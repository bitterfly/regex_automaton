package regex

type NDFA struct {
	initialState int
	numStates    int
	finalState   int
	delta        *MultipleDeltaTransitions
}

func NewNDFA(initialState, numStates, finalState int, delta *MultipleDeltaTransitions) *NDFA {
	return &NDFA{
		initialState: initialState,
		numStates:    numStates,
		finalState:   finalState,
		delta:        delta,
	}
}
