package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
)

type (
	Test struct {
		Rule     io.Reader
		Data     io.Reader
		Expected io.Reader
	}

	Tests []Test
)

func convertInterfaceToReader(i any) io.Reader {
	var result bytes.Buffer

	encoder := json.NewEncoder(&result)
	err := encoder.Encode(i)
	if err != nil {
		panic(err)
	}

	return &result
}

func GetScenariosFromOfficialTestSuite() Tests {
	var tests Tests

	response, err := http.Get("http://jsonlogic.com/tests.json")
	if err != nil {
		log.Fatal(err)

		return tests
	}

	buffer, _ := io.ReadAll(response.Body)

	response.Body.Close()

	var scenarios []any

	err = json.Unmarshal(buffer, &scenarios)
	if err != nil {
		log.Fatal(err)

		return tests
	}

	// add missing but relevant scenarios
	var rule []any

	scenarios = append(scenarios,
		append(rule,
			make(map[string]any),
			make(map[string]any),
			make(map[string]any)))

	for _, scenario := range scenarios {
		if reflect.ValueOf(scenario).Kind() == reflect.String {
			continue
		}

		tests = append(tests, Test{
			Rule:     convertInterfaceToReader(scenario.([]any)[0]),
			Data:     convertInterfaceToReader(scenario.([]any)[1]),
			Expected: convertInterfaceToReader(scenario.([]any)[2]),
		})
	}

	return tests
}
