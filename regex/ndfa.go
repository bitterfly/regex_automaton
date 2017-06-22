package regex

import (
	"fmt"
	"os"

	"github.com/bitterfly/pka/common"
)

type NDFA struct {
	finalStates  map[int]struct{}
	initialState int
	numStates    int
	transitions  map[int][]common.Transition
}

func NewNDFA(initialState int, numStates int, finalStates map[int]struct{}, transitions map[int][]common.Transition) *NDFA {
	return &NDFA{
		finalStates:  finalStates,
		initialState: initialState,
		transitions:  transitions,
	}
}

func (n *NDFA) GetInitialState() int {
	return n.initialState
}

func (n *NDFA) GetNumStates() int {
	return n.numStates
}

func (n *NDFA) HasFinal(states map[int]struct{}) bool {
	for state, _ := range states {
		if n.isFinal(state) {
			return true
		}
	}
	return false
}

func (n *NDFA) isFinal(state int) bool {
	_, ok := n.finalStates[state]
	return ok
}

func (n *NDFA) GetDestinations(states map[int]struct{}) map[rune]map[int]struct{} {
	destinations := make(map[rune]map[int]struct{})
	for state, _ := range states {
		//get one state and find all its children
		for _, transition := range n.transitions[state] {
			_, ok := destinations[transition.GetLetter()]
			if !ok {
				destinations[transition.GetLetter()] = make(map[int]struct{})
			}
			destinations[transition.GetLetter()][transition.GetState()] = struct{}{}
		}
	}
	return destinations
}

//==============================================

func (n *NDFA) Dot(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()
	fmt.Fprintf(f, "digraph gs {\n")
	for state, trs := range n.transitions {
		for _, tr := range trs {
			if tr.GetLetter() == 0 {
				fmt.Fprintf(f, "%d -> %d [label=\"ε\"];\n", state, tr.GetState())

			} else {

				fmt.Fprintf(f, "%d -> %d [label=\"%c\"];\n", state, tr.GetState(), tr.GetLetter())
			}
		}

	}
	for state, _ := range n.finalStates {
		fmt.Fprintf(f, "%d [style=filled,color=\"0.2 0.9 0.85\"];\n", state)
	}

	fmt.Fprintf(f, "}\n")
}
