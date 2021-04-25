package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
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
	fmt.Println("Hello, world.")
	img, err := readImage("./wallpaper.jpg")
	if err != nil {
		panic(err)
	}
	mod := image.NewRGBA(img.Bounds())
	// kernel size
	radius := 1
	point := img.Bounds().Size()
	width := point.X
	height := point.Y
	var constant uint16 = 0

	fmt.Println(fmt.Sprintf("Width : %d, Height : %d", width, height))
	// Iterate every pixel
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			topLeftX := max(0, i-radius)
			topLeftY := max(0, j-radius)
			botRightX := min(width-1, i+radius)
			botRightY := min(height-1, j+radius)

			// Starts with 1 as this starts off with the values at i, j
			var amount int = 1
			col := img.At(i, j)
			r, g, b, _ := col.RGBA()
			for y := topLeftY; y <= botRightY; y++ {
				for x := topLeftX; x <= botRightX; x++ {
					if x == i && y == j {
						continue
					}
					// The values of the ones at the surrounding pixels within the kernel
					sCol := img.At(x, y)
					sR, sG, sB, _ := sCol.RGBA()
					r += sR
					g += sG
					b += sB
					amount++
				}
			}
			// fmt.Println("Amount : ", amount)

			r8 := clampRGB(uint16(r/uint32(amount)/256) + constant)
			g8 := clampRGB(uint16(g/uint32(amount)/256) + constant)
			b8 := clampRGB(uint16(b/uint32(amount)/256) + constant)
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
