package regex

import (
	"fmt"
	"os"

	"github.com/bitterfly/pka/automaton"
)

type NDFA struct {
	automaton *automaton.FA
	delta     *MultipleDeltaTransitions
}

func EmptyExpressionNDFA(initialState int) *NDFA {
	automaton := automaton.NewFA(initialState+1, 2, 2, []int{initialState + 1})
	transitions := make(map[MultipleTransition]struct{})
	transitions[*NewMultipleTransition(initialState, 0, initialState+1)] = struct{}{}

	return &NDFA{
		automaton: automaton,
		delta:     NewMultipleDeltaTransitions(transitions),
	}
}

func LetterExpressionNDFA(initialState int, letter rune) *NDFA {
	automaton := automaton.NewFA(initialState+1, 2, 2, []int{initialState + 1})
	transitions := make(map[MultipleTransition]struct{})
	transitions[*NewMultipleTransition(initialState, letter, initialState+1)] = struct{}{}

	return &NDFA{
		automaton: automaton,
		delta:     NewMultipleDeltaTransitions(transitions),
	}
}

func UnionExpressionsNDFA(initialState int, first, second *NDFA) *NDFA {
	// /func NewFA(maxState int, numStates int, eqClasses int, _finalStates []int)
	numberOfStates := 2 + first.automaton.NumStates + second.automaton.NumStates
	maxState := initialState + numberOfStates - 1

	automaton := automaton.NewFA(maxState, numberOfStates, 0, []int{maxState})
	delta := NewMultipleEmptyTransition()

	delta.addTransition(initialState, 0, first.automaton.GetInitialState())
	delta.addTransition(initialState, 0, second.automaton.GetInitialState())

	delta.addTransitions(first.delta)
	delta.addTransitions(second.delta)

	delta.addFinalStates(first.automaton.FinalStates, maxState)
	delta.addFinalStates(second.automaton.FinalStates, maxState)

	return &NDFA{
		automaton: automaton,
		delta:     delta,
	}
}

func (n *NDFA) MoveToInitial(initialState int) {
	offset := n.automaton.GetInitialState() - initialState

	fmt.Printf("Offset: %d", offset)
	n.automaton.MaxState -= offset

	newFinalStates := make(map[int]struct{}, len(n.automaton.FinalStates))
	for state, _ := range n.automaton.FinalStates {
		newFinalStates[state-offset] = struct{}{}
	}

	n.automaton.FinalStates = newFinalStates

	newTransitions := make(map[MultipleTransition]struct{}, len(n.delta.transitions))

	for transition, _ := range n.delta.transitions {
		newTransitions[*NewMultipleTransition(transition.initialState-offset, transition.letter, transition.goalState-offset)] = struct{}{}
	}
	n.delta.transitions = newTransitions
}

func ConcatenateExpressionsNDFA(first, second *NDFA) *NDFA {
	delta := NewMultipleEmptyTransition()
	second.MoveToInitial(first.automaton.GetFinalState())
	delta.addTransitions(first.delta)
	delta.addTransitions(second.delta)

	// /func NewFA(maxState int, numStates int, eqClasses int, _finalStates []int)
	numStates := first.automaton.NumStates + second.automaton.NumStates
	automaton := automaton.NewFA(second.automaton.MaxState, numStates, 0, []int{second.automaton.GetFinalState()})

	return &NDFA{
		automaton: automaton,
		delta:     delta,
	}
}

//=================================================

func (n *NDFA) Print() {
	fmt.Printf("====NDFA====\n")
	fmt.Printf("Max: %d, NumState: %d, NumEQClasses: %d,  Final: %v\n", n.automaton.MaxState, n.automaton.NumStates, n.automaton.NumEqClasses, n.automaton.SortedFinalStates())
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
	for finalState, _ := range n.automaton.FinalStates {
		fmt.Fprintf(f, "%d [style=filled,color=\"0.2 0.9 0.85\"];\n", finalState)
	}
	fmt.Fprintf(f, "}\n")

}
