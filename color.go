package migateway

import (
	"fmt"
	mcolor "github.com/bingbaba/tool/color"
	"image/color"
)

func RGBNumber(c color.Color) uint32 {
	r, g, b, a := c.RGBA()
	r = r << 24 >> 24
	g = g << 24 >> 24
	b = b << 24 >> 24
	a = a << 24 >> 24
	if a >= 100 {
		a = 100
	}

	fmt.Printf("r:%d,g:%d,b:%d,a:%d\n", r, g, b, a)
	if r|g|b == 0 {
		return 0
	} else {
		var v uint32 = uint32(100 - a)
		return v<<24 |
			r<<16 |
			g<<8 |
			b
	}
}

func NewRGB(n uint32) *mcolor.RGBA {
	fmt.Printf("******parse rgb number:%d\n", n)
	a := 100 - n>>24
	if a < 0 {
		a = 0
	}

	fmt.Printf("Bri: %d\n", n>>24)
	return &mcolor.RGBA{
		uint8(n << 8 >> 24),
		uint8(n << 16 >> 24),
		uint8(n << 24 >> 24),
		uint8(a),
	}
}
