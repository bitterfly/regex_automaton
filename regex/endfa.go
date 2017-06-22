package regex

import (
	"fmt"
	"os"
)

type ENDFA struct {
	initialState int
	numStates    int
	finalState   int
	delta        *MultipleDeltaTransitions
}

func NewENDFA(initialState, numStates, finalState int, delta *MultipleDeltaTransitions) *ENDFA {
	return &ENDFA{
		initialState: initialState,
		numStates:    numStates,
		finalState:   finalState,
		delta:        delta,
	}
}

func (n *ENDFA) EpsilonClosure(states map[int]struct{}) (map[int]struct{}, bool) {
	epsilonClosure := make(map[int]struct{})
	stack := make([]int, len(states))

	i := 0
	for state, _ := range states {
		epsilonClosure[state] = struct{}{}
		stack[i] = state
		i += 1
	}

	for len(stack) != 0 {
		state := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		for _, destinations := range n.delta.transitions[state] {
			if destinations.GetLetter() == 0 {
				_, ok := epsilonClosure[destinations.GetState()]
				if !ok {
					epsilonClosure[destinations.GetState()] = struct{}{}
					stack = append(stack, destinations.GetState())
				}
			}
		}
	}

	_, ok := epsilonClosure[n.finalState]
	return epsilonClosure, ok
}

func (n *ENDFA) RemoveEpsilonTransitions() *NDFA {
	delta := NewMultipleEmptyTransition()
	finalStates := make(map[int]struct{})
	states := make(map[int]struct{})

	stack := make([]int, 1)
	stack[0] = n.initialState
	states[n.initialState] = struct{}{}

	for len(stack) != 0 {
		state := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		enclosure, isFinal := n.EpsilonClosure(map[int]struct{}{state: struct{}{}})
		if isFinal {
			finalStates[state] = struct{}{}
		}

		for otherState, _ := range enclosure {
			destinations, _ := n.delta.transitions[otherState]
			for _, destination := range destinations {
				if destination.GetLetter() != 0 {
					delta.addTransition(state, destination.GetLetter(), destination.GetState())
					_, ok := states[destination.GetState()]
					if !ok {

						stack = append(stack, destination.GetState())
						states[destination.GetState()] = struct{}{}
					}
				}
			}
		}
	}

	return NewNDFA(n.initialState, len(states), finalStates, delta.transitions)
}

func (n *ENDFA) GetNonEpsilonTransitions(states map[int]struct{}) map[rune]map[int]struct{} {
	transitions := make(map[rune]map[int]struct{})

	for state, _ := range states {
		for _, tr := range n.delta.transitions[state] {
			if tr.GetLetter() != 0 {
				_, ok := transitions[tr.GetLetter()]
				if !ok {
					transitions[tr.GetLetter()] = make(map[int]struct{})
				}
				transitions[tr.GetLetter()][tr.GetState()] = struct{}{}
			}
		}
	}

	return transitions
}

//=================================================

func (n *ENDFA) GetInitialState() int {
	return n.initialState
}

func (n *ENDFA) GetNumStates() int {
	return n.numStates
}

func (n *ENDFA) Print() {
	fmt.Printf("====ENDFA====\n")
	fmt.Printf("Initial: %d, NumStates: %d, FinalState: %d\n", n.initialState, n.numStates, n.finalState)
	n.PrintFunction()
	fmt.Printf("\n====NAFD====\n")
}

func (n *ENDFA) PrintFunction() {
	fmt.Printf("(p, a, q)\n\n")
	for transition, _ := range n.delta.triple {
		if transition.letter == 0 {
			fmt.Printf("(%d, ε) -> %d)\n", transition.initialState, transition.goalState)
		} else {
			fmt.Printf("(%d, %c) -> %d)\n", transition.initialState, transition.letter, transition.goalState)
		}
	}

	fmt.Printf("p ->  []\n\n")
	for s, tr := range n.delta.transitions {
		fmt.Printf(" %d -> [", s)
		for _, t := range tr {
			if t.GetLetter() == 0 {
				fmt.Printf("(%d, ε), ", t.GetState())
			} else {
				fmt.Printf("(%d, %c), ", t.GetState(), t.GetLetter())
			}
		}
		fmt.Printf("]\n\n")
	}
}

func (n *ENDFA) Dot(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()
	fmt.Fprintf(f, "digraph gs {\n")
	for transition, _ := range n.delta.triple {
		if transition.letter == 0 {
			fmt.Fprintf(f, "%d -> %d [label=\"ε\"];\n", transition.initialState, transition.goalState)

		} else {

			fmt.Fprintf(f, "%d -> %d [label=\"%c\"];\n", transition.initialState, transition.goalState, transition.letter)
		}
	}
	fmt.Fprintf(f, "%d [style=filled,color=\"0.2 0.9 0.85\"];\n", n.finalState)
	fmt.Fprintf(f, "}\n")
}
