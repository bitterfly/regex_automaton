package regex

type RegexParser struct {
	maxState   int
	regexStack *Stack
}

func NewRegexParser() *RegexParser {
	return &RegexParser{
		maxState:   0,
		regexStack: NewStack(),
	}
}

func (p *RegexParser) NewState() int {
	p.maxState += 1
	return p.maxState
}

func (p *RegexParser) Parse(regex string) *ENDFA {
	for _, symbol := range regex {
		switch symbol {
		case '|':
			second := p.regexStack.Pop()
			first := p.regexStack.Pop()

			initialState := p.NewState()
			finalState := p.NewState()

			p.regexStack.Push(UnionExpressionsENDFA(initialState, finalState, first, second))
		case '.':
			second := p.regexStack.Pop()
			first := p.regexStack.Pop()

			p.regexStack.Push(ConcatenateExpressionsENDFA(first, second))
		case '*':
			initialState := p.NewState()
			finalState := p.NewState()

			last := p.regexStack.Pop()

			p.regexStack.Push(KleeneExpressionENDFA(initialState, finalState, last))
		case '?':
			initialState := p.NewState()
			finalState := p.NewState()

			p.regexStack.Push(LetterExpressionENDFA(initialState, finalState, 0))

		default:
			initialState := p.NewState()
			finalState := p.NewState()

			p.regexStack.Push(LetterExpressionENDFA(initialState, finalState, symbol))
		}
	}

	if p.regexStack.Len() == 0 {
		return &ENDFA{}
	} else {
		return p.regexStack.Pop()
	}
}
