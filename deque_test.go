package deque

import (
	"fmt"
	"slices"
	"testing"
	"unicode"
)

func TestEmpty(t *testing.T) {
	q := new(Deque[string])
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expect 0")
	}
	if q.Cap() != 0 {
		t.Error("expected q.Cap() == 0")
	}
	idx := q.Index(func(item string) bool {
		return true
	})
	if idx != -1 {
		t.Error("should return -1 index for nil deque")
	}
	idx = q.RIndex(func(item string) bool {
		return true
	})
	if idx != -1 {
		t.Error("should return -1 index for nil deque")
	}
}

func TestNil(t *testing.T) {
	var q *Deque[int]
	if q.Len() != 0 {
		t.Error("expected q.Len() == 0")
	}
	if q.Cap() != 0 {
		t.Error("expected q.Cap() == 0")
	}
	q.Rotate(5)
	idx := q.Index(func(item int) bool {
		return true
	})
	if idx != -1 {
		t.Error("should return -1 index for nil deque")
	}
	idx = q.RIndex(func(item int) bool {
		return true
	})
	if idx != -1 {
		t.Error("should return -1 index for nil deque")
	}
}

func TestFrontBack(t *testing.T) {
	var q Deque[string]
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
	var q Deque[int]
	const size = minCapacity * 2

	for i := range size {
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
	var q Deque[int]
	const size = minCapacity * 2

	for i := range size {
		if q.Len() != i {
			t.Error("q.Len() =", q.Len(), "expected", i)
		}
		q.PushBack(i)
	}
	bufLen := len(q.buf)

	// Remove from Front
	for i := range size {
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
	var q Deque[int]

	for i := range minCapacity {
		q.PushBack(i)
	}
	if q.Front() != 0 {
		t.Fatalf("expected 0 at front, got %d", q.Front())
	}
	if q.Back() != minCapacity-1 {
		t.Fatalf("expected %d at back, got %d", minCapacity-1, q.Back())
	}

	for i := range minCapacity {
		if q.Front() != i {
			t.Error("peek", i, "had value", q.Front())
		}
		x := q.PopFront()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}

	q.Clear()
	for i := range minCapacity {
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
	var q Deque[int]

	for i := range minCapacity {
		q.PushBack(i)
	}

	for i := range 3 {
		q.PopFront()
		q.PushBack(minCapacity + i)
	}

	for i := range minCapacity {
		if q.Front() != i+3 {
			t.Error("peek", i, "had value", q.Front())
		}
		q.PopFront()
	}
}

func TestBufferWrapReverse(t *testing.T) {
	var q Deque[int]

	for i := range minCapacity {
		q.PushFront(i)
	}
	for i := range 3 {
		q.PopBack()
		q.PushFront(minCapacity + i)
	}

	for i := range minCapacity {
		if q.Back() != i+3 {
			t.Error("peek", i, "had value", q.Front())
		}
		q.PopBack()
	}
}

func TestLen(t *testing.T) {
	var q Deque[int]

	if q.Len() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := range 1000 {
		q.PushBack(i)
		if q.Len() != i+1 {
			t.Error("adding: queue with", i, "elements has length", q.Len())
		}
	}
	for i := range 1000 {
		q.PopFront()
		if q.Len() != 1000-i-1 {
			t.Error("removing: queue with", 1000-i-i, "elements has length", q.Len())
		}
	}
}

func TestBack(t *testing.T) {
	var q Deque[int]

	for i := range minCapacity + 5 {
		q.PushBack(i)
		if q.Back() != i {
			t.Errorf("Back returned wrong value")
		}
	}
}

func TestGrow(t *testing.T) {
	var q Deque[int]
	assertPanics(t, "should panic when calling Groe with invalid size", func() {
		q.Grow(-1)
	})

	q.Grow(35)
	if q.Cap() != 64 {
		t.Fatal(t, "did not grow to expected capacity")
	}

	q.Grow(55)
	if q.Cap() != 64 {
		t.Fatal(t, "expected no capacity change")
	}

	q.Grow(77)
	if q.Cap() != 128 {
		t.Fatal(t, "did not grow to expected capacity")
	}

	for i := range 127 {
		q.PushBack(i)
	}
	if q.Cap() != 128 {
		t.Fatal(t, "expected no capacity change")
	}

	q.Grow(2)
	if q.Cap() != 256 {
		t.Fatal(t, "did not grow to expected capacity")
	}
}

func TestNew(t *testing.T) {
	const minCap = 64
	q := &Deque[string]{}
	q.SetBaseCap(minCap)
	if q.Cap() != 0 {
		t.Fatal("should not have allowcated mem yet")
	}
	q.PushBack("foo")
	q.PopFront()
	if q.Len() != 0 {
		t.Fatal("Len() should return 0")
	}
	if q.Cap() != minCap {
		t.Fatalf("worng capactiy expected %d, got %d", minCap, q.Cap())
	}
	curCap := 128
	q = new(Deque[string])
	q.SetBaseCap(minCap)
	q.Grow(curCap)
	if q.Cap() != curCap {
		t.Fatalf("Cap() should return %d, got %d", curCap, q.Cap())
	}
	if q.Len() != 0 {
		t.Fatalf("Len() should return 0")
	}
	q.PushBack("foo")
	if q.Cap() != curCap {
		t.Fatalf("Cap() should return %d, got %d", curCap, q.Cap())
	}
}

func checkRotate(t *testing.T, size int) {
	var q Deque[int]
	for i := range size {
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
		if q.Back() != i {
			t.Fatal("wrong value during rotation")
		}
	}
	for i := q.Len() - 1; i >= 0; i-- {
		q.Rotate(-1)
		if q.Front() != i {
			t.Fatal("wrong value during reverse rotation")
		}
	}
}

func TestRotate(t *testing.T) {
	checkRotate(t, 10)
	checkRotate(t, minCapacity)
	checkRotate(t, minCapacity+minCapacity/2)

	var q Deque[int]
	for i := range 10 {
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
	var q Deque[int]

	for i := range 1000 {
		q.PushBack(i)
	}

	// Front to back.
	for j := 0; j < q.Len(); j++ {
		if q.At(j) != j {
			t.Errorf("wrong item at index %d", j)
		}
	}

	// Back to front
	for j := q.Len() - 1; j >= 0; j-- {
		if q.At(j) != j {
			t.Errorf("wrong item at index %d", j)
		}
	}
}

func TestSet(t *testing.T) {
	var q Deque[int]

	for i := range 1000 {
		q.PushBack(i)
		q.Set(i, i+50)
	}

	// Front to back.
	for j := 0; j < q.Len(); j++ {
		if q.At(j) != j+50 {
			t.Errorf("index %d doesn't contain %d", j, j+50)
		}
	}
}

func TestClear(t *testing.T) {
	var q Deque[int]

	for i := range 100 {
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
		if q.buf[i] != 0 {
			t.Error("queue has non-nil deleted elements after Clear()")
			break
		}
	}

	for i := range 128 {
		q.PushBack(i)
	}
	q.Clear()
	// Check that there are no remaining references after Clear()
	for i := 0; i < len(q.buf); i++ {
		if q.buf[i] != 0 {
			t.Error("queue has non-nil deleted elements after Clear()")
			break
		}
	}
}

func TestIndex(t *testing.T) {
	var q Deque[rune]
	for _, x := range "Hello, 世界" {
		q.PushBack(x)
	}
	idx := q.Index(func(item rune) bool {
		c := item
		return unicode.Is(unicode.Han, c)
	})
	if idx != 7 {
		t.Fatal("Expected index 7, got", idx)
	}
	idx = q.Index(func(item rune) bool {
		c := item
		return c == 'H'
	})
	if idx != 0 {
		t.Fatal("Expected index 0, got", idx)
	}
	idx = q.Index(func(item rune) bool {
		return false
	})
	if idx != -1 {
		t.Fatal("Expected index -1, got", idx)
	}
}

func TestRIndex(t *testing.T) {
	var q Deque[rune]
	for _, x := range "Hello, 世界" {
		q.PushBack(x)
	}
	idx := q.RIndex(func(item rune) bool {
		c := item
		return unicode.Is(unicode.Han, c)
	})
	if idx != 8 {
		t.Fatal("Expected index 8, got", idx)
	}
	idx = q.RIndex(func(item rune) bool {
		c := item
		return c == 'H'
	})
	if idx != 0 {
		t.Fatal("Expected index 0, got", idx)
	}
	idx = q.RIndex(func(item rune) bool {
		return false
	})
	if idx != -1 {
		t.Fatal("Expected index -1, got", idx)
	}
}

func TestInsert(t *testing.T) {
	q := new(Deque[rune])
	for _, x := range "ABCDEFG" {
		q.PushBack(x)
	}
	q.Insert(4, 'x') // ABCDxEFG
	if q.At(4) != 'x' {
		t.Error("expected x at position 4, got", q.At(4))
	}

	q.Insert(2, 'y') // AByCDxEFG
	if q.At(2) != 'y' {
		t.Error("expected y at position 2")
	}
	if q.At(5) != 'x' {
		t.Error("expected x at position 5")
	}

	q.Insert(0, 'b') // bAByCDxEFG
	if q.Front() != 'b' {
		t.Error("expected b inserted at front, got", q.Front())
	}

	q.Insert(q.Len(), 'e') // bAByCDxEFGe

	for i, x := range "bAByCDxEFGe" {
		if q.PopFront() != x {
			t.Error("expected", x, "at position", i)
		}
	}

	qs := new(Deque[string])
	qs.Grow(16)

	for i := 0; i < qs.Cap(); i++ {
		qs.PushBack(fmt.Sprint(i))
	}
	// deque: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
	// buffer: [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]
	for i := 0; i < qs.Cap()/2; i++ {
		qs.PopFront()
	}
	// deque: 8 9 10 11 12 13 14 15
	// buffer: [_,_,_,_,_,_,_,_,8,9,10,11,12,13,14,15]
	for i := 0; i < qs.Cap()/4; i++ {
		qs.PushBack(fmt.Sprint(qs.Cap() + i))
	}
	// deque: 8 9 10 11 12 13 14 15 16 17 18 19
	// buffer: [16,17,18,19,_,_,_,_,8,9,10,11,12,13,14,15]

	at := qs.Len() - 2
	qs.Insert(at, "x")
	// deque: 8 9 10 11 12 13 14 15 16 17 x 18 19
	// buffer: [16,17,x,18,19,_,_,_,8,9,10,11,12,13,14,15]
	if qs.At(at) != "x" {
		t.Error("expected x at position", at)
	}
	if qs.At(at) != "x" {
		t.Error("expected x at position", at)
	}

	qs.Insert(2, "y")
	// deque: 8 9 y 10 11 12 13 14 15 16 17 x 18 19
	// buffer: [16,17,x,18,19,_,_,8,9,y,10,11,12,13,14,15]
	if qs.At(2) != "y" {
		t.Error("expected y at position 2")
	}
	if qs.At(at+1) != "x" {
		t.Error("expected x at position 5")
	}

	qs.Insert(0, "b")
	// deque: b 8 9 y 10 11 12 13 14 15 16 17 x 18 19
	// buffer: [16,17,x,18,19,_,b,8,9,y,10,11,12,13,14,15]
	if qs.Front() != "b" {
		t.Error("expected b inserted at front, got", qs.Front())
	}

	qs.Insert(qs.Len(), "e")
	if qs.Cap() != qs.Len() {
		t.Fatal("Expected full buffer")
	}
	// deque: b 8 9 y 10 11 12 13 14 15 16 17 x 18 19 e
	// buffer: [16,17,x,18,19,e,b,8,9,y,10,11,12,13,14,15]
	for i, x := range []string{"16", "17", "x", "18", "19", "e", "b", "8", "9", "y", "10", "11", "12", "13", "14", "15"} {
		if qs.buf[i] != x {
			t.Error("expected", x, "at buffer position", i)
		}
	}
	for i, x := range []string{"b", "8", "9", "y", "10", "11", "12", "13", "14", "15", "16", "17", "x", "18", "19", "e"} {
		if qs.Front() != x {
			t.Error("expected", x, "at position", i, "got", qs.Front())
		}
		qs.PopFront()
	}
}

func TestRemove(t *testing.T) {
	q := new(Deque[rune])
	for _, x := range "ABCDEFG" {
		q.PushBack(x)
	}

	if q.Remove(4) != 'E' { // ABCDFG
		t.Error("expected E from position 4")
	}

	if q.Remove(2) != 'C' { // ABDFG
		t.Error("expected C at position 2")
	}
	if q.Back() != 'G' {
		t.Error("expected G at back")
	}

	if q.Remove(0) != 'A' { // BDFG
		t.Error("expected to remove A from front")
	}
	if q.Front() != 'B' {
		t.Error("expected G at back")
	}

	if q.Remove(q.Len()-1) != 'G' { // BDF
		t.Error("expected to remove G from back")
	}
	if q.Back() != 'F' {
		t.Error("expected F at back")
	}

	if q.Len() != 3 {
		t.Error("wrong length")
	}
}

func TestSwap(t *testing.T) {
	var q Deque[string]

	q.PushBack("a")
	q.PushBack("b")
	q.PushBack("c")
	q.PushBack("d")
	q.PushBack("e")

	q.Swap(0, 4)
	if q.Front() != "e" {
		t.Fatal("wrong value at front")
	}
	if q.Back() != "a" {
		t.Fatal("wrong value at back")
	}
	q.Swap(3, 1)
	if q.At(1) != "d" {
		t.Fatal("wrong value at index 1")
	}
	if q.At(3) != "b" {
		t.Fatal("wrong value at index 3")
	}
	q.Swap(2, 2)
	if q.At(2) != "c" {
		t.Fatal("wrong value at index 2")
	}

	assertPanics(t, "should panic when removing out of range", func() {
		q.Swap(1, 5)
	})
	assertPanics(t, "should panic when removing out of range", func() {
		q.Swap(5, 1)
	})
}

func TestIter(t *testing.T) {
	var q Deque[int]

	for range q.Iter() {
		t.Fatal("iterated when empty")
	}

	q.Grow(50)
	for i := range 50 {
		q.PushBack(i)
	}

	// Front to back.
	var i int
	for item := range q.Iter() {
		if item != i {
			t.Errorf("index %d contains %d", i, item)
		}
		if i == 40 {
			break
		}
		i++
	}

	assertPanics(t, "Iter must panic when deque modified during iteration", func() {
		for item := range q.Iter() {
			if item == 42 {
				q.PushBack(51)
				q.PopFront()
			}
		}
	})

	for i := 50; i < 60; i++ {
		q.PushFront(i)
	}
	items := slices.Sorted(q.Iter())
	if !slices.IsSorted(items) {
		t.Fatal("expected sorted slice")
	}
}

func TestRIter(t *testing.T) {
	var q Deque[int]

	for range q.RIter() {
		t.Fatal("iterated when empty")
	}

	q.Grow(50)
	for i := range 50 {
		q.PushBack(i)
	}

	// Back to front
	i := q.Len() - 1
	for item := range q.RIter() {
		if item != i {
			t.Fatalf("index %d contains %d", i, item)
		}
		if i == 10 {
			break
		}
		i--
	}

	assertPanics(t, "RIter must panic when deque modified during iteration", func() {
		for item := range q.RIter() {
			if item == 42 {
				q.PushBack(51)
				q.PopFront()
			}
		}
	})
}

func TestIterPopBack(t *testing.T) {
	const (
		baseCap = 32
		size    = baseCap * 5
	)
	var q Deque[int]
	q.SetBaseCap(baseCap)

	for i := range q.IterPopBack() {
		t.Fatalf("popped %d when nothig to pop", i)
	}

	q.Grow(size)
	for i := range size {
		q.PushBack(i)
	}

	last := q.Front()
	var removed, lastRm int
	for i := range q.IterPopBack() {
		removed++
		lastRm = i
	}

	if last != lastRm {
		t.Fatal("did not expose expected item list")
	}
	if removed != size {
		t.Fatal("wrong removed count")
	}
	if lastRm != 0 {
		t.Fatal("wrong last item removed")
	}
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expected 0")
	}
	if q.Cap() != baseCap {
		t.Error("capacity should return to configured base capacity")
	}

	q.Grow(size)
	for i := range size {
		q.PushBack(i)
	}
	last = q.Front()
	removed = 0
	for i := range q.IterPopBack() {
		removed++
		lastRm = i
	}
	if last != lastRm {
		t.Fatal("did not expose expected item list")
	}
	if removed != size {
		t.Fatal("wrong removed count, got", removed, " expected", size)
	}
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expected 0")
	}
	for range q.IterPopBack() {
		t.Fatal("iteration with 0 items")
	}

	for i := range 5 {
		q.PushBack(i)
	}
	c := q.Cap()
	removed = 0
	for range q.IterPopBack() {
		removed++
	}
	if removed != 5 {
		t.Fatal("wrong removed count")
	}
	if q.Cap() != c {
		t.Fatal("unexpected capacity change")
	}

	for i := range 65 {
		q.PushBack(i)
	}
	c = q.Cap()
	removed = 0
	for range q.IterPopBack() {
		removed++
		if removed == 3 {
			break
		}
	}
	if q.Cap() != c {
		t.Fatal("unexpected capacity change")
	}
}

func TestIterPopFront(t *testing.T) {
	const (
		baseCap = 32
		size    = baseCap * 5
	)
	var q Deque[int]
	q.SetBaseCap(baseCap)

	q.Grow(size)
	for i := range size {
		q.PushBack(i)
	}
	last := q.Back()
	var lastRm, removed int
	for i := range q.IterPopFront() {
		removed++
		lastRm = i
	}
	if last != lastRm {
		t.Fatal("did not expose expected item list")
	}
	if removed != size {
		t.Fatal("wrong removed count")
	}
	if lastRm != size-1 {
		t.Fatal("wrong last item removed")
	}
	if q.Len() != 0 {
		t.Error("q.Len() =", q.Len(), "expected 0")
	}
	if q.Cap() != baseCap {
		t.Error("capacity should return to configured base capacity")
	}

	for range q.IterPopFront() {
		t.Fatal("iteration with 0 items")
	}

	for i := range 5 {
		q.PushBack(i)
	}
	c := q.Cap()
	removed = 0
	for range q.IterPopFront() {
		removed++
	}
	if removed != 5 {
		t.Fatal("wrong removed count")
	}
	if q.Cap() != c {
		t.Fatal("unexpected capacity change")
	}

	for i := range 65 {
		q.PushBack(i)
	}
	c = q.Cap()
	removed = 0
	for range q.IterPopFront() {
		removed++
		if removed == 3 {
			break
		}
	}
	for range q.IterPopFront() {
		if q.Len() == 32 {
			break
		}
	}
	if q.Cap() == c {
		t.Fatal("expected capacity change")
	}
}

func TestFrontBackOutOfRangePanics(t *testing.T) {
	const msg = "should panic when peeking empty queue"
	var q Deque[int]
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
	var q Deque[int]

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
	var q Deque[int]

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
	var q Deque[int]

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

func TestSetOutOfRangePanics(t *testing.T) {
	var q Deque[int]

	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	assertPanics(t, "should panic when negative index", func() {
		q.Set(-4, 1)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		q.Set(4, 1)
	})
}

func TestInsertOutOfRangePanics(t *testing.T) {
	q := new(Deque[string])

	q.Insert(1, "A")
	if q.Front() != "A" {
		t.Fatal("expected A at front")
	}
	if q.Back() != "A" {
		t.Fatal("expected A at back")
	}

	q.Insert(-1, "B")
	if q.Front() != "B" {
		t.Fatal("expected B at front")
	}

	q.Insert(999, "C")
	if q.Back() != "C" {
		t.Fatal("expected C at back")
	}
}

func TestRemoveOutOfRangePanics(t *testing.T) {
	q := new(Deque[string])

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

func TestSetMBaseapacity(t *testing.T) {
	var q Deque[string]
	q.SetBaseCap(200)
	q.PushBack("A")
	if q.minCap != 256 {
		t.Fatal("wrong minimum capacity")
	}
	if q.Cap() != 256 {
		t.Fatal("wrong buffer size")
	}
	q.PopBack()
	if q.minCap != 256 {
		t.Fatal("wrong minimum capacity")
	}
	if q.Cap() != 256 {
		t.Fatal("wrong buffer size")
	}
	q.SetBaseCap(0)
	if q.minCap != minCapacity {
		t.Fatal("wrong minimum capacity")
	}
}

func TestCopyInSliceCopyOutSlice(t *testing.T) {
	var q Deque[int]

	out := make([]int, 5)
	n := q.CopyOutSlice(out)
	if n != 0 {
		t.Fatal("expected nothing replaced")
	}
	content1 := []int{1, 2, 3}
	q.CopyInSlice(content1)
	if q.Front() != 1 {
		t.Fatal("wrong value for Fromt")
	}
	if q.Back() != 3 {
		t.Fatal("wrong value for Back")
	}

	n = q.CopyOutSlice(out)
	if !slices.Equal(content1, out[:n]) {
		t.Fatal("did not get expected previous content, got:", out[:n])
	}
	content2 := []int{4, 5}
	q.CopyInSlice(content2)
	if q.Len() != len(content2) {
		t.Fatalf("q should have %d items", len(content2))
	}
	for i, v := range content2 {
		if q.At(i) != v {
			t.Fatal("unexpected content after replace")
		}
	}

	n = q.CopyOutSlice(out)
	if !slices.Equal(content2, out[:n]) {
		t.Fatal("did not get expected previous content, got:", out[:n])
	}
	content3 := []int{6, 7, 8, 9, 10, 11, 12}
	q.CopyInSlice(content3)
	if q.Len() != len(content3) {
		t.Fatalf("q should have %d items", len(content3))
	}
	for i, v := range content3 {
		if q.At(i) != v {
			t.Fatal("unexpected content after replace")
		}
	}

	n = q.CopyOutSlice(out)
	if n != len(out) {
		t.Fatalf("expected %d replaced items", len(out))
	}
	if !slices.Equal(content3[:n], out) {
		t.Fatal("did not get expected previous content, got:", out)
	}
	q.CopyInSlice(nil)
	if q.Len() != 0 {
		t.Fatal("q should have 0 items")
	}

	n = q.CopyOutSlice(out)
	if n != 0 {
		t.Fatal("expected 0 replaced items")
	}
	content4 := []int{13, 14}
	q.CopyInSlice(content4)
	if q.Len() != len(content4) {
		t.Fatalf("q should have %d items", len(content4))
	}

	n = q.CopyOutSlice(nil)
	if n != 0 {
		t.Fatal("expected 0 replaced items")
	}
	if q.Len() != len(content4) {
		t.Fatalf("q should have %d items", len(content4))
	}

	q.Clear()
	cap := q.Cap()
	content5 := make([]int, 0, cap+cap/2)
	for i := range cap {
		q.PushBack(i)
		content5 = append(content5, 200)
	}
	for i := 0; i < cap/2; i++ {
		q.PopFront()
		q.PushBack(i + 100)
		content5 = append(content5, 200)
	}

	buf := make([]int, q.Len()*2)
	n = q.CopyOutSlice(buf)
	if buf[0] != q.Front() {
		t.Fatal("wrong first element in output")
	}
	if buf[n-1] != q.At(n-1) {
		t.Fatal("wrong last element in output")
	}

	first := q.Front()
	last := q.At(len(out) - 1)
	n = q.CopyOutSlice(out)
	if n != len(out) {
		t.Fatalf("expected %d replaced items", len(out))
	}
	if out[0] != first {
		t.Fatalf("expected previous[0] to be %d, got %d", first, out[0])
	}
	if out[len(out)-1] != last {
		t.Fatalf("expected last previous to be %d, got %d", last, out[len(out)-1])
	}

	q.CopyInSlice(content5)
	if q.Len() != len(content5) {
		t.Fatalf("q should have %d items", len(content5))
	}
	for i, v := range content5 {
		if q.At(i) != v {
			t.Fatal("unexpected content after replace")
		}
	}
}

func TestAppendToSlice(t *testing.T) {
	var q Deque[int]
	content := []int{1, 2, 3}
	q.CopyInSlice(content)

	out := q.AppendToSlice(nil)
	if !slices.Equal(content, out) {
		t.Fatal("did not get expected content")
	}

	out = q.AppendToSlice(out)
	if len(out) != 2*len(content) {
		t.Fatal("worng output length")
	}
	if !slices.Equal(content, out[:len(content)]) {
		t.Fatal("did not get expected content at start of out")
	}
	if !slices.Equal(content, out[len(content):]) {
		t.Fatal("did not get expected content at end of out")
	}

	q.Clear()
	for i := 0; i < q.Cap(); i++ {
		q.PushBack(i)
	}
	for i := 0; i < q.Cap()/2; i++ {
		q.PopFront()
		q.PushBack(i + 100)
	}

	out = q.AppendToSlice(out[:0])
	if len(out) != q.Len() {
		t.Fatal("wrong output length")
	}
	if out[0] != q.Front() {
		t.Fatal("worng first element")
	}
	if out[len(out)-1] != q.Back() {
		t.Fatal("worng last element")
	}

	q.Clear()
	newOut := q.AppendToSlice(out)
	if !slices.Equal(out, newOut) {
		t.Fatal("expected same slice returned")
	}
}

func TestCopy(t *testing.T) {
	var a, b Deque[int]

	n := b.Copy(a)
	if n != 0 {
		t.Fatal("expected nothing copied")
	}

	content1 := []int{1, 2, 3}
	a.CopyInSlice(content1)

	n = b.Copy(a)
	if n != a.Len() {
		t.Fatal("Copy returned wrong length")
	}
	if !equalInt(a, b) {
		t.Fatal("different contents after copy")
	}

	a.PopFront()
	a.PopBack()
	n = b.Copy(a)
	if n != a.Len() {
		t.Fatal("Copy returned wrong length")
	}
	if !equalInt(a, b) {
		t.Fatal("different contents after copy")
	}

	a.Clear()
	n = b.Copy(a)
	if n != 0 {
		t.Fatal("expected nothing copied")
	}
	if b.Len() != 0 {
		t.Fatal("Destination should be empty")
	}
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
	var q Deque[int]
	for i := 0; i < b.N; i++ {
		q.PushFront(i)
	}
}

func BenchmarkPushBack(b *testing.B) {
	var q Deque[int]
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
}

func BenchmarkSerial(b *testing.B) {
	var q Deque[int]
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	for i := 0; i < b.N; i++ {
		q.PopFront()
	}
}

func BenchmarkSerialReverse(b *testing.B) {
	var q Deque[int]
	for i := 0; i < b.N; i++ {
		q.PushFront(i)
	}
	for i := 0; i < b.N; i++ {
		q.PopBack()
	}
}

func BenchmarkRotate(b *testing.B) {
	q := new(Deque[int])
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
	q := new(Deque[int])
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Insert(q.Len()/2, -i)
	}
}

func BenchmarkRemove(b *testing.B) {
	q := new(Deque[int])
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Remove(q.Len() / 2)
	}
}

func BenchmarkYoyo(b *testing.B) {
	var q Deque[int]
	for i := 0; i < b.N; i++ {
		for j := range 65536 {
			q.PushBack(j)
		}
		for range 65536 {
			q.PopFront()
		}
	}
}

func BenchmarkYoyoFixed(b *testing.B) {
	var q Deque[int]
	q.SetBaseCap(64000)
	for i := 0; i < b.N; i++ {
		for j := range 65536 {
			q.PushBack(j)
		}
		for range 65536 {
			q.PopFront()
		}
	}
}

func BenchmarkClearContiguous(b *testing.B) {
	var src Deque[int]
	for i := range 1<<15 + 1<<14 {
		src.PushBack(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := src
		q.buf = slices.Clone(q.buf)
		q.Clear()
	}
}

func BenchmarkClearSplit(b *testing.B) {
	var src Deque[int]
	for i := range 1<<15 + 1<<14 {
		src.PushFront(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := src
		q.buf = slices.Clone(q.buf)
		q.Clear()
	}
}

func equalInt(a, b Deque[int]) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if a.At(i) != b.At(i) {
			return false
		}
	}
	return true
}
