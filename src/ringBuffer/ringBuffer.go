package ringBuffer

import (
	"fmt"
)

// RingBuffer is simple FIFO implementation of ring array
type RingBuffer struct {
	size   int
	data   []*int
	start  int
	end    int
	isFull bool
}

// NewBuffer returns new instance of RingBuffer.
// In case of incorrect size returns error
func NewBuffer(size int) (*RingBuffer, error) {
	if size < 1 {
		return nil, fmt.Errorf("wrong buffer size: must be greater than zero")
	}
	return &RingBuffer{
		size:   size,
		data:   make([]*int, size),
		start:  0,
		end:    0,
		isFull: false,
	}, nil
}

//IsEmpty tells about if buffer currently empty
func (r *RingBuffer) IsEmpty() bool {
	return r.start == r.end && !r.isFull
}

// IsFull tells about if buffer currently full
func (r *RingBuffer) IsFull() bool {
	return r.isFull
}

// Read extracts the value from the buffer.
// Returns error if buffer is empty
func (r *RingBuffer) Read() (int, error) {
	if r.IsEmpty() {
		return 0, fmt.Errorf("buffer is empty")
	}

	// read logic
	f := func(el *int) (*int, bool) {
		// increment start index
		// if index is on the edge of underlying array
		// we set start index to beginning
		if r.start != r.size-1 {
			r.start++
		} else {
			r.start = 0
		}
		return el, el == nil
	}

	// supply traversal function
	// with our logic
	el := r.traverse(f)

	if r.isFull {
		r.isFull = false
	}

	return *el, nil
}

// Write writes provided value to the buffer.
// Returns error if buffer is already full
func (r *RingBuffer) Write(v int) error {
	if r.isFull {
		return fmt.Errorf("buffer is full")
	}

	r.data[r.end] = &v

	// increment end index
	// if index is on the edge of underlying array
	// we set the end index to beginning
	if r.end < r.size-1 {
		r.end++
	} else {
		r.end = 0
	}

	if r.end == r.start {
		r.isFull = true
	}

	return nil
}

// Print prints buffer contents to console
func (r *RingBuffer) Print() {
	if r.IsEmpty() {
		fmt.Println("buffer is empty")
		return
	}

	// print logic
	f := func(el *int) (*int, bool) {
		if el == nil {
			return el, false
		}
		fmt.Printf("%d\t", *el)
		return el, true
	}

	// supply traversal function
	// with our logic
	r.traverse(f)

	fmt.Printf("\n")
}

// traverse iterates through all elements
// of the buffer from start index to the end index
// and applies provided function to each element.
// Provided function must return pointer to processed value
// (whatever this function has done with it)
// and a flag that determines whether further traversal is needed
func (r *RingBuffer) traverse(f func(*int) (*int, bool)) *int {

	steps := r.steps()
	var (
		el *int
		ok bool
	)
	for i := r.start; steps > 0; i, steps = i+1, steps-1 {

		if el, ok = f(r.data[i]); !ok {
			return el
		}

		if i == r.size-1 {
			i = -1
		}
	}
	return el
}

// steps calculates the amount of steps
// needed to iterate through the buffer
// elements. It depends on current positions
// of start and end indexes
func (r *RingBuffer) steps() int {

	s := r.size - r.start + r.end
	if r.start < r.end {
		s = r.end - r.start
	}

	return s
}
