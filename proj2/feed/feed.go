// The feed package simulates a Twitter feed, which supports user's add post, remove post, and check if
// contains a post operations. The post is uniquely identified by its timestamp.
package feed

import (
	"proj2/lock"
	"sync"
)

// Feed represents a user's twitter feed
// You will add to this interface the implementations as you complete them.
type Feed interface {
	Add(body string, timestamp float64)
	Remove(timestamp float64) bool
	Contains(timestamp float64) bool
	Feeds() []any
}

// feed is the internal representation of a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation. You can assume the feed will not have duplicate posts
type feed struct {
	start *post // a pointer to the beginning post
	lock  *lock.RWLock
}

// post is the internal representation of a post on a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation.
type post struct {
	body      string  // the text of the post
	timestamp float64 // Unix timestamp of the post
	next      *post   // the next post in the feed
}

// NewPost creates and returns a new post value given its body and timestamp
func newPost(body string, timestamp float64, next *post) *post {
	return &post{body, timestamp, next}
}

// NewFeed creates a empy Twitter feed
func NewFeed() Feed {
	var mutex sync.Mutex
	condVar := sync.NewCond(&mutex)
	return &feed{start: nil, lock: lock.NewRWLock(&mutex, condVar, 0)}
}

// Return all posts in the Feed in the format
// [{"body":body of post1,"timestamp":time of post 1},
// {"body":body of post2,"timestamp":time of post 2},
// ...]
func (f *feed) Feeds() []any {
	f.lock.RLock()
	defer f.lock.RUnlock()
	var s []any
	p := f.start
	for p != nil {
		s = append(s, map[string]any{"body": p.body, "timestamp": p.timestamp})
		p = p.next
	}
	return s
}

// Add inserts a new post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new post somewhere in the feed because
// the given timestamp may not be the most recent.
func (f *feed) Add(body string, timestamp float64) {
	// Create a new post
	newPost := newPost(body, timestamp, nil)
	// Lock the feed
	f.lock.Lock()
	// The feed is empty
	if f.start == nil {
		f.start = newPost
		f.lock.Unlock()
		return
	} else {
		// The feed is not empty, find the right position to add post
		var prevPost *post
		curPost := f.start
		for timestamp < curPost.timestamp {
			prevPost = curPost
			curPost = curPost.next
			if curPost == nil {
				break
			}
		}
		newPost.next = curPost
		if curPost == f.start {
			f.start = newPost
		} else {
			prevPost.next = newPost
		}
		f.lock.Unlock()
		return
	}
}

// Remove deletes the post with the given timestamp. If the timestamp
// is not included in a post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *feed) Remove(timestamp float64) bool {
	// Lock the feed
	f.lock.Lock()

	if f.start == nil {
		f.lock.Unlock()
		return false
	} else if f.start.timestamp == timestamp {
		f.start = f.start.next
		f.lock.Unlock()
		return true
	} else {
		var prevPost *post
		curPost := f.start
		for curPost.timestamp != timestamp {
			prevPost = curPost
			curPost = curPost.next
			if curPost == nil {
				f.lock.Unlock()
				return false
			}
		}
		prevPost.next = curPost.next
		f.lock.Unlock()
		return true
	}
}

// Contains determines whether a post with the given timestamp is
// inside a feed. The function returns true if there is a post
// with the timestamp, otherwise, false.
func (f *feed) Contains(timestamp float64) bool {
	// Lock the feed in reader mode
	f.lock.RLock()

	curPost := f.start
	for curPost != nil {
		if curPost.timestamp == timestamp {
			f.lock.RUnlock()
			return true
		}
		curPost = curPost.next
	}
	f.lock.RUnlock()
	return false
}
