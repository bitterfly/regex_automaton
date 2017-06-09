package dfa

import "fmt"

//states are consecutive numbers
//start state is always 1
type DFA struct {
	maxState    int
	finalStates map[int]struct{}
	delta       DeltaTransitions
}

func NewDFA(maxState int, _finalStates []int, delta map[Transition]int) *DFA {
	finalStates := make(map[int]struct{})
	for state := range _finalStates {
		finalStates[state] = struct{}{}
	}

	return &DFA{
		maxState:    maxState,
		finalStates: finalStates,
		delta:       *NewDeltaTransitions(delta),
	}
}

func EmptyAutomaton() *DFA {
	return &DFA{
		maxState:    1,
		finalStates: make(map[int]struct{}),
		delta:       *NewDeltaTransitions(make(map[Transition]int)),
	}
}

func BuildDFAFromDict(dict []string) *DFA {
	checked := make(map[int]struct{})
	dfa := EmptyAutomaton()

	for _, word := range dict {
		//fmt.Printf("word: %s\n", word)
		remaining, lastState := dfa.delta.commonPrefix(word)
		//fmt.Printf("remaining: %s\n", remaining)
		//fmt.Printf("last_state: %d\n", lastState)

		if dfa.delta.hasChildren(lastState) {
			dfa.reduce(lastState, &checked)
		}
		dfa.AddWord(lastState, remaining)
		//fmt.Printf("DFA: \n")
		//dfa.Print()
		//fmt.Printf("==============\n\n")

	}
	dfa.reduce(1, &checked)
	return dfa
}

func (d *DFA) reduce(state int, checked *map[int]struct{}) {
	//fmt.Printf("Enter reduce with %d\n", state)
	child := d.delta.stateToTransitions[state].lastChild
	if d.delta.hasChildren(child.state) {
		d.reduce(child.state, checked)
	}

	found_equivalent := false

	//fmt.Printf("Checked %v\n", checked)
	for checked_state, _ := range *checked {
		if d.checkEquivalentStates(checked_state, child.state) {
			//fmt.Printf("found euqivalent: %d %d\n", checked_state, child.state)
			found_equivalent = true
			d.delta.removeTransition(state, child.letter, child.state, *NewTransition(checked_state, child.letter))
			d.removeState(child.state)
			d.delta.addTransition(state, child.letter, checked_state)
		}
	}
	if !found_equivalent {
		//fmt.Printf("Adding child: %d\n", child.state)
		(*checked)[child.state] = struct{}{}
	}
}

func (d *DFA) AddWord(state int, word string) {
	d.addNewStates(len(word))
	d.finalStates[d.maxState] = struct{}{}
	d.delta.addWord(state, d.maxState-len(word)+1, word)
}

func (d *DFA) isFinal(state int) bool {
	_, ok := d.finalStates[state]
	return ok
}

func (d *DFA) checkEquivalentStates(first int, second int) bool {
	return (d.isFinal(first) == d.isFinal(second)) &&
		(d.delta.compareOutgoing(first, second))
}

func (d *DFA) addNewStates(number int) {
	d.maxState += number
}

func (d *DFA) removeState(state int) {
	if d.isFinal(state) {
		delete(d.finalStates, state)
	}
}

//===========================Human Friendly======================================

func (d *DFA) Print() {
	fmt.Printf("====DFA====\n")
	fmt.Printf("Max: %d, Final: %v\n", d.maxState, d.finalStates)
	d.PrintFunction()
	fmt.Printf("\n====AFD====\n")
}

func (d *DFA) PrintFunction() {
	fmt.Printf("(p, a) -> q\n\n")
	for transition, goalState := range d.delta.transitionToState {
		fmt.Printf("(%d, %c) -> %d)\n", transition.state, transition.letter, goalState)
	}
	fmt.Printf("\np -> (a, q)\n\n")
	for initialState, children := range d.delta.stateToTransitions {
		fmt.Printf("%d -> ", initialState)
		for child, _ := range children.children {
			fmt.Printf("%s ", child.String())
		}
		fmt.Printf("\n")
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

func (d *DFA) Check() {
	fmt.Printf("Equiv 1:1 %v\n", d.checkEquivalentStates(1, 1))
	fmt.Printf("Equiv 1:2 %v\n", d.checkEquivalentStates(1, 2))
}
