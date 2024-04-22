package scheduler

import (
	"encoding/json"
	"os"
	"proj3/png"
	"strings"
	"sync/atomic"
)

// Initialize the unique list.
func InitializeUniqueList() *List {
	tasklist := &List{Head: nil}
	var TAS atomic.Bool
	(&TAS).Store(false)
	tasklist.Lock = &TAS
	return tasklist
}

// Append the task to the unique list.
func EnqueueUniqueList(tasklist *List, task *Node) {
	if tasklist.Head == nil {
		tasklist.Head = task
	} else {
		node := tasklist.Head
		for node.Next != nil {
			node = node.Next
		}
		node.Next = task
	}
}

// Dequeue the unique task list.
func DequeueUniqueList(tasklist *List) *Node {
	task := tasklist.Head
	tasklist.Head = task.Next
	return task
}

// Process the image tasks sequentially without any parallelization.
func RunSequential(config Config) {
	// Initialize the unique task list
	tasklist := InitializeUniqueList()

	effectsFile, _ := os.Open("../data/effects.txt")
	reader := json.NewDecoder(effectsFile)
	dir := strings.Split(config.DataDirs, "+")
	// Get tasks and append to the list
	for {
		var m map[string]interface{}
		var inFilePath, outFilePath string
		var effects []interface{}

		// An error indicates the end of JSON file, return.
		if err := reader.Decode(&m); err != nil {
			return
		}
		// Get task info.
		for k, v := range m {
			switch k {
			case "inPath":
				inFilePath = v.(string)
			case "outPath":
				outFilePath = v.(string)
			case "effects":
				effects = v.([]interface{})
			}
		}
		//Iterate through data directories.
		for _, dataDir := range dir {
			_, err := png.Load("../data/in/" + dataDir + "/" + inFilePath)
			if err != nil {
				panic(err)
			}
			// Create task node.
			task := &Node{
				InputFile:  "../data/in/" + dataDir + "/" + inFilePath,
				OutputFile: "../data/out/" + dataDir + "_" + outFilePath,
				Effects:    effects,
				Next:       nil}
			// Enqueue task to the list.
			EnqueueUniqueList(tasklist, task)
		}
		//Dequeue list and process tasks
		for tasklist.Head != nil {
			taskToProcess := DequeueUniqueList(tasklist)
			ProcessTask(taskToProcess)
		}
	}
}
