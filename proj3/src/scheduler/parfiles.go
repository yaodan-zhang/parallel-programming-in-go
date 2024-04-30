package scheduler

import (
	"encoding/json"
	"math/rand"
	"os"
	"proj3/png"
	"strings"
	"sync"
	"sync/atomic"
)

// The image task node.
type Node struct {
	InputFile  string
	OutputFile string
	Effects    []any
	Next       *Node
}

// The task list for each thread.
type List struct {
	Head *Node
	Lock *atomic.Bool // TAS lock for performing dequeue on the list
}

// Process the image task.
func ProcessTask(task *Node) {
	pngImg, err := png.Load(task.InputFile)
	if err != nil {
		panic(err)
	}
	yMin, yMax, xMin, xMax := pngImg.GetBounds()
	bd := (yMax - yMin) * (xMax - xMin)
	for _, s := range task.Effects {
		switch s.(string) {
		case "G":
			pngImg.Grayscale(0, bd)
		case "S":
			pngImg.Sharpen(0, bd)
		case "E":
			pngImg.EdgeDetection(0, bd)
		case "B":
			pngImg.Blur(0, bd)
		}
		// swap the in and out image pointer for applying the next effect.
		pngImg.Swap()
	}
	// Counteract the last swap.
	pngImg.Swap()
	// Save the image.
	err = pngImg.Save(task.OutputFile)
	// Check save errors.
	if err != nil {
		panic(err)
	}
}

// Initialize task lists for threads.
func InitializeTaskLists(ThreadCount int) []*List {
	taskLists := make([]*List, ThreadCount*16) // padding
	for i := 0; i < ThreadCount; i++ {
		taskLists[i*16] = &List{Head: nil}
		var TAS atomic.Bool
		(&TAS).Store(false)
		taskLists[i*16].Lock = &TAS
	}
	return taskLists
}

// Append the task to the index list.
func EnqueueList(tasklists []*List, index int, task *Node) {
	if tasklists[index*16].Head == nil {
		tasklists[index*16].Head = task
	} else {
		node := tasklists[index*16].Head
		for node.Next != nil {
			node = node.Next
		}
		node.Next = task
	}
}

// Dequeue the list
func DequeueList(tasklists []*List, index int) *Node {
	task := tasklists[index*16].Head
	tasklists[index*16].Head = task.Next
	return task
}

// Return a stealed task if success, if failed, return nil
func StealTask(tasklists []*List, index int) *Node {
	lock := (tasklists[index*16].Lock).Swap(true)
	// Acquired the lock
	if !lock {
		defer tasklists[index*16].Lock.Store(false)
		if tasklists[index*16].Head != nil {
			task := DequeueList(tasklists, index)
			return task
		}
		return nil
	}
	// Acquire lock failed
	return nil
}

// Dequeue the list only if the thread get the TAS lock.
// Otherwise, spin until it gets the lock.
// If get the lock,
// 1. the list is not empty: dequeue, release the lock, and process the task.
// 2. the list is empty: release the lock, and steal task from a random thread.
// Return if all tasks are done.
func SpinList(config Config, tasklists []*List, index int, counter *atomic.Int64, wg *sync.WaitGroup) {
	// Spin until all tasks are done.
	for counter.Load() != 0 {
		lock := (tasklists[index*16].Lock).Swap(true)
		// Somebody else is dequeueing, i.e., stealing the task
		if lock {
			continue
		} else { // successfully get the lock
			var task *Node
			if tasklists[index*16].Head == nil {
				// release the lock
				tasklists[index*16].Lock.Store(false)
				// Steal task from another thread
				task = StealTask(tasklists, rand.Intn(config.ThreadCount))
			} else {
				// dequeue list
				task = DequeueList(tasklists, index)
				// release the lock
				tasklists[index*16].Lock.Store(false)
			}
			if task != nil {
				counter.Add(-1)
				ProcessTask(task)
			}
		}
	}
	wg.Done()
}

// Create the task lists as an array each entry store the head of a list.
// Each thread spins on its own list using the thread id.
// Return the pointer to the array.
func CreateTaskLists(config Config, counter *atomic.Int64) []*List {
	effectsFile, e := os.Open("../data/effects.txt")
	// Check open file error.
	if e != nil {
		panic(e)
	}
	reader := json.NewDecoder(effectsFile)
	dir := strings.Split(config.DataDirs, "+")
	// Initiate each list.
	taskLists := InitializeTaskLists(config.ThreadCount)
	// Generate tasks.
	for {
		var m map[string]any
		var inFilePath, outFilePath string
		var effects []any

		// A non-nil error indicates the end of JSON file.
		if err := reader.Decode(&m); err != nil {
			break
		}
		// Get info for a task.
		for k, v := range m {
			switch k {
			case "inPath":
				inFilePath = v.(string)
			case "outPath":
				outFilePath = v.(string)
			case "effects":
				effects = v.([]any)
			}
		}
		// Iterate through directories.
		for _, dataDir := range dir {
			// Create task node.
			task := &Node{
				InputFile:  "../data/in/" + dataDir + "/" + inFilePath,
				OutputFile: "../data/out/" + dataDir + "_" + outFilePath,
				Effects:    effects,
				Next:       nil}
			listIndex := rand.Intn(config.ThreadCount)
			// Append the task to a random list in tasklists.
			EnqueueList(taskLists, listIndex, task)
			counter.Add(1)
		}
	}
	return taskLists
}

// Run multiple image tasks in parallel, but each image task is
// processed wholely within one thread.
func RunParallelFiles(config Config) {
	// Atomic counter records the number of tasks.
	var counter atomic.Int64
	// Create task lists for threads. Increment counter for each task added.
	ThreadLists := CreateTaskLists(config, &counter)
	// Spawn threads to dequeue and process image tasks using the workstealing confinement.
	// Specifically, each thread has its own queue of tasks.
	// When all tasks are done, return.
	var wg sync.WaitGroup
	for i := 0; i < config.ThreadCount; i++ {
		wg.Add(1)
		go SpinList(config, ThreadLists, i, &counter, &wg)
		wg.Wait()
	}
}
