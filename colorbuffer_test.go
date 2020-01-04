package main

import "testing"

func TestClear(t *testing.T) {
	cb := NewColorBuffer()

	for x := 0; x < WindowWidth; x++ {
		for y := 0; y < WindowHeight; y++ {
			if cb[WindowWidth*y+x] != 0x00 {
				t.Errorf("Image is not uniform black, %b", cb[WindowWidth*y+x])
			}
		}
	}
}

func TestSet(t *testing.T) {
	cb := NewColorBuffer()
	cb.Set(10, 10, 0xFFFFFFFF)

	if cb[WindowWidth*10*4+10*4+0] != 0xFF {
		t.Error("Error setting the correct color")
	}
	if cb[WindowWidth*10*4+10*4+1] != 0xFF {
		t.Error("Error setting the correct color")
	}
	if cb[WindowWidth*10*4+10*4+2] != 0xFF {
		t.Error("Error setting the correct color")
	}
	if cb[WindowWidth*10*4+10*4+3] != 0xFF {
		t.Error("Error setting the correct color")
	}
}
