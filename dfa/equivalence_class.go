package dfa

type EquivalenceClass struct {
	isFinal  bool
	children []Transition
}

func NewEquivalenceClass(isFinal bool, children []Transition) *EquivalenceClass {
	return &EquivalenceClass{
		isFinal:  isFinal,
		children: children,
	}
}

func CompareEquivalenceClasses(first, second *EquivalenceClass) int {
	//the final one is bigger
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

	return compareTransitionSlices(first.children, second.children)
}
