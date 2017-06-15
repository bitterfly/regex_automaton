package dfa

import (
	"fmt"
	"os"
	"sort"
)

//states are consecutive numbers
//start state is always 1
type DFA struct {
	maxState     int
	finalStates  map[int]struct{}
	delta        DeltaTransitions
	NumStates    int
	NumEqClasses int
}

func NewDFA(maxState int, _finalStates []int, delta map[Transition]int) *DFA {
	finalStates := make(map[int]struct{})
	for _, state := range _finalStates {
		finalStates[state] = struct{}{}
	}

	return &DFA{
		maxState:     maxState,
		finalStates:  finalStates,
		delta:        *NewDeltaTransitions(delta),
		NumStates:    1,
		NumEqClasses: 1,
	}
}

func EmptyAutomaton() *DFA {
	return &DFA{
		maxState:     1,
		finalStates:  make(map[int]struct{}),
		delta:        *NewDeltaTransitions(make(map[Transition]int)),
		NumStates:    0,
		NumEqClasses: 0,
	}
}

func BuildDFAFromDict(dict <-chan string) *DFA {
	var checked *EquivalenceTree
	dfa := EmptyAutomaton()

	//	i := 0
	for word := range dict {
		// if i%1000 == 0 {
		// 	fmt.Printf("%d\n", i)
		// }
		// i += 1

		runeWord := make([]rune, 0, len(word))
		for _, r := range word {
			runeWord = append(runeWord, r)
		}

		remaining, lastState := dfa.delta.commonPrefix(runeWord)

		if dfa.delta.hasChildren(lastState) {
			dfa.reduce(lastState, &checked)
		}

		if len(remaining) == 0 {
			dfa.makeFinal(lastState)
		} else {
			dfa.AddWord(lastState, remaining)
		}
	}
	dfa.reduce(1, &checked)
	return dfa
}

func (d *DFA) reduce(state int, checked **EquivalenceTree) {
	children := d.delta.stateToTransitions[state]
	child := children[len(children)-1]
	if d.delta.hasChildren(child.state) {
		d.reduce(child.state, checked)
	}

	childEquivalenceClass := *NewEquivalenceClass(d.isFinal(child.state), d.delta.getChildren(child.state))
	childEquivalenceNode := *NewEquivalenceNode(child.state, childEquivalenceClass)

	checked_state, ok := (*checked).find(childEquivalenceNode)
	if checked_state == child.state {
		return
	}

	if ok {
		d.delta.removeTransition(state, child.letter, child.state)
		d.removeState(child.state)

		d.delta.addTransition(state, child.letter, checked_state)
	} else {
		(*checked) = insert((*checked), &childEquivalenceNode)
		//fmt.Printf("Tree after insert:\n %s \n", (*checked).print())
		d.NumEqClasses += 1
	}
}

func (d *DFA) AddWord(state int, word []rune) {
	d.addNewStates(len(word))
	d.makeFinal(d.maxState)
	d.delta.addWord(state, d.maxState-len(word)+1, word)
	d.NumStates += len(word)
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

func (d *DFA) makeFinal(state int) {
	d.finalStates[state] = struct{}{}
}

func (d *DFA) removeState(state int) {
	if d.isFinal(state) {
		delete(d.finalStates, state)
	}
	d.delta.removeTransitionsFor(state)
	d.NumStates -= 1
}

//===========================Human Friendly======================================
func (d *DFA) CountStates() {
	states := make(map[int]struct{})

	for tr, i := range d.delta.transitionToState {
		states[tr.state] = struct{}{}
		states[i] = struct{}{}
	}

	states_2 := make(map[int]struct{})
	for i, tr_2 := range d.delta.stateToTransitions {
		states_2[i] = struct{}{}
		for _, t := range tr_2 {
			states_2[t.state] = struct{}{}
		}
	}

	fmt.Printf("First function unique states: %v\n", states)
	fmt.Printf("Second function unique states: %v\n", states_2)
}

func (d *DFA) sortedFinalStates() []int {
	var states []int
	for k, _ := range d.finalStates {
		states = append(states, k)
	}
	sort.Ints(states)
	return states
}

func (d *DFA) Print() {
	fmt.Printf("====DFA====\n")
	fmt.Printf("Max: %d, Final: %v\n", d.maxState, d.sortedFinalStates())
	d.PrintFunction()
	fmt.Printf("\n====AFD====\n")
}

func (d *DFA) DotGraph(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()
	fmt.Fprintf(f, "digraph gs {\n")
	for transition, goalState := range d.delta.transitionToState {
		fmt.Fprintf(f, "%d -> %d [label=\"%c\"];\n", transition.state, goalState, transition.letter)
	}
	for finalState, _ := range d.finalStates {
		fmt.Fprintf(f, "%d [style=filled,color=\"0.2 0.9 0.85\"];\n", finalState)
	}
	fmt.Fprintf(f, "}\n")

}

func (d *DFA) PrintFunction() {
	fmt.Printf("(p, a) -> q\n\n")

	for transition, goalState := range d.delta.transitionToState {
		fmt.Printf("(%d, %c) -> %d)\n", transition.state, transition.letter, goalState)
	}

	fmt.Printf("\np -> (a, q)\n\n")
	for initialState, children := range d.delta.stateToTransitions {
		fmt.Printf("%d -> [", initialState)
		for _, child := range children {
			fmt.Printf("(%c, %d),", child.letter, child.state)
		}
		fmt.Printf("]\n")
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

func (d *DFA) CheckLanguage(dict <-chan string) bool {
	for word := range dict {
		ok, state := d.delta.traverse(word)
		if !ok {
			fmt.Printf("No transition: %s\n", word)
			return false
		}
		if !d.isFinal(state) {
			fmt.Printf("First failing word: %s\n", word)
			return false
		}
	}
	return true
}

func (d *DFA) GetMaxState() int {
	return d.maxState
}
