package dfa

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/bitterfly/regex_automata/common"
)

type DeltaTransitions struct {
	transitionToState  map[common.Transition]int
	stateToTransitions map[int][]common.Transition
}

func NewDeltaTransitions(transitionToState map[common.Transition]int) *DeltaTransitions {
	stateToTransitions := make(map[int][]common.Transition)

	for transition, goalState := range transitionToState {
		children := stateToTransitions[transition.GetState()]
		children = append(children, *common.NewTransition(goalState, transition.GetLetter()))
	}

	for _, children := range stateToTransitions {
		sort.Slice(children, func(i, j int) bool {
			return children[i].GetLetter() < children[j].GetLetter() || children[i].GetState() < children[j].GetState()
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

func (dt *DeltaTransitions) getChildren(state int) []common.Transition {
	if dt.stateToTransitions[state] == nil {
		return []common.Transition{}
	}

	return (*dt).stateToTransitions[state]
}

func (dt *DeltaTransitions) addTransition(initialState int, letter rune, goalState int) {
	// if goalState == 0 {
	// 	panic("WTF")
	// }
	dt.transitionToState[*common.NewTransition(initialState, letter)] = goalState

	children := dt.stateToTransitions[initialState]
	if len(children) > 0 && children[len(children)-1].GetLetter() == letter {
		dt.stateToTransitions[initialState][len(children)-1].SetState(goalState)
	} else {
		dt.stateToTransitions[initialState] = append(children, *common.NewTransition(goalState, letter))
	}

	if len(dt.stateToTransitions[initialState]) > 1 {
		if common.CompareTransition(dt.stateToTransitions[initialState][len(dt.stateToTransitions[initialState])-2], *common.NewTransition(goalState, letter)) == 1 {
			fmt.Printf("Adding things: (%d, %c, %d)", initialState, letter, goalState)
			fmt.Printf("Previous in other map: %v\n", dt.stateToTransitions[initialState])
			fmt.Printf("Previous: %v\n", dt.stateToTransitions[initialState])
			fmt.Printf("New: %v\n", *common.NewTransition(goalState, letter))
			panic("New transition isn't bigger than previous")

		}
	}
}

func (dt *DeltaTransitions) removeTransition(initialState int, letter rune, goalState int) {

	state, ok := dt.transitionToState[*common.NewTransition(initialState, letter)]
	if ok {
		if state != goalState {
			panic("We are deleting the wrong thing")
		}
	}
	delete(dt.transitionToState, *common.NewTransition(initialState, letter))

	outgoing_transitions := dt.stateToTransitions[initialState]

	if common.CompareTransition(outgoing_transitions[len(outgoing_transitions)-1], *common.NewTransition(goalState, letter)) != 0 {
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
		delete(dt.transitionToState, *common.NewTransition(state, child.GetLetter()))
	}
	delete(dt.stateToTransitions, state)
}

func (dt *DeltaTransitions) traverse(word string) (bool, int) {
	state := 1
	ok := true
	for _, letter := range word {
		state, ok = (dt.transitionToState[*common.NewTransition(state, letter)])
		if !ok {
			return false, -1
		}
	}
	return true, state
}

func (dt *DeltaTransitions) commonPrefix(word []rune) ([]rune, int) {
	last_state := 1
	for index, letter := range word {
		state, ok := (dt.transitionToState[*common.NewTransition(last_state, letter)])
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
