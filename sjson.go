package sjson

import (
	"fmt"
	"strings"

	// "github.com/Jeffail/gabs" // default import
	"github.com/Jeffail/gabs/v2" // support go module when go version >= go.11 and has open GO111MODULE="on"
	"github.com/tidwall/gjson"
	vsj "github.com/tidwall/sjson"
)

const (
	PREFIX = "__sjson__"
)

type Json struct {
	data          []byte
	wildcardPaths []string
	container     *gabs.Container
	sjdata        string
	sj            string
}

func redefinePath(path string) string {
	return PREFIX + "." + path
}

func redefineJson(json []byte) []byte {
	js := fmt.Sprintf(`{"%s":%s}`, PREFIX, string(json))
	return []byte(js)

}

func newJson(json []byte) *Json {
	j := &Json{}
	j.data = redefineJson(json)
	j.sjdata = string(j.data)
	var err error
	j.container, err = gabs.ParseJSON(j.data)
	if err != nil {
		panic(err)
	}
	j.sj = `{}`
	j.wildcardPaths = make([]string, 0, 1)
	return j

}

func (j *Json) addPrefix(path string) string {
	return PREFIX + "." + path
}

func (j *Json) IsArray(path string) bool {
	_, ok := j.pathData(path).([]interface{})
	return ok
}

func (j *Json) IsObject(path string) bool {
	_, ok := j.pathData(path).(map[string]interface{})
	return ok
}

func (j *Json) pathData(path string) interface{} {
	return j.container.Path(j.addPrefix(path)).Data()
}

func (j *Json) IsString(path string) bool {
	v := j.pathData(path)
	_, ok := v.(string)
	return ok
}

func (j *Json) IsFloat(path string) bool {
	v := j.pathData(path)
	_, ok := v.(float64)
	return ok

}

func (j *Json) IsNumber(path string) bool {
	v := j.pathData(path)
	_, ok := v.(float64)
	return ok
}

func (j *Json) IsScalar(path string) bool {
	if j.IsArray(path) || j.IsObject(path) {
		return false
	}
	return true
}

func (j *Json) IsBool(path string) bool {
	v := j.pathData(path)
	_, ok := v.(bool)
	return ok
}

// func (j *Json) GetPaths(path string) []string {
// ps := make([]string, 0, 10)
// return ps
// }

// check path is contain '#' or '*'
// the '#' is used in path that expression [{item1} , {item2} , {itemN...} ] which index of "0,1,2 ..."
// the '*' is used in path that expression {"name1":{item1} , "name2":{item2}  , "nameN":{itemN}} which index of "name1,name2,nameN..."
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

func getByBytes(json []byte, paths []string) ([]byte, error) {
	j := newJson(json)
	if j.checkIsAll(paths) {
		return json, nil
	}
	//
	for _, path := range paths {
		path = redefinePath(path)
		if j.IsCommonPath(path) {
			j.pathGabsSet(path)
			continue
		}
		j.wildcardPathGet(path)
	}
	gjson, _ := gabs.ParseJSON([]byte(j.sj))
	return gjson.Search(PREFIX).Bytes(), nil
}

func (j *Json) buildWildcardPaths(path, c string) (s []string, err error) {
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
			j.wildcardPaths, err = j.buildWildcardPaths(line, p)
			if err != nil {
				return err
			}
			continue
		}

		var newPaths []string
		newWildcardPs := make([]string, 0, 1)
		for _, pth := range j.wildcardPaths {
			newPaths, err = j.buildWildcardPaths(pth, p)
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
	j.resetWildcardPaths()
	return nil
}

// init wildcardPaths to used next path
func (j *Json) resetWildcardPaths() {
	j.wildcardPaths = make([]string, 0, 1) // init
}

// check get path value is empty
func (j *Json) isPathNil(path string) bool {
	return gjson.Get(j.sjdata, path).Value() == nil
}

//
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

// check split char is '#' and '*'
func (j *Json) isNotWildcard(char string) bool {
	switch char {
	case
		"#",
		"*":
		return false
	}
	return true
}

// if you input a slice path and one value is '*' then return origin json
// that is mean the path should be exists field
// the all data return example  :
//  ["*"]
//  ["path1","*","path2"]
func (j *Json) checkIsAll(paths []string) bool {
	for _, p := range paths {
		if p == "*" {
			return true
		}
	}
	return false
}
