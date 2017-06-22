package regex

type NormalNDFA struct {
	finalStates  map[int]struct{}
	initialState int
	transitions  map[int][]common.Transition
}

func NewNormalNDFA(finalStates map[int]struct{}, initialState int, transitions map[int][]common.Transitions) *NormalNDFA {
	return &NormalNDFA{
		finalStates:  finalStates,
		initialState: initialState,
		transitions:  transitions,
	}
}

func (n *NormalNDFA) GetInitialState() int {
	return n.initialState
}

func (n *NormalNDFA) IsFinal(state int) bool {
	_, ok := n.finalStates[state]
	return ok
}

func (n *NormalNDFA) GetTransitions(state int) []common.Transition {
	return n.transitions[state]
}

func (n *NormalNDFA) GetDestinations(state int) map[rune][]int {
	destinations := make(map[rune]int)

	for _, transition := range n.transitions[state] {
		destinations[transition.GetLetter()] = append(destinations[transition.GetLetter()], destinations.GetState())
	}
	return destinations
}
