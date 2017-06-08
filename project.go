package main

import (
	"github.com/bitterfly/pka/dfa"
)

func main() {
	alphabet := "abc"
	states := 3
	finishStates := []int{3}

	delta := map[dfa.Transition]int{
		*dfa.NewTransition(1, 'a'): 2,
		*dfa.NewTransition(1, 'b'): 1,
		*dfa.NewTransition(2, 'c'): 3,
	}

	first := dfa.NewDFA(alphabet, states, finishStates, delta)
	first.PrintFunction()
	// first.Traverse("bba")
	// first.Traverse("bbac")
	// first.Traverse("bbaca")
	first.FindCommonPrefix("baba")
	first.FindCommonPrefix("pliok")
	first.FindCommonPrefix("aca")

}
