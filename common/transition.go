package common

import "fmt"

type Transition struct {
	letter rune
	state  int
}

func (t *Transition) GetLetter() rune {
	return t.letter
}

func (t *Transition) GetState() int {
	return t.state
}

func (t *Transition) SetState(state int) {
	t.state = state
}

func NewTransition(state int, letter rune) *Transition {
	return &Transition{letter: letter, state: state}
}

func CompareTransition(first, second Transition) int {
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

func CompareTransitionSlices(first []Transition, second []Transition) int {
	for i, _ := range first {
		compareResult := CompareTransition(first[i], second[i])
		if compareResult != 0 {
			return compareResult
		}
	}
	return 0
}

func (t *Transition) String() string {
	return fmt.Sprintf("(%c, %d)", t.letter, t.state)
}
