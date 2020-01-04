package main

import (
	"testing"
)

const (
	width  = 20
	height = 20
)

func TestClear(t *testing.T) {
	cb := NewColorBuffer(width, height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if cb.Pixels[width*y+x] != 0x00 {
				t.Errorf("Image is not uniform black, %b", cb.Pixels[width*y+x])
			}
		}
	}
}

func TestSet(t *testing.T) {
	cb := NewColorBuffer(width, height)

	cb.Set(10, 10, 0xFFFFFFFF)

	if cb.Pixels[width*10*4+10*4+0] != 0xFF {
		t.Error("Error setting the correct color")
	}
	if cb.Pixels[width*10*4+10*4+1] != 0xFF {
		t.Error("Error setting the correct color")
	}
	if cb.Pixels[width*10*4+10*4+2] != 0xFF {
		t.Error("Error setting the correct color")
	}
	if cb.Pixels[width*10*4+10*4+3] != 0xFF {
		t.Error("Error setting the correct color")
	}
}

func TestAt(t *testing.T) {
	cb := NewColorBuffer(width, height)
	cb.Set(10, 10, 0xEEEEFFFF)

	if cb.At(10, 10) != 0xEEEEFFFF {
		t.Error("Error setting the correct color")
	}
}
