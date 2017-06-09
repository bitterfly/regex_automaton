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

	// dict := []string{"babite", "babo", "babu"}
	// test := dfa.BuildDFAFromDict(dict)
	// test.Print()

	children := *dfa.NewChildren(*dfa.NewTransition(2, 'a'))
	ec := *dfa.NewEquivalenceClass(true, children)
	en := *dfa.NewEquivalenceNode(1, ec)

	root := dfa.NewEquivalenceTree(en)

	children_2 := *dfa.NewChildren(*dfa.NewTransition(3, 'a'))
	ec_2 := *dfa.NewEquivalenceClass(true, children_2)
	en_2 := *dfa.NewEquivalenceNode(2, ec_2)

	ec_3 := *dfa.NewEquivalenceClass(false, children_2)
	en_3 := *dfa.NewEquivalenceNode(3, ec_3)

	// fmt.Printf("Expect -1, found: %d\n", ec.Compare(ec_2))
	// fmt.Printf("Expect 0, found: %d\n", ec.Compare(ec))
	// fmt.Printf("Expect 1, found: %d\n", ec_2.Compare(ec_3))

	dfa.Insert(&root, en_2)
	dfa.Insert(&root, en_3)

	children_3 := *dfa.NewChildren(*dfa.NewTransition(2, 'b'))
	ec_4 := *dfa.NewEquivalenceClass(true, children_3)
	en_4 := *dfa.NewEquivalenceNode(4, ec_4)

	fmt.Printf("Found? %v\n", root.Find(en_2))
	fmt.Printf("Found? %v\n", root.Find(en_3))
	fmt.Printf("Found? %v\n", root.Find(en))
	fmt.Printf("Found? %v\n", root.Find(en_4))
}
