package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

func readImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// TODO : Think about other image formats
	image, err := jpeg.Decode(f)
	// image, _, err := image.Decode(f)
	return image, err
}

// Given a uint16 (not a uint8) such that it does not overload,
// clamp it back into a uint8, mainly used to add constants to colors
func clampRGB(x uint16) uint8 {
	if x >= 255 {
		return 255
	}
	return uint8(x)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func main() {
	img, err := readImage("./wallpaper.jpg")
	if err != nil {
		panic(err)
	}
	mod := image.NewRGBA(img.Bounds())
	// kernel size
	radius := 7
	var sigma float64 = math.Max(float64(radius/2), 1.0)
	var kernelWidth int = (radius * 2) + 1 // left + right + center

	// Create a 2d array of length kernelWidth that holds []float64
	var kernel = make([][]float64, kernelWidth)
	// Fill each of the indices with []float64 arrays
	for i := range kernel {
		kernel[i] = make([]float64, kernelWidth)
	}

	var twoPi float64 = 2 * math.Pi
	var sum float64 = 0

	// Generate the gaussian kernel values
	for x := -radius; x < radius; x++ {
		for y := -radius; y < radius; y++ {
			var expNumerator float64 = -(float64(x*x) + float64(y*y))
			var expDenominator float64 = 2 * (sigma * sigma)
			var e float64 = math.Pow(math.E, float64(expNumerator/expDenominator))
			var value float64 = e / float64(twoPi*sigma*sigma)
			kernel[x+radius][y+radius] = value
			sum += value
		}
	}
	// Normalize the kernel, so that all the values of the kernel adds up to 1
	for x := 0; x < kernelWidth; x++ {
		for y := 0; y < kernelWidth; y++ {
			kernel[x][y] /= sum
		}
	}

	point := img.Bounds().Size()
	width := point.X
	height := point.Y
	// TODO : contrast
	// var contrast uint16 = 0
	fmt.Println("kernel width : ", kernelWidth)

	fmt.Println(fmt.Sprintf("Width : %d, Height : %d", width, height))
	// Iterate every pixel
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			topLeftX := max(0, i-radius)
			topLeftY := max(0, j-radius)
			botRightX := min(width-1, i+radius)
			botRightY := min(height-1, j+radius)

			var r float64 = 0
			var g float64 = 0
			var b float64 = 0
			for y := topLeftY; y <= botRightY; y++ {
				for x := topLeftX; x <= botRightX; x++ {
					sCol := img.At(x, y)
					sR, sG, sB, _ := sCol.RGBA()

					dx := (x - i) + radius
					dy := (y - j) + radius

					value := kernel[dx][dy]

					r += float64(sR) * value
					g += float64(sG) * value
					b += float64(sB) * value
				}
			}

			r8 := uint8(uint16(r) / 256)
			g8 := uint8(uint16(g) / 256)
			b8 := uint8(uint16(b) / 256)
			mod.Set(i, j, color.RGBA{
				R: r8,
				G: g8,
				B: b8,
				A: 255,
			})
		}
	}

	writer, _ := os.Create("test.jpg")
	jpeg.Encode(writer, mod, nil)
}
