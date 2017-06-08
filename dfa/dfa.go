package dfa

import "fmt"

type Transition struct {
	letter rune
	state  int
}

func NewTransition(state int, letter rune) *Transition {
	return &Transition{letter: letter, state: state}
}

type DeltaTransitions struct {
	data map[Transition]int
}

func NewDeltaTransitions(data map[Transition]int) *DeltaTransitions {
	return &DeltaTransitions{
		data: data,
	}
}

func (dt *DeltaTransitions) traverse(word string) (bool, int) {
	state := 1
	ok := true
	for _, letter := range word {
		state, ok = (dt.data[*NewTransition(state, letter)])
		if !ok {
			return false, -1
		}
	}
	return true, state
}

func (dt *DeltaTransitions) commonPrefix(word string) (string, int) {
	last_state := 1
	for index, letter := range word {
		state, ok := (dt.data[*NewTransition(last_state, letter)])
		if !ok {
			return word[index:], last_state
		}
		last_state = state
	}
	return "", last_state
}

//states are consecutive numbers
//start state is always 1
type DFA struct {
	alphabet    string
	maxState    int
	finalStates []int
	delta       DeltaTransitions
}

func NewDFA(alphabet string, maxState int, finalStates []int, delta map[Transition]int) *DFA {
	return &DFA{
		alphabet:    alphabet,
		maxState:    maxState,
		finalStates: finalStates,
		delta:       *NewDeltaTransitions(delta),
	}
}

func (d *DFA) PrintFunction() {
	for k, v := range d.delta.data {
		fmt.Printf("(%d, %c, %v)\n", k.state, k.letter, v)
	}
}

func (d *DFA) Traverse(word string) {
	ok, state := d.delta.traverse(word)
	if !ok {
		fmt.Printf("Not in the automation - %s\n", word)
	} else {
		fmt.Printf("%s leads to %d\n", word, state)
	}
}

func (d *DFA) FindCommonPrefix(word string) {
	remaining, state := d.delta.commonPrefix(word)
	fmt.Printf("Word: %s\nRemaining: %s, last_state: %d\n\n", word, remaining, state)
}
