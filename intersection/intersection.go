package intersection

import (
	"fmt"

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

func (i Intersector) Intersect(ndfaStates map[int]struct{}, dfaState int, wordSoFar string) {
	ndfaStates, isFinal := i.ndfa.EpsilonClosure(ndfaStates)

	if isFinal && i.dfa.IsFinal(dfaState) {
		fmt.Printf("%s\n", wordSoFar)
	}

	ndfaTransitions := i.ndfa.GetNonEpsilonTransitions(ndfaStates)
	for _, transitions := range i.dfa.GetTransitions(dfaState) {
		ndfaDestinations, ok := ndfaTransitions[transitions.GetLetter()]
		if ok {
			i.Intersect(ndfaDestinations, transitions.GetState(), wordSoFar+string(transitions.GetLetter()))
		}
	}
}
