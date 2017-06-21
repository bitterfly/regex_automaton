package intersection

import (
	"github.com/bitterfly/pka/dfa"
	"github.com/bitterfly/pka/regex"
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
		wordSoFar = make([]rune, 0)
		i.intersect(map[int]struct{}{i.ndfa.GetInitialState(): struct{}{}}, 1, wordSoFar, matched)
		close(matched)
	}()
	return matched
}

func (i *Intersector) intersect(ndfaStates map[int]struct{}, dfaState int, wordSoFar *[]rune, matched chan string) {
	ndfaStates, isFinal := i.ndfa.EpsilonClosure(ndfaStates)

	if isFinal && i.dfa.IsFinal(dfaState) {
		matched <- String(wordSoFar)
	}

	ndfaTransitions := i.ndfa.GetNonEpsilonTransitions(ndfaStates)
	for _, transitions := range i.dfa.GetTransitions(dfaState) {
		ndfaDestinations, ok := ndfaTransitions[transitions.GetLetter()]
		if ok {
			wordSoFar = append(wordSoFar, transitions.GetLetter())
			i.intersect(ndfaDestinations, transitions.GetState(), wordSoFar, matched)
		}
	}

	wordSoFar = wordSoFar[0 : len(wordSoFar)-1]
}
