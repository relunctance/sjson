package sjson

import (
	"testing"

	"github.com/relunctance/goutils/dump"
)

func TestDemo(t *testing.T) {
	/*
		json, err := ioutil.ReadFile("./a.json")
		if err != nil {
			panic(err)
		}
		fields := []string{
			"name.first",
			"age",
			"friends.#.last",
		}
	*/

	json := `{
				"a": [
					{
						"b1": [
							{
								"c": 111
							}
						]
					}, 
					{
						"b1": [
							{
							"c": 222
							}
						]
					}, 
					{
						"b1": [
							{
							"c": 333
							}
						]
					}
				]
}`
	fields := []string{
		"a.#.b1.#.c",
	}
	data, err := getByBytes([]byte(json), fields)
	if err != nil {
		panic(err)
	}
	dump.Println(string(data))
}
