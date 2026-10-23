package utils

import (
	"reflect"
	"testing"
)

func TestToValues(t *testing.T) {
	cases := [][]int{
		{1, 0, 2, 0, 3, 0, 4},
		{1, 2, 3, 4, 5, 6, 7},
		{2, 1, 3},
	}

	for _, input := range cases {
		root := ToTreeNode(input)
		values := ToValues(root)
		if !reflect.DeepEqual(input, values) {
			t.Errorf("input: %v, values: %v", input, values)
		}
	}
}
