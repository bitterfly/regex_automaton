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
	// if goalState == 0 {
	// 	panic("WTF")
	// }
	dt.transitionToState[*NewTransition(initialState, letter)] = goalState

	children := dt.stateToTransitions[initialState]
	if len(children) > 0 && children[len(children)-1].letter == letter {
		dt.stateToTransitions[initialState][len(children)-1].state = goalState
	} else {
		dt.stateToTransitions[initialState] = append(children, *NewTransition(goalState, letter))
	}

	if len(dt.stateToTransitions[initialState]) > 1 {
		if compareTransition(dt.stateToTransitions[initialState][len(dt.stateToTransitions[initialState])-2], *NewTransition(goalState, letter)) == 1 {
			fmt.Printf("Adding things: (%d, %c, %d)", initialState, letter, goalState)
			fmt.Printf("Previous in other map: %v\n", dt.stateToTransitions[initialState])
			fmt.Printf("Previous: %v\n", dt.stateToTransitions[initialState])
			fmt.Printf("New: %v\n", *NewTransition(goalState, letter))
			panic("New transition isn't bigger than previous")

		}
	}
}

func (dt *DeltaTransitions) removeTransition(initialState int, letter rune, goalState int) {

	state, ok := dt.transitionToState[*NewTransition(initialState, letter)]
	if ok {
		if state != goalState {
			panic("We are deleting the wrong thing")
		}
	}
	delete(dt.transitionToState, *NewTransition(initialState, letter))

	outgoing_transitions := dt.stateToTransitions[initialState]

	if compareTransition(outgoing_transitions[len(outgoing_transitions)-1], *NewTransition(goalState, letter)) != 0 {
		panic("We aren't removing last transition\n")
	}

	dt.stateToTransitions[initialState] = outgoing_transitions[:len(outgoing_transitions)-1]
	if len(dt.stateToTransitions[initialState]) == 0 {
		delete(dt.stateToTransitions, initialState)
	}
}

func (dt *DeltaTransitions) removeTransitionsFor(state int) {
	children := dt.stateToTransitions[state]
	for _, child := range children {
		delete(dt.transitionToState, *NewTransition(state, child.letter))
	}
	delete(dt.stateToTransitions, state)
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

func (dt *DeltaTransitions) commonPrefix(word []rune) ([]rune, int) {
	last_state := 1
	for index, letter := range word {
		state, ok := (dt.transitionToState[*NewTransition(last_state, letter)])
		if !ok {
			remaining := make([]rune, len(word[index:]))
			copy(remaining, word[index:])
			return remaining, last_state
		}
		last_state = state
	}
	return nil, last_state
}

func (dt *DeltaTransitions) addWord(initialState int, firstNewState int, word []rune) {
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
