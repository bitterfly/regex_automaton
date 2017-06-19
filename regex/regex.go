package regex

import (
	"fmt"
	"os"
)

func EmptyExpressionNDFA(initialState, finalState int) *NDFA {
	delta := NewMultipleEmptyTransition()
	//            ε
	//    -> o --------> (o)
	//
	delta.addTransition(initialState, 0, finalState)

	return NewNDFA(initialState, 2, finalState, delta)
}

func LetterExpressionNDFA(initialState, finalState int, letter rune) *NDFA {
	delta := NewMultipleEmptyTransition()
	//            a
	//    -> o --------> (o)
	//
	delta.addTransition(initialState, letter, finalState)

	return NewNDFA(initialState, 2, finalState, delta)
}

func UnionExpressionsNDFA(initialState, finalState int, first, second *NDFA) *NDFA {
	delta := NewMultipleEmptyTransition()
	//                 ____________
	//                /            \   ε
	//           ε,-----> o     o -------- -,
	//           /    \____________/         \
	//   ----> o       ____________          (o)
	//           \    /            \   ε     /
	//           ε'-----> o     o  ---------'
	//                \____________/

	// new initial to old initials
	delta.addTransition(initialState, 0, first.initialState)
	delta.addTransition(initialState, 0, second.initialState)

	//all previous transitions
	delta.addTransitions(first.delta)
	delta.addTransitions(second.delta)

	// olf final to new final
	delta.addTransition(first.finalState, 0, finalState)
	delta.addTransition(second.finalState, 0, finalState)

	numStates := 2 + first.numStates + second.numStates
	return NewNDFA(initialState, numStates, finalState, delta)
}

func ConcatenateExpressionsNDFA(first, second *NDFA) *NDFA {
	delta := NewMultipleEmptyTransition()
	//       ________________           _______________
	//      /                \    ε    /               \
	//   -----> o first   o---|---------> o  second (o)|
	//		\________________/         \_______________/
	//

	//first final to second initial
	delta.addTransition(first.finalState, 0, second.initialState)

	// all previous transitions
	delta.addTransitions(first.delta)
	delta.addTransitions(second.delta)

	numStates := first.numStates + second.numStates

	return NewNDFA(first.initialState, numStates, second.finalState, delta)
}

func KleeneExpressionsNDFA(initialState, finalState int, ndfa *NDFA) *NDFA {
	delta := NewMultipleEmptyTransition()

	//
	//                      ___ε____
	//                     /________\__
	//			 ε   	 /V          \ \  ε
	//    ---> o -----> | o   NDFA    o | ----> (o)
	//          \        \_____________/
	//           \____________________________/
	//				           ε
	//

	// new initial to old initial
	delta.addTransition(initialState, 0, ndfa.initialState)

	//new initial to new final
	delta.addTransition(initialState, 0, ndfa.finalState)

	//olf final to new final
	delta.addTransition(ndfa.finalState, 0, finalState)

	// all old transitions
	delta.addTransitions(first.delta)
	delta.addTransitions(second.delta)

	//old final to old initial
	delta.addTransition(ndfa.finalState, 0, ndfa.initialState)

	numStates := first.numStates + second.numStates

	return NewNDFA(initialState, numStates, finalState, delta)
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
