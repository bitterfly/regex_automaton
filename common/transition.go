package common

import "fmt"

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
