package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("remove test", func(t *testing.T) {
		l := NewList()

		x := l.PushBack("x")
		l.PushBack("y")
		l.PushBack("z")

		require.Equal(t, 3, l.Len())

		c := l.PushFront("c")
		l.PushFront("b")
		l.PushFront("a")

		require.Equal(t, 6, l.Len())

		l.Remove(x)
		l.Remove(c)

		require.Equal(t, 4, l.Len())

		elems := make([]string, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(string))
		}
		require.Equal(t, []string{"a", "b", "y", "z"}, elems)
	})

	t.Run("reverse test", func(t *testing.T) {
		l := NewList()

		c := l.PushBack("c")
		b := l.PushBack("b")
		a := l.PushBack("a")

		l.MoveToFront(c)
		l.MoveToFront(b)
		l.MoveToFront(a)

		require.Equal(t, 3, l.Len())

		elems := make([]string, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(string))
		}
		require.Equal(t, []string{"a", "b", "c"}, elems)
	})

	t.Run("Move last one to front", func(t *testing.T) {
		l := NewList()

		l.PushBack("a")
		b := l.PushBack("b")
		c := l.PushBack("c")

		require.Equal(t, 3, l.Len())
		require.Equal(t, l.Back(), c)

		l.MoveToFront(c)
		require.Equal(t, 3, l.Len())
		require.Equal(t, l.Back(), b)
		require.Equal(t, l.Front(), c)

		elems := make([]string, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(string))
		}
		require.Equal(t, []string{"c", "a", "b"}, elems)
	})
}
