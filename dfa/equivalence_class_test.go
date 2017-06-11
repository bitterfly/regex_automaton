package dfa

import "testing"

func TestCompareNonFinalGreater(t *testing.T) {
	a := map[Transition]struct{}{
		*NewTransition(4, 'e'):  struct{}{},
		*NewTransition(14, 'i'): struct{}{},
	}
	b := map[Transition]struct{}{
		*NewTransition(4, 'e'):   struct{}{},
		*NewTransition(114, 'a'): struct{}{},
	}

	c1 := &Children{
		children:  a,
		lastChild: *NewTransition(4, 'e'),
	}

	c2 := &Children{
		children:  b,
		lastChild: *NewTransition(14, 'i'),
	}

	egc1 := *NewEquivalenceClass(true, *c1)
	egc2 := *NewEquivalenceClass(true, *c2)

	if egc1.Compare(egc2) != 1 {
		t.Errorf("Different non-final states are not properly ordered")
	}
}
