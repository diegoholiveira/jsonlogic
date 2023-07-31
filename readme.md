# Go JSON Logic (CIDgravity custom version)

Forked from [https://github.com/diegoholiveira/jsonlogic](https://github.com/diegoholiveira/jsonlogic)
For original documentation, view the original repository

Added function to generate the json logic with solved variables

```go
func solveVarsBackToJsonLogic(rule, data interface{}) ([]byte, error) {
	ruleMap := rule.(map[string]interface{})
	result := make(map[string]interface{})

	for operator, values := range ruleMap {
		result[operator] = solveVars(values, data)
	}

	body, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return body, nil
}

```

For example, if you specify the folowing rules model :


```json
{
  "and":[
    { "==":[{ "var":"VariableA" }, true] },
    { "==":[{ "var":"VariableB" }, true] },
    { ">=":[{ "var":"VariableC" }, 17179869184] },
    { "==":[{ "var":"VariableD" }, "0"] },
    { "<":[{ "var":"VariableE" }, 20] }
  ]
}

```

You will get as output, the folowing response (using a specific data, all variables will be replaced with matching values) :

```json
{
  "and":[
    { "==":[false, true] },
    { "==":[true, true] },
    { ">=":[34359738368, 17179869184] },
    { "==":[12, "0"] },
    { "<":[14, 20] }
  ]
}

```