package sjson

import (
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

func GetKeys(json []byte, path string) ([]string, error) {
	j, err := gabs.ParseJSON(json)
	if err != nil {
		return nil, err
	}
	ms := j.ChildrenMap()
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
	proc := newProcess(json)
	for _, path := range paths {
		if proc.j.IsCommonPath(path) {
			proc.commonPathGet(path)
		} else {
			if strings.Index(path, "#") != -1 {
				proc.numbersignPathGet(path)
			}
		}
	}
	return proc.j.json.Search(PREFIX).Bytes(), nil
}

type process struct {
	json []byte
	j    *Json
}

func newProcess(json []byte) *process {
	p := &process{}
	p.json = redefineJson(json)
	p.j = NewJson()
	return p
}

func (p *process) gJsonResult(path string) gjson.Result {
	return gjson.GetBytes(p.json, path)
}

func (p *process) commonPathGet(path string) error {
	path = redefinePath(path)
	result := p.gJsonResult(path)
	_, err := p.j.json.SetP(result.Value(), path)
	return err
}

func (p *process) wildcardPathGet(path string) error {
	return nil
}

func (p *process) numbersignPathGet(path string) error {
	path = redefinePath(path)
	result := p.gJsonResult(path)
	return p.j.setPath(path, result)
}
