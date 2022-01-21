package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bitterfly/regex_automata/dfa"
	"github.com/bitterfly/regex_automata/intersection"
	"github.com/bitterfly/regex_automata/regex"
	"github.com/bitterfly/regex_automata/rpn"
)

func printVerbose(message string, verbose bool) {
	if verbose {
		fmt.Println(message)
	}
}

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
	//f, err := os.Create("pka.prof")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()
	//=========== END =================
	//=========== Read Arguments=======

	infixPtr := flag.Bool("infix", true, "If infix is true convert expression to rpn first.")
	verbosePtr := flag.Bool("verbose", false, "If true debug messages are print.")

	var outputFile string
	flag.StringVar(&outputFile, "output", "", "Write matched words into file.")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Printf("usage: pka [OPTIONS] filename\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	dict := readWord(flag.Args()[0])

	//========= READ DICTIONARY ========
	var startTime time.Time
	var elapsed time.Duration

	startTime = time.Now()
	printVerbose("Building dictionary automaton.", *verbosePtr)
	dfa := dfa.BuildDFAFromDict(dict)
	elapsed = time.Since(startTime)
	printVerbose(fmt.Sprintf("Time elapsed: %s", elapsed), *verbosePtr)
	// dfa.DotGraph("dfa.dot")

	dict = readWord(flag.Args()[0])
	printVerbose(
		fmt.Sprintf(
			`
		Correct language: %v
		Number of states: %d
		Number of eq classes: %d
		`, dfa.CheckLanguage(dict),
			dfa.GetNumStates(),
			dfa.GetNumEqClasses()),
		*verbosePtr)
	//============= END ================

	//============ READ REGEX ==========
	scanner := bufio.NewScanner(os.Stdin)
	printVerbose("Enter regular expression.", *verbosePtr)
	for scanner.Scan() {
		expression := scanner.Text()
		if *infixPtr {
			printVerbose("Converting to rpn...", *verbosePtr)
			expression = rpn.ConvertToRpn(expression)
		}

		parser := regex.NewRegexParser(*verbosePtr)
		printVerbose("Building Regex Automaton...", *verbosePtr)
		startTime = time.Now()
		endfa := parser.Parse(expression)
		// endfa.Dot("endfa.dot")

		elapsed = time.Since(startTime)
		printVerbose(fmt.Sprintf("Time: %s\n\n", elapsed), *verbosePtr)
		printVerbose("Removing epsilon transitions...", *verbosePtr)
		startTime = time.Now()
		ndfa := endfa.RemoveEpsilonTransitions()
		elapsed = time.Since(startTime)
		printVerbose(fmt.Sprintf("Time: %s\n\n", elapsed), *verbosePtr)
		// ndfa.Dot("eps.dot")
		//============= END ================

		intersector := intersection.NewIntersector(ndfa, dfa)
		printVerbose("Running intersection...", *verbosePtr)
		startTime = time.Now()
		matched := intersector.Intersect()

		number := 0
		if outputFile != "" {
			output, err := os.Create(outputFile)

			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			for word := range matched {
				fmt.Fprintf(output, "%s\n", word)
				number += 1
			}
		} else {
			for word := range matched {
				fmt.Printf("%s\n", word)
				number += 1
			}
		}
		elapsed = time.Since(startTime)
		printVerbose(
			fmt.Sprintf(
				`Found words: %d
			Time: %s`,
				number,
				elapsed), *verbosePtr)
	}
}
