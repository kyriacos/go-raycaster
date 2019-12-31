package main

import (
	"image"
)

// ColorBuffer - the buffer that will store all the information for the pixels on the screen will be copied into an SDL Texture
type ColorBuffer struct{ image.RGBA }

// NewColorBuffer - Create new colorbuffer and initialize it
func NewColorBuffer() *ColorBuffer {
	return new(ColorBuffer).Init()
}

// Init - Initialize the ColorBuffer with a rectangle for the image based on the window width and height
func (cb *ColorBuffer) Init() *ColorBuffer {
	cb.RGBA = *image.NewRGBA(image.Rect(0, 0, WindowWidth, WindowHeight))
	return cb
}

// Clear - Clear the color buffer. Set the value to a default or whatever is passed in
func (cb *ColorBuffer) Clear(c ...uint32) {
	var col uint32 = 0x000000FF // default color if nothing is passed in

	if len(c) > 0 {
		col = c[0]
	}

	for x := 0; x < WindowWidth; x++ {
		for y := 0; y < WindowHeight; y++ {
			cb.SetRGBA(x, y, uint32ToColorRGBA(col))
		}
	}
}
