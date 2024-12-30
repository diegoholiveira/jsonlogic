package jsonlogic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func BenchmarkEngine_normal(b *testing.B) {
	rules := `{"<": [{"var": "temp"}, 110]}`

	for i := 1; i <= 1000; i++ {
		tempData := fmt.Sprintf(`{"temp": %d}`, i)
		data := strings.NewReader(tempData)

		var result strings.Builder
		Apply(strings.NewReader(rules), data, &result)
	}

}

func BenchmarkEngine_Apply(b *testing.B) {

	engine := NewEngine()
	rules := `{"<": [{"var": "temp"}, 110]}`

	hash := sha256.Sum256([]byte(rules))
	hashKey := hex.EncodeToString(hash[:]) // Convert hash to a string

	// Call Build with both ruleKey and hashKey
	engine.Build(rules, hashKey)

	for i := 1; i <= 1000; i++ {
		tempData := fmt.Sprintf(`{"temp": %d}`, i)
		data := strings.NewReader(tempData)
		var result strings.Builder
		engine.Apply(strings.NewReader(rules), data, &result)
	}
}
