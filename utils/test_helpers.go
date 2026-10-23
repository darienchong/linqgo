package utils

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

const TestCaseFailureMessageTemplate = `
Test case %q failed:
- Expected: %+v
- Actual:   %+v
`

type TestCase[TInput, TOutput any] struct {
	Input    TInput
	Expected TOutput
	Name     string
}

func (tc *TestCase[TInput, TOutput]) Execute(t *testing.T, actual TOutput) {
	t.Helper()

	if reflect.DeepEqual(actual, tc.Expected) {
		return
	}

	if tc.Name == "" {
		tc.Name = fmt.Sprintf("(Input: %+v, Expected: %+v)", tc.Input, tc.Expected)
	}

	t.Errorf(TestCaseFailureMessageTemplate, tc.Name, tc.Expected, actual)
}

func (tc *TestCase[TInput, TOutput]) FuzzyExecute(t *testing.T, actual, tolerance float64) {
	t.Helper()

	if math.Abs(any(tc.Expected).(float64)-actual) <= tolerance {
		return
	}

	if tc.Name == "" {
		tc.Name = fmt.Sprintf("(Input: %+v, Expected: %+v)", tc.Input, tc.Expected)
	}

	t.Errorf(TestCaseFailureMessageTemplate, tc.Name, tc.Expected, actual)
}
