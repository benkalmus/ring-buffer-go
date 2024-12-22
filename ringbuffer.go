package ringbuffer

import "errors"

// A ring buffer is a fixed queue, FIFO like structure.
// Pushing to the buffer, will add a new item to the front of the buffer
// unless the buffer is already full.
// Popping from the buffer will drop the oldest item, from the back unless
// buffer is empty.
type RingBuffer[T any] struct {
	data        []T
	front, back int
	isFull      bool
}

var (
	ErrBufferFull  = errors.New("ring buffer full")
	ErrBufferEmpty = errors.New("ring buffer empty")
)

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data: make([]T, capacity),
	}
}

func (r *RingBuffer[T]) Push(elem T) error {
	if r.front == r.back && r.isFull {
		return ErrBufferFull
	}
	r.data[r.front] = elem
	increment(&r.front, len(r.data)-1)
	if r.front == r.back {
		r.isFull = true
	}
	return nil
}

func (r *RingBuffer[T]) Pop() (T, error) {
	if r.front == r.back && !r.isFull {
		return *new(T), ErrBufferEmpty
	}
	elem := r.data[r.back]
	increment(&r.back, len(r.data)-1)
	// if pointers meet again, buffer must be empty
	if r.back == r.front {
		r.isFull = false
	}
	return elem, nil
}

func (r *RingBuffer[T]) PopAll() []T {
	result := make([]T, 0, len(r.data)-1)
	for {
		val, err := r.Pop()
		if err != nil {
			return result
		}
		result = append(result, val)
	}
}

func (r *RingBuffer[T]) Peek() (T, error) {
	if r.front == r.back && !r.isFull {
		return *new(T), ErrBufferEmpty
	}
	elem := r.data[r.back]
	return elem, nil
}

// a function that moves a pointer forward around the circle buffer
func increment(val *int, max int) {
	// if reached end of slice, return to 0
	if *val == max {
		*val = 0
		return
	}
	*val++
}
