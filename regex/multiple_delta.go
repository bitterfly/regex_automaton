package regex

type MultipleTransition struct {
	initialState int
	letter       rune
	goalState    int
}

func NewMultipleTransition(initialState int, letter rune, goalState int) *MultipleTransition {
	return &MultipleTransition{
		initialState: initialState,
		letter:       letter,
		goalState:    goalState,
	}
}

type MultipleDeltaTransitions struct {
	transitions map[MultipleTransition]struct{}
}

func NewMultipleEmptyTransition() *MultipleDeltaTransitions {
	return &MultipleDeltaTransitions{
		transitions: make(map[MultipleTransition]struct{}),
	}
}

func NewMultipleDeltaTransitions(transitions map[MultipleTransition]struct{}) *MultipleDeltaTransitions {
	return &MultipleDeltaTransitions{
		transitions: transitions,
	}
}

func (mdt *MultipleDeltaTransitions) addTransition(initialState int, letter rune, goalState int) {
	mdt.transitions[*NewMultipleTransition(initialState, letter, goalState)] = struct{}{}
}

func (mdt *MultipleDeltaTransitions) removeTransition(initialState int, letter rune, goalState int) {
	delete(mdt.transitions, *NewMultipleTransition(initialState, letter, goalState))
}

func (mdt *MultipleDeltaTransitions) addTransitions(other *MultipleDeltaTransitions) {
	for k, v := range other.transitions {
		mdt.transitions[k] = v
	}
}
