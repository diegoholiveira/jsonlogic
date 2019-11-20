package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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

func convertInterfaceToReader(i interface{}) io.Reader {
	var result bytes.Buffer

	encoder := json.NewEncoder(&result)
	encoder.Encode(i)

	return &result

}

func GetScenariosFromOfficialTestSuite() Tests {
	var tests Tests

	response, err := http.Get("http://jsonlogic.com/tests.json")
	if err != nil {
		log.Fatal(err)

		return tests
	}

	buffer, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()

	var scenarios []interface{}

	err = json.Unmarshal(buffer, &scenarios)
	if err != nil {
		log.Fatal(err)

		return tests
	}

    // add missing but relevant scenarios
    var rule []interface{}

    scenarios = append(scenarios,
        append(rule,
            make(map[string]interface {}, 0),
            make(map[string]interface {}, 0),
            make(map[string]interface {}, 0)))

	for _, scenario := range scenarios {
		if reflect.ValueOf(scenario).Kind() == reflect.String {
			continue
		}

		tests = append(tests, Test{
			Rule:     convertInterfaceToReader(scenario.([]interface{})[0]),
			Data:     convertInterfaceToReader(scenario.([]interface{})[1]),
			Expected: convertInterfaceToReader(scenario.([]interface{})[2]),
		})
	}

	return tests
}
