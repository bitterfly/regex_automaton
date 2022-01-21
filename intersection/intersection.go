package intersection

import (
	"github.com/bitterfly/regex_automata/dfa"
	"github.com/bitterfly/regex_automata/regex"
)

type Intersector struct {
	ndfa *regex.NDFA
	dfa  *dfa.DFA
}

func NewIntersector(ndfa *regex.NDFA, dfa *dfa.DFA) *Intersector {
	return &Intersector{
		ndfa: ndfa,
		dfa:  dfa,
	}
}

func (i *Intersector) Intersect() chan string {
	matched := make(chan string)
	go func() {
		wordSoFar := make([]rune, 0)
		i.intersect(map[int]struct{}{i.ndfa.GetInitialState(): struct{}{}}, 1, &wordSoFar, matched)
		close(matched)
	}()
	return matched
}

func (i *Intersector) intersect(ndfaStates map[int]struct{}, dfaState int, wordSoFar *[]rune, matched chan string) {
	if i.ndfa.HasFinal(ndfaStates) && i.dfa.IsFinal(dfaState) {
		matched <- string(*wordSoFar)
	}

	ndfaTransitions := i.ndfa.GetDestinations(ndfaStates)

	for _, transitions := range i.dfa.GetTransitions(dfaState) {
		ndfaDestinations, ok := ndfaTransitions[transitions.GetLetter()]
		if ok {
			*wordSoFar = append(*wordSoFar, transitions.GetLetter())
			i.intersect(ndfaDestinations, transitions.GetState(), wordSoFar, matched)
		}
	}
	if len(*wordSoFar) > 0 {
		*wordSoFar = (*wordSoFar)[0 : len(*wordSoFar)-1]
	}
}
