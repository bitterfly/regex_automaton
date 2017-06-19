package regex

import "github.com/bitterfly/pka/common"

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
	triple      map[MultipleTransition]struct{}
	transitions map[common.Transition][]int
}

func NewMultipleEmptyTransition() *MultipleDeltaTransitions {
	return &MultipleDeltaTransitions{
		triple:      make(map[MultipleTransition]struct{}),
		transitions: make(map[common.Transition][]int),
	}
}

func NewMultipleDeltaTransitions(triple map[MultipleTransition]struct{}) *MultipleDeltaTransitions {
	transitions := make(map[common.Transition][]int)
	for k, _ := range triple {
		transitions[*common.NewTransition(k.initialState, k.letter)] = append(transitions[*common.NewTransition(k.initialState, k.letter)], k.goalState)
	}

	return &MultipleDeltaTransitions{
		triple:      triple,
		transitions: transitions,
	}
}

func (mdt *MultipleDeltaTransitions) addTransition(initialState int, letter rune, goalState int) {
	mdt.triple[*NewMultipleTransition(initialState, letter, goalState)] = struct{}{}
	mdt.transitions[*common.NewTransition(initialState, letter)] = append(mdt.transitions[*common.NewTransition(initialState, letter)], goalState)
}

// func (mdt *MultipleDeltaTransitions) removeTransition(initialState int, letter rune, goalState int) {
// 	delete(mdt.triple, *NewMultipleTransition(initialState, letter, goalState))
// }

func (mdt *MultipleDeltaTransitions) addTransitions(other *MultipleDeltaTransitions) {
	for k, v := range other.triple {
		mdt.triple[k] = v
		mdt.transitions[*common.NewTransition(k.initialState, k.letter)] = append(mdt.transitions[*common.NewTransition(k.initialState, k.letter)], k.goalState)
	}
}
