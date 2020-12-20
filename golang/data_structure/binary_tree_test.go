package data_structure

import "testing"

// TestBinaryTreeOrder 测试二叉树前中后序遍历结果
func TestBinaryTreeOrder(t *testing.T) {
	r := CreateBinaryTree()

	BinaryTreePreOrderBinaryTree(r)
	t.Log("pre-order:", GetResultOrder())

	BinaryTreeInOrder(r)
	t.Log("in-order:", GetResultOrder())

	BinaryTreePostOrder(r)
	t.Log("post-order:", GetResultOrder())
}
