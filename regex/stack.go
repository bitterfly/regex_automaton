package regex

type (
	Stack struct {
		top    *node
		length int
	}
	node struct {
		value *ENDFA
		prev  *node
	}
)

func NewStack() *Stack {
	return &Stack{
		top:    nil,
		length: 0,
	}
}

// Return the number of items in the stack
func (this *Stack) Len() int {
	return this.length
}

// View the top item on the stack
func (this *Stack) Peek() *ENDFA {
	if this.length == 0 {
		return &ENDFA{}
	}
	return this.top.value
}

// Pop the top item of the stack and return it
func (this *Stack) Pop() *ENDFA {
	if this.length == 0 {
		return &ENDFA{}
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the stack
func (this *Stack) Push(value *ENDFA) {
	n := &node{value, this.top}
	this.top = n
	this.length++
}
