package regex

import (
	"fmt"
	"os"
)

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

func (n *NDFA) EpsilonClosure(states map[int]struct{}) (map[int]struct{}, bool) {
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

func (n *NDFA) GetNonEpsilonTransitions(states map[int]struct{}) map[rune]map[int]struct{} {
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

func (n *NDFA) GetInitialState() int {
	return n.initialState
}

func (n *NDFA) Print() {
	fmt.Printf("====NDFA====\n")
	fmt.Printf("Initial: %d, NumStates: %d, FinalState: %d\n", n.initialState, n.numStates, n.finalState)
	n.PrintFunction()
	fmt.Printf("\n====NAFD====\n")
}

func (n *NDFA) PrintFunction() {
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

func (n *NDFA) Dot(filename string) {
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
