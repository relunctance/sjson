# SJSON
Go json select with field
* Sjson can used Precise API field control
* Sjson can used Authority control



Getting Started
===============


Installing
----------

To start using DJSON, install Go and run `go get`:

```sh
go get -u -v github.com/relunctance/sjson
```

This will retrieve the library.


Example
----------

```go
 json := `{
		"data":	{
			"a" : {
				"name" : "v1",
				"age" : 18,
				"pass": "abc"
			},
			"d" : [
				{
					"name" : "v1",
					"age" : 18,
					"pass": "abc",
					"eee" : {
						"xx" : { "name":"gql" },
						"ff" : { "name":"gql" },
						"gg" : { "name":"gql" }
					}
				},
				{
					"name" : "v2",
					"age" : 19,
					"pass": "def",
					"eee" : { 
                        "cc" : { "name":"gql" } 
                    }
				}
			]
		}
	}`
    // with slice path , return string
    newStrJson , err := sjson.Select(json , []string{"data.*.name","data.d.#.eee.*.name"});
    // alias sjson.Select()
    newStrJson , err := sjson.SelectString(json , []string{"data.*.name","data.d.#.eee.*.name"});

    // with string path , return string
    newStrJson , err := sjson.SelectPath(json , "data.*.name");

    // use byte json to select
    newByteJson , err := sjson.SelectBytes([]byte(json) , []string{"data.*.name","data.d.#.eee.*.name"}));

    // use byte json with single path
    newByteJson , err := sjson.SelectBytesPath([]byte(json), "data.*.name")


```



```go
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
```

```go
	json := []byte(demojson)
	fields := []string{
		"name.first",
		"lname",
		"age",
		"friends.#.last",
	}

	data, _ := sjson.SelectBytes([]byte(json), fields)
	//data = `{"age":37,"friends":[{"last":"Murphy"},{"last":"Craig"},{"last":"Murphy"},{"last":"Craig"}],"lname":{"first":"Tom1","last":"Anderson1"},"name":{"first":"Tom"}}`
```

```go
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
	fields := []string{
		"data.*.name",
		"data.d.#.eee.*.name",
	}
	data, _ := sjson.SelectBytes([]byte(json), fields)
	//data = `{"data":{"a":{"name":"v1"},"b":{"name":"v2"},"c":{"name":"v3"},"d":[{"eee":{"ff":{"name":"gql"},"gg":{"name":"gql"},"xx":{"name":"gql"}}},{"eee":{"cc":{"name":"gql"}}}]}}`
```

```go
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
	data, _ := sjson.SelectBytes([]byte(json), fields)
	//data = `[{"b1":[{"c":{"gql":"cde"}}]},{"b1":[{"c":222}]},{"b1":[{"c":333}]}]`

```

```go
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
	data, _ := sjson.SelectBytes([]byte(json), fields)
	//data = `[{"b1":[{"c":111}]},{"b1":[{"c":222}]},{"b1":[{"c":333}]}]`

```

```go
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

	data, _ := sjson.SelectBytes([]byte(json), fields)
	// `{"a":[{"b1":[{"c":111}]},{"b1":[{"c":222}]},{"b1":[{"c":333}]}]}`
```



Thanks
----------

[github.com/Jeffail/gabs](https://github.com/Jeffail/gabs)
