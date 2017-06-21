package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/bitterfly/pka/dfa"
	"github.com/bitterfly/pka/intersection"
	"github.com/bitterfly/pka/regex"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func readWord(fileName string) chan string {
	dict := make(chan string, 1000)
	go func() {
		defer close(dict)

		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			dict <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()
	return dict
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if len(os.Args) != 2 {
		fmt.Printf("usage: pka filename\n")
		return
	}

	//===================================

	parser := regex.NewRegexParser()
	fmt.Printf("(b.a).(a|b|t|d|n)*\n")
	ndfa := parser.Parse("abtdn||||*ab..")
	// fmt.Printf("(b.a).(c|d)*\n")
	// ndfa := parser.Parse("cd|*ab..")
	ndfa.Dot("ndfa.dot")
	// ndfa.Print()
	// ec, f := ndfa.EpsilonClosure(map[int]struct{}{5: struct{}{}})
	// fmt.Printf("Find epsilon closure for: %d - %v\n", ndfa.GetInitialState(), ec)
	// fmt.Printf("Does it contain final state? %v\n", f)

	// epsilon := regex.EmptyExpressionNDFA(3, 4)
	// epsilon.Print()

	// letter := regex.LetterExpressionNDFA(5, 6, 'a')
	// letter.Print()

	// kleene := regex.KleeneExpressionNDFA(4, 7, letter)
	// kleene.Print()
	// kleene.Dot("a.dot")

	// union := regex.UnionExpressionsNDFA(2, epsilon, letter)

	// epsilon2 := regex.EmptyExpressionNDFA(8)
	// doubleUnion := regex.UnionExpressionsNDFA(1, epsilon2, union)

	// epsilon3 := regex.EmptyExpressionNDFA(11)
	// concatenation := regex.ConcatenateExpressionsNDFA(doubleUnion, epsilon3)
	// concatenation.Print()
	// concatenation.Dot("a.dot")

	//=====================

	dict := readWord(os.Args[1])

	start_time := time.Now()

	dfa := dfa.BuildDFAFromDict(dict)
	elapsed := time.Since(start_time)
	dfa.DotGraph("dfa.dot")
	//dfa.Print()

	// dict = readWord(os.Args[1])
	fmt.Printf("Correct language: %v\nTime: %s\n", dfa.CheckLanguage(dict), elapsed)
	// //fmt.Printf("Is minimal? %v\n", (i == eq_c))
	// fmt.Printf("Number of states: %d\n", dfa.GetNumStates())
	// fmt.Printf("Number of eq classes: %d\n", dfa.GetNumEqClasses())

	// fmt.Printf("Check real minimality: %v\n", dfa.CheckMinimal())

	//==================================
	intersector := intersection.NewIntersector(ndfa, dfa)
	intersector.Intersect()
}
