package rpn

func isOperator(symbol rune) bool {
	return symbol == '|' || symbol == '*' || symbol == '.'
}

func isBigger(first, second rune) bool {
	if first == '|' {
		return true
	}

	if second == '|' {
		return false
	}

	if first == '.' {
		return true
	}

	if second == '.' {
		return true
	}

	return true
}

func ConvertToRpn(infixExpression string) string {
	stack := NewStack()
	rpnExpression := make([]rune, 0)

	for _, symbol := range infixExpression {
		switch symbol {
		case '(':
			stack.Push(symbol)
		case ')':
			for isOperator(stack.Peek()) {
				rpnExpression = append(rpnExpression, stack.Pop())
			}
			stack.Pop()
		case '|':
			fallthrough
		case '*':
			fallthrough
		case '.':
			for isOperator(stack.Peek()) && isBigger(stack.Peek(), symbol) {
				rpnExpression = append(rpnExpression, stack.Pop())
			}
			stack.Push(symbol)
		default:
			rpnExpression = append(rpnExpression, symbol)
		}
	}

	for stack.Peek() != 0 {
		rpnExpression = append(rpnExpression, stack.Pop())
	}

	return string(rpnExpression)
}

//if the token is an operator, then:
// while there is an operator at the top of the operator stack with
// 	greater than or equal to precedence:
// 		pop operators from the operator stack, onto the output queue;
// push the read operator onto the operator stack.
