package vardump

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntStack(t *testing.T) {
	var val int
	var ok bool
	s := newIntStack()

	assert.Equal(t, []int{}, s.Array())
	assert.False(t, s.Swap(0))

	assert.Equal(t, 0, s.Count())
	s.Push(42)
	s.Push(1337)
	assert.Equal(t, 2, s.Count())
	assert.Equal(t, []int{42, 1337}, s.Array())

	val, ok = s.Pop()
	assert.Equal(t, 1, s.Count())
	assert.Equal(t, []int{42}, s.Array())
	assert.Equal(t, 1337, val)
	assert.True(t, ok)

	s.Push(21)
	assert.Equal(t, 2, s.Count())
	assert.True(t, s.Swap(13))
	assert.Equal(t, 2, s.Count())

	val, ok = s.Pop()
	assert.Equal(t, 1, s.Count())
	assert.Equal(t, 13, val)
	assert.True(t, ok)

	val, ok = s.Pop()
	assert.Equal(t, 0, s.Count())
	assert.Equal(t, []int{}, s.Array())
	assert.Equal(t, 42, val)
	assert.True(t, ok)

	val, ok = s.Pop()
	assert.Equal(t, 0, s.Count())
	assert.False(t, ok)
	assert.False(t, s.Swap(0))
}

func TestLargeIntStack(t *testing.T) {
	s := newIntStack()
	for i := 0; i < 100; i++ {
		s.Push(i)
		assert.Equal(t, i+1, s.Count())
	}

	for i := 99; i >= 0; i-- {
		val, ok := s.Pop()
		assert.Equal(t, i, s.Count())
		assert.Equal(t, i, val)
		assert.True(t, ok)
	}
}
