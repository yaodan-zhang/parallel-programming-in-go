package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj1/png"
	"strings"
)

// Process the image tasks sequentially without any parallelization.
func RunSequential(config Config) {
	effectsPathFile := fmt.Sprintf("../data/effects.txt")
	effectsFile, _ := os.Open(effectsPathFile)
	reader := json.NewDecoder(effectsFile)
	dir := strings.Split(config.DataDirs, "+")

	for {
		var m map[string]interface{}
		var inFilePath, outFilePath string
		var effects []interface{}

		// An error indicates the end of JSON file.
		// Return, since all tasks are done.
		if err := reader.Decode(&m); err != nil {
			return
		}

		// Get image task info.
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

		// Process image task.
		for _, dataDir := range dir {
			pngImg, err := png.Load("../data/in/" + dataDir + "/" + inFilePath)

			if err != nil {
				panic(err)
			}
			// Get image bounds.
			yMin, yMax, xMin, xMax := pngImg.GetBounds()
			// Set start position to 0 and end position to the end of the image.
			bd := (yMax - yMin) * (xMax - xMin)

			for _, s := range effects {
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
			err = pngImg.Save("../data/out/" + dataDir + "_" + outFilePath)

			//Checks save errors.
			if err != nil {
				panic(err)
			}
		}
	}
}
