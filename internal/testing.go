package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
)

type (
	Test struct {
		Rule     any
		Data     any
		Expected any
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

// This gets the tests.json file that we've proposed become the new official one in
// https://github.com/jwadhams/json-logic/pull/48 but that hasn't merged yet.
func GetScenariosFromProposedOfficialTestSuite() Tests {
	var err error
	buffer, err := os.ReadFile("internal/json_logic_pr_48_tests.json")
	if err != nil {
		log.Fatal(err)
	}

	return getScenariosFromFile(buffer)
}

func GetScenariosFromOfficialTestSuite() Tests {
	response, err := http.Get("http://jsonlogic.com/tests.json")
	if err != nil {
		log.Fatal(err)
	}

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()
	return getScenariosFromFile(buffer)
}

func getScenariosFromFile(buffer []byte) Tests {
	var (
		tests     Tests
		scenarios []any
		err       = json.Unmarshal(buffer, &scenarios)
	)
	if err != nil {
		log.Fatal(err)
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
			Rule:     scenario.([]any)[0],
			Data:     scenario.([]any)[1],
			Expected: scenario.([]any)[2],
		})
	}

	return tests
}
