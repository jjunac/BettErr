package betterr

import (
	"encoding/json"
	"regexp"
	"testing"
)

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("\nExpected: %v\nActual: %v", expected, actual)
	}
}

func assertRegexp(t *testing.T, pattern string, value string) {
	t.Helper()
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		t.Fatalf("\nInvalid regexp pattern '%s': %v", pattern, err)
	}
	if !matched {
		t.Errorf("\nExpected pattern:\n%s\n\nActual value:\n%s\n", pattern, value)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("\nExpected no error, got: %v", err)
	}
}

func assertJSONEq(t *testing.T, expected, actual string) {
	t.Helper()
	var expectedJSON, actualJSON interface{}
	
	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		t.Fatalf("\nInvalid expected JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		t.Fatalf("\nInvalid actual JSON: %v", err)
	}
	
	expectedBytes, _ := json.Marshal(expectedJSON)
	actualBytes, _ := json.Marshal(actualJSON)
	if string(expectedBytes) != string(actualBytes) {
		t.Errorf("\nJSONs not equal:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
	}
}