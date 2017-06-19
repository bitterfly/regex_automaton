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

//=================================================

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

	fmt.Printf("(p, a) ->  []\n\n")
	for tr, v := range n.delta.transitions {
		if tr.GetLetter() == 0 {
			fmt.Printf("(%d, ε) -> ", tr.GetState())
		} else {
			fmt.Printf("(%d, %c) -> ", tr.GetState(), tr.GetLetter())
		}

		fmt.Printf("[")
		for _, s := range v {
			fmt.Printf("%d, ", s)
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
