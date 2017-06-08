package dfa

import (
	"fmt"
)

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

func EmptyAutomaton() *DFA {
	return &DFA{
		maxState:    1,
		finalStates: nil,
		delta:       *NewDeltaTransitions(make(map[Transition]int)),
	}
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

func (dt *DeltaTransitions) addWord(initialState int, firstNewState int, word string) {
	currentState := firstNewState
	for index, letter := range word {
		if index == 0 {
			dt.data[*NewTransition(initialState, letter)] = currentState
		} else {
			dt.data[*NewTransition(currentState, letter)] = currentState + 1
			currentState += 1
		}
	}
}

//states are consecutive numbers
//start state is always 1
type DFA struct {
	maxState    int
	finalStates []int
	delta       DeltaTransitions
}

func NewDFA(maxState int, finalStates []int, delta map[Transition]int) *DFA {
	return &DFA{
		maxState:    maxState,
		finalStates: finalStates,
		delta:       *NewDeltaTransitions(delta),
	}
}

// func BuildDFAFromDictionary(dictionary []string) {
// 	var checked []int = nil
// 	dfa := NewDFA(0, nil, nil)
// }

func (d *DFA) AddWord(state int, word string) {
	d.addNewStates(len(word))
	d.finalStates = append(d.finalStates, d.maxState)
	d.delta.addWord(state, d.maxState-len(word)+1, word)
}

//===========================Human Friendly======================================

func (d *DFA) Print() {
	fmt.Printf("====DFA====\n")
	fmt.Printf("Max: %d, Final: %v\n", d.maxState, d.finalStates)
	d.PrintFunction()
	fmt.Printf("\n====AFD====\n")
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

func (d *DFA) addNewStates(number int) {
	d.maxState += number
}
