package sjson

import (
	//	sj "github.com/guyannanfei25/go-simplejson"

	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/relunctance/goutils/dump"
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

func runMmap(mmap map[string]interface{}, ps []string, i *int, ret gjson.Result) interface{} {
	key := ps[*i]
	next := *i + 1
	if *i == len(ps)-1 { // 最后一个值设置
		fmt.Printf("endkey:%v\n", key)
		mmap[key] = ret.Value()
		return mmap[key]
	}

	nkey := ps[next]
	if nkey == "#" {
		l := gjsonLength(ret)
		tmpSlice := make([]interface{}, 0, l) // 也是一个map
		for n := 0; n < l; n++ {
			tmpSlice = append(tmpSlice, map[string]interface{}{})
		}
		mmap[key] = tmpSlice
		*i += 1 // 游标再次位移
	}
	fmt.Printf("key:%v , nkey:%v\n", key, nkey)

	if mmap[key] == nil {
		mmap[key] = map[string]interface{}{}

	}
	return mmap[key] // 重点:  vm 被指向了下一级
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
		}
		dump.Println(path, result.String())

		if strings.Index(path, "#") != -1 { // 处理#号情况

			err := j.SetPath(path, result)
			if err != nil {
				return nil, err
			}
		}

	}
	return j.Bytes(), nil
}
