package regex

type (
	Stack struct {
		top    *node
		length int
	}
	node struct {
		value *NDFA
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
func (this *Stack) Peek() *NDFA {
	if this.length == 0 {
		return &NDFA{}
	}
	return this.top.value
}

// Pop the top item of the stack and return it
func (this *Stack) Pop() *NDFA {
	if this.length == 0 {
		return &NDFA{}
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the stack
func (this *Stack) Push(value *NDFA) {
	n := &node{value, this.top}
	this.top = n
	this.length++
}
