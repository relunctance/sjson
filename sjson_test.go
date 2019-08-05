package sjson

import (
	"fmt"
	"testing"
)

const (
	demojson = `{
    "name": {"first": "Tom", "last": "Anderson"},
    "lname": {"first": "Tom1", "last": "Anderson1"},
    "age":37,
    "children": ["Sara","Alex","Jack"],
    "fav.movie": "Deer Hunter",
    "friends": [
        {"last": "Murphy"},
        {"first": "Roger", "last": "Craig","name":"abc"},
        {"first": "James1", "last": "Murphy"},
        {
            "first": "Roger2",
            "last": "Craig",
            "name":"abc",
            "c":[

                {
                    "kw": {
                        "signature": [
                            "c1"
                        ]
                    },
                    "tabname": "tab1"
               },
                {
                    "kw": {
                        "signature": [
                            "c2"
                        ]
                    },
                    "tabname": "tab2"
               },
                {
                    "kw": {
                        "signature": [
                            "c3"
                        ]
                    },
                    "tabname": "tab3"
               }
            ]
        }
    ]
	}`
)

func TestStar(t *testing.T) {
	json := `{
		"data":	{
			"a" : {
				"name" : "v1",
				"age" : 18,
				"pass": "abc"
			},
			"b" : {
				"name" : "v2",
				"age" : 19,
				"pass": "abc1234"
			},
			"c" : {
				"name" : "v3",
				"age" : 18,
				"pass": "defg23423"
			}
		}
	}`
	fields := []string{
		"data.*.name",
	}
	data, err := getByBytes([]byte(json), fields)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func TestMapJson(t *testing.T) {
	json := `[
					{
						"b1": [
							{
								"c": {"gql":"cde"},
								"d": {"gql":"cde"},
								"e": {"gql":"cde"}
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
	expectJson := `[{"b1":[{"c":{"gql":"cde"}}]},{"b1":[{"c":222}]},{"b1":[{"c":333}]}]`
	if string(data) != expectJson {
		t.Fatalf("should be == `%s`", expectJson)

	}
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
	expectJson := `[{"b1":[{"c":111}]},{"b1":[{"c":222}]},{"b1":[{"c":333}]}]`
	if string(data) != expectJson {
		t.Fatalf("should be == `%s`", expectJson)

	}
}

func TestMutliJson(t *testing.T) {
	json := []byte(demojson)
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
	expectJson := `{"age":37,"friends":[{"last":"Murphy"},{"last":"Craig"},{"last":"Murphy"},{"last":"Craig"}],"lname":{"first":"Tom1","last":"Anderson1"},"name":{"first":"Tom"}}`
	if string(data) != expectJson {
		t.Fatalf("should be == `%s`", expectJson)

	}
}

func TestDemo(t *testing.T) {
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
	expectJson := `{"a":[{"b1":[{"c":111}]},{"b1":[{"c":222}]},{"b1":[{"c":333}]}]}`
	if string(data) != expectJson {
		t.Fatalf("should be == `%s`", expectJson)

	}
}
