package ndfa

import "fmt"

type StateLetter struct {
	letter string
	state  int
}

func NewStateLetter(state int, letter string) *StateLetter {
	return &StateLetter{letter: letter, state: state}
}

type Ndfa struct {
	Alphabet     []string
	NumStates    int
	StartStates  []int
	FinishStates []int
	Transition   map[StateLetter][]int
}

func NewNdfa(alphabet []string, numStates int, startStates []int, finishStates []int, transition map[StateLetter][]int) *Ndfa {

	return &Ndfa{
		Alphabet:     alphabet,
		NumStates:    numStates,
		StartStates:  startStates,
		FinishStates: finishStates,
		Transition:   transition,
	}
}

func remapArray(states []int, offset int) {
	for i := range states {
		states[i] = states[i] + offset
	}
}

func remapTransition(transition map[StateLetter][]int, offset int) map[StateLetter][]int {
	offset_transition := make(map[StateLetter][]int)
	for k, v := range transition {
		remapArray(v, offset)
		offset_transition[*NewStateLetter(k.state+offset, k.letter)] = v
	}

	return offset_transition
}

func EmptyNdfa() *Ndfa {
	return &Ndfa{
		Alphabet:     nil,
		NumStates:    0,
		StartStates:  nil,
		FinishStates: nil,
		Transition:   nil,
	}
}

func Union(first *Ndfa, second *Ndfa) *Ndfa {
	fmt.Println("First: ")
	first.PrintFunction()

	fmt.Println("Second: ")
	second.PrintFunction()

	second.Transition = remapTransition(second.Transition, first.NumStates)

	fmt.Println("Change: ")
	second.PrintFunction()

	return EmptyNdfa()
}

func (n *Ndfa) PrintFunction() {
	for k, v := range n.Transition {
		fmt.Printf("(%d, %s, %v)\n", k.state, k.letter, v)
	}
}
