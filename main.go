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

	fmt.Println(fmt.Sprintf("Width : %d, Height : %d", width, height))
	// Iterate every pixel
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			col := img.At(i, j)
			r, g, b, _ := col.RGBA()
			r8 := uint8(r / 256)
			g8 := uint8(g / 256)
			b8 := uint8(b / 256)
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
