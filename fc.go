package sjson

import (
	"fmt"
	"strings"
)

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
	return PREFIX + "." + path
}

func redefineJson(json []byte) []byte {
	js := fmt.Sprintf(`{"%s":%s}`, PREFIX, string(json))
	return []byte(js)

}

func copySlice(arr, arr2 []string) []string {
	if len(arr2) > len(arr) {
		return arr2
	}
	copy(arr, arr2)
	return arr
}
