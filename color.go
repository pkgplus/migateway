package migateway

import (
	"image/color"
)

var (
	COLOR_RED    = &RGB{0xff, 0x00, 0x00}
	COLOR_YELLOW = &RGB{0xff, 0xff, 0x00}
	COLOR_BLUE   = &RGB{0x00, 0x00, 0xff}
	COLOR_GREEN  = &RGB{0x00, 0xff, 0x00}
	COLOR_CYAN   = &RGB{0x00, 0xff, 0xff}
	COLOR_ORANGE = &RGB{0xff, 0x61, 0x00}
	COLOR_PURPLE = &RGB{0x99, 0x33, 0xFA}
	COLOR_PINK   = &RGB{0xff, 0x14, 0x93}
	COLOR_GRAY   = &RGB{0xc0, 0xc0, 0xc0}
	COLOR_BLACK  = &RGB{0x00, 0x00, 0x00}

	COLOR_ALL = []color.Color{
		COLOR_RED,
		COLOR_YELLOW,
		COLOR_BLUE,
		COLOR_GREEN,
		COLOR_CYAN,
		COLOR_ORANGE,
		COLOR_PURPLE,
		COLOR_PINK,
		COLOR_GRAY,
		COLOR_BLACK,
	}
)

type RGB struct {
	r uint32
	g uint32
	b uint32
}

func (rgb *RGB) RGBA() (r, g, b, a uint32) {
	return 0xff<<24 + rgb.r<<16,
		0xff<<24 + rgb.g<<8,
		0xff<<24 + rgb.b,
		0
}

func RGBNumber(c color.Color) uint32 {
	r, g, b, a := c.RGBA()
	if r|g|b|a == 0xff<<24 {
		return 0
	} else {
		return r | g | b | a
	}
}
