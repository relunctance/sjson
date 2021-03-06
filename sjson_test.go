package sjson

import (
	"testing"
)

const (
	demojson = `{
	"ccc":false,
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

func TestGetPaths(t *testing.T) {
	j := newJson([]byte(demojson))
	if !j.IsObject("name") {
		t.Fatalf("should be true")
	}
	if ok := j.IsArray("children"); ok != true {
		t.Fatalf("should be true")
	}

	// not exists key 'childrenxxxxxx'
	if ok := j.IsArray("childrenxxxxxx"); ok != false {
		t.Fatalf("should be false")
	}
	if ok := j.IsObject("children"); ok != false {
		t.Fatalf("should be false")
	}

	if ok := j.IsString("name.first"); ok != true {
		t.Fatalf("should be false")
	}
	if ok := j.IsNumber("age"); ok != true {
		t.Fatalf("should be true")
	}
	if !j.IsBool("ccc") {
		t.Fatalf("should be true")

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
		t.Fatalf("should be == `%s` , but not is `%s`", expectJson, string(data))

	}
}

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
			},
			"d" : [
				{
					"name" : "v1",
					"age" : 18,
					"pass": "abc",
					"eee" : {
						"xx" : {
							"name":"gql" 
						},
						"ff" : {
							"name":"gql" 
						},
						"gg" : {
							"name":"gql" 
						}
					}
				},
				{
					"name" : "v2",
					"age" : 19,
					"pass": "def",
					"eee" : {
						"cc" : {
							"name":"gql"
						}
					}
				}
			]
		}
	}`
	/*
	 */
	fields := []string{
		"data.*.name",
		"data.d.#.eee.*.name",
	}
	data, err := getByBytes([]byte(json), fields)
	if err != nil {
		panic(err)
	}
	expectJson := `{"data":{"a":{"name":"v1"},"b":{"name":"v2"},"c":{"name":"v3"},"d":[{"eee":{"ff":{"name":"gql"},"gg":{"name":"gql"},"xx":{"name":"gql"}}},{"eee":{"cc":{"name":"gql"}}}]}}`
	if string(data) != expectJson {
		t.Fatalf("should be == `%s` but not is `%s`", expectJson, string(data))
	}
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
