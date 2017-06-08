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

func (t *Transition) String() string {
	return fmt.Sprintf("(%c, %d)", t.letter, t.state)
}

type DeltaTransitions struct {
	transitionToState map[Transition]int
	stateToTransition map[int][]Transition
}

func NewDeltaTransitions(transitionToState map[Transition]int) *DeltaTransitions {
	stateToTransition := make(map[int][]Transition)

	for transition, goalState := range transitionToState {
		transitions, ok := stateToTransition[transition.state]
		newTransition := *NewTransition(goalState, transition.letter)
		if !ok {
			stateToTransition[transition.state] = []Transition{newTransition}
		} else {
			stateToTransition[transition.state] = append(transitions, newTransition)
		}
	}

	return &DeltaTransitions{
		transitionToState: transitionToState,
		stateToTransition: stateToTransition,
	}
}

func (dt *DeltaTransitions) AddTransition(initialState int, letter rune, goalState int) {
	dt.transitionToState[*NewTransition(initialState, letter)] = goalState

	transitions, ok := dt.stateToTransition[initialState]
	newTransition := *NewTransition(goalState, letter)
	if !ok {
		dt.stateToTransition[initialState] = []Transition{newTransition}
	} else {
		dt.stateToTransition[initialState] = append(transitions, newTransition)
	}

}

func (dt *DeltaTransitions) traverse(word string) (bool, int) {
	state := 1
	ok := true
	for _, letter := range word {
		state, ok = (dt.transitionToState[*NewTransition(state, letter)])
		if !ok {
			return false, -1
		}
	}
	return true, state
}

func (dt *DeltaTransitions) commonPrefix(word string) (string, int) {
	last_state := 1
	for index, letter := range word {
		state, ok := (dt.transitionToState[*NewTransition(last_state, letter)])
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
			dt.AddTransition(initialState, letter, currentState)
		} else {
			dt.AddTransition(currentState, letter, currentState+1)
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

func EmptyAutomaton() *DFA {
	return &DFA{
		maxState:    1,
		finalStates: nil,
		delta:       *NewDeltaTransitions(make(map[Transition]int)),
	}
}

func BuildDFAFromDict(dict []string) {
	// var checked []int = nil
	dfa := EmptyAutomaton()
	for _, word := range dict {
		remaining, lastState := dfa.delta.commonPrefix(word)
		dfa.AddWord(lastState, remaining)
	}
}

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
	fmt.Printf("(p, a) -> q\n")
	for transition, goalState := range d.delta.transitionToState {
		fmt.Printf("(%d, %c) -> %d)\n", transition.state, transition.letter, goalState)
	}
	fmt.Printf("\np -> (a, q)\n")
	for initialState, transition := range d.delta.stateToTransition {
		fmt.Printf("%d -> %v\n", initialState, transition)
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
