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


Thanks
----------

[github.com/Jeffail/gabs](https://github.com/Jeffail/gabs)