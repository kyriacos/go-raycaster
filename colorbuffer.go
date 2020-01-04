package main

// ColorBuffer - the buffer that will store all the information for the pixels on the screen will be copied into an SDL Texture
type ColorBuffer [WindowWidth * WindowHeight]uint32

// NewColorBuffer - Create new colorbuffer and initialize it
func NewColorBuffer() *ColorBuffer {
	return new(ColorBuffer).Init()
}

// Init - Initialize the ColorBuffer with a rectangle for the image based on the window width and height
func (cb *ColorBuffer) Init() *ColorBuffer {
	// cb.NRGBA = *image.NewNRGBA(image.Rect(0, 0, WindowWidth, WindowHeight))
	return cb
}

// Clear - Clear the color buffer. Set the value to a default or whatever is passed in
func (cb *ColorBuffer) Clear(c ...uint32) {
	// var col color.NRGBA = ColorBlack // default color if nothing is passed in
	var col uint32 = 0x000000FF

	if len(c) > 0 {
		// col = uint32ToColorNRGBA(c[0])
	}

	// draw.Draw(cb, cb.Bounds(), &image.Uniform{color.NRGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)
	for x := 0; x < WindowWidth; x++ {
		for y := 0; y < WindowHeight; y++ {
			// cb.SetNRGBA(x, y, col)
			cb[WindowWidth*y+x] = col
		}
	}
}
