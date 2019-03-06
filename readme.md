# Go JSON Logic

[![Build Status](https://travis-ci.org/diegoholiveira/jsonlogic.svg)](https://travis-ci.org/diegoholiveira/jsonlogic)
[![codecov](https://codecov.io/gh/diegoholiveira/jsonlogic/branch/master/graph/badge.svg)](https://codecov.io/gh/diegoholiveira/jsonlogic)

Implementation of [JSON Logic](http://jsonlogic.com) in Go Lang.


## What's JSON Logic?

JSON Logic is a DSL to write logic decisions in JSON. It's has a great specification and is very simple to learn.
The official website has a great documentation with examples: http://jsonlogic.com


## How to use it

The use of this library is very straightforward. Here's a simple example:


```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/diegoholiveira/jsonlogic"
)

func main() {
	var logic interface{}

	err := json.Unmarshal([]byte(`{"==": [1, 1]}`), &logic)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	var result interface{}

	jsonlogic.Apply(
		logic,
		nil,
		&result,
	)

	fmt.Println(result)
}
```

This will output `true` in your console.

Here's another example, but this time using variables passed in the `data` parameter:


```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/diegoholiveira/jsonlogic"
)

func main() {
	var logic interface{}
	var data interface{}
	var result interface{}

	json.Unmarshal([]byte(`{
		"filter": [
			{"var": "users"},
			{">=": [
				{"var": ".age"},
				18
			]}
		]
	}`), &logic)

	json.Unmarshal([]byte(`{
		"users": [
			{"name": "Diego", "age": 33, "location": "Florianópolis"},
			{"name": "Jack", "age": 12, "location": "London"},
			{"name": "Pedro", "age": 19, "location": "Lisbon"},
			{"name": "Leopoldina", "age": 30, "location": "Rio de Janeiro"}
		]
	}`), &data)

	err := jsonlogic.Apply(logic, data, &result)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println("Users older than 18:")
	for _, _user := range result.([]interface{}) {
		user := _user.(map[string]interface{})

		fmt.Printf("    - %s\n", user["name"].(string))
	}
}
```

## Limitations

The `Apply` function have three params as input:

- the first one is the logic to be executed;
- next you have the data that can be used by the logic;
- and the last one is the variable to store the result.

The type of those params must be `interface{}` to be easy to reflect of it and use it.
Also, all values passed in any of the two input params must be one of those:

    bool for JSON booleans,
    float64 for JSON numbers,
    string for JSON strings, and
    nil for JSON null.

This is the same values that `encoding/json` work with.

Here's an example of an invalid data:

```go
func main() {
	var rules interface{}
	var result interface{}

	json.Unmarshal([]byte(`{
		"filter": [
			{"var": "users"},
			{">=": [
				{"var": ".age"},
				18
			]}
		]
	}`), &rules)

	data := interface{}(map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"name":     string("Diego"),
				"age":      int(33),
				"location": string("Florianópolis"),
			},
		},
	})

	err := jsonlogic.Apply(rules, data, &result)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println("Users older than 18:")
	for _, _user := range result.([]interface{}) {
		user := _user.(map[string]interface{})

		fmt.Printf("    - %s\n", user["name"].(string))
	}
}
```

This will produce this error: `panic: interface conversion: interface {} is int, not float64`.
So, to avoid this error, make sure to always work with types compatible with `encoding/json`.


# License

This project is licensed under the MIT License - see the LICENSE file for details
