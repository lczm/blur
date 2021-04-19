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
func clamp(x uint16) uint8 {
	if x >= 255 {
		return 255
	}
	return uint8(x)
}

func main() {
	fmt.Println("Hello, world.")
	img, err := readImage("./wallpaper.jpg")
	if err != nil {
		panic(err)
	}
	mod := image.NewRGBA(img.Bounds())

	// kernel size
	// radius := 3
	point := img.Bounds().Size()
	width := point.X
	height := point.Y
	var constant uint16 = 255

	fmt.Println(fmt.Sprintf("Width : %d, Height : %d", width, height))
	// Iterate every pixel
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			col := img.At(i, j)
			r, g, b, _ := col.RGBA()
			r8 := clamp(uint16(r/256) + constant)
			g8 := clamp(uint16(g/256) + constant)
			b8 := clamp(uint16(b/256) + constant)
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
