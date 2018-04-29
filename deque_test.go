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

func TestGrowShrink(t *testing.T) {
	var q Deque
	for i := 0; i < minCapacity*2; i++ {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		q.PushBack(i)
	}
	// Check that all values are as expected.
	for i := 0; i < minCapacity*2; i++ {
		if q.Get(i) != i {
			t.Errorf("q.Get(%d) = %d, expected %d", i, q.Get(i), i)
		}
	}
	bufLen := len(q.buf)

	// Remove from back.
	for i := minCapacity * 2; i > 0; i-- {
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

	// Fill up queue again.
	for i := 0; i < minCapacity*2; i++ {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		q.PushBack(i)
	}
	bufLen = len(q.buf)

	// Remove from Front
	for i := 0; i < minCapacity*2; i++ {
		if q.Len() != minCapacity*2-i {
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

func TestGet(t *testing.T) {
	var q Deque

	for i := 0; i < 1000; i++ {
		q.PushBack(i)
	}

	// Front to back.
	for j := 0; j < q.Len(); j++ {
		if q.Get(j).(int) != j {
			t.Errorf("index %d doesn't contain %d", j, j)
		}
	}

	// Back to front
	for j := 1; j <= q.Len(); j++ {
		if q.Get(q.Len()-j).(int) != q.Len()-j {
			t.Errorf("index %d doesn't contain %d", q.Len()-j, q.Len()-j)
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

func TestCopy(t *testing.T) {
	var q Deque
	a := make([]interface{}, minCapacity)
	if q.Copy(a) != 0 {
		t.Error("Copied wrong size, expected 0")
	}

	for i := 0; i < minCapacity/2; i++ {
		q.PushBack(i)
		q.PopFront()
	}
	for i := 0; i < minCapacity; i++ {
		q.PushBack(i)
	}
	q.Copy(a)
	for i := range a {
		if a[i].(int) != i {
			t.Error("Copy has wrong value at position", i)
		}
	}

	a = []interface{}{}
	if q.Copy(a) != 0 {
		t.Error("Copied wrong size, expected 0")
	}

	a = make([]interface{}, q.Len()/2)
	if q.Copy(a) != len(a) {
		t.Error("Copied wrong size, expected", len(a))
	}

	a = make([]interface{}, q.Len()*2)
	if q.Copy(a) != q.Len() {
		t.Error("Copied wrong size", q.Len())
	}
}

func TestRotate(t *testing.T) {
	var q Deque
	for i := 0; i < 10; i++ {
		q.PushBack(i)
	}

	a := make([]interface{}, q.Len())
	for i := 0; i < q.Len(); i++ {
		q.Copy(a)
		x := i
		for n := range a {
			if a[n] != x {
				t.Fatalf("a[%d] != %d after rotate and copy", n, x)
			}
			x++
			if x == q.Len() {
				x = 0
			}
		}

		v := q.PopFront()
		if v.(int) != i {
			t.Fatal("wrong value during rotation")
		}
		q.PushBack(v)

	}
	for i := q.Len() - 1; i >= 0; i-- {
		v := q.PopBack()
		if v.(int) != i {
			t.Fatal("wrong value during reverse rotation")
		}
		q.PushFront(v)
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
	var q Deque
	q.PushBack("A")
	q.PushBack("B")
	q.PushBack("C")
	q.PushBack("D")
	q.PushBack("E")
	q.PushBack("F")
	q.PushBack("G")

	q.Insert(4, "x")
	if q.Get(4) != "x" {
		t.Error("expected x at position 4")
	}

	q.Insert(2, "y")
	if q.Get(2) != "y" {
		t.Error("expected y at position 2")
	}

	if q.Get(5) != "x" {
		t.Error("expected x at position 5")
	}

	q.Insert(0, "b")
	if q.Front() != "b" {
		t.Error("expected b inserted at front")
	}

	q.Insert(q.Len(), "e")
	if q.Back() != "e" {
		t.Error("expected e inserted at back")
	}
}

func TestRemove(t *testing.T) {
	var q Deque
	q.PushBack("A")
	q.PushBack("B")
	q.PushBack("C")
	q.PushBack("D")
	q.PushBack("E")
	q.PushBack("F")
	q.PushBack("G")

	if q.Remove(4) != "E" {
		t.Error("expected E from position 4")
	}
	if q.Get(4) != "F" {
		t.Error("expected F at position 4")
	}

	if q.Remove(2) != "C" {
		t.Error("expected C at position 2")
	}
	if q.Get(2) != "D" {
		t.Error("expected D at position 4")
	}

	if q.Get(4) != "G" {
		t.Error("expected G at position 4")
	}

	if q.Remove(0) != "A" {
		t.Error("expected to remove A from front")
	}

	if q.Remove(q.Len()-1) != "G" {
		t.Error("expected to remove G from back")
	}
}

func TestReplace(t *testing.T) {
	var q Deque
	q.PushBack("a")
	q.PushBack("b")
	q.PushBack("c")

	q.Replace(0, "A")
	if q.Front() != "A" {
		t.Error("expected A at front")
	}

	q.Replace(q.Len()-1, "C")
	if q.Back() != "C" {
		t.Error("expected C at back")
	}

	q.Replace(1, "-")
	if q.Get(1) != "-" {
		t.Error("expected - at position 1")
	}
}

func TestGetOutOfRangePanics(t *testing.T) {
	var q Deque

	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	assertPanics(t, "should panic when negative index", func() {
		q.Get(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		q.Get(4)
	})
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

func TestInsertOutOfRangePanics(t *testing.T) {
	var q Deque

	assertPanics(t, "should panic when inserting out of range", func() {
		q.Insert(1, "X")
	})

	q.PushBack("A")

	assertPanics(t, "should panic when inserting at negative index", func() {
		q.Insert(-1, "Y")
	})

	assertPanics(t, "should panic when inserting out of range", func() {
		q.Insert(2, "B")
	})
}

func TestRemoveOutOfRangePanics(t *testing.T) {
	var q Deque

	assertPanics(t, "should panic when removing from empty queue", func() {
		q.Remove(0)
	})

	q.PushBack("A")

	assertPanics(t, "should panic when removing at negative index", func() {
		q.Remove(-1)
	})

	assertPanics(t, "should panic when removing out of range", func() {
		q.Remove(1)
	})
}

func TestReplaceOutOfRangePanics(t *testing.T) {
	var q Deque

	assertPanics(t, "should panic when replacing in empty queue", func() {
		q.Replace(0, "x")
	})

	q.PushBack("A")

	assertPanics(t, "should panic when replacing at negative index", func() {
		q.Replace(-1, "Z")
	})

	assertPanics(t, "should panic when replacing out of range", func() {
		q.Replace(1, "Y")
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
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	// N complete rotations on length N.
	for i := 0; i < b.N; i++ {
		for j := 0; j < b.N; j++ {
			q.PushBack(q.PopFront())
		}
	}
}

func BenchmarkInsert(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Insert(q.Len()/2, -i)
	}
}

func BenchmarkRemove(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Remove(q.Len() / 2)
	}
}
