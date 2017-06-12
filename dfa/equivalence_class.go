package dfa

import (
	"reflect"
	"sort"
)

type EquivalenceClass struct {
	isFinal  bool
	children []Transition
}

func NewEquivalenceClass(isFinal bool, children Children) *EquivalenceClass {
	c := make([]Transition, len(children.children))
	i := 0
	for k, _ := range children.children {
		c[i] = k
		i += 1
	}

	sort.Slice(c, func(i, j int) bool { return c[i].letter < c[j].letter || c[i].state < c[j].state })
	return &EquivalenceClass{
		isFinal:  isFinal,
		children: c,
	}
}

func CompareEquivalenceClasses(first, second *EquivalenceClass) int {
	//the final one is bigger
	if (*first).isFinal != (*second).isFinal {
		if (*first).isFinal {
			return 1
		}
		return -1
	}
	if (*first).children == nil && (*second).children == nil {
		return 0
	}

	if len((*first).children) != len((*second).children) {
		if len((*first).children) > len((*second).children) {
			return 1
		}
		return -1
	}

	return compareTransitionSlices(first.children, second.children)
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

	compare_result := CompareEquivalenceClasses(&t.data.equivalenceClass, &needle.equivalenceClass)
	if compare_result == 0 {

		if !reflect.DeepEqual(needle.equivalenceClass.children, t.data.equivalenceClass.children) {
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
		compare_result := CompareEquivalenceClasses(&(*node).data.equivalenceClass, &needle.equivalenceClass)
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
	compare_result := CompareEquivalenceClasses(&(*node).data.equivalenceClass, &needle.equivalenceClass)
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
