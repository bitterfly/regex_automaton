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
		i.intersect(map[int]struct{}{i.ndfa.GetInitialState(): struct{}{}}, 1, "", matched)
		close(matched)
	}()
	return matched
}

func (i *Intersector) intersect(ndfaStates map[int]struct{}, dfaState int, wordSoFar string, matched chan string) {
	ndfaStates, isFinal := i.ndfa.EpsilonClosure(ndfaStates)

	if isFinal && i.dfa.IsFinal(dfaState) {
		matched <- wordSoFar
	}

	ndfaTransitions := i.ndfa.GetNonEpsilonTransitions(ndfaStates)
	for _, transitions := range i.dfa.GetTransitions(dfaState) {
		ndfaDestinations, ok := ndfaTransitions[transitions.GetLetter()]
		if ok {
			i.intersect(ndfaDestinations, transitions.GetState(), wordSoFar+string(transitions.GetLetter()), matched)
		}
	}
}
