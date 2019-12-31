package main

import (
	"image/color"
	"math"
)

func distanceBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}

func normalizeAngle(angle float64) float64 {
	rAngle := math.Remainder(angle, TwoPI) // get the remainder - dont go beyond 360 degrees

	if rAngle < 0 { // if the angle is negative add 2 * PI
		rAngle = rAngle + TwoPI
	}

	return rAngle
}

// Get back values we can use with Go SDL2
func minimapScale(val interface{}) int32 {
	v := 0.0
	switch val.(type) {
	case int8, int32, int64:
		v = float64(val.(int))
	case int:
		v = float64(val.(int))
	default:
		v = val.(float64)
	}

	return int32(MinimapScaleFactor * v)
}

// Convert from Uint32 to RGBA color values
func uint32ToColorRGBA(h uint32) color.RGBA {
	return color.RGBA{
		R: uint8(h >> 24),
		G: uint8(h >> 16),
		B: uint8(h >> 8),
		A: uint8(h),
	}
}

// func clamp()

/*
	Calculating the pitch
	[https://stackoverflow.com/questions/37643392/getting-all-pixel-valuesrgba]

	RGBA struct in Go:
	type RGBA struct {
		// Pix holds the image's pixels, in R, G, B, A order. The pixel at
		// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
		Pix []uint8
		// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
		Stride int
		// Rect is the image's bounds.
		Rect Rectangle
	}
	Since Pix holds the individual pixels in a uint8 array, the pitch as
	defined in the SDL documentation [https://wiki.libsdl.org/SDL\_UpdateTexture]:
	`pitch: the number of bytes in a row of pixel data, including padding between lines`

	The number of bytes in a row of pixel data for each row for us is `WindowWidth * 4`.
	Since every pixel (well rgba) is 4 bytes so each row has WindowWidth length and each RGBA value
	for the single pixel will be 4 bytes long.

	Long description since this tripped my up a few times so i am making a lengthy note here.

	UPDATE: and i just did not simply see there is a Stride which is the pitch basically and we can use
			that instead! :) ooops
*/
func calculatePitch() int {
	// var a uint32
	// pitch := int(WindowWidth * unsafe.Sizeof(a))
	return WindowWidth * 4
}
