package sjson

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/tidwall/gjson"
)

func (j *Json) SetPathSlice(path string, ret gjson.Result) {
	lp := newLpath(path)
	fmt.Println("lp:", lp)
	ret.ForEach(func(key, value gjson.Result) bool {
		kbool := key.String() == ""

		fmt.Printf("key:%v , kbool:%v  , ktype:%d  , valType:%d, val:%v\n", key.String(), kbool, key.Type, value.Type, value.String())

		if value.IsArray() {
			//arr := value.Array()
			if !j.json.ExistsP(lp.left) {
				j.json.ArrayP(lp.left)
			}
			j.json.ArrayAppendP(j.buildRight(lp.right, value), lp.left)
		}

		if value.Type == gjson.JSON {

			//j.SetPathSlice()
		}
		return true // keep iterating
	})

}
func (j *Json) getValue(ret gjson.Result) (v interface{}) {

	switch ret.Type {
	case
		gjson.Null,
		gjson.False,
		gjson.Number,
		gjson.String,
		gjson.True:
		v = ret.Value()
	case gjson.JSON:
		length := gjsonLength(ret)
		if ret.IsArray() {
			vj := make([]interface{}, 0, length)
			ret.ForEach(func(key, value gjson.Result) bool {
				vj = append(vj, value.Value())
				return true
			})
			v = vj
		} else if ret.IsObject() {
			vj := make(map[interface{}]interface{}, length)
			ret.ForEach(func(key, value gjson.Result) bool {
				k := key.String()
				vj[k] = value.Value()
				return true
			})
			v = vj
		}
	default:
		panic(fmt.Errorf("unknow type"))
	}
	return
}
func (j *Json) buildRight(path string, ret gjson.Result) (v interface{}) {
	lp := newLpath(path)
	vj := gabs.New()
	vj.Set(j.getValue(ret), lp.left)
	return vj
	/*
		ret.ForEach(func(key, value gjson.Result) bool {

		})
		return j.buildRight(lp.right)
	*/

}
