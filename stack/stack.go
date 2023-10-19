package stack

type Stack[T any] []T

func (s Stack[T]) Top() T {
	var top T
	if len(s) > 0 {
		top = s[len(s)-1]
	}
	return top
}

func (s Stack[T]) Push(ts ...T) Stack[T] {
	return append(s, ts...)
}

func (s Stack[T]) Pop() Stack[T] {
	if len(s) == 0 {
		return s
	}
	return s[:len(s)-1]
}

func (s Stack[T]) Size() int {
	return len(s)
}
