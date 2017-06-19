package automaton

import "sort"

type FA struct {
	MaxState     int
	FinalStates  map[int]struct{}
	NumStates    int
	NumEqClasses int
}

func NewFA(maxState int, numStates int, eqClasses int, _finalStates []int) *FA {
	finalStates := make(map[int]struct{})
	for _, state := range _finalStates {
		finalStates[state] = struct{}{}
	}

	return &FA{
		MaxState:     maxState,
		FinalStates:  finalStates,
		NumStates:    numStates,
		NumEqClasses: eqClasses,
	}
}

func EmptyAutomaton() *FA {
	return &FA{
		MaxState:     1,
		FinalStates:  make(map[int]struct{}),
		NumStates:    1,
		NumEqClasses: 1,
	}
}

func (f *FA) RemoveFinalState(state int) {
	delete(f.FinalStates, state)
}

func (f *FA) SortedFinalStates() []int {
	var states []int
	for k, _ := range f.FinalStates {
		states = append(states, k)
	}
	sort.Ints(states)
	return states
}

func (f *FA) GetInitialState() int {
	return f.MaxState - f.NumStates + 1
}

func (f *FA) GetNumState() int {
	for k, _ := range f.FinalStates {
		return k
	}
	return -1
}
