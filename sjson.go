package sjson

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/tidwall/gjson"
	vsj "github.com/tidwall/sjson"
)

type Json struct {
	json          *gabs.Container
	data          []byte
	wildcardPaths []string
	s             []string
	container     *gabs.Container
	finishp       string
	sjdata        string
	sj            string
}

func NewJson(json []byte) *Json {
	j := &Json{}
	j.json = gabs.New()

	j.data = redefineJson(json)
	j.sjdata = string(j.data)
	var err error
	j.container, err = gabs.ParseJSON(j.data)
	if err != nil {
		panic(err)
	}
	j.sj = `{}`
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
			j.pathGabsSet(path)
			continue
			//j.commonPathGet(path)
			//continue
		}
		/*
			if j.isAllNumbersign(path) {
				j.numbersignPathGet(path)
				continue
			}
		*/
		j.wildcardPathGet(path)
	}
	j.json, _ = gabs.ParseJSON([]byte(j.sj))
	return j.json.Search(PREFIX).Bytes(), nil
}

/*
[
	"__sjson__.data.d.0.eee",
	"__sjson__.data.d.1.eee",
]
=>
[
	"__sjson__.data.d.0.eee.xx",
	"__sjson__.data.d.0.eee.ff",
	"__sjson__.data.d.0.eee.gg",
	"__sjson__.data.d.1.eee.cc",
]
*/

func (j *Json) initWildcardPaths(path, c string) (s []string, err error) {
	path = strings.TrimRight(path, ".")
	switch c {
	case "#":
		num := gjson.Get(string(j.data), path+".#").Int()
		l := int(num)
		s = make([]string, 0, l)
		for i := 0; i < l; i++ {
			s = append(s, fmt.Sprintf("%s.%d", path, i))
		}
	case "*":
		keys, err := j.findMapKeys(path)
		if err != nil {
			return nil, err
		}
		s = make([]string, 0, len(keys))
		for _, key := range keys {
			s = append(s, fmt.Sprintf("%s.%s", path, key))
		}

	}
	return s, nil
}

func (j *Json) wildcardPathGet(path string) (err error) {
	ps := strings.Split(path, ".")
	var line string
	for _, p := range ps {
		if j.isNotWildcard(p) {
			if len(j.wildcardPaths) == 0 {
				line += p + "."
			} else {
				for key, pth := range j.wildcardPaths {
					j.wildcardPaths[key] = pth + "." + p
				}
			}
			continue
		}
		if len(j.wildcardPaths) == 0 {
			j.wildcardPaths, err = j.initWildcardPaths(line, p)
			if err != nil {
				return err
			}
			continue
		}

		var newPaths []string
		newWildcardPs := make([]string, 0, 1)
		for _, pth := range j.wildcardPaths {
			newPaths, err = j.initWildcardPaths(pth, p)
			if err != nil {
				return err
			}
			newWildcardPs = append(newWildcardPs, newPaths...)
		}
		j.wildcardPaths = newWildcardPs
	}
	for _, path := range j.wildcardPaths {
		j.pathGabsSet(path)
	}
	j.wildcardPaths = make([]string, 0, 1) // init
	return nil
}

func (j *Json) isPathNil(path string) bool {
	return gjson.Get(j.sjdata, path).Value() == nil
}

func (j *Json) pathGabsSet(path string) {
	if j.isPathNil(path) {
		return
	}
	value := j.container.Path(path).Data()
	var err error
	j.sj, err = vsj.Set(j.sj, path, value)
	if err != nil {
		panic(err)
	}
	return
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
