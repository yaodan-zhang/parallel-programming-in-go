// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color"
)

// Grayscale applies a grayscale filtering effect to the image slice from start to end position, concatenating row-wise.
func (img *Image) Grayscale(start int, end int) {
	yMin, _, xMin, xMax := img.GetBounds() // Get image bounds.
	n := xMax - xMin                       // Get the number of pixel in a row.

	for i := start; i < end; i++ {
		// Perform position matching
		y := i/n + yMin
		x := i%n + xMin
		//Returns the pixel (i.e., RGBA) value at a (x,y) position
		// Note: These get returned as int32 so based on the math you'll
		// be performing you'll need to do a conversion to float64(..)
		r, g, b, a := img.in.At(x, y).RGBA()

		//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
		//For certain computations (i.e., convolution) the values might fall outside this
		// range so you need to clamp them between those values.
		greyC := clamp(float64(r+g+b) / 3.0)

		//Note: The values need to be stored back as uint16.
		img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
	}
}

// Perform the designated kernel to the image slice from start to end position, concatenating row-wise.
func (img *Image) PerformKernel(kernel *[9]float64, start int, end int) {
	yMin, yMax, xMin, xMax := img.GetBounds()
	n := xMax - xMin
	for k := start; k < end; k++ {
		y := k/n + yMin
		x := k%n + xMin
		kIndex := 0
		rNew, gNew, bNew, aNew := float64(0), float64(0), float64(0), float64(0)
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if (yMin <= (y + i)) && ((y + i) < yMax) && (xMin <= (x + j)) && ((x + j) < xMax) {
					r, g, b, a := img.in.At(x+j, y+i).RGBA()
					if j == 0 && i == 0 {
						aNew = float64(a)
					}
					kMult := kernel[kIndex]
					rNew += float64(r) * kMult
					gNew += float64(g) * kMult
					bNew += float64(b) * kMult
				}
				kIndex += 1
			}
		}
		img.out.Set(x, y, color.RGBA64{clamp(rNew), clamp(gNew), clamp(bNew), clamp(aNew)})
	}
}

// Sharpen effect on the image.
func (img *Image) Sharpen(start int, end int) {
	kernel := [9]float64{0, -1, 0, -1, 5, -1, 0, -1, 0}
	img.PerformKernel(&kernel, start, end)
}

// Edge detection on the image.
func (img *Image) EdgeDetection(start int, end int) {
	kernel := [9]float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
	img.PerformKernel(&kernel, start, end)
}

// Blur the image.
func (img *Image) Blur(start int, end int) {
	kernel := [9]float64{1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
	img.PerformKernel(&kernel, start, end)
}
