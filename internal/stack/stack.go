package stack

// Stack is a stack interface for integers.
type Stack struct {
	Push   func(int)
	Pop    func() int
	Length func() int
}

// NewStack returns a new stack of integers.
func NewStack() Stack {
	slice := make([]int, 0)
	return Stack{
		Push: func(i int) {
			slice = append(slice, i)
		},
		Pop: func() int {
			res := slice[len(slice)-1]
			slice = slice[:len(slice)-1]
			return res
		},
		Length: func() int {
			return len(slice)
		},
	}
}
