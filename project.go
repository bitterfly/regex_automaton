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
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func readWord() chan string {
	dict := make(chan string, 1000)
	go func() {
		defer close(dict)

		file, err := os.Open("/tmp/dict.txt")
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

	dict := readWord()

	start := time.Now()
	test := dfa.BuildDFAFromDict(dict)
	elapsed := time.Since(start)
	//test.Print()

	dict = readWord()
	fmt.Printf("Correct language: %v\nTime: %s\n", test.CheckLanguage(dict), elapsed)
	//fmt.Printf("Is minimal? %v\n", (i == eq_c))
	fmt.Printf("Number of states: %d\n", test.NumStates)
	fmt.Printf("Number of eq classes: %d\n", test.NumEqClasses)
}
