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
		if q.PeekAt(i) != i {
			t.Errorf("q.PeekAt(%d) = %d, expected %d", i, q.PeekAt(i), i)
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

func TestDequeSimple(t *testing.T) {
	q := new(Deque)

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

func TestDequeWrap(t *testing.T) {
	q := new(Deque)

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

func TestDequeWrapReverse(t *testing.T) {
	q := new(Deque)

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

func TestDequeLen(t *testing.T) {
	q := new(Deque)

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

func TestDequePeekAt(t *testing.T) {
	q := new(Deque)

	for i := 0; i < 1000; i++ {
		q.PushBack(i)
		for j := 0; j < q.Len(); j++ {
			if q.PeekAt(j).(int) != j {
				t.Errorf("index %d doesn't contain %d", j, j)
			}
		}
	}
}

func TestDequePeekAtNegative(t *testing.T) {
	q := new(Deque)

	for i := 0; i < 1000; i++ {
		q.PushBack(i)
		for j := 1; j <= q.Len(); j++ {
			if q.PeekAt(-j).(int) != q.Len()-j {
				t.Errorf("index %d doesn't contain %d", -j, q.Len()-j)
			}
		}
	}
}

func TestDequeBack(t *testing.T) {
	q := new(Deque)

	for i := 0; i < minCapacity+5; i++ {
		q.PushBack(i)
		if q.Back() != i {
			t.Errorf("Back returned wrong value")
		}
	}
}

func TestCopy(t *testing.T) {
	q := new(Deque)
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
	q := new(Deque)
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

func TestDequeClear(t *testing.T) {
	q := new(Deque)

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

func TestPeekAtOutOfRangePanics(t *testing.T) {
	q := new(Deque)

	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	assertPanics(t, "should panic when negative index", func() {
		q.PeekAt(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		q.PeekAt(4)
	})
}

func TestFrontBackOutOfRangePanics(t *testing.T) {
	const msg = "should panic when peeking empty queue"
	q := new(Deque)
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
	q := new(Deque)

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
	q := new(Deque)

	assertPanics(t, "should panic when removing empty queue", func() {
		q.PopBack()
	})

	q.PushBack(1)
	q.PopBack()

	assertPanics(t, "should panic when removing emptied queue", func() {
		q.PopBack()
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

// Size (number of items) of Deque to use for benchmarks.
const size = minCapacity + (minCapacity / 2)

func BenchmarkPushFront(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			q.PushFront(n)
		}
	}
}

func BenchmarkPushBack(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			q.PushBack(n)
		}
	}
}

func BenchmarkSerial(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			q.PushBack(i)
		}
		for n := 0; n < size; n++ {
			x := q.Front()
			if q.PopFront() != x {
				panic("bad PopFront()")
			}
		}
	}
}

func BenchmarkSerialReverse(b *testing.B) {
	var q Deque
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			q.PushFront(i)
		}
		for n := 0; n < size; n++ {
			x := q.Back()
			if q.PopBack() != x {
				panic("bad PopBack()")
			}
		}
	}
}

func BenchmarkRotate(b *testing.B) {
	q := new(Deque)
	for i := 0; i < size; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < size; j++ {
			v := q.PopFront()
			q.PushBack(v)
		}
	}
}

func BenchmarkRotateReverse(b *testing.B) {
	q := new(Deque)
	for i := 0; i < size; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < size; j++ {
			v := q.PopBack()
			q.PushFront(v)
		}
	}
}

func BenchmarkDequePeekAt(b *testing.B) {
	q := new(Deque)
	for i := 0; i < size; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < size; j++ {
			q.PeekAt(j)
		}
	}
}

func BenchmarkDequePushPop(b *testing.B) {
	q := new(Deque)
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			q.PushBack(nil)
			q.PopFront()
		}
	}
}
