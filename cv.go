package cv

import (
	"image"
	"image/color"
)

func RGB2BGR(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			dst.Set(x, y, rgb(b, g, r))
		}
	}
	return dst
}

func RGB2Gray(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			v := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
			dst.SetGray16(x, y, gray(v))
		}
	}
	return dst
}

func rgb(r, g, b uint32) color.RGBA64 {
	return color.RGBA64{uint16(r), uint16(g), uint16(b), 65535}
}

func gray(y float64) color.Gray16 {
	return color.Gray16{uint16(y)}
}
