// The server package implements a producer-consumer environment, where the producer
// extracts tasks from a given source and enqueues them as requests to a lock-free queue,
// and different consumers dequeue to process the requests and return the result (success or not, etc.).
package server

import (
	"encoding/json"
	"proj2/feed"
	"proj2/queue"
	"sync"
)

type Config struct {
	Encoder *json.Encoder // Represents the buffer to encode Responses
	Decoder *json.Decoder // Represents the buffer to decode Requests
	Mode    string        // Represents whether the server should execute
	// sequentially or in parallel
	// If Mode == "s"  then run the sequential version
	// If Mode == "p"  then run the parallel version
	// These are the only values for Version
	ConsumersCount int // Represents the number of consumers to spawn
}

// The comsumer function. If dequeue succeeds, it processes the task;
// otherwise, wait until there is an available job.
func consumer(q *queue.LockFreeQueue, f feed.Feed, config *Config, wg *sync.WaitGroup) {
	for {
		if task := q.Dequeue(); task == nil {
			q.Mutex.Lock()
			q.Cond.Wait()
			q.Mutex.Unlock()
		} else {
			cmd := task.Command
			if cmd == "ADD" {
				f.Add(task.Body, task.Timestamp)
				config.Encoder.Encode(map[string]any{"success": true, "id": task.Id})
			} else if cmd == "REMOVE" {
				config.Encoder.Encode(map[string]any{"success": f.Remove(task.Timestamp), "id": task.Id})
			} else if cmd == "CONTAINS" {
				config.Encoder.Encode(map[string]any{"success": f.Contains(task.Timestamp), "id": task.Id})
			} else if cmd == "FEED" {
				config.Encoder.Encode(map[string]any{"id": task.Id, "feed": f.Feeds()})
			} else if cmd == "DONE" {
				wg.Done()
				return
			}
		}
	}
}

// The producer function. It extracts the request one at a time and enqueues it into the lock-free queue,
// after it enqueues the "DONE" request, it will return.
func producer(q *queue.LockFreeQueue, config *Config) {
	for {
		var m map[string]any
		config.Decoder.Decode(&m)
		var task *queue.Request
		switch m["command"].(string) {
		case "ADD":
			task = &queue.Request{nil, "ADD", int(m["id"].(float64)), m["body"].(string), m["timestamp"].(float64)}
		case "REMOVE":
			task = &queue.Request{nil, "REMOVE", int(m["id"].(float64)), "", m["timestamp"].(float64)}
		case "CONTAINS":
			task = &queue.Request{nil, "CONTAINS", int(m["id"].(float64)), "", m["timestamp"].(float64)}
		case "FEED":
			task = &queue.Request{nil, "FEED", int(m["id"].(float64)), "", float64(0)}
		case "DONE":
			task = &queue.Request{nil, "DONE", int(0), "", float64(0)}
			q.Enqueue(task)
			q.Mutex.Lock()
			q.Cond.Broadcast()
			q.Mutex.Unlock()
			return
		}
		q.Enqueue(task)
		q.Mutex.Lock()
		q.Cond.Signal()
		q.Mutex.Unlock()
	}
}

// Run starts up the twitter server based on the configuration
// information provided and only returns when the server is fully
// shutdown.
func Run(config Config) {
	// Make a new lock-free queue and twitter feed
	taskQueue := queue.NewLockFreeQueue()
	myFeed := feed.NewFeed()
	// Parallel mode
	if config.Mode == "p" {
		var wg sync.WaitGroup
		// Spawn consumer goroutines
		for i := 0; i < config.ConsumersCount; i++ {
			wg.Add(1)
			go consumer(taskQueue, myFeed, &config, &wg)
		}
		// Call the producer function
		producer(taskQueue, &config)
		// Wait until all consumers are done
		wg.Wait()
	} else if config.Mode == "s" {
		// Sequential mode
		for {
			var m map[string]any
			config.Decoder.Decode(&m)
			switch m["command"].(string) {
			case "ADD":
				myFeed.Add(m["body"].(string), m["timestamp"].(float64))
				config.Encoder.Encode(map[string]any{"success": true, "id": int(m["id"].(float64))})
			case "REMOVE":
				config.Encoder.Encode(map[string]any{"success": myFeed.Remove(m["timestamp"].(float64)), "id": int(m["id"].(float64))})
			case "CONTAINS":
				config.Encoder.Encode(map[string]any{"success": myFeed.Contains(m["timestamp"].(float64)), "id": int(m["id"].(float64))})
			case "FEED":
				config.Encoder.Encode(map[string]any{"id": int(m["id"].(float64)), "feed": myFeed.Feeds()})
			case "DONE":
				return
			}
		}
	}
}
