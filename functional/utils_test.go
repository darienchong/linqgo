package functional

import (
	"cmp"
	"fmt"
	"testing"
)

type TestCase[In, Out any] struct {
	Input    In
	Expected Out
}

func (t *TestCase[In, Out]) String() string {
	return fmt.Sprintf("(%v, %v)", t.Input, t.Expected)
}

func TestUtils_Reverse(t *testing.T) {
	tcs := []TestCase[[]int, int]{
		{
			Input:    []int{1, 2},
			Expected: 1, // 1 < 2, so the reverse comparator should return 1
		},
	}
	reverseComparer := reverse(cmp.Compare[int])
	for _, tc := range tcs {
		if actual := reverseComparer(tc.Input[0], tc.Input[1]); actual != tc.Expected {
			t.Fatalf("%v: expected %d, but got %d\n", tc, tc.Expected, actual)
		}
	}
}
