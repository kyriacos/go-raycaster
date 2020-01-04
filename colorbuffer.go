package main

// ColorBuffer - Stores all the information for the pixels on the screen will be copied into an SDL Texture
type ColorBuffer struct {
	Pixels        []byte
	Stride        int // Width * 4
	Width, Height int
}

// NewColorBuffer - Create new ColorBuffer based on the supplied width and height. Set the Stride at the time we initialize the ColorBuffer
func NewColorBuffer(w, h int) *ColorBuffer {
	buf := make([]byte, w*h*4, w*h*4)
	return &ColorBuffer{
		Pixels: buf,
		Stride: w * 4,
		Width:  w,
		Height: h,
	}
}

// GetPitch - Returns the Stride. I just remember pitch instead of Stride so i added this method
func (cb *ColorBuffer) GetPitch() int {
	return cb.Stride
}

// PixOffset returns the index of the first element of Pix that corresponds to the pixel at (x, y).
func (cb *ColorBuffer) PixOffset(x, y int) int {
	// Width*4*y + x*4
	return cb.Stride*y + x*4
}

// Set - set the values in the bytes buffer
func (cb *ColorBuffer) Set(x, y int, c uint32) {
	// cb[WindowWidth*y*4+x*4+3] = byte(c)

	i := cb.PixOffset(x, y)
	s := cb.Pixels[i : i+4 : i+4]
	s[0] = byte(c >> 24)
	s[1] = byte(c >> 16)
	s[2] = byte(c >> 8)
	s[3] = byte(c)
}

// At - retrieve uint32 color value (ordered as rgba)
func (cb *ColorBuffer) At(x, y int) uint32 {
	offset := cb.PixOffset(x, y)
	return uint32(cb.Pixels[offset])<<24 |
		uint32(cb.Pixels[offset+1])<<16 |
		uint32(cb.Pixels[offset+2])<<8 |
		uint32(cb.Pixels[offset+3])
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
