package migateway

import (
	mcolor "github.com/bingbaba/tool/color"
	"image/color"
)

func RGBNumber(c color.Color) uint32 {
	r, g, b, a := c.RGBA()
	if a >= 100 {
		a = 100
	}

	if r|g|b == 0 {
		return 0
	} else {
		v := 0xff * (100 - a) / 100
		return v<<24 |
			r<<24>>24<<16 |
			g<<24>>24<<8 |
			b<<24>>24
	}
}

func NewRGB(n uint32) *mcolor.RGBA {
	return &mcolor.RGBA{
		uint8(n << 8 >> 24),
		uint8(n << 16 >> 24),
		uint8(n << 24 >> 24),
		uint8(n >> 24),
	}
}
