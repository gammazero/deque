package deque

import "testing"

func TestEmpty(t *testing.T) {
	var q Deque
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expect 0")
	}
}

func TestFrontBack(t *testing.T) {
	var q Deque
	q.PushBack("foo")
	q.PushBack("bar")
	q.PushBack("baz")
	if q.Front() != "foo" {
		t.Error("wrong value at front of queue")
	}
	if q.Back() != "baz" {
		t.Error("wrong value at back of queue")
	}

	if q.PopFront() != "foo" {
		t.Error("wrong value removed from front of queue")
	}
	if q.Front() != "bar" {
		t.Error("wrong value remaining at front of queue")
	}
	if q.Back() != "baz" {
		t.Error("wrong value remaining at back of queue")
	}

	if q.PopBack() != "baz" {
		t.Error("wrong value removed from back of queue")
	}
	if q.Front() != "bar" {
		t.Error("wrong value remaining at front of queue")
	}
	if q.Back() != "bar" {
		t.Error("wrong value remaining at back of queue")
	}
}

func TestGrowShrinkBack(t *testing.T) {
	var q Deque
	size := minCapacity * 2

	for i := 0; i < size; i++ {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		q.PushBack(i)
	}
	bufLen := len(q.buf)

	// Remove from back.
	for i := size; i > 0; i-- {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		x := q.PopBack()
		if x != i-1 {
			t.Error("q.PopBack() =", x, "expected", i-1)
		}
	}
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expected 0")
	}
	if len(q.buf) == bufLen {
		t.Error("queue buffer did not shrink")
	}
}

func TestGrowShrinkFront(t *testing.T) {
	var q Deque
	size := minCapacity * 2

	for i := 0; i < size; i++ {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		q.PushBack(i)
	}
	bufLen := len(q.buf)

	// Remove from Front
	for i := 0; i < size; i++ {
		if q.Len() != size-i {
			t.Error("q.Len() =", q.Len(), "expected", minCapacity*2-i)
		}
		x := q.PopFront()
		if x != i {
			t.Error("q.PopBack() =", x, "expected", i)
		}
	}
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expected 0")
	}
	if len(q.buf) == bufLen {
		t.Error("queue buffer did not shrink")
	}
}

func TestSimple(t *testing.T) {
	var q Deque

	for i := 0; i < minCapacity; i++ {
		q.PushBack(i)
	}
	for i := 0; i < minCapacity; i++ {
		if q.Front() != i {
			t.Error("peek", i, "had value", q.Front())
		}
		x := q.PopFront()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}

	q.Clear()
	for i := 0; i < minCapacity; i++ {
		q.PushFront(i)
	}
	for i := minCapacity - 1; i >= 0; i-- {
		x := q.PopFront()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}
}

func TestBufferWrap(t *testing.T) {
	var q Deque

	for i := 0; i < minCapacity; i++ {
		q.PushBack(i)
	}

	for i := 0; i < 3; i++ {
		q.PopFront()
		q.PushBack(minCapacity + i)
	}

	for i := 0; i < minCapacity; i++ {
		if q.Front().(int) != i+3 {
			t.Error("peek", i, "had value", q.Front())
		}
		q.PopFront()
	}
}

func TestBufferWrapReverse(t *testing.T) {
	var q Deque

	for i := 0; i < minCapacity; i++ {
		q.PushFront(i)
	}
	for i := 0; i < 3; i++ {
		q.PopBack()
		q.PushFront(minCapacity + i)
	}

	for i := 0; i < minCapacity; i++ {
		if q.Back().(int) != i+3 {
			t.Error("peek", i, "had value", q.Front())
		}
		q.PopBack()
	}
}

func TestLen(t *testing.T) {
	var q Deque

	if q.Len() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := 0; i < 1000; i++ {
		q.PushBack(i)
		if q.Len() != i+1 {
			t.Error("adding: queue with", i, "elements has length", q.Len())
		}
	}
	for i := 0; i < 1000; i++ {
		q.PopFront()
		if q.Len() != 1000-i-1 {
			t.Error("removing: queue with", 1000-i-i, "elements has length", q.Len())
		}
	}
}

func TestBack(t *testing.T) {
	var q Deque

	for i := 0; i < minCapacity+5; i++ {
		q.PushBack(i)
		if q.Back() != i {
			t.Errorf("Back returned wrong value")
		}
	}
}

func checkRotate(t *testing.T, size int) {
	var q Deque
	for i := 0; i < size; i++ {
		q.PushBack(i)
	}

	for i := 0; i < q.Len(); i++ {
		x := i
		for n := 0; n < q.Len(); n++ {
			if q.At(n) != x {
				t.Fatalf("a[%d] != %d after rotate and copy", n, x)
			}
			x++
			if x == q.Len() {
				x = 0
			}
		}
		q.Rotate(1)
		if q.Back().(int) != i {
			t.Fatal("wrong value during rotation")
		}
	}
	for i := q.Len() - 1; i >= 0; i-- {
		q.Rotate(-1)
		if q.Front().(int) != i {
			t.Fatal("wrong value during reverse rotation")
		}
	}
}

func TestRotate(t *testing.T) {
	checkRotate(t, 10)
	checkRotate(t, minCapacity)
	checkRotate(t, minCapacity+minCapacity/2)

	var q Deque
	for i := 0; i < 10; i++ {
		q.PushBack(i)
	}
	q.Rotate(11)
	if q.Front() != 1 {
		t.Error("rotating 11 places should have been same as one")
	}
	q.Rotate(-21)
	if q.Front() != 0 {
		t.Error("rotating -21 places should have been same as one -1")
	}
	q.Rotate(q.Len())
	if q.Front() != 0 {
		t.Error("should not have rotated")
	}
	q.Clear()
	q.PushBack(0)
	q.Rotate(13)
	if q.Front() != 0 {
		t.Error("should not have rotated")
	}
}

func TestAt(t *testing.T) {
	var q Deque

	for i := 0; i < 1000; i++ {
		q.PushBack(i)
	}

	// Front to back.
	for j := 0; j < q.Len(); j++ {
		if q.At(j).(int) != j {
			t.Errorf("index %d doesn't contain %d", j, j)
		}
	}

	// Back to front
	for j := 1; j <= q.Len(); j++ {
		if q.At(q.Len()-j).(int) != q.Len()-j {
			t.Errorf("index %d doesn't contain %d", q.Len()-j, q.Len()-j)
		}
	}
}

func TestClear(t *testing.T) {
	var q Deque

	for i := 0; i < 100; i++ {
		q.PushBack(i)
	}
	if q.Len() != 100 {
		t.Error("push: queue with 100 elements has length", q.Len())
	}
	cap := len(q.buf)
	q.Clear()
	if q.Len() != 0 {
		t.Error("empty queue length not 0 after clear")
	}
	if len(q.buf) != cap {
		t.Error("queue capacity changed after clear")
	}

	// Check that there are no remaining references after Clear()
	for i := 0; i < len(q.buf); i++ {
		if q.buf[i] != nil {
			t.Error("queue has non-nil deleted elements after Clear()")
			break
		}
	}
}

func TestInsert(t *testing.T) {
	q := new(Deque)
	for _, x := range "ABCDEFG" {
		q.PushBack(x)
	}
	insert(q, 4, 'x') // ABCDxEFG
	if q.At(4) != 'x' {
		t.Error("expected x at position 4")
	}

	insert(q, 2, 'y') // AByCDxEFG
	if q.At(2) != 'y' {
		t.Error("expected y at position 2")
	}
	if q.At(5) != 'x' {
		t.Error("expected x at position 5")
	}

	insert(q, 0, 'b') // bAByCDxEFG
	if q.Front() != 'b' {
		t.Error("expected b inserted at front")
	}

	insert(q, q.Len(), 'e') // bAByCDxEFGe

	for i, x := range "bAByCDxEFGe" {
		if q.PopFront() != x {
			t.Error("expected", x, "at position", i)
		}
	}
}

func TestRemove(t *testing.T) {
	q := new(Deque)
	for _, x := range "ABCDEFG" {
		q.PushBack(x)
	}

	if remove(q, 4) != 'E' { // ABCDFG
		t.Error("expected E from position 4")
	}

	if remove(q, 2) != 'C' { // ABDFG
		t.Error("expected C at position 2")
	}
	if q.Back() != 'G' {
		t.Error("expected G at back")
	}

	if remove(q, 0) != 'A' { // BDFG
		t.Error("expected to remove A from front")
	}
	if q.Front() != 'B' {
		t.Error("expected G at back")
	}

	if remove(q, q.Len()-1) != 'G' { // BDF
		t.Error("expected to remove G from back")
	}
	if q.Back() != 'F' {
		t.Error("expected F at back")
	}

	if q.Len() != 3 {
		t.Error("wrong length")
	}
}

func TestFrontBackOutOfRangePanics(t *testing.T) {
	const msg = "should panic when peeking empty queue"
	var q Deque
	assertPanics(t, msg, func() {
		q.Front()
	})
	assertPanics(t, msg, func() {
		q.Back()
	})

	q.PushBack(1)
	q.PopFront()

	assertPanics(t, msg, func() {
		q.Front()
	})
	assertPanics(t, msg, func() {
		q.Back()
	})
}

func TestPopFrontOutOfRangePanics(t *testing.T) {
	var q Deque

	assertPanics(t, "should panic when removing empty queue", func() {
		q.PopFront()
	})

	q.PushBack(1)
	q.PopFront()

	assertPanics(t, "should panic when removing emptied queue", func() {
		q.PopFront()
	})
}

func TestPopBackOutOfRangePanics(t *testing.T) {
	var q Deque

	assertPanics(t, "should panic when removing empty queue", func() {
		q.PopBack()
	})

	q.PushBack(1)
	q.PopBack()

	assertPanics(t, "should panic when removing emptied queue", func() {
		q.PopBack()
	})
}

func TestAtOutOfRangePanics(t *testing.T) {
	var q Deque

	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	assertPanics(t, "should panic when negative index", func() {
		q.At(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		q.At(4)
	})
}

func TestInsertOutOfRangePanics(t *testing.T) {
	q := new(Deque)

	assertPanics(t, "should panic when inserting out of range", func() {
		insert(q, 1, "X")
	})

	q.PushBack("A")

	assertPanics(t, "should panic when inserting at negative index", func() {
		insert(q, -1, "Y")
	})

	assertPanics(t, "should panic when inserting out of range", func() {
		insert(q, 2, "B")
	})
}

func TestRemoveOutOfRangePanics(t *testing.T) {
	q := new(Deque)

	assertPanics(t, "should panic when removing from empty queue", func() {
		remove(q, 0)
	})

	q.PushBack("A")

	assertPanics(t, "should panic when removing at negative index", func() {
		remove(q, -1)
	})

	assertPanics(t, "should panic when removing out of range", func() {
		remove(q, 1)
	})
}

func assertPanics(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: didn't panic as expected", name)
		}
	}()

	f()
}

func BenchmarkPushFront(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushFront(i)
	}
}

func BenchmarkPushBack(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
}

func BenchmarkSerial(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	for i := 0; i < b.N; i++ {
		q.PopFront()
	}
}

func BenchmarkSerialReverse(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushFront(i)
	}
	for i := 0; i < b.N; i++ {
		q.PopBack()
	}
}

func BenchmarkRotate(b *testing.B) {
	q := new(Deque)
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	// N complete rotations on length N - 1.
	for i := 0; i < b.N; i++ {
		q.Rotate(b.N - 1)
	}
}

func BenchmarkInsert(b *testing.B) {
	q := new(Deque)
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insert(q, q.Len()/2, -i)
	}
}

func BenchmarkRemove(b *testing.B) {
	q := new(Deque)
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		remove(q, q.Len()/2)
	}
}

// insert is used to insert an element into the middle of the queue, before the
// element at the specified index.  insert(0,e) is the same as PushFront(e) and
// insert(Len(),e) is the same as PushBack(e).  Complexity is constant plus
// linear in the lesser of the distances between i and either of the ends of
// the queue.  Accepts only non-negative index values, and panics if index is
// out of range.
func insert(q *Deque, i int, elem interface{}) {
	if i < 0 || i > q.Len() {
		panic("deque: Insert() called with index out of range")
	}
	if i == 0 {
		q.PushFront(elem)
		return
	}
	if i == q.Len() {
		q.PushBack(elem)
		return
	}
	if i <= q.Len()/2 {
		q.Rotate(i)
		q.PushFront(elem)
		q.Rotate(-i)
	} else {
		rots := q.Len() - i
		q.Rotate(-rots)
		q.PushBack(elem)
		q.Rotate(rots)
	}
}

// remove removes and returns an element from the middle of the queue, at the
// specified index.  remove(0) is the same as PopFront() and remove(Len()-1) is
// the same as PopBack().  Complexity is constant plus linear in the lesser of
// the distances between i and either of the ends of the queue.  Accepts only
// non-negative index values, and panics if index is out of range.
func remove(q *Deque, i int) interface{} {
	if i < 0 || i >= q.Len() {
		panic("deque: Remove() called with index out of range")
	}
	if i == 0 {
		return q.PopFront()
	}
	if i == q.Len()-1 {
		return q.PopBack()
	}
	if i <= q.Len()/2 {
		q.Rotate(i)
		elem := q.PopFront()
		q.Rotate(-i)
		return elem
	}
	rots := q.Len() - 1 - i
	q.Rotate(-rots)
	elem := q.PopBack()
	q.Rotate(rots)
	return elem
}
