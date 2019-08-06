package sjson

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/tidwall/gjson"
)

type Json struct {
	json    *gabs.Container
	data    []byte
	s       []string
	finishp string
}

func NewJson(json []byte) *Json {
	j := &Json{}
	j.json = gabs.New()
	j.data = redefineJson(json)
	j.s = make([]string, 0, 1)
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

func (j *Json) GetKeys(json []byte, path string) ([]string, error) {
	if !j.IsCommonPath(path) {
		return nil, fmt.Errorf("path should be common")
	}

	js, err := gabs.ParseJSON(json)
	if err != nil {
		return nil, err
	}
	ms := js.ChildrenMap()
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
	for _, path := range paths {
		if j.IsCommonPath(path) {
			j.commonPathGet(path)
		} else {
			if strings.Index(path, "#") != -1 {
				j.numbersignPathGet(path)
			}
		}
	}
	return j.json.Search(PREFIX).Bytes(), nil
}

func (j *Json) gJsonResult(path string) gjson.Result {
	return gjson.GetBytes(j.data, path)
}

func (j *Json) commonPathGet(path string) error {
	path = redefinePath(path)
	result := j.gJsonResult(path)
	_, err := j.json.SetP(result.Value(), path)
	return err
}

func (j *Json) wildcardPathGet(path string) error {
	return nil
}

func (j *Json) numbersignPathGet(path string) error {
	path = redefinePath(path)
	result := j.gJsonResult(path)
	return j.setPath(path, result)
}
