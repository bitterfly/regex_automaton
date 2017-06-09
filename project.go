package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

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

	file, err := os.Open("/tmp/dict.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var dict []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dict = append(dict, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//dict := []string{"babite", "babo", "babu", "kaka", "kapaci"}
	start := time.Now()
	test := dfa.BuildDFAFromDict(dict)
	elapsed := time.Since(start)
	test.Print()

	fmt.Printf("Correct language: %v\n time: %s\n", test.CheckLanguage(dict), elapsed)
}
