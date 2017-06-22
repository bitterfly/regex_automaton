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
	"github.com/bitterfly/pka/rpn"
)

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

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
	f, err := os.Create("pka.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	//=========== END =================
	//=========== Read Arguments=======

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("usage: pka [OPTIONS] filename, %d\n", len(os.Args))
		return
	}

	infix := flag.Bool("infix", false, "If infix is true convert expression to rpn first.")
	flag.Parse()

	var dict chan string
	if *infix {
		dict = readWord(os.Args[2])
	} else {
		dict = readWord(os.Args[1])
	}

	//========= GET DICTIONARY ========
	var startTime time.Time
	var elapsed time.Duration

	startTime = time.Now()
	fmt.Printf("Building dictionary automaton.\n")

	dfa := dfa.BuildDFAFromDict(dict)
	elapsed = time.Since(startTime)
	fmt.Printf("Time: %s\n", elapsed)
	dfa.DotGraph("dfa.dot")

	if *infix {
		dict = readWord(os.Args[2])
	} else {
		dict = readWord(os.Args[1])
	}
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
		if *infix {
			fmt.Printf("Converting to rpn...\n")
			expression = rpn.ConvertToRpn(expression)
		}

		startTime = time.Now()
		parser := regex.NewRegexParser()

		fmt.Printf("\nBuilding Regex Automaton...\n")
		ndfa := parser.Parse(expression)
		epsilonless := ndfa.RemoveEpsilonTransitions()
		epsilonless.Dot("eps.dot")

		elapsed = time.Since(startTime)
		fmt.Printf("Time: %s\n\n", elapsed)
		ndfa.Dot("ndfa.dot")
		//============= END ================

		//intersector := intersection.NewIntersector(epsilonless, dfa)
		intersector := intersection.NewIntersector(ndfa, dfa)
		startTime = time.Now()
		fmt.Printf("\nRunning intersection...\n")
		matched := intersector.Intersect()

		fmt.Printf("Matching words: \n")
		number := 0
		for word := range matched {
			fmt.Printf("%s\n", word)
			number += 1
		}
		elapsed = time.Since(startTime)
		fmt.Printf("Found words: %d\n", number)
		fmt.Printf("Time: %s\n\n", elapsed)
		fmt.Printf("=====================\n")
		fmt.Printf("Enter regular expression: \n")
	}
}
