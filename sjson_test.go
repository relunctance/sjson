package sjson

import (
	"io/ioutil"
	"testing"

	"github.com/relunctance/goutils/dump"
)

func TestMapJson(t *testing.T) {
	json := `[
					{
						"b1": [
							{
								"c": {"gaoqilin":"cde"}
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
	]`
	fields := []string{
		"#.b1.#.c",
	}
	data, err := getByBytes([]byte(json), fields)
	if err != nil {
		panic(err)
	}
	dump.Println(string(data))
}

func TestSliceJson(t *testing.T) {
	json := `[
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
	]`
	fields := []string{
		"#.b1.#.c",
	}

	data, err := getByBytes([]byte(json), fields)

	if err != nil {
		panic(err)
	}
	dump.Println(string(data))
}

func TestMutliJson(t *testing.T) {
	json, err := ioutil.ReadFile("./a.json")
	if err != nil {
		panic(err)
	}
	fields := []string{
		"name.first",
		"lname",
		"age",
		"friends.#.last",
	}

	data, err := getByBytes([]byte(json), fields)
	if err != nil {
		panic(err)
	}
	dump.Println(string(data))
}

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
