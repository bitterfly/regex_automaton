package dfa

import (
	"reflect"
	"sort"
)

type EquivalenceClass struct {
	isFinal  bool
	children Children
}

func NewEquivalenceClass(isFinal bool, children Children) *EquivalenceClass {
	return &EquivalenceClass{
		isFinal:  isFinal,
		children: children,
	}
}

func compare(first []int, second []int) int {
	for i, value := range first {
		if value > second[i] {
			return 1
		}
		if value < second[i] {
			return -1
		}
	}
	return 0
}

func compareRuneSlices(first []rune, second []rune) int {
	for i, value := range first {
		if value > second[i] {
			return 1
		}
		if value < second[i] {
			return -1
		}
	}
	return 0
}

func Compare(state int, first EquivalenceClass, other EquivalenceClass) int {
	//the final one is bigger
	if first.isFinal != other.isFinal {
		if first.isFinal {
			return 1
		}
		return -1
	}
	if first.children.children == nil && other.children.children == nil {
		return 0
	}

	if len(first.children.children) != len(other.children.children) {
		if len(first.children.children) > len(other.children.children) {
			return 1
		}
		return -1
	}

	// both are final/non final and have the same number of children
	first_labels := make([]rune, len(first.children.children))
	first_states := make([]int, len(first.children.children))

	i := 0
	for k, _ := range first.children.children {
		first_labels[i] = k.letter
		first_states[i] = k.state
		i += 1
	}

	sort.Slice(first_labels, func(i, j int) bool { return first_labels[i] < first_labels[j] })
	sort.Ints(first_states)

	second_labels := make([]rune, len(other.children.children))
	second_states := make([]int, len(other.children.children))

	i = 0
	for k, _ := range other.children.children {
		second_labels[i] = k.letter
		second_states[i] = k.state
		i += 1
	}

	sort.Slice(second_labels, func(i, j int) bool { return second_labels[i] < second_labels[j] })
	sort.Ints(second_states)

	// -1 0 1
	labels_difference := compareRuneSlices(first_labels, second_labels)
	states_difference := compare(first_states, second_states)

	if labels_difference != 0 {
		return labels_difference
	}

	if states_difference != 0 {
		return states_difference
	}

	return 0
}

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
	parent *EquivalenceTree
	left   *EquivalenceTree
	right  *EquivalenceTree
	data   EquivalenceNode
}

func NewEquivalenceTree(data EquivalenceNode) *EquivalenceTree {
	return &EquivalenceTree{
		parent: nil,
		left:   nil,
		right:  nil,
		data:   data,
	}
}

func (t *EquivalenceTree) Find(needle EquivalenceNode) (int, bool) {
	if t == nil {
		return -1, false
	}

	compare_result := Compare(needle.state, t.data.equivalenceClass, needle.equivalenceClass)
	if compare_result == 0 {

		if !reflect.DeepEqual(needle.equivalenceClass.children.sortedChildren(), t.data.equivalenceClass.children.sortedChildren()) {
			panic("compare said unequal things are equal")
		}
		return t.data.state, true
	}

	if compare_result == 1 {
		return t.left.Find(needle)
	} else {
		return t.right.Find(needle)
	}
}

func Insert(node **EquivalenceTree, needle EquivalenceNode) {
	if (*node) == nil {
		(*node) = NewEquivalenceTree(needle)
	} else {
		compare_result := Compare(needle.state, (*node).data.equivalenceClass, needle.equivalenceClass)
		if compare_result == 1 {
			Insert(&(*node).left, needle)
			(*node).left.parent = (*node)
		} else {
			Insert(&(*node).right, needle)
			(*node).right.parent = (*node)
		}
	}
}

func FindMin(node *EquivalenceTree) *EquivalenceTree {
	currentNode := node
	for currentNode.left != nil {
		currentNode = currentNode.left
	}
	return currentNode
}

func ReplaceNode(node **EquivalenceTree, newNode *EquivalenceTree) {
	if (*node).parent != nil {
		if (*node) == (*node).parent.left {
			(*node).parent.left = newNode
		} else {
			(*node).parent.right = newNode
		}
	}
	if newNode != nil {
		newNode.parent = (*node).parent
	}
}

func Delete(node **EquivalenceTree, needle EquivalenceNode) {
	compare_result := Compare(needle.state, (*node).data.equivalenceClass, needle.equivalenceClass)
	if compare_result == -1 {
		Delete(&(*node).right, needle)
	} else if compare_result == 1 {
		Delete(&(*node).left, needle)
	} else {
		if (*node).left != nil && (*node).right != nil {
			succesor := FindMin((*node).right)
			(*node).data = succesor.data
			Delete(&succesor, succesor.data)
		} else if (*node).left != nil {
			ReplaceNode(node, (*node).left)
		} else if (*node).right != nil {
			ReplaceNode(node, (*node).right)
		} else {
			ReplaceNode(node, nil)
		}
	}
}

func Update(node **EquivalenceTree, oldNode EquivalenceNode, newNode EquivalenceNode) {
	Delete(node, oldNode)
	Insert(node, newNode)
}
