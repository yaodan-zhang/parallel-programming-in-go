// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

import (
	"sync"
)

type RWLock struct {
	mutex       *sync.Mutex
	cond        *sync.Cond
	readerCount int
}

// Create a new RWLock
func NewRWLock(mutex *sync.Mutex, cond *sync.Cond, readerCount int) *RWLock {
	return &RWLock{mutex, cond, readerCount}
}

// Read lock
func (l *RWLock) RLock() {
	l.mutex.Lock()
	for l.readerCount >= 32 {
		l.cond.Wait()
	}
	l.readerCount++
	l.mutex.Unlock()
}

// Read unlock
func (l *RWLock) RUnlock() {
	l.mutex.Lock()
	l.readerCount--
	if l.readerCount == 0 { //|| l.readerCount < 32 {
		l.cond.Broadcast()
	}
	l.mutex.Unlock()
}

// Write lock
func (l *RWLock) Lock() {
	l.mutex.Lock()
	for l.readerCount > 0 {
		l.cond.Wait()
	}
}

// Write unlock
func (l *RWLock) Unlock() {
	l.mutex.Unlock()
}
