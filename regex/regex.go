package regex

func LetterExpressionENDFA(initialState, finalState int, letter rune) *ENDFA {
	//if letter == 0, then it's ε

	delta := NewMultipleEmptyTransition()
	//            a
	//    -> o --------> (o)
	//
	delta.addTransition(initialState, letter, finalState)

	return NewENDFA(initialState, 2, finalState, delta)
}

func UnionExpressionsENDFA(initialState, finalState int, first, second *ENDFA) *ENDFA {
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
	return NewENDFA(initialState, numStates, finalState, delta)
}

func ConcatenateExpressionsENDFA(first, second *ENDFA) *ENDFA {
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

	return NewENDFA(first.initialState, numStates, second.finalState, delta)
}

func KleeneExpressionENDFA(initialState, finalState int, ndfa *ENDFA) *ENDFA {
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
	delta.addTransition(initialState, 0, finalState)

	//olf final to new final
	delta.addTransition(ndfa.finalState, 0, finalState)

	// all old transitions
	delta.addTransitions(ndfa.delta)

	//old final to old initial
	delta.addTransition(ndfa.finalState, 0, ndfa.initialState)

	numStates := ndfa.numStates + 2

	return NewENDFA(initialState, numStates, finalState, delta)
}
