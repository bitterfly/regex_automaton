package dfa

import (
	"fmt"
	"reflect"
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

type Children struct {
	children  map[Transition]struct{}
	lastChild Transition
}

func NewChildren(child Transition) *Children {
	return &Children{
		children:  map[Transition]struct{}{child: struct{}{}},
		lastChild: child,
	}
}

func (c *Children) addChild(child Transition) {
	c.children[child] = struct{}{}
	c.lastChild = child
}

type DeltaTransitions struct {
	transitionToState  map[Transition]int
	stateToTransitions map[int]*Children
}

func NewDeltaTransitions(transitionToState map[Transition]int) *DeltaTransitions {
	stateToTransitions := make(map[int]*Children)

	for transition, goalState := range transitionToState {
		_, ok := stateToTransitions[transition.state]
		newTransition := *NewTransition(goalState, transition.letter)
		if !ok {
			stateToTransitions[transition.state] = NewChildren(newTransition)
		} else {
			stateToTransitions[transition.state].addChild(newTransition)
		}
	}

	return &DeltaTransitions{
		transitionToState:  transitionToState,
		stateToTransitions: stateToTransitions,
	}
}

func (dt *DeltaTransitions) hasChildren(state int) bool {
	children, ok := dt.stateToTransitions[state]
	if !ok {
		return false
	}

	if children.children == nil {
		return false
	}

	return len(children.children) != 0
}

func (dt *DeltaTransitions) addTransition(initialState int, letter rune, goalState int) {
	dt.transitionToState[*NewTransition(initialState, letter)] = goalState

	_, ok := dt.stateToTransitions[initialState]
	newTransition := *NewTransition(goalState, letter)
	if !ok {
		dt.stateToTransitions[initialState] = NewChildren(newTransition)
	} else {
		dt.stateToTransitions[initialState].addChild(newTransition)
	}
}

func (dt *DeltaTransitions) removeTransition(initialState int, letter rune, goalState int, newLastChild Transition) {
	delete(dt.transitionToState, *NewTransition(initialState, letter))
	outgoing_transitions := dt.stateToTransitions[initialState]
	delete(outgoing_transitions.children, *NewTransition(goalState, letter))
	outgoing_transitions.lastChild = newLastChild
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
	transitions, ok := dt.stateToTransitions[state]
	if !ok {
		return 0
	} else {
		return len(transitions.children)
	}
}

func (dt *DeltaTransitions) compareOutgoing(first int, second int) bool {
	first_transitions, first_ok := dt.stateToTransitions[first]
	second_transitions, second_ok := dt.stateToTransitions[second]

	if first_ok != second_ok {
		return false
	}

	return reflect.DeepEqual(first_transitions, second_transitions)
}
