// Package queue implements a lock-free queue that supports the enqueue and dequeue
// operation on a series of requests.
package queue

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// The node in LockFreeQueue
type Request struct {
	Next      unsafe.Pointer
	Command   string
	Id        int
	Body      string
	Timestamp float64
}

// LockfreeQueue represents a FIFO structure with operations to enqueue
// and dequeue tasks represented as Request
type LockFreeQueue struct {
	Mutex *sync.Mutex
	Head  unsafe.Pointer
	Tail  unsafe.Pointer
	Cond  *sync.Cond
}

// NewLockFreeQueue creates and returns a LockFreeQueue
func NewLockFreeQueue() *LockFreeQueue {
	var mutex sync.Mutex
	condVar := sync.NewCond(&mutex)
	aux := &Request{Next: nil, Command: ""}
	head := unsafe.Pointer(aux)
	tail := unsafe.Pointer(aux)
	q := &LockFreeQueue{Head: head, Tail: tail, Mutex: &mutex, Cond: condVar}
	return q
}

// Enqueue adds a series of Request to the queue
func (q *LockFreeQueue) Enqueue(task *Request) {
	oldTail := atomic.LoadPointer(&q.Tail)
	for !atomic.CompareAndSwapPointer(&(*Request)(oldTail).Next, nil, unsafe.Pointer(task)) {
		oldTail = atomic.LoadPointer(&q.Tail)
	}
	atomic.CompareAndSwapPointer(&q.Tail, oldTail, unsafe.Pointer(task))
}

// Dequeue removes a Request from the queue
func (q *LockFreeQueue) Dequeue() *Request {
	oldHead := atomic.LoadPointer(&q.Head)
	if (*Request)(oldHead).Next == nil {
		if (*Request)(oldHead).Command == "DONE" {
			return &Request{Command: "DONE"}
		}
		return nil
	}
	nextHead := (*Request)(oldHead).Next // next to head
	for !atomic.CompareAndSwapPointer(&q.Head, oldHead, nextHead) {
		oldHead = atomic.LoadPointer(&q.Head)
		nextHead = (*Request)(oldHead).Next // next to head
		if nextHead == nil {
			if (*Request)(oldHead).Command == "DONE" {
				return &Request{Command: "DONE"}
			}
			return nil
		}
	}
	return (*Request)(nextHead)
}
