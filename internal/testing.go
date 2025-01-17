package internal

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
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

var testFile = flag.String("jsonlogic-test-file", "", "tests.json file to use instead of http://jsonlogic.com/tests.json")

func GetScenariosFromOfficialTestSuite() Tests {
	var tests Tests

	var buffer []byte
	if *testFile != "" {
		fmt.Printf("reading from local file\n")
		var err error
		buffer, err = os.ReadFile(*testFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		response, err := http.Get("http://jsonlogic.com/tests.json")
		if err != nil {
			log.Fatal(err)

			return tests
		}

		buffer, err = io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		response.Body.Close()
	}

	var scenarios []any
	var err = json.Unmarshal(buffer, &scenarios)
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
