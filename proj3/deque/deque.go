package deque

import (
	"sync/atomic"
)

// BoundDeque represents a deque with one-ended push and two-ended pop.
type BoundDeque struct {
	top    Index
	bottom int32
	tasks  []interface{}
}

// NewBoundDeque creates and initializes a BoundDeque
func NewBoundDeque(size int) *BoundDeque {
	return &BoundDeque{
		top:    *NewIndex(0, 0),
		bottom: 0,
		tasks:  make([]interface{}, size),
	}
}

// PushBottom adds a Task to the deque from the back.
func (deque *BoundDeque) PushBottom(task interface{}) {
	deque.tasks[deque.bottom] = task
	atomic.AddInt32(&deque.bottom, 1)
}

// PopBottom removes and returns a Task from the back of the deque.
func (deque *BoundDeque) PopBottom() interface{} {
	if deque.bottom == 0 {
		return nil
	}
	atomic.AddInt32(&deque.bottom, -1)
	task := deque.tasks[deque.bottom]
	oldTop, oldStamp := deque.top.Get()
	var newTop int32 = 0
	newStamp := oldStamp + 1
	if deque.bottom > oldTop {
		return task
	}
	if deque.bottom == oldTop {
		deque.bottom = 0
		if deque.top.CompareAndSwap(oldTop, newTop, oldStamp, newStamp) {
			return task
		}
	}
	deque.top.Set(newTop, newStamp)
	deque.bottom = 0
	return nil
}

// PopTop removes and returns a Task from the front of the deque.
func (deque *BoundDeque) PopTop() interface{} {
	oldTop, oldStamp := deque.top.Get()
	newTop := oldTop + 1
	newStamp := oldStamp + 1
	if deque.bottom <= oldTop {
		return nil
	}
	task := deque.tasks[oldTop]
	if deque.top.CompareAndSwap(oldTop, newTop, oldStamp, newStamp) {
		return task
	}
	return nil
}

// IsEmpty checks if the deque is empty.
func (deque *BoundDeque) IsEmpty() bool {
	top := deque.top.Value()
	return deque.bottom <= top
}
