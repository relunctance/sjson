package sjson

import (
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
