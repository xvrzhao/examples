package data_structure

// binaryTreeNode 是二叉树节点
type binaryTreeNode struct {
	value interface{}
	left  *binaryTreeNode
	right *binaryTreeNode
}

func newBinaryTreeNode(v interface{}) *binaryTreeNode {
	return &binaryTreeNode{value: v}
}

// CreateBinaryTree 创建二叉树并返回根节点
func CreateBinaryTree() (root *binaryTreeNode) {
	root = newBinaryTreeNode("a")
	b := newBinaryTreeNode("b")
	c := newBinaryTreeNode("c")
	d := newBinaryTreeNode("d")
	e := newBinaryTreeNode("e")
	f := newBinaryTreeNode("f")
	g := newBinaryTreeNode("g")

	root.left, root.right = b, c
	b.left, b.right = d, e
	c.left, c.right = f, g

	return
}

var order []interface{}

func init() {
	order = make([]interface{}, 0, 7)
}

// GetResultOrder 返回排序结果并清空内部状态
func GetResultOrder() []interface{} {
	defer func() {
		if order != nil {
			order = order[:0]
		}
	}()
	return order
}

// BinaryTreePreOrderBinaryTree 根据根节点前序遍历二叉树
func BinaryTreePreOrderBinaryTree(root *binaryTreeNode) {
	order = append(order, root.value)

	if left := root.left; left != nil {
		BinaryTreePreOrderBinaryTree(left)
	}

	if right := root.right; right != nil {
		BinaryTreePreOrderBinaryTree(right)
	}
}

// BinaryTreeInOrder 根据根节点中序遍历二叉树
func BinaryTreeInOrder(root *binaryTreeNode) {
	if left := root.left; left != nil {
		BinaryTreeInOrder(left)
	}

	order = append(order, root.value)

	if right := root.right; right != nil {
		BinaryTreeInOrder(right)
	}
}

// BinaryTreePostOrder 根据根节点后序遍历二叉树
func BinaryTreePostOrder(root *binaryTreeNode) {
	if left := root.left; left != nil {
		BinaryTreePostOrder(left)
	}

	if right := root.right; right != nil {
		BinaryTreePostOrder(right)
	}

	order = append(order, root.value)
}
