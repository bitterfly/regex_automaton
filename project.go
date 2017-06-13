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

	start := time.Now()
	test := dfa.BuildDFAFromDict(dict)
	elapsed := time.Since(start)
	//test.Print()

	fmt.Printf("Correct language: %v\n time: %s\n", test.CheckLanguage(dict), elapsed)
	fmt.Printf("Number of states: %d\n", test.CountStates())
	fmt.Printf("Enters function: %d\n", dfa.GetTimes())
	fmt.Printf("Enters reduce: %d\n", dfa.GetTimesReduce())

}
