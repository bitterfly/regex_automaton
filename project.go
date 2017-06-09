package main

import (
	"fmt"

	"github.com/bitterfly/pka/dfa"
)

func main() {
	// states := 3
	// finishStates := []int{3}

	// delta := map[dfa.Transition]int{
	// 	*dfa.NewTransition(1, 'a'): 2,
	// 	*dfa.NewTransition(1, 'b'): 1,
	// 	*dfa.NewTransition(2, 'c'): 3,
	// }

	// first := dfa.NewDFA(states, finishStates, delta)
	// first.PrintFunction()
	// // first.Traverse("bba")
	// // first.Traverse("bbac")
	// // first.Traverse("bbaca")
	// first.FindCommonPrefix("baba")
	// first.FindCommonPrefix("pliok")
	// first.FindCommonPrefix("aca")

	// test := dfa.EmptyAutomaton()
	// test.Print()
	// test.AddWord(1, "bla")
	// test.AddWord(2, "gs")
	// test.Print()
	// test.Check()

	dict := []string{"babite", "babo", "babu", "kaka", "kapaci"}
	test := dfa.BuildDFAFromDict(dict)
	test.Print()

	fmt.Printf("Correct language: %v", test.CheckLanguage(dict))
}
