package regex

import "github.com/bitterfly/pka/automaton"

type NDFA struct {
	automaton *automaton.FA
	delta     *MultipleDeltaTransitions
}
