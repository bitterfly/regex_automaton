package dfa

import "github.com/bitterfly/regex_automata/common"

type EquivalenceClass struct {
	isFinal  bool
	children []common.Transition
}

func NewEquivalenceClass(isFinal bool, children []common.Transition) *EquivalenceClass {
	return &EquivalenceClass{
		isFinal:  isFinal,
		children: children,
	}
}

func CompareEquivalenceClasses(first, second *EquivalenceClass) int {
	if first.isFinal != second.isFinal {
		if first.isFinal {
			return 1
		}
		return -1
	}
	if first.children == nil && second.children == nil {
		return 0
	}

	if len(first.children) != len(second.children) {
		if len(first.children) > len(second.children) {
			return 1
		}
		return -1
	}

	return common.CompareTransitionSlices(first.children, second.children)
}
