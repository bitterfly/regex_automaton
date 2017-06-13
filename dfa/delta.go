package dfa

import (
	"fmt"
	"reflect"
	"sort"
)

type Transition struct {
	letter rune
	state  int
}

func NewTransition(state int, letter rune) *Transition {
	return &Transition{letter: letter, state: state}
}

func compareTransition(first, second Transition) int {
	if first.letter != second.letter {
		if first.letter > second.letter {
			return 1
		} else {
			return -1
		}
	}

	if first.state != second.state {
		if first.state > second.state {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

func compareTransitionSlices(first []Transition, second []Transition) int {
	for i, _ := range first {
		compareResult := compareTransition(first[i], second[i])
		if compareResult != 0 {
			return compareResult
		}
	}
	return 0
}

func (t *Transition) String() string {
	return fmt.Sprintf("(%c, %d)", t.letter, t.state)
}

type DeltaTransitions struct {
	transitionToState  map[Transition]int
	stateToTransitions map[int][]Transition
}

func NewDeltaTransitions(transitionToState map[Transition]int) *DeltaTransitions {
	stateToTransitions := make(map[int][]Transition)

	for transition, goalState := range transitionToState {
		fmt.Printf("transition: %v\n", transition)
		children := stateToTransitions[transition.state]
		children = append(children, *NewTransition(goalState, transition.letter))
	}

	for _, children := range stateToTransitions {
		sort.Slice(children, func(i, j int) bool {
			return children[i].letter < children[j].letter || children[i].state < children[j].state
		})
	}

	return &DeltaTransitions{
		transitionToState:  transitionToState,
		stateToTransitions: stateToTransitions,
	}
}

func (dt *DeltaTransitions) hasChildren(state int) bool {
	return len(dt.stateToTransitions[state]) != 0
}

func (dt *DeltaTransitions) getChildren(state int) []Transition {
	if dt.stateToTransitions[state] == nil {
		return []Transition{}
	}

	return (*dt).stateToTransitions[state]
}

func (dt *DeltaTransitions) addTransition(initialState int, letter rune, goalState int) {
	dt.transitionToState[*NewTransition(initialState, letter)] = goalState

	children := dt.stateToTransitions[initialState]
	dt.stateToTransitions[initialState] = append(children, *NewTransition(goalState, letter))
}

func (dt *DeltaTransitions) removeTransition(initialState int, letter rune, goalState int, newLastChild Transition) {
	delete(dt.transitionToState, *NewTransition(initialState, letter))
	outgoing_transitions := dt.stateToTransitions[initialState]

	if compareTransition(outgoing_transitions[len(outgoing_transitions)-1], *NewTransition(goalState, letter)) != 0 {
		panic("We aren't removing last transition\n")
	}

	outgoing_transitions = outgoing_transitions[:len(outgoing_transitions)-1]
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
			dt.addTransition(initialState, letter, currentState)
		} else {
			dt.addTransition(currentState, letter, currentState+1)
			currentState += 1
		}
	}
}

func (dt *DeltaTransitions) numOutgoing(state int) int {
	return len(dt.stateToTransitions[state])
}

func (dt *DeltaTransitions) compareOutgoing(first int, second int) bool {
	first_transitions, first_ok := dt.stateToTransitions[first]
	second_transitions, second_ok := dt.stateToTransitions[second]

	if first_ok != second_ok {
		return false
	}

	return reflect.DeepEqual(first_transitions, second_transitions)
}
