package deque

import (
	"fmt"
	"testing"
)

func TestIteratorDetectsClear(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		for _, tt := range []struct {
			size   int
			rotate int
		}{
			{15, 0},
			{16, 0},
			{32, 0},
			{16, 3},
		} {
			t.Run(fmt.Sprintf("reverse=%v/size=%d/rotate=%d", reverse, tt.size, tt.rotate), func(t *testing.T) {
				var q Deque[int]
				for i := range tt.size {
					q.PushBack(i + 1)
				}
				q.Rotate(tt.rotate)
				items := q.Iter()
				if reverse {
					items = q.RIter()
				}
				seen := 0
				assertPanics(t, "iterator must detect Clear", func() {
					for range items {
						seen++
						if seen == 1 {
							q.Clear()
						}
					}
				})
				if seen != 1 {
					t.Fatalf("iterator yielded %d items after Clear, want one before detection", seen)
				}
			})
		}
	}
}

func TestIterNil(t *testing.T) {
	var q *Deque[int]
	for range q.Iter() {
		t.Fatal("Iter yielded from nil deque")
	}
	for range q.RIter() {
		t.Fatal("RIter yielded from nil deque")
	}
}

func TestIterModifyAfterBreak(t *testing.T) {
	var q Deque[int]
	for i := range 8 {
		q.PushBack(i)
	}

	for range q.Iter() {
		q.PushBack(0)
		break
	}
	for range q.RIter() {
		q.PushBack(0)
		break
	}
	if q.Len() != 10 {
		t.Fatalf("expected length 10, got %d", q.Len())
	}
}
