package ringBuffer

import "testing"

func TestNewBuffer(t *testing.T) {
	_, err := NewBuffer(0)
	if err == nil {
		t.Error("incorrect result: expected error, got nil")
	}
	_, err = NewBuffer(5)
	if err != nil {
		t.Errorf("incorrect result: expected buffer instance, got error: %s", err)
	}
}

func TestWrite(t *testing.T) {
	buf, err := NewBuffer(1)
	if err != nil {
		t.Fatalf("incorrect result: expected buffer instance, got error: %s", err)
	}

	err = buf.Write(1)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}

	err = buf.Write(2)
	if err == nil {
		t.Fatal("incorrect result: expected error, got nil")
	}
}

func TestRead(t *testing.T) {
	buf, err := NewBuffer(2)
	if err != nil {
		t.Fatalf("incorrect result: expected buffer instance, got error: %s", err)
	}

	x, y := 1, 2

	err = buf.Write(x)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}

	err = buf.Write(y)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}

	v, err := buf.Read()
	if err != nil {
		t.Fatalf("incorrect result: expected read from buffer, got error: %s", err)
	}
	if v != x {
		t.Errorf("incorrect result: expected value: %d, got value: %d", x, v)
	}

	v, err = buf.Read()
	if err != nil {
		t.Fatalf("incorrect result: expected read from buffer, got error: %s", err)
	}
	if v != y {
		t.Errorf("incorrect result: expected value: %d, got value: %d", y, v)
	}
}

func TestIsEmpty(t *testing.T) {
	buf, err := NewBuffer(2)
	if err != nil {
		t.Fatalf("incorrect result: expected buffer instance, got error: %s", err)
	}

	if !buf.IsEmpty() {
		t.Fatal("incorrect result: expected that buffer is empty")
	}

	err = buf.Write(1)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}
	err = buf.Write(2)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}

	_, err = buf.Read()
	if err != nil {
		t.Fatalf("incorrect result: expected read from buffer, got error: %s", err)
	}

	_, err = buf.Read()
	if err != nil {
		t.Fatalf("incorrect result: expected read from buffer, got error: %s", err)
	}

	if !buf.IsEmpty() {
		t.Fatal("incorrect result: expected that buffer is empty")
	}

}

func TestIsFull(t *testing.T) {
	buf, err := NewBuffer(2)
	if err != nil {
		t.Fatalf("incorrect result: expected buffer instance, got error: %s", err)
	}

	if buf.IsFull() {
		t.Fatal("incorrect result: expected that buffer is not full")
	}

	err = buf.Write(1)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}
	err = buf.Write(2)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}

	if !buf.IsFull() {
		t.Fatal("incorrect result: expected that buffer is full")
	}

}

func Test_traverse(t *testing.T) {
	buf, err := NewBuffer(3)
	if err != nil {
		t.Fatalf("incorrect result: expected buffer instance, got error: %s", err)
	}

	err = buf.Write(1)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}
	err = buf.Write(2)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}
	err = buf.Write(3)
	if err != nil {
		t.Fatalf("incorrect result: expected write to buffer, got error: %s", err)
	}

	ints := make([]int, 0, 3)

	f := func(el *int) (*int, bool) {

		if el != nil {
			ints = append(ints, *el)
		}

		return el, el != nil
	}

	buf.traverse(f)

	if len(ints) != 3 {
		t.Fatalf("incorrect result: expected size of test slice is 3, got %d", len(ints))
	}

	if ints[0] != 1 && ints[1] != 2 && ints[2] != 3 {
		t.Fatalf("incorrect result: expected slice values: %d, %d, %d; got %d, %d, %d",
			1, 2, 3, ints[0], ints[1], ints[2])
	}

}
