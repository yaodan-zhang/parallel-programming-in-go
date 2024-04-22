package scheduler

import (
	"math"
	"proj1/png"
	"sync"
)

// Process the image slice.
func SliceChunk(s any, img *png.Image, chunkStart int, chunkEnd int, wg *sync.WaitGroup) {
	switch s.(string) {
	case "G":
		img.Grayscale(chunkStart, chunkEnd)
	case "S":
		img.Sharpen(chunkStart, chunkEnd)
	case "E":
		img.EdgeDetection(chunkStart, chunkEnd)
	case "B":
		img.Blur(chunkStart, chunkEnd)
	}
	wg.Done()
}

// Run parallel image slice tasks for each image,
// continue to the next image task until the current task is done.
func RunParallelSlices(config Config) {
	// Create the image tasks list from effects.txt.
	taskList := CreateTaskList(config)

	// Dequeue and process one image at a time.
	for taskList.head != nil {
		// Dequeue, get an image task.
		task := taskList.head
		taskList.head = task.next

		// Load image
		pngImg, err := png.Load(task.inputFile)
		if err != nil {
			panic(err)
		}

		// Get image bounds
		yMin, yMax, xMin, xMax := pngImg.GetBounds()

		// Get number of pixels per thread
		numPixels := (yMax - yMin) * (xMax - xMin)
		numThreads := config.ThreadCount
		numPixelsPerThread := int(math.Ceil(float64(numPixels) / float64(numThreads)))

		// Process each effect sequantially
		for _, s := range task.effects {
			var wg sync.WaitGroup
			// Spawn a go routine for each image slice.
			for i := 0; i < numThreads; i++ {
				chunkStart := numPixelsPerThread * i
				if chunkStart > numPixels {
					chunkStart = numPixels
				}
				chunkEnd := chunkStart + numPixelsPerThread
				if chunkEnd > numPixels {
					chunkEnd = numPixels
				}
				wg.Add(1)
				go SliceChunk(s, pngImg, chunkStart, chunkEnd, &wg)
			}
			// Wait until all slices are done.
			wg.Wait()
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
}
