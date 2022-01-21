package regex

import "fmt"

type RegexParser struct {
	maxState   int
	regexStack *Stack
	verbose    bool
}

func NewRegexParser(verbose bool) *RegexParser {
	return &RegexParser{
		maxState:   0,
		regexStack: NewStack(),
		verbose:    verbose,
	}
}

func (p *RegexParser) NewState() int {
	p.maxState += 1
	return p.maxState
}

func (p *RegexParser) Parse(regex string) *ENDFA {
	if p.verbose {
		fmt.Printf("Regex string: %s\n", regex)
	}
	for _, symbol := range regex {
		switch symbol {
		case '|':
			if p.verbose {
				fmt.Printf("Union\nPop 2 from stack.\n")
			}
			second := p.regexStack.Pop()
			first := p.regexStack.Pop()

			initialState := p.NewState()
			finalState := p.NewState()

			union := UnionExpressionsENDFA(initialState, finalState, first, second)

			if p.verbose {
				fmt.Printf("Pushing union into stack\n")
			}
			p.regexStack.Push(union)
		case '.':
			if p.verbose {
				fmt.Printf("Concatenate\nPop 2 from stack.\n")
			}

			second := p.regexStack.Pop()
			first := p.regexStack.Pop()

			if p.verbose {
				fmt.Printf("Pushing concatenation into stack\n")
			}
			concatenation := ConcatenateExpressionsENDFA(first, second)
			p.regexStack.Push(concatenation)
		case '*':
			if p.verbose {
				fmt.Printf("Kleene.\nPop 1 from stack\n")
			}

			initialState := p.NewState()
			finalState := p.NewState()

			last := p.regexStack.Pop()
			kleene := KleeneExpressionENDFA(initialState, finalState, last)

			if p.verbose {
				fmt.Printf("Pushing kleene into stack\n")
			}
			p.regexStack.Push(kleene)
		case '?':
			if p.verbose {
				fmt.Printf("Epsilon.\n Pushing into stack\n")
			}

			initialState := p.NewState()
			finalState := p.NewState()

			eps := EmptyExpressionENDFA(initialState, finalState)
			p.regexStack.Push(eps)

		default:
			initialState := p.NewState()
			finalState := p.NewState()

			if p.verbose {
				fmt.Printf("Pushing into stack %c\n", symbol)
			}
			letter := LetterExpressionENDFA(initialState, finalState, symbol)
			p.regexStack.Push(letter)
		}
	}

	if p.regexStack.Len() == 0 {
		return &ENDFA{}
	} else {
		return p.regexStack.Pop()
	}
}
