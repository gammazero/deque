/*
Package deque provides a fast ring-buffer deque (double-ended queue)
implementation.  The ring-buffer automatically re-sizes by powers of two,
growing when additional capacity is needed and shrinking when only a quarter of
the capacity is used, and uses bitwise arithmetic for all calculations.

The ring-buffer implementation significantly improves memory and time
performance with fewer GC pauses, compared to implementations based on slices
and linked lists.

For maximum speed, this deque implementation leaves concurrency safety up to
the application to provide, however the application chooses, if needed at all.

Queue (FIFO) operations are supported using PushBack() and PopFront().  Stack
(LIFO) operations are supported using PushBack() and PopBack().

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
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	// calculate new tail using bitwise modulus
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

// PushFront prepends an element to the front of the queue.
func (q *Deque) PushFront(elem interface{}) {
	if q.count == len(q.buf) {
		q.resize()
	}

	// calculate new head using bitwise modulus
	q.head = (q.head - 1) & (len(q.buf) - 1)
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
	// Calculate new head using bitwise modulus.
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.buf) > minCapacity && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack().  If the queue is empty, the call
// panics.
func (q *Deque) PopBack() interface{} {
	if q.count <= 0 {
		panic("deque: PopBack() called on empty queue")
	}

	// Calculate new tail using bitwise modulus.
	q.tail = (q.tail - 1) & (len(q.buf) - 1)

	// Remove value at tail.
	ret := q.buf[q.tail]
	q.buf[q.tail] = nil
	q.count--

	// Resize down if buffer 1/4 full.
	if len(q.buf) > minCapacity && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret
}

// Front returns the element at the front of the queue.  This is the element
// that would be returned by PopFront().  This call panics if the queue is
// empty.
func (q *Deque) Front() interface{} {
	if q.count <= 0 {
		panic("deque: Front() called on empty queue")
	}
	return q.buf[q.head]
}

// Back returns the element at the back of the queue.  This is the element
// that would be returned by PopBack().  This call panics if the queue is
// empty.
func (q *Deque) Back() interface{} {
	if q.count <= 0 {
		panic("deque: Back() called on empty queue")
	}
	// bitwise modulus
	return q.buf[(q.tail-1)&(len(q.buf)-1)]
}

// Get returns the element at index i in the queue.  This method accepts only
// non-negative index values.  Get(0) refers to the first element and is the
// same as Front().  Get(Len()-1) refers to the last element and is the same
// as Back().  If the index is invalid, the call panics.
func (q *Deque) Get(i int) interface{} {
	if i < 0 || i >= q.count {
		panic("deque: Get() called with index out of range")
	}
	// bitwise modulus
	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Insert is used to insert an element into the middle of the queue.
// Insert(0,e) is the same as PushFront(e) and Insert(Len(),e) is the same as
// PushBack(e).  Accepts only non-negative index values, and panics if index is
// out of range.
func (q *Deque) Insert(i int, elem interface{}) {
	if i < 0 || i > q.count {
		panic("deque: Insert() called with index out of range")
	}
	if i == 0 {
		q.PushFront(elem)
		return
	}
	if i == q.count {
		q.PushBack(elem)
		return
	}
	if i <= q.count/2 {
		// If inserting closer to front, rotate front to back i places, put
		// element at front, then rotate back to front i places.
		for j := 0; j < i; j++ {
			q.PushBack(q.PopFront())
		}
		q.PushFront(elem)
		for j := 0; j < i; j++ {
			q.PushFront(q.PopBack())
		}
	} else {
		// If inserting closer to back, rotate back to front Len() - i places,
		// put element at back, then rotate front to back same amount.
		rots := q.count - i
		for j := 0; j < rots; j++ {
			q.PushFront(q.PopBack())
		}
		q.PushBack(elem)
		for j := 0; j < rots; j++ {
			q.PushBack(q.PopFront())
		}
	}
}

// Remove removes and returns an element from the middle of the queue.
// Remove(0) is the same as PopFront() and Remove(Len()-1) is the same as
// PopBack().  Accepts only non-negative index values, and panics if index is
// out of range.
func (q *Deque) Remove(i int) interface{} {
	if i < 0 || i >= q.count {
		panic("deque: Remove() called with index out of range")
	}
	if i == 0 {
		return q.PopFront()
	}
	if i == q.count-1 {
		return q.PopBack()
	}
	var elem interface{}
	if i <= q.count/2 {
		// If removing closer to front, rotate front to back i places, remove
		// element at front, then rotate back to front i places.
		for j := 0; j < i; j++ {
			q.PushBack(q.PopFront())
		}
		elem = q.PopFront()
		for j := 0; j < i; j++ {
			q.PushFront(q.PopBack())
		}
	} else {
		// If removing closer to back, rotate back to front Len() - 1 - i
		// places, remove element at back, then rotate front to back same
		// amount.
		rots := (q.count - 1) - i
		for j := 0; j < rots; j++ {
			q.PushFront(q.PopBack())
		}
		elem = q.PopBack()
		for j := 0; j < rots; j++ {
			q.PushBack(q.PopFront())
		}
	}
	return elem
}

// Replace replaces the element at position i with the given element.  Accepts
// only non-negative index values, and panics if index is out of range
func (q *Deque) Replace(i int, elem interface{}) {
	if i < 0 || i >= q.count {
		panic("deque: Remove() called with index out of range")
	}
	// bitwise modulus
	q.buf[(q.head+i)&(len(q.buf)-1)] = elem
}

// Clear removes all elements from the queue, but retains the current capacity.
// The queue will not be resized smaller as long as items are only added.
// Only when items are removed is the queue subject to getting resized smaller.
func (q *Deque) Clear() {
	// bitwise modulus
	mbits := len(q.buf) - 1
	for h := q.head; h != q.tail; h = (h + 1) & mbits {
		q.buf[h] = nil
	}
	q.head = 0
	q.tail = 0
	q.count = 0
}

// Copy copies elements from the queue to the destination slice, from front to
// back.  Copy returns the number of elements copied, which will be the minimum
// of q.Len() and len(dst).
func (q *Deque) Copy(dst []interface{}) int {
	count := q.count
	if len(dst) < q.count {
		count = len(dst)
	}
	if count == 0 {
		return 0
	}
	if q.head+count <= len(q.buf) {
		copy(dst, q.buf[q.head:q.head+count])
	} else {
		n := copy(dst, q.buf[q.head:])
		copy(dst[n:], q.buf[:count-n])
	}
	return count
}

// resize resizes the deque to fit exactly twice its current contents.
// This results in shrinking if the queue is less than half-full, or growing
// the queue when it is full.
func (q *Deque) resize() {
	if len(q.buf) == 0 {
		q.buf = make([]interface{}, minCapacity)
		return
	}

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
