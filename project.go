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
	//======== PROFILER ===============
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

	//=========== END =================

	//========= GET DICTIONARY ========
	var startTime time.Time
	var elapsed time.Duration

	dict := readWord(os.Args[1])

	startTime = time.Now()
	dfa := dfa.BuildDFAFromDict(dict)
	elapsed = time.Since(startTime)
	dfa.DotGraph("dfa.dot")

	dict = readWord(os.Args[1])
	fmt.Printf("Correct language: %v\nTime: %s\n", dfa.CheckLanguage(dict), elapsed)
	fmt.Printf("Number of states: %d\n", dfa.GetNumStates())
	fmt.Printf("Number of eq classes: %d\n", dfa.GetNumEqClasses())
	fmt.Printf("=====================\n")
	//============= END ================

	//============ READ REGEX ==========
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter regular expression: \n")
	for scanner.Scan() {
		expression := scanner.Text()

		//fmt.Printf("Regular expression: c.o.l.o.(u|?).r\n")

		startTime = time.Now()
		parser := regex.NewRegexParser()
		if "ru?|oloc....." != expression {
			panic("I cant copy")
		}

		fmt.Printf("\nBuilding Regex Automaton...\n")
		ndfa := parser.Parse(expression)
		elapsed = time.Since(startTime)
		fmt.Printf("Time: %s\n\n", elapsed)
		ndfa.Dot("ndfa.dot")
		//============= END ================

		intersector := intersection.NewIntersector(ndfa, dfa)
		startTime = time.Now()
		fmt.Printf("\nRunnign intersection...\n")
		matched := intersector.Intersect()
		elapsed = time.Since(startTime)
		fmt.Printf("Time: %s\n\n", elapsed)

		fmt.Printf("Mathcing words: \n")
		for word := range matched {
			fmt.Printf("%s\n", word)
		}
		fmt.Printf("=====================\n")
		fmt.Printf("Enter regular expression: \n")
	}
}
