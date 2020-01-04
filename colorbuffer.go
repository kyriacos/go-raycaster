package main

// ColorBuffer - the buffer that will store all the information for the pixels on the screen will be copied into an SDL Texture
type ColorBuffer [WindowWidth * WindowHeight * 4]byte

// NewColorBuffer - Create new colorbuffer and initialize it
func NewColorBuffer() *ColorBuffer {
	return new(ColorBuffer).Init()
}

// Init - Initialize the ColorBuffer with a rectangle for the image based on the window width and height
func (cb *ColorBuffer) Init() *ColorBuffer {
	return cb
}

// Set - set the values in the bytes buffer
func (cb *ColorBuffer) Set(x, y int, c uint32) {
	cb[WindowWidth*y*4+x*4+0] = byte(c >> 24)
	cb[WindowWidth*y*4+x*4+1] = byte(c >> 16)
	cb[WindowWidth*y*4+x*4+2] = byte(c >> 8)
	cb[WindowWidth*y*4+x*4+3] = byte(c)
}

// At - retrieve uint32 color value (ordered as rgba)
func (cb *ColorBuffer) At(x, y int) uint32 {
	return uint32(cb[WindowWidth*y*4+x*4+0])<<24 | uint32(cb[WindowWidth*y*4+x*4+1])<<16 | uint32(cb[WindowWidth*y*4+x*4+2])<<8 | uint32(cb[WindowWidth*y*4+x*4+3])
}

// Clear - Clear the color buffer. Set the value to a default or whatever is passed in
func (cb *ColorBuffer) Clear(c ...uint32) {
	var col uint32 = 0x00000000

	if len(c) > 0 {
		col = c[0]
	}

	for x := 0; x < WindowWidth; x++ {
		for y := 0; y < WindowHeight; y++ {
			cb.Set(x, y, col)
		}
	}
}
