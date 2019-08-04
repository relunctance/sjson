package sjson

import (
	"fmt"
	"strings"

	"github.com/relunctance/goutils/dump"
	"github.com/tidwall/gjson"
)

const (
	slice_char = "-"
)

// TODO 结尾不可以是#
func (j *Json) checkPath(path string) error {
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
	sp := strings.TrimRight(p, slice_char)
	sp = strings.TrimRight(sp, ".")
	if strings.Index(sp, slice_char) != -1 {
		panic("should not be exists '-'")
	}
	j.json.SetP([]interface{}{}, sp)
}
func (j *Json) replaceSliceCharToIndex(p string, i int) string {
	return fmt.Sprintf(strings.Replace(p, slice_char, "%d", -1), i)
}
func (j *Json) buildSlice(p string, n int) {
	dump.Println("n:", n)
	if j.finishp == "" {
		j.initSlice(p)
	} else {
		for _, line := range j.s {
			newP := strings.Replace(p, j.finishp, line, -1)
			j.initSlice(newP)
		}
	}

	for i := 0; i < n; i++ {
		if j.finishp == "" {
			j.s = append(j.s, j.replaceSliceCharToIndex(p, i))
			//dump.Println("jsssss 00:", j.s) // jsssss 00: [a.0 a.1 a.2]
			j.json.SetP(map[string]interface{}{}, p)
		} else {
			for key, line := range j.s {
				newP := strings.Replace(p, j.finishp, line, -1) // newP: a.1.b1.- n: 1
				j.json.SetP(map[string]interface{}{}, newP)
				j.s[key] = j.replaceSliceCharToIndex(newP, i)
			}
			//dump.Println("jsssss:", j.s)	// jsssss: [a.0.b1.0, a.1.b1.0,a.2.b1.0]
			//os.Exit(1)
		}
	}
	j.finishp = p
	for _, line := range j.s {
		j.json.SetP(map[string]interface{}{}, line)
	}
	dump.Println("j.s:", j.s)
}

func (j *Json) setSlice(p string, ret gjson.Result) {

}
func (j *Json) SetPath(path string, ret gjson.Result) error {
	if err := j.checkPath(path); err != nil {
		return err
	}
	ps := j.buildPaths(path)
	for _, p := range ps {
		fmt.Println("sssp:", p)
		n := strings.Count(p, slice_char) - 1
		if IsEndChar(p, slice_char) {
			num := InterfaceSliceLength(ret.Value(), n)
			j.buildSlice(p, num) // has finish
		} else {
			j.setSlice(p, ret) // todo
		}
	}
	/*
		l := len(ps)
		var nextP string
		arr := ret.Array()
		dump.Println("l:", len(arr))
		var nextL int
		for i, p := range ps {
			if i == l-1 { // 这里是end
			}

			nextP = strings.TrimRight(nextP, ".")
			fmt.Println("nextP:", nextP)
			for k := 0; k < nextL; k++ {

			}

			if i == 0 {
				j.json.SetP(make([]interface{}, len(arr)), ps[0])
			}
			nextP = nextP + p + ".-."
		}
	*/
	return nil
}
