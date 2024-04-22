package scheduler

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"
	"proj3/png"
	"strings"
	"sync"
	"sync/atomic"
)

type ChunkTask struct {
	Effect     string
	PngImg     *png.Image
	ChunkStart int
	ChunkEnd   int
	Next       *ChunkTask
}

type ListC struct {
	Head *ChunkTask
	Lock *atomic.Bool // TAS lock for performing dequeue on the list
}

// Process the image chunk.
func SliceChunk(effect string, img *png.Image, chunkStart int, chunkEnd int) {
	switch effect {
	case "G":
		img.Grayscale(chunkStart, chunkEnd)
	case "S":
		img.Sharpen(chunkStart, chunkEnd)
	case "E":
		img.EdgeDetection(chunkStart, chunkEnd)
	case "B":
		img.Blur(chunkStart, chunkEnd)
	}
}

// Initialize the task lists
func InitializeTaskListsChunkVer(ThreadCount int) []*ListC {
	taskLists := make([]*ListC, ThreadCount*16) // padding
	for i := 0; i < ThreadCount; i++ {
		taskLists[i*16] = &ListC{Head: nil}
		var TAS atomic.Bool
		(&TAS).Store(false)
		taskLists[i*16].Lock = &TAS
	}
	return taskLists
}

// Append the chunk task to the index list.
func EnqueueListChunkVer(tasklists []*ListC, index int, task *ChunkTask) {
	// busy waiting until get the lock for the list
	for {
		if (tasklists[index*16].Lock).Swap(true) {
			continue
		} else {
			// Acquired the lock
			if tasklists[index*16].Head == nil {
				tasklists[index*16].Head = task
			} else {
				node := tasklists[index*16].Head
				for node.Next != nil {
					node = node.Next
				}
				node.Next = task
			}
			tasklists[index*16].Lock.Store(false)
			return
		}
	}
}

// Dequeue from an indexed list
func DequeueListChunkVer(tasklists []*ListC, index int) *ChunkTask {
	if tasklists[index*16].Head != nil {
		task := tasklists[index*16].Head
		tasklists[index*16].Head = task.Next
		return task
	}
	return nil
}

// Steal task from an indexed list.
func StealTaskChunkVer(tasklists []*ListC, index int) *ChunkTask {
	lock := (tasklists[index*16].Lock).Swap(true)
	// Acquired the lock
	if !lock {
		defer tasklists[index*16].Lock.Store(false)
		return DequeueListChunkVer(tasklists, index)
	}
	// Failed to acquire lock
	return nil
}

// Dequeue the list only if the thread get the TAS lock.
// Otherwise, spin until it gets the lock.
// If get the lock,
// 1. the list is not empty: dequeue, release the lock, and process the task.
// 2. the list is empty: release the lock, and steal task from a random thread.
// Return if all tasks are done.
func SpinListChunVer(config Config, tasklists []*ListC, index int, cond *sync.Cond, end *bool, counter *atomic.Int64, wg *sync.WaitGroup) {
	for !(*end) {
		lock := (tasklists[index*16].Lock).Swap(true)
		// Somebody else is dequeueing, i.e., stealing the task
		if lock {
			continue
		} else { // successfully get the lock
			var task *ChunkTask
			if tasklists[index*16].Head == nil {
				// release the lock
				tasklists[index*16].Lock.Store(false)
				// Steal task from another thread
				task = StealTaskChunkVer(tasklists, rand.Intn(config.ThreadCount))
			} else {
				// dequeue list
				task = DequeueListChunkVer(tasklists, index)
				// release the lock
				tasklists[index*16].Lock.Store(false)
			}
			if task != nil {
				// Process task
				SliceChunk(task.Effect, task.PngImg, task.ChunkStart, task.ChunkEnd)
				counter.Add(-1)
				//fmt.Print(counter.Load())
				// Signal to the barrier that all slices are done for this effect
				if counter.Load() == 0 {
					cond.L.Lock()
					cond.Signal()
					cond.L.Unlock()
				}
			}
		}
	}
	wg.Done()
}

// Run parallel image slice tasks for each image,
// continue to the next image task until the current task is done.
// Each effects should be processed sequentially using a barrier for synchronization.
func RunParallelSlices(config Config) {
	// Initialize condition variable as barrier between different effects
	m := sync.Mutex{}
	c := sync.NewCond(&m)

	// Initiate task lists for threads.
	taskLists := InitializeTaskListsChunkVer(config.ThreadCount)
	end := false

	// Spawn goroutines to spin on its own task list.
	var wg sync.WaitGroup
	var counter atomic.Int64
	for i := 0; i < config.ThreadCount; i++ {
		wg.Add(1)
		go SpinListChunVer(config, taskLists, i, c, &end, &counter, &wg)
	}

	// Dequeue and process one image at a time.
	effectsFile, e := os.Open("../data/effects.txt")
	if e != nil {
		panic(e)
	}
	reader := json.NewDecoder(effectsFile)
	dir := strings.Split(config.DataDirs, "+")

	for {
		var m map[string]any
		var inFilePath, outFilePath string
		var effects []any

		// A non-nil error indicates the end of JSON file.
		if err := reader.Decode(&m); err != nil {
			break
		}

		// Get info for an image task.
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
			// Load image.
			pngImg, err := png.Load("../data/in/" + dataDir + "/" + inFilePath)
			if err != nil {
				panic(err)
			}
			// Get image bounds
			yMin, yMax, xMin, xMax := pngImg.GetBounds()

			// Generate 4 tasks for each thread.
			numPixels := (yMax - yMin) * (xMax - xMin)
			auxiNumThreads := config.ThreadCount * 4
			numPixelsPerThread := int(math.Ceil(float64(numPixels) / float64(auxiNumThreads)))

			// Sequentially process the effects using a barrier.
			for _, effect := range effects {
				// Create tasks and add to random queue
				for i := 0; i < auxiNumThreads; i++ {
					chunkStart := numPixelsPerThread * i
					if chunkStart > numPixels {
						chunkStart = numPixels
					}
					chunkEnd := chunkStart + numPixelsPerThread
					if chunkEnd > numPixels {
						chunkEnd = numPixels
					}
					// Randomly add the task to a queue
					chunkTask := &ChunkTask{Effect: effect.(string), PngImg: pngImg, ChunkStart: chunkStart, ChunkEnd: chunkEnd}
					EnqueueListChunkVer(taskLists, rand.Intn(config.ThreadCount), chunkTask)
					counter.Add(1)
				}
				// Wait until all slices are done.
				if counter.Load() != 0 {
					c.L.Lock()
					c.Wait()
					c.L.Unlock()
				}
				// swap the in and out image pointer for applying the next effect.
				pngImg.Swap()
			}
			// Counteract the last swap.
			pngImg.Swap()
			// Save the image.
			err = pngImg.Save("../data/out/" + dataDir + "_" + outFilePath)
			// Check save errors.
			if err != nil {
				panic(err)
			}
		}
	}
	//Signal the end for goroutines.
	end = true
	wg.Wait()
}
