package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"time"
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

// Clamp rgb values
// Takes in an integer and outputs a uint8(0-255)
// takes in a signed integer as the values can be negative
// i.e. negative contrast, goes below 0
func clampRGB(x int) uint8 {
	if x >= 255 {
		return 255
	} else if x < 0 {
		return 0
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

// Cli arguments
var inFile string
var outFile string
var technique string
var radius int
var contrast int
var verbose bool

// Start timer for verbose output
var start time.Time

func init() {
	// Input file
	flag.StringVar(&inFile, "i", "", "Input Image")
	flag.StringVar(&inFile, "input", "", "Input Image")
	// Output file
	flag.StringVar(&outFile, "o", "", "Output Image")
	flag.StringVar(&outFile, "output", "", "Output Image")
	// Technique
	flag.StringVar(&technique, "t", "gaussian", "Technique")
	flag.StringVar(&technique, "technique", "gaussian", "Technique")
	// Radius
	flag.IntVar(&radius, "r", 1, "Radius")
	flag.IntVar(&radius, "radius", 1, "Radius")
	// Contrast
	flag.IntVar(&contrast, "c", 0, "Contrast")
	flag.IntVar(&contrast, "contrast", 0, "Contrast")
	// Verbosity
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.BoolVar(&verbose, "verbose", false, "Verbose")
}

func main() {
	// Parse arguments
	flag.Parse()

	if verbose {
		start = time.Now()
	}

	if inFile == "" {
		fmt.Println("Must have an input image, use -i or --input")
		os.Exit(0)
	}
	if outFile == "" {
		fmt.Println("Must have an output image, use -o or --output")
		os.Exit(0)
	}
	if radius < 0 {
		fmt.Println("Radius must be at least 1")
		os.Exit(0)
	}
	if contrast < -255 || contrast > 255 {
		fmt.Println("Contrast must be within 1 and 255")
		os.Exit(0)
	}

	img, err := readImage(inFile)
	if err != nil {
		panic(err)
	}
	// Create new image
	mod := image.NewRGBA(img.Bounds())

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

	size := img.Bounds().Size()
	width := size.X
	height := size.Y

	if verbose {
		fmt.Println("Image width:", width)
		fmt.Println("Image height:", height)
		fmt.Println("Kernel width:", kernelWidth)
	}

	var calculateRGB = func(i, j int) {
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
		r8 := clampRGB(int(r/256) + contrast)
		g8 := clampRGB(int(g/256) + contrast)
		b8 := clampRGB(int(b/256) + contrast)
		mod.Set(i, j, color.RGBA{
			R: r8,
			G: g8,
			B: b8,
			A: 255,
		})
		// fmt.Println("calculating")
	}

	// Iterate every pixel
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			// calculateRGB(i, j)
			go calculateRGB(i, j)
		}
	}

	// fmt.Println("done")

	writer, _ := os.Create(outFile)
	jpeg.Encode(writer, mod, nil)

	if verbose {
		elapsed := time.Since(start)
		fmt.Printf("Time taken: %s\n", elapsed)
	}
}
