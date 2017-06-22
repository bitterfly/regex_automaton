package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bitterfly/pka/dfa"
	"github.com/bitterfly/pka/intersection"
	"github.com/bitterfly/pka/regex"
	"github.com/bitterfly/pka/rpn"
)

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
	//======== PROFILER ===============
	// f, err := os.Create("pka.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	//=========== END =================
	//=========== Read Arguments=======

	infixPtr := flag.Bool("infix", false, "If infix is true convert expression to rpn first.")
	dotPtr := flag.Bool("dot", false, "If dot is true make dot file for svg automaton.")
	var outputFile string
	flag.StringVar(&outputFile, "output", "", "puke words here inseat of stdin")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Printf("usage: pka [OPTIONS] filename, %d\n", len(os.Args))
		os.Exit(1)
	}

	dict := readWord(flag.Args()[0])

	//========= GET DICTIONARY ========
	var startTime time.Time
	var elapsed time.Duration

	startTime = time.Now()
	fmt.Printf("Building dictionary automaton.\n")

	dfa := dfa.BuildDFAFromDict(dict)
	elapsed = time.Since(startTime)
	fmt.Printf("Time: %s\n", elapsed)
	if *dotPtr {
		dfa.DotGraph("dfa.dot")
	}

	dict = readWord(flag.Args()[0])

	fmt.Printf("Correct language: %v\n", dfa.CheckLanguage(dict))
	fmt.Printf("Number of states: %d\n", dfa.GetNumStates())
	fmt.Printf("Number of eq classes: %d\n", dfa.GetNumEqClasses())
	fmt.Printf("=====================\n")
	//============= END ================

	//============ READ REGEX ==========
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter regular expression: \n")
	for scanner.Scan() {
		expression := scanner.Text()
		if *infixPtr {
			fmt.Printf("Converting to rpn...\n")
			expression = rpn.ConvertToRpn(expression)
		}

		parser := regex.NewRegexParser()
		fmt.Printf("\nBuilding Regex Automaton...\n")
		startTime = time.Now()
		endfa := parser.Parse(expression)

		elapsed = time.Since(startTime)
		fmt.Printf("Time: %s\n\n", elapsed)

		fmt.Printf("\nRemoving epsilon transitions...\n")
		startTime = time.Now()
		ndfa := endfa.RemoveEpsilonTransitions()
		elapsed = time.Since(startTime)
		fmt.Printf("Time: %s\n\n", elapsed)

		if *dotPtr {
			ndfa.Dot("ndfa.dot")
			endfa.Dot("endfa.dot")
		}
		//============= END ================

		intersector := intersection.NewIntersector(ndfa, dfa)

		fmt.Printf("\nRunning intersection...\n")
		startTime = time.Now()
		matched := intersector.Intersect()

		number := 0
		if outputFile != "" {
			output, err := os.Create(outputFile)

			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			fmt.Printf("Writing into file: %s\n", outputFile)

			for word := range matched {
				fmt.Fprintf(output, "%s\n", word)
				number += 1
			}
		} else {
			fmt.Printf("Matching words: \n")
			for word := range matched {
				fmt.Printf("%s\n", word)
				number += 1
			}
		}
		elapsed = time.Since(startTime)
		fmt.Printf("Found words: %d\n", number)
		fmt.Printf("Time: %s\n\n", elapsed)
		fmt.Printf("=====================\n")
		fmt.Printf("Enter regular expression: \n")
	}
}
