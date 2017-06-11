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
	file, err := os.Open("/tmp/small_dict.txt")
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

}
