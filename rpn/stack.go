package rpn

import "fmt"

type (
	Stack struct {
		top    *node
		length int
	}
	node struct {
		value rune
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
func (this *Stack) Peek() rune {
	if this.length == 0 {
		return 0
	}
	return this.top.value
}

// Pop the top item of the stack and return it
func (this *Stack) Pop() rune {
	if this.length == 0 {
		return 0
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the stack
func (this *Stack) Push(value rune) {
	n := &node{value, this.top}
	this.top = n
	this.length++
}

func (this *Stack) Print() {
	n := this.top
	for n != nil {
		fmt.Printf("%c", n.value)
		n = n.prev
	}
	fmt.Printf("\n")
}
