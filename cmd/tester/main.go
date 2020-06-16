package main

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/lusingander/cv"
)

const (
	input  = ""
	output = ""
)

func run(args []string) error {
	src, err := loadImage(input)
	if err != nil {
		return err
	}
	dst := cv.Mean(src, 3)
	return saveImage(output, dst)
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func saveImage(path string, dst image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, dst)
}

func main() {
	if err := run(os.Args); err != nil {
		panic(err)
	}
}
