package cv

import (
	"image"
	"image/color"
	"math"
	"sort"
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

func Gaussian(src image.Image, sigma float64, kernelSize int) image.Image {
	kernel := make([][]float64, kernelSize)
	total := 0.
	c := kernelSize / 2
	for y := 0; y < kernelSize; y++ {
		kernel[y] = make([]float64, kernelSize)
		for x := 0; x < kernelSize; x++ {
			dy := float64(y - c)
			dx := float64(x - c)
			kernel[y][x] = (1. / (2 * math.Pi * sq(sigma))) * math.Exp(-(sq(dx)+sq(dy))/(2*sq(sigma)))
			total += kernel[y][x]
		}
	}
	for y := 0; y < kernelSize; y++ {
		for x := 0; x < kernelSize; x++ {
			kernel[y][x] /= total
		}
	}
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rr, gg, bb := 0., 0., 0.
			for yy := -c; yy <= c; yy++ {
				for xx := -c; xx <= c; xx++ {
					if x+xx >= bounds.Min.X && y+yy >= bounds.Min.Y && x+xx < bounds.Max.X && y+yy < bounds.Max.Y {
						r, g, b, _ := src.At(x+xx, y+yy).RGBA()
						rr += kernel[yy+c][xx+c] * float64(r)
						gg += kernel[yy+c][xx+c] * float64(g)
						bb += kernel[yy+c][xx+c] * float64(b)
					}
				}
			}
			dst.Set(x, y, rgb(uint32(rr), uint32(gg), uint32(bb)))
		}
	}
	return dst
}

func Median(src image.Image, kernelSize int) image.Image {
	c := kernelSize / 2
	m := kernelSize * kernelSize / 2
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rs := make([]uint32, 0)
			gs := make([]uint32, 0)
			bs := make([]uint32, 0)
			for yy := -c; yy <= c; yy++ {
				for xx := -c; xx <= c; xx++ {
					if x+xx >= bounds.Min.X && y+yy >= bounds.Min.Y && x+xx < bounds.Max.X && y+yy < bounds.Max.Y {
						r, g, b, _ := src.At(x+xx, y+yy).RGBA()
						rs = append(rs, r)
						gs = append(gs, g)
						bs = append(bs, b)
					} else {
						rs = append(rs, 0)
						gs = append(gs, 0)
						bs = append(bs, 0)
					}
				}
			}
			sortUint32(rs)
			sortUint32(gs)
			sortUint32(bs)
			dst.Set(x, y, rgb(rs[m], gs[m], bs[m]))
		}
	}
	return dst
}

func Mean(src image.Image, kernelSize int) image.Image {
	c := kernelSize / 2
	k := uint32(kernelSize * kernelSize)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rr := uint32(0)
			gg := uint32(0)
			bb := uint32(0)
			for yy := -c; yy <= c; yy++ {
				for xx := -c; xx <= c; xx++ {
					if x+xx >= bounds.Min.X && y+yy >= bounds.Min.Y && x+xx < bounds.Max.X && y+yy < bounds.Max.Y {
						r, g, b, _ := src.At(x+xx, y+yy).RGBA()
						rr += r
						gg += g
						bb += b
					}
				}
			}
			dst.Set(x, y, rgb(rr/k, gg/k, bb/k))
		}
	}
	return dst
}

func sortUint32(s []uint32) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}

func rgb(r, g, b uint32) color.RGBA64 {
	return color.RGBA64{uint16(r), uint16(g), uint16(b), 65535}
}

func sq(v float64) float64 {
	return v * v
}
