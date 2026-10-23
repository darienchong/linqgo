package utils

import (
	"fmt"
)

func ToPosTreeNode(nodeValues []int) *TreeNode {
	if len(nodeValues) == 0 {
		return nil
	}

	root := &TreeNode{Val: nodeValues[0]}
	var q []*TreeNode
	q = append(q, root)
	ptr := 0
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		var left, right *TreeNode
		if ptr+1 < len(nodeValues) && nodeValues[ptr+1] >= 0 {
			left = &TreeNode{Val: nodeValues[ptr+1]}
		}
		if ptr+2 < len(nodeValues) && nodeValues[ptr+2] >= 0 {
			right = &TreeNode{Val: nodeValues[ptr+2]}
		}
		curr.Left = left
		curr.Right = right
		if left != nil {
			q = append(q, left)
		}
		if right != nil {
			q = append(q, right)
		}
		ptr += 2
	}
	return root
}

func ToTreeNode(nodeValues []int) *TreeNode {
	if len(nodeValues) == 0 {
		return nil
	}

	root := &TreeNode{Val: nodeValues[0]}
	var q []*TreeNode
	q = append(q, root)
	ptr := 0
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		var left, right *TreeNode
		if ptr+1 < len(nodeValues) && nodeValues[ptr+1] > 0 {
			left = &TreeNode{Val: nodeValues[ptr+1]}
		}
		if ptr+2 < len(nodeValues) && nodeValues[ptr+2] > 0 {
			right = &TreeNode{Val: nodeValues[ptr+2]}
		}
		curr.Left = left
		curr.Right = right
		if left != nil {
			q = append(q, left)
		}
		if right != nil {
			q = append(q, right)
		}
		ptr += 2
	}
	return root
}

// e.g. 1
//      2   3
//      4 5 6 7
// = [1 2 3 4 5 6 7]
// = [1 | 2 3 | 4 5 6 7]

func ToValues(root *TreeNode) []int {
	var values []int
	var q []*TreeNode
	q = append(q, root)
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]
		if curr == nil {
			values = append(values, 0)
			continue
		}
		values = append(values, curr.Val)
		q = append(q, curr.Left)
		q = append(q, curr.Right)
	}
	// Prune trailing zeroes
	idx := len(values) - 1
	for values[idx] == 0 {
		idx--
	}
	return values[:idx+1]
}

func walk(node *TreeNode) {
	if node == nil {
		fmt.Printf("<nil>, ")
		return
	}
	fmt.Printf("%d, ", node.Val)
	walk(node.Left)
	walk(node.Right)
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}
