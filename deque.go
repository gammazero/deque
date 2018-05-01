/*
Package deque provides a fast ring-buffer deque (double-ended queue)
implementation.

Deque generalizes a queue and a stack, to efficiently add and remove items at
either end with O(1) performance.  Queue (FIFO) operations are supported using
PushBack() and PopFront().  Stack (LIFO) operations are supported using
PushBack() and PopBack().

Ring-buffer Performance

The ring-buffer automatically resizes by
powers of two, growing when additional capacity is needed and shrinking when
only a quarter of the capacity is used, and uses bitwise arithmetic for all
calculations.

The ring-buffer implementation significantly improves memory and time
performance with fewer GC pauses, compared to implementations based on slices
and linked lists.

For maximum speed, this deque implementation leaves concurrency safety up to
the application to provide, however the application chooses, if needed at all.

Reading Empty Deque

Since it is OK for the deque to contain a nil value, it is necessary to either
panic or return a second boolean value to indicate the deque is empty, when
reading or removing an element.  This deque panics when reading from an empty
deque.  This is a run-time check to help catch programming errors, which may be
missed if a second return value is ignored.  Simply check Deque.Len() before
reading from the deque.

*/
package deque

// minCapacity is the smallest capacity that deque may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

// Deque represents a single instance of the deque data structure.
type Deque struct {
	buf   []interface{}
	head  int
	tail  int
	count int
}

// Len returns the number of elements currently stored in the queue.
func (q *Deque) Len() int {
	return q.count
}

// PushBack appends an element to the back of the queue.  Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (q *Deque) PushBack(elem interface{}) {
	q.growIfFull()

	q.buf[q.tail] = elem
	// Calculate new tail position.
	q.tail = q.next(q.tail)
	q.count++
}

// PushFront prepends an element to the front of the queue.
func (q *Deque) PushFront(elem interface{}) {
	q.growIfFull()

	// Calculate new head position.
	q.head = q.prev(q.head)
	q.buf[q.head] = elem
	q.count++
}

// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *Deque) PopFront() interface{} {
	if q.count <= 0 {
		panic("deque: PopFront() called on empty queue")
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// Calculate new head position.
	q.head = q.next(q.head)
	q.count--

	q.shrinkIfExcess()
	return ret
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *Deque) PopBack() interface{} {
	if q.count <= 0 {
		panic("deque: PopBack() called on empty queue")
	}

	// Calculate new tail position
	q.tail = q.prev(q.tail)

	// Remove value at tail.
	ret := q.buf[q.tail]
	q.buf[q.tail] = nil
	q.count--

	q.shrinkIfExcess()
	return ret
}

// Front returns the element at the front of the queue.  This is the element
// that would be returned by PopFront().  This call panics if the queue is
// empty.
func (q *Deque) Front() interface{} {
	if q.count <= 0 {
		panic("deque: Front() called when empty")
	}
	return q.buf[q.head]
}

// Back returns the element at the back of the queue.  This is the element
// that would be returned by PopBack().  This call panics if the queue is
// empty.
func (q *Deque) Back() interface{} {
	if q.count <= 0 {
		panic("deque: Back() called when empty")
	}
	return q.buf[q.prev(q.tail)]
}

// Clear removes all elements from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse.  The queue will not be resized smaller as long as items are
// only added.  Only when items are removed is the queue subject to getting
// resized smaller.
func (q *Deque) Clear() {
	// bitwise modulus
	modBits := len(q.buf) - 1
	for h := q.head; h != q.tail; h = (h + 1) & modBits {
		q.buf[h] = nil
	}
	q.head = 0
	q.tail = 0
	q.count = 0
}

// Rotate rotates the deque n steps front-to-back.  If n is negative, rotates
// back-to-front.  Having Deque provide Rotate() avoids resizing, that could
// not be may happen if implementing rotation using Pop and Push methods.
func (q *Deque) Rotate(n int) {
	if q.count <= 1 {
		return
	}
	// Rotating a multiple of q.count is same as no rotation.
	n %= q.count
	if n == 0 {
		return
	}

	modBits := len(q.buf) - 1
	// If no empty space in buffer, only move head and tail indexes.
	if q.head == q.tail {
		// Calculate new head and tail using bitwise modulus.
		q.head = (q.head + n) & modBits
		q.tail = (q.tail + n) & modBits
		return
	}

	if n < 0 {
		// Rotate back to front.
		for ; n < 0; n++ {
			// Calculate new head and tail using bitwise modulus.
			q.head = (q.head - 1) & modBits
			q.tail = (q.tail - 1) & modBits
			// Put tail value at head and remove value at tail.
			q.buf[q.head] = q.buf[q.tail]
			q.buf[q.tail] = nil
		}
		return
	}

	// Rotate front to back.
	for ; n > 0; n-- {
		// Put head value at tail and remove value at head.
		q.buf[q.tail] = q.buf[q.head]
		q.buf[q.head] = nil
		// Calculate new head and tail using bitwise modulus.
		q.head = (q.head + 1) & modBits
		q.tail = (q.tail + 1) & modBits
	}
}

// prev returns the previous buffer position wrapping around buffer.
func (q *Deque) prev(i int) int {
	return (i - 1) & (len(q.buf) - 1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (q *Deque) next(i int) int {
	return (i + 1) & (len(q.buf) - 1) // bitwise modulus
}

// growIfFull resizes up if the buffer is full.
func (q *Deque) growIfFull() {
	if len(q.buf) == 0 {
		q.buf = make([]interface{}, minCapacity)
		return
	}
	if q.count == len(q.buf) {
		q.resize()
	}
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (q *Deque) shrinkIfExcess() {
	if len(q.buf) > minCapacity && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

// resize resizes the deque to fit exactly twice its current contents.
// This results in shrinking if the queue is less than half-full, or growing
// the queue when it is full.
func (q *Deque) resize() {
	newBuf := make([]interface{}, q.count<<1)
	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}
