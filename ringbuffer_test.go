package ringbuffer

import (
	"reflect"
	"testing"
)

const capacity = 5

func TestRingBuffer(t *testing.T) {
	t.Run("push pop", func(t *testing.T) {
		elem := 6
		buffer := NewRingBuffer[int](capacity)

		err := buffer.Push(elem)
		AssertEqual(t, nil, err)

		got, err := buffer.Pop()
		AssertEqual(t, nil, err)

		AssertEqual(t, elem, got)
	})
	t.Run("push until full and pop until empty", func(t *testing.T) {
		elems := []int{1, 2, 3, 4, 5}
		buffer := NewRingBuffer[int](len(elems))

		for i, val := range elems {
			t.Logf("push %d", i)
			err := buffer.Push(val)
			AssertEqual(t, nil, err)
		}
		// pushing one more time should give an error
		err := buffer.Push(88)
		AssertIsError(t, err)

		for i, want := range elems {
			t.Logf("pop %d", i)
			got, err := buffer.Pop()
			AssertEqual(t, nil, err)
			AssertEqual(t, want, got)
		}
		// popping one more time should give an error
		_, err = buffer.Pop()
		AssertIsError(t, err)
	})

	t.Run("push until about halfway and pop remaining", func(t *testing.T) {
		elems := []int{1, 2, 3}
		buffer := NewRingBuffer[int](5)

		for i, val := range elems {
			t.Logf("push %d", i)
			err := buffer.Push(val)
			AssertEqual(t, nil, err)
		}

		for i, want := range elems {
			t.Logf("pop %d", i)
			got, err := buffer.Pop()
			AssertEqual(t, nil, err)
			AssertEqual(t, want, got)
		}
		// popping one more time should give an error
		_, err := buffer.Pop()
		AssertIsError(t, err)
	})

	t.Run("push twice, pop once, push again, pop remaining", func(t *testing.T) {
		elems := []int{1, 2}
		buffer := NewRingBuffer[int](5)

		for i, val := range elems {
			t.Logf("push %d", i)
			err := buffer.Push(val)
			AssertEqual(t, nil, err)
		}
		// buffer now has [1,2]

		got, err := buffer.Pop()
		AssertEqual(t, nil, err)
		AssertEqual(t, 1, got)
		// buffer now has [2]
		err = buffer.Push(3)
		AssertEqual(t, nil, err)
		err = buffer.Push(4)
		AssertEqual(t, nil, err)
		// buffer now has [2, 3, 4]

		// pop remaining
		got, err = buffer.Pop()
		AssertEqual(t, nil, err)
		AssertEqual(t, 2, got)
		got, err = buffer.Pop()
		AssertEqual(t, nil, err)
		AssertEqual(t, 3, got)
		got, err = buffer.Pop()
		AssertEqual(t, nil, err)
		AssertEqual(t, 4, got)
		// popping one more time should give an error
		_, err = buffer.Pop()
		AssertIsError(t, err)
	})
	t.Run("push pop", func(t *testing.T) {
		elems := []int{1, 2, 3, 4, 5}
		buffer := NewRingBuffer[int](len(elems))

		for i, val := range elems {
			t.Logf("push %d", i)
			err := buffer.Push(val)
			AssertEqual(t, nil, err)

			got, err := buffer.Pop()
			AssertEqual(t, nil, err)
			AssertEqual(t, val, got)
		}
		// pushing one more time should give an error
		_, err := buffer.Pop()
		AssertIsError(t, err)
	})
}

func TestRingBufferPopAll(t *testing.T) {
	t.Run("pop all one item", func(t *testing.T) {
		buffer := NewRingBuffer[int](capacity)

		err := buffer.Push(5)
		AssertEqual(t, nil, err)

		got := buffer.PopAll()
		AssertEqual(t, nil, err)
		AssertEqual(t, []int{5}, got)
	})
	t.Run("pop all on full buffer", func(t *testing.T) {
		elems := []int{1, 2, 3, 4, 5}
		buffer := NewRingBuffer[int](len(elems))

		for _, val := range elems {
			err := buffer.Push(val)
			AssertEqual(t, nil, err)
		}
		got := buffer.PopAll()
		AssertEqual(t, elems, got)
	})
}

func AssertEqual[T any](t testing.TB, want, got T) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Got %v\nWant %v\n", got, want)
	}
}

func AssertIsError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected err but is nil")
	}
}
