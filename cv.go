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

func RGB2Gray16(src image.Image) *image.Gray16 {
	bounds := src.Bounds()
	dst := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			v := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
			dst.SetGray16(x, y, color.Gray16{uint16(v)})
		}
	}
	return dst
}

func RGB2Gray(src image.Image) *image.Gray {
	bounds := src.Bounds()
	dst := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			v := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
			dst.SetGray(x, y, color.Gray{uint8(v / 255.)})
		}
	}
	return dst
}

func Binalize(src image.Image, th uint8) *image.Gray {
	gray := RGB2Gray(src)

	bounds := gray.Bounds()
	dst := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			v := gray.GrayAt(x, y).Y
			if v < th {
				dst.SetGray(x, y, color.Gray{0})
			} else {
				dst.SetGray(x, y, color.Gray{255})
			}
		}
	}
	return dst
}

func OtsuBinalize(src image.Image) *image.Gray {
	gray := RGB2Gray(src)

	bounds := gray.Bounds()
	var sigmaMax float64
	var sigmaMaxT uint8
	for t := 0; t < 255; t++ {
		o0, o1 := 0, 0
		totalM0, totalM1 := 0, 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				v := gray.GrayAt(x, y).Y
				if v <= uint8(t) {
					o0++
					totalM0 += int(v)
				} else {
					o1++
					totalM1 += int(v)
				}
			}
		}
		// should consider zero div...
		M0 := float64(totalM0) / float64(o0)
		M1 := float64(totalM1) / float64(o1)
		sigma := (float64(o0*o1) / sq(float64(o0+o1))) * sq(M0-M1)
		if sigma > sigmaMax {
			sigmaMax = sigma
			sigmaMaxT = uint8(t)
		}
	}
	return Binalize(gray, sigmaMaxT)
}

func rgb(r, g, b uint32) color.RGBA64 {
	return color.RGBA64{uint16(r), uint16(g), uint16(b), 65535}
}

func sq(v float64) float64 {
	return v * v
}
