package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj1/png"
	"strings"
	"sync"
)

// The image task node.
type Node struct {
	inputFile  string
	outputFile string
	effects    []any
	next       *Node
}

// The image tasks linked list.
type List struct {
	head  *Node
	lock  bool // TAS lock for performing dequeue on the list
	count int
}

// Process the image task.
func ProcessTask(task *Node, wg *sync.WaitGroup) {
	pngImg, err := png.Load(task.inputFile)
	if err != nil {
		panic(err)
	}
	yMin, yMax, xMin, xMax := pngImg.GetBounds()
	bd := (yMax - yMin) * (xMax - xMin)
	for _, s := range task.effects {
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
	err = pngImg.Save(task.outputFile)
	// Check save errors.
	if err != nil {
		panic(err)
	}
}

// Dequeue the list only if the thread get the TAS lock.
// Otherwise, spin until it gets the lock or the list becomes empty.
// If get the lock, dequeue, release the lock, and process the task.
func DequeueList(list *List, wg *sync.WaitGroup) {
	// Spin until the list is empty.
	for list.head != nil {
		lock := list.lock
		list.lock = true

		// Somebody else is dequeing
		if lock == true {
			continue
		} else { // dequeue success
			task := list.head
			list.head = task.next
			list.lock = false
			ProcessTask(task, wg)
		}
	}
	wg.Done()
}

// Create the image task linked list from effects.txt.
// Return a pointer to the list.
func CreateTaskList(config Config) *List {
	effectsPathFile := fmt.Sprintf("../data/effects.txt")
	effectsFile, _ := os.Open(effectsPathFile)
	reader := json.NewDecoder(effectsFile)
	dir := strings.Split(config.DataDirs, "+")

	taskList := &List{head: nil, lock: false, count: 0}

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
		for _, dataDir := range dir {
			// Append the task to the end of the queue.
			task := &Node{
				inputFile:  "../data/in/" + dataDir + "/" + inFilePath,
				outputFile: "../data/out/" + dataDir + "_" + outFilePath,
				effects:    effects,
				next:       nil}
			if taskList.head == nil {
				taskList.head = task
			} else {
				node := taskList.head
				for node.next != nil {
					node = node.next
				}
				node.next = task
			}
			taskList.count += 1
		}
	}
	return taskList
}

// Run multiple image tasks in parallel, but each image task must be
// processed within one thread.
func RunParallelFiles(config Config) {
	taskList := CreateTaskList(config)
	wg := sync.WaitGroup{}
	numThreads := config.ThreadCount
	if taskList.count < config.ThreadCount {
		numThreads = taskList.count
	}

	// Spawn threads to dequeue and process image tasks.
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go DequeueList(taskList, &wg)
	}
	// Wait until all tasks are done.
	wg.Wait()
}
