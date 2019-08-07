package sjson

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/relunctance/goutils/dump"
	"github.com/tidwall/gjson"
)

type Json struct {
	json          *gabs.Container
	data          []byte
	wildcardPaths []string
	s             []string
	container     *gabs.Container
	finishp       string
}

func NewJson(json []byte) *Json {
	j := &Json{}
	j.json = gabs.New()

	j.data = redefineJson(json)
	var err error
	j.container, err = gabs.ParseJSON(j.data)
	if err != nil {
		panic(err)
	}
	j.s = make([]string, 0, 1)
	j.wildcardPaths = make([]string, 0, 1)
	return j

}

func (j *Json) Bytes() []byte {
	return j.json.Bytes()
}

func (j *Json) String() string {
	return j.json.String()
}

func (j *Json) IsCommonPath(path string) bool {
	if strings.Index(path, "#") != -1 {
		return false
	}
	if strings.Index(path, "*") != -1 {
		return false
	}
	return true

}

func (j *Json) findMapKeys(path string) ([]string, error) {
	if !j.IsCommonPath(path) {
		return nil, fmt.Errorf("path should be common")
	}
	ms := j.container.Path(path).ChildrenMap()
	ret := make([]string, 0, len(ms))
	for key, _ := range ms {
		ret = append(ret, key)
	}
	return ret, nil

}

// data.#.a.*.name
// data.*.name
// data.*.name.#.c
func getByBytes(json []byte, paths []string) ([]byte, error) {
	j := NewJson(json)
	if j.checkIsAll(paths) {
		return json, nil
	}
	for _, path := range paths {
		path = redefinePath(path)
		if j.IsCommonPath(path) {
			j.commonPathGet(path)
			continue
		}
		if j.isAllNumbersign(path) {
			j.numbersignPathGet(path)
			continue
		}
		//j.tmpdata = j.data
		j.wildcardPathGet(path)
	}
	return j.json.Search(PREFIX).Bytes(), nil
}

func (j *Json) initWildcardPaths(data []byte, path, c string) (s []string, err error) {
	s = make([]string, 0, 1)
	path = strings.TrimRight(path, ".")
	switch c {
	case "#":
		num := gjson.Get(string(data), path+".#").Int()
		for i := 0; i < int(num); i++ {
			s = append(s, fmt.Sprintf("%s.%d", path, i))
		}
	case "*":
		keys, err := j.findMapKeys(path)
		if err != nil {
			return s, err
		}
		for _, key := range keys {
			s = append(s, fmt.Sprintf("%s.%s", path, key))
		}
	}
	return s, nil
}

// [1]data.#.a.*.name
// [2]data.*.name
// [3]data.*.name.#.c
func (j *Json) wildcardPathGet(path string) (err error) {
	ps := strings.Split(path, ".")
	var line string
	for _, p := range ps {
		if j.isNotWildcard(p) && len(j.wildcardPaths) == 0 {
			line += p + "."
			continue
		}
		if len(j.wildcardPaths) == 0 {
			j.wildcardPaths, err = j.initWildcardPaths(j.data, line, p)
			if err != nil {
				return err
			}
			continue
		}
		if j.isNotWildcard(p) {
			// 进入到这里  如果是[2] 则p == name
			// 进入到这里  如果是[3] 则p == name
			for key, pth := range j.wildcardPaths {
				j.wildcardPaths[key] = pth + "." + p
			}
		} else {
			for _, pth := range j.wildcardPaths {
				j.wildcardPaths, err = j.initWildcardPaths(j.data, pth, p)
				if err != nil {
					return err
				}
			}
		}
	}
	dump.Println("xxxx:", j.wildcardPaths)
	return nil
}

func (j *Json) isAllNumbersign(path string) bool {
	return strings.Index(path, "#") != -1 && strings.Index(path, "*") == -1
}

func (j *Json) isNotWildcard(char string) bool {
	switch char {
	case
		"#",
		"*":
		return false
	}
	return true
}
func (j *Json) checkIsAll(paths []string) bool {
	for _, p := range paths {
		if p == "*" {
			return true
		}
	}
	return false
}

func (j *Json) gJsonResult(path string) gjson.Result {
	return gjson.GetBytes(j.data, path)
}

func (j *Json) commonPathGet(path string) error {
	result := j.gJsonResult(path)
	_, err := j.json.SetP(result.Value(), path)
	return err
}

func (j *Json) numbersignPathGet(path string) error {
	result := j.gJsonResult(path)
	return j.setPath(path, result)
}
