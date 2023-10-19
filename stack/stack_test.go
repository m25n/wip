package stack_test

import (
	"github.com/m25n/wip/stack"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStack(t *testing.T) {
	t.Run("stacks start empty", func(t *testing.T) {
		var s stack.Stack[int]

		require.Zero(t, s.Size())
	})

	t.Run("stacks can grow", func(t *testing.T) {
		var s stack.Stack[int]
		s = s.Push(1).Push(2).Push(3)

		require.Equal(t, 3, s.Size())
	})

	t.Run("stacks can shrink", func(t *testing.T) {
		var s stack.Stack[int]
		s = s.Push(1).Push(2).Push(3).Pop().Pop().Pop()

		require.Zero(t, s.Size())
	})

	t.Run("popping an empty stack does nothing", func(t *testing.T) {
		var s stack.Stack[int]
		s = s.Pop()

		require.Zero(t, s.Size())
	})

	t.Run("the top of a stack is the most recently pushed value", func(t *testing.T) {
		var s stack.Stack[int]
		s = s.Push(1).Push(2).Push(3)

		require.Equal(t, 3, s.Top())
	})

	t.Run("the top of a stack is the empty value by default", func(t *testing.T) {
		var s stack.Stack[int]

		require.Empty(t, s.Top())
	})
}
