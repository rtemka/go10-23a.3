package worker

import (
	"testing"
)

func TestFilter(t *testing.T) {
	f := func(x int) (int, bool) {
		return x, x > 100
	}

	w := NewIntProcessor(1, f)

	in := make(chan int)
	defer close(in)

	done := make(chan struct{})
	defer close(done)

	out := w.Filter(done, in)

	x := 101
	in <- x

	v, ok := <-out
	if !ok {
		t.Fatal("incorrect result: expected to recieve value, got nothing")
	}
	if v != x {
		t.Fatalf("incorrect result: expected to recieve value: %d, got %d", x, v)
	}

	x = 99
	in <- x
	done <- struct{}{}

	if v, ok := <-out; ok {
		t.Fatalf("incorrect result: expected to recieve nothing, got %d", v)
	}

}
