package regex

import (
	"fmt"
	"os"
)

func EmptyExpressionNDFA(initialState int) *NDFA {
	transitions := make(map[MultipleTransition]struct{})
	transitions[*NewMultipleTransition(initialState, 0, initialState+1)] = struct{}{}

	return NewNDFA(initialState, 2, initialState+1, NewMultipleDeltaTransitions(transitions))
}

func LetterExpressionNDFA(initialState int, letter rune) *NDFA {
	transitions := make(map[MultipleTransition]struct{})
	transitions[*NewMultipleTransition(initialState, letter, initialState+1)] = struct{}{}

	return NewNDFA(initialState, 2, initialState+1, NewMultipleDeltaTransitions(transitions))
}

func UnionExpressionsNDFA(initialState int, first, second *NDFA) *NDFA {
	numStates := 2 + first.numStates + second.numStates
	newFinalState := initialState + numStates - 1

	delta := NewMultipleEmptyTransition()

	delta.addTransition(initialState, 0, first.initialState)
	delta.addTransition(initialState, 0, second.initialState)

	delta.addTransitions(first.delta)
	delta.addTransitions(second.delta)

	delta.addTransition(first.finalState, 0, newFinalState)
	delta.addTransition(second.finalState, 0, newFinalState)

	return NewNDFA(initialState, numStates, newFinalState, delta)
}

// func (n *NDFA) MoveToInitial(initialState int) {
// 	offset := n.automaton.GetInitialState() - initialState

// 	fmt.Printf("Offset: %d", offset)
// 	n.automaton.MaxState -= offset

// 	newFinalStates := make(map[int]struct{}, len(n.automaton.FinalStates))
// 	for state, _ := range n.automaton.FinalStates {
// 		newFinalStates[state-offset] = struct{}{}
// 	}

// 	n.automaton.FinalStates = newFinalStates

// 	newTransitions := make(map[MultipleTransition]struct{}, len(n.delta.transitions))

// 	for transition, _ := range n.delta.transitions {
// 		newTransitions[*NewMultipleTransition(transition.initialState-offset, transition.letter, transition.goalState-offset)] = struct{}{}
// 	}
// 	n.delta.transitions = newTransitions
// }

func ConcatenateExpressionsNDFA(first, second *NDFA) *NDFA {
	delta := NewMultipleEmptyTransition()

	delta.addTransition(first.finalState, 0, second.initialState)

	delta.addTransitions(first.delta)
	delta.addTransitions(second.delta)

	numStates := first.numStates + second.numStates

	return NewNDFA(first.initialState, numStates, second.finalState, delta)
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

	for transition, _ := range n.delta.transitions {
		if transition.letter == 0 {
			fmt.Printf("(%d, ε) -> %d)\n", transition.initialState, transition.goalState)
		} else {
			fmt.Printf("(%d, %c) -> %d)\n", transition.initialState, transition.letter, transition.goalState)
		}
	}
}

func (n *NDFA) Dot(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()
	fmt.Fprintf(f, "digraph gs {\n")
	for transition, _ := range n.delta.transitions {
		if transition.letter == 0 {
			fmt.Fprintf(f, "%d -> %d [label=\"ε\"];\n", transition.initialState, transition.goalState)

		} else {

			fmt.Fprintf(f, "%d -> %d [label=\"%c\"];\n", transition.initialState, transition.goalState, transition.letter)
		}
	}
	fmt.Fprintf(f, "%d [style=filled,color=\"0.2 0.9 0.85\"];\n", n.finalState)
	fmt.Fprintf(f, "}\n")

}
