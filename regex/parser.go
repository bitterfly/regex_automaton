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
			// fmt.Println("Union\n")
			// fmt.Printf("Pop 2 from stack")

			second := p.regexStack.Pop()
			first := p.regexStack.Pop()

			initialState := p.NewState()
			finalState := p.NewState()

			union := UnionExpressionsENDFA(initialState, finalState, first, second)
			//union.Print()

			// fmt.Printf("Pushing into stack\n")
			p.regexStack.Push(union)
		case '.':
			// fmt.Println("Concatenate\n")
			// fmt.Printf("Pop 2 from stack\n")

			second := p.regexStack.Pop()
			first := p.regexStack.Pop()

			// fmt.Printf("Pushing into stack\n")
			concatenation := ConcatenateExpressionsENDFA(first, second)
			//concatenation.Print()

			p.regexStack.Push(concatenation)
		case '*':
			// fmt.Println("Kleene\n")
			// fmt.Printf("Pop 1 from stack\n")

			initialState := p.NewState()
			finalState := p.NewState()

			last := p.regexStack.Pop()
			kleene := KleeneExpressionENDFA(initialState, finalState, last)
			//kleene.Print()

			// fmt.Printf("Pushing into stack\n")
			p.regexStack.Push(kleene)
		case '?':
			// fmt.Println("Epsilon\n")
			// fmt.Printf("Pushing into stack\n")

			initialState := p.NewState()
			finalState := p.NewState()

			eps := EmptyExpressionENDFA(initialState, finalState)
			//eps.Print()
			p.regexStack.Push(eps)

		default:
			initialState := p.NewState()
			finalState := p.NewState()

			// fmt.Printf("Pushing into stack %c\n", symbol)
			letter := LetterExpressionENDFA(initialState, finalState, symbol)
			//letter.Print()
			p.regexStack.Push(letter)
		}
	}

	if p.regexStack.Len() == 0 {
		return &ENDFA{}
	} else {
		return p.regexStack.Pop()
	}
}
