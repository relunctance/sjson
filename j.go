package sjson

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	slice_char = "-"
	prefix     = "__sjson__"
)

func (j *Json) checkPath(path string) error {
	if path[len(path)-1] == '#' {
		return fmt.Errorf("path end char can not be '#'")
	}
	if strings.Index(path, "-") != -1 {
		return fmt.Errorf("can not allow key exists '-' , check path:[%s]", path)
	}
	return nil
}

func (j *Json) buildPaths(path string) []string {
	path = strings.Replace(path, "#", "-", -1)
	ps := strings.Split(path, "-")
	for k, p := range ps {
		ps[k] = strings.Trim(p, ".")
	}

	return j.createGabsPath(ps)
}

func (j *Json) createGabsPath(ps []string) []string {
	l := len(ps)
	data := make([]string, 0, l)
	var line string
	for i, p := range ps {
		if i < l-1 {
			line = line + "." + p + "." + slice_char
		}
		if i == l-1 {
			line += "." + p
		}
		line = strings.Trim(line, ".")
		data = append(data, line)

	}
	return data
}

func (j *Json) initSlice(p string) {
	sp := strings.TrimRight(strings.TrimRight(p, slice_char), ".")
	if strings.Index(sp, slice_char) != -1 {
		panic("should not be exists '-'")
	}
	j.json.SetP([]interface{}{}, sp)
}

func (j *Json) replaceSliceCharToIndex(p string, i int) string {
	return fmt.Sprintf(strings.Replace(p, slice_char, "%d", -1), i)
}

func (j *Json) renewPath(p, line string) string {
	return strings.Replace(p, j.finishp, line, -1)
}

func (j *Json) buildSlice(p string, n int) {
	if j.finishp == "" {
		j.initSlice(p)
	} else {
		for _, line := range j.s {
			newP := j.renewPath(p, line)
			j.initSlice(newP)
		}
	}

	for i := 0; i < n; i++ {
		if j.finishp == "" {
			j.s = append(j.s, j.replaceSliceCharToIndex(p, i))
			j.json.SetP(map[string]interface{}{}, p)
		} else {
			for key, line := range j.s {
				newP := j.renewPath(p, line)
				j.json.SetP(map[string]interface{}{}, newP)
				j.s[key] = j.replaceSliceCharToIndex(newP, i)
			}
		}
	}
	j.finishp = p
	for _, line := range j.s {
		j.json.SetP(map[string]interface{}{}, line)
	}
}

func (j *Json) getValueWithPath(p, newp string, ret gjson.Result) interface{} {
	pslice := strings.Split(p, ".")
	nslice := strings.Split(newp, ".")
	data := ret.Array()
	var item gjson.Result
	for key, path := range pslice {
		if path != "-" {
			continue
		}
		k := nslice[key]
		index, _ := strconv.Atoi(k)
		item = data[index]
	}
	return item.Value()
}

func (j *Json) setSlice(p string, ret gjson.Result) {
	for _, line := range j.s {
		newp := j.renewPath(p, line)
		v := j.getValueWithPath(p, newp, ret)
		j.json.SetP(v, newp)
	}
}

func (j *Json) setPath(path string, ret gjson.Result) error {
	if err := j.checkPath(path); err != nil {
		return err
	}
	ps := j.buildPaths(path)
	for _, p := range ps {
		n := strings.Count(p, slice_char) - 1
		if IsEndChar(p, slice_char) {
			num := InterfaceSliceLength(ret.Value(), n)
			j.buildSlice(p, num)
		} else {
			j.setSlice(p, ret)
		}
	}
	return nil
}
