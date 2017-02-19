package migateway

import (
	"image/color"
)

var (
	COLOR_RED    = &RGB{0xff, 0xff, 0x00, 0x00}
	COLOR_YELLOW = &RGB{0xff, 0xff, 0xff, 0x00}
	COLOR_BLUE   = &RGB{0xff, 0x00, 0x00, 0xff}
	COLOR_GREEN  = &RGB{0xff, 0x00, 0xff, 0x00}
	COLOR_CYAN   = &RGB{0xff, 0x00, 0xff, 0xff}
	COLOR_ORANGE = &RGB{0xff, 0xff, 0x61, 0x00}
	COLOR_PURPLE = &RGB{0xff, 0x99, 0x33, 0xFA}
	COLOR_PINK   = &RGB{0xff, 0xff, 0x14, 0x93}
	COLOR_GRAY   = &RGB{0xff, 0xc0, 0xc0, 0xc0}
	COLOR_BLACK  = &RGB{0x00, 0x00, 0x00, 0x00}
	COLOR_WHITE  = &RGB{0xff, 0xff, 0xff, 0xff}

	COLOR_ALL = []color.Color{
		COLOR_WHITE,
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
	v uint32
	r uint32
	g uint32
	b uint32
}

func (rgb *RGB) RGBA() (r, g, b, a uint32) {
	if rgb.v == 0 {
		return 0, 0, 0, 0
	} else {
		return rgb.v<<24 + rgb.r<<16,
			rgb.v<<24 + rgb.g<<8,
			rgb.v<<24 + rgb.b,
			0
	}
}

func (rgb *RGB) Brightness(v float64) {
	if v < 0 || v > 1 {
		return
	}

	if v == 1 {
		rgb.v = 255
	} else {
		rgb.v = uint32(v*float64(100)) / 2
	}

	/*
		var max uint32
		if rgb.r > max {
			max = rgb.r
		}
		if rgb.g > max {
			max = rgb.g
		}
		if rgb.g > max {
			max = rgb.g
		}

		LOGGER.Warn("%d %d %d", rgb.r, rgb.g, rgb.b)
		var cur_v float64 = float64(max) / float64(255)
		rgb.r = uint32(float64(rgb.r) / cur_v * v)
		rgb.g = uint32(float64(rgb.g) / cur_v * v)
		rgb.b = uint32(float64(rgb.b) / cur_v * v)
		LOGGER.Warn("%d %d %d", rgb.r, rgb.g, rgb.b)
	*/
}

func NewRGB(n uint32) *RGB {
	return &RGB{
		v: n >> 24,
		r: n << 8 >> 24,
		g: n << 16 >> 24,
		b: n << 24 >> 24,
	}
}

func RGBNumber(c color.Color) uint32 {
	r, g, b, a := c.RGBA()
	return r | g | b | a
}
