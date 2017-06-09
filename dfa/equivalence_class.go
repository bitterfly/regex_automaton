package dfa

import "sort"

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

func (e *EquivalenceClass) Compare(other EquivalenceClass) int {
	if e.isFinal != other.isFinal {
		if e.isFinal {
			return 1
		}
		return -1
	} else {
		if len(e.children.children) == len(other.children.children) {
			first_labels := make([]rune, len(e.children.children))
			first_states := make([]int, len(e.children.children))

			i := 0
			for k, _ := range e.children.children {
				first_labels[i] = k.letter
				first_states[i] = k.state
			}

			sort.Slice(first_labels, func(i, j int) bool { return first_labels[i] < first_labels[j] })
			sort.Ints(first_states)

			second_labels := make([]rune, len(e.children.children))
			second_states := make([]int, len(e.children.children))

			i = 0
			for k, _ := range other.children.children {
				second_labels[i] = k.letter
				second_states[i] = k.state
			}

			sort.Slice(second_labels, func(i, j int) bool { return second_labels[i] < second_labels[j] })
			sort.Ints(second_states)

			labels_difference := compareRuneSlices(first_labels, second_labels)
			states_difference := compare(first_states, second_states)

			if labels_difference == 0 {
				// same labes
				if states_difference == 0 {
					//same labels and states
					return 0
				} else {
					//return who has lexicographically bigger state
					return states_difference
				}
			} else {
				// different labes - return who has lexicographically bigger label
				return labels_difference
			}
			// children are not the same lenght - > has more children
		} else {
			if len(e.children.children) > len(other.children.children) {
				return 1
			} else {
				return -1
			}
		}
	}
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
	left  *EquivalenceTree
	right *EquivalenceTree
	data  EquivalenceNode
}

func NewEquivalenceTree(data EquivalenceNode) *EquivalenceTree {
	return &EquivalenceTree{
		left:  nil,
		right: nil,
		data:  data,
	}
}

func (t *EquivalenceTree) Find(needle EquivalenceNode) bool {
	if t == nil {
		return false
	}

	compare_result := t.data.equivalenceClass.Compare(needle.equivalenceClass)
	if compare_result == 0 {
		return true
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
		compare_result := (*node).data.equivalenceClass.Compare(needle.equivalenceClass)
		if compare_result == 1 {
			Insert(&(*node).left, needle)
		} else {
			Insert(&(*node).right, needle)
		}
	}
}
