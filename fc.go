package sjson

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

func gjsonLength(r gjson.Result) int {
	v := r.Value()
	switch val := v.(type) {
	case []interface{}:
		return len(val)
	case map[string]interface{}:
		return len(val)
	}
	return 0
}

func IsEndChar(v, char string) bool {
	return strings.LastIndex(v, char) == len(v)-1
}

func InterfaceSliceLength(v interface{}, n int) int {
	if n < 0 {
		panic("n should be >= 0")
	}
	varr := v.([]interface{})
	if n == 0 {
		return len(varr)
	}
	return InterfaceSliceLength(varr[0], n-1)
}

func IsPrefixSlice(path string) bool {
	return path[0] == '#'
}

func redefinePath(path string) string {
	return prefix + "." + path
}

func redefineJson(json []byte) []byte {
	js := fmt.Sprintf(`{"%s":%s}`, prefix, string(json))
	return []byte(js)

}
