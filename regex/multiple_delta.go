package regex

import "github.com/bitterfly/regex_automata/common"

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
	transitions map[int][]common.Transition
}

func NewMultipleEmptyTransition() *MultipleDeltaTransitions {
	return &MultipleDeltaTransitions{
		triple:      make(map[MultipleTransition]struct{}),
		transitions: make(map[int][]common.Transition),
	}
}

func NewMultipleDeltaTransitions(triple map[MultipleTransition]struct{}) *MultipleDeltaTransitions {
	transitions := make(map[int][]common.Transition)
	for k, _ := range triple {
		transitions[k.initialState] = append(transitions[k.initialState], *common.NewTransition(k.goalState, k.letter))
	}

	return &MultipleDeltaTransitions{
		triple:      triple,
		transitions: transitions,
	}
}

func (mdt *MultipleDeltaTransitions) addTransition(initialState int, letter rune, goalState int) {
	mdt.triple[*NewMultipleTransition(initialState, letter, goalState)] = struct{}{}
	mdt.transitions[initialState] = append(mdt.transitions[initialState], *common.NewTransition(goalState, letter))
}

func (mdt *MultipleDeltaTransitions) addTransitions(other *MultipleDeltaTransitions) {
	for k := range other.triple {
		mdt.addTransition(k.initialState, k.letter, k.goalState)
	}
}
