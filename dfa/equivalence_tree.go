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
		left:   nil,
		right:  nil,
		height: 1,
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
		(*root).height = max(height((*root).left), height((*root).right)) + 1
	}

	compareResult := CompareEquivalenceClasses(&(*root).data.equivalenceClass, &data.equivalenceClass)

	//data < root.data
	if compareResult == 1 {
		insert(&(*root).left, data)
		if height((*root).left)-height((*root).right) == 2 {
			if CompareEquivalenceClasses(&(*root).left.data.equivalenceClass, &data.equivalenceClass) == 1 {
				(*root) = rightRotate((*root))
			} else {
				(*root) = leftRightRotate((*root))
			}
		}
	}

	if compareResult == -1 {
		insert(&(*root).right, data)
		if height((*root).right)-height((*root).left) == 2 {
			if CompareEquivalenceClasses(&(*root).right.data.equivalenceClass, &data.equivalenceClass) == -1 {
				(*root) = leftRotate((*root))
			} else {
				(*root) = rightLeftRotate((*root))
			}
		}
	}

	(*root).height = max(height((*root).left), height((*root).right)) + 1
}

func (t *EquivalenceTree) find(needle EquivalenceNode) (int, bool) {
	if t == nil {
		return -1, false
	}

	compare_result := CompareEquivalenceClasses(&t.data.equivalenceClass, &needle.equivalenceClass)
	if compare_result == 0 {
		return t.data.state, true
	}

	if compare_result == 1 {
		return t.left.find(needle)
	} else {
		return t.right.find(needle)
	}
}
