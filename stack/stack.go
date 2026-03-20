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
	var empty T
	s[len(s)-1] = empty
	return s[:len(s)-1]
}

func (s Stack[T]) Size() int {
	return len(s)
}

func (s Stack[T]) Move(from int, to int) Stack[T] {
	var buf Stack[T]
	var moved T
	for i := 0; i <= max(from, to); i++ {
		if i == from {
			moved = s.Top()
		} else {
			buf = buf.Push(s.Top())
		}
		s = s.Pop()
	}
	for i := max(from, to); i >= 0; i-- {
		if i == to {
			s = s.Push(moved)
			continue
		}
		s = s.Push(buf.Top())
		buf = buf.Pop()
	}
	return s
}
