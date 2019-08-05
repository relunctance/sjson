package sjson

import (
	//	sj "github.com/guyannanfei25/go-simplejson"

	"strings"

	"github.com/Jeffail/gabs"
	"github.com/tidwall/gjson"
)

type Json struct {
	json    *gabs.Container
	s       []string
	finishp string
}

func NewJson() *Json {
	j := &Json{}
	j.json = gabs.New()
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

func getByBytes(json []byte, paths []string) ([]byte, error) {
	j := NewJson()

	js, err := gabs.ParseJSONFile("./a.json")
	if err != nil {
		panic(err)
	}
	js.Path("")
	for _, path := range paths {
		result := gjson.GetBytes(json, path)
		if j.IsCommonPath(path) {
			j.json.SetP(result.Value(), path)
		} else {
			if strings.Index(path, "#") != -1 { // 处理#号情况

				err := j.SetPath(path, result)
				if err != nil {
					return nil, err
				}
			}
		}

	}
	return j.Bytes(), nil
}