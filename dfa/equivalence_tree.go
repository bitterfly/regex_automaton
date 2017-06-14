package dfa

type EquivalenceNode struct {
	state            int
	equivalenceClass EquivalenceClass
}

func NewEquivalenceNode(state int, equivalenceClass EquivalenceClass) *EquivalenceNode {
	return &EquivalenceNode{
		state:            state,
		equivalenceClass: equivalenceClass,
	}
}

type EquivalenceTree struct {
	data        EquivalenceNode
	height      int
	left, right *EquivalenceTree
}

func NewEquivalenceTree(data EquivalenceNode) *EquivalenceTree {
	return &EquivalenceTree{
		parent: nil,
		left:   nil,
		right:  nil,
		data:   data,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func height(root *EquivalenceTree) int {
	if root != nil {
		return root.height
	} else {
		return -1
	}
}

func leftRotate(root *EquivalenceTree) *EquivalenceTree {
	node := root.right
	root.right = node.left
	node.left = root

	root.height = max(height(root.left), height(root.right)) + 1
	node.height = max(height(node.left), height(node.right)) + 1
	return node
}

func rightRotate(root *EquivalenceTree) *EquivalenceTree {
	node := root.left
	root.left = node.right
	node.right = root

	root.height = max(height(root.left), height(root.right)) + 1
	node.height = max(height(node.left), height(node.right)) + 1
	return node
}

func leftRightRotate(root *EquivalenceTree) *EquivalenceTree {
	root.left = leftRotate(root.left)
	root = rightRotate(root)
	return root
}

func rightLeftRotate(root *EquivalenceTree) *EquivalenceTree {
	root.right = rightRotate(root.right)
	root = leftRotate(root)
	return root
}

func insert(root **EquivalenceTree, data *EquivalenceNode) {
	if (*root) == nil {
		(*root) = NewEquivalenceTree(*data)
		(*root).height = max(height(root.left), height(root.right)) + 1
	}

	compareResult := CompareEquivalenceClasses(root.data.equivalenceClass, data.equivalenceClass)

	//data < root.data
	if compareResult == 1 {
		(*root).left = insert((*root).left, data)
		if height((*root).left)-height((*root).right) == 2 {
			if CompareEquivalenceClasses((*root).left.data.equivalenceClass, data.equivalenceClass) == 1 {
				(*root) = rightRotate((*root))
			} else {
				(*root) = leftRightRotate((*root))
			}
		}
	}

	if compareResult == -1 {
		(*root).right = insert((*root).right, data)
		if height((*root).right)-height((*root).left) == 2 {
			if CompareEquivalenceClasses((*root).right.data.equivalenceClass, data.equivalenceClass) == -1 {
				(*root) = leftRotate((*root))
			} else {
				(*root) = rightLeftRotate((*root))
			}
		}
	}

	(*root).height = max(height((*root).left), height((*root).right)) + 1
}
