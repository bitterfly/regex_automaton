package main

import (
	"fmt"

	"github.com/bitterfly/project/ndfa"
)

func main() {
	alphabet := []string{"a", "b", "c"}
	states := 3
	startStates := []int{1, 2}
	finishStates := []int{3}

	transition := map[ndfa.StateLetter][]int{
		*ndfa.NewStateLetter(1, "a"): []int{1, 2},
		*ndfa.NewStateLetter(1, "b"): []int{1, 2, 3},
		*ndfa.NewStateLetter(2, "c"): []int{3},
	}

	first := ndfa.NewNdfa(alphabet, states, startStates, finishStates, transition)

	alphabet = []string{"a", "b", "c"}
	states = 5
	startStates = []int{2, 4}
	finishStates = []int{5}

	transition = map[ndfa.StateLetter][]int{
		*ndfa.NewStateLetter(4, "a"): []int{1, 2, 3},
		*ndfa.NewStateLetter(3, "b"): []int{1, 2, 3, 4},
		*ndfa.NewStateLetter(5, "c"): []int{3},
	}

	second := ndfa.NewNdfa(alphabet, states, startStates, finishStates, transition)

	fmt.Println("First: ")
	first.PrintFunction()
	fmt.Println("Second: ")
	second.PrintFunction()
	fmt.Println("Union ")

	ndfa.Union(first, second)
}
