package main

import (
	"encoding/base64"
	"image/color"
	"image/png"
	"io"
	"strings"
	"testing"
)

// This image is a base64 encoding of a 'broken' png i had. It has an RGBA model instead of the
// more common NRGBA model that i usually got from the other pngs i saved
const wall = `
iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAIAAAAlC+aJAAAAeElEQVR4nOzVIQoAIADFUBXvf2UFs30M9oJB0zD8NeT2O8
/vaSru9T9QAK0Amj5Arx2gFUArgKYP0GsHaAXQCqDpA
/TaAVoBtAJo+gC9doBWAK0Amj5Arx2gFUArgKYP0GsHaAXQCqDpA
/TaAVoBtAJo+oAbAAD
//5UCQHtHbgkvAAAAAElFTkSuQmCC`

func base64PNG() io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(wall))
}

func TestDecodeImage(t *testing.T) {
	cfg, _ := png.DecodeConfig(base64PNG())
	img, _ := decodeImage(base64PNG())
	// the original file is in RGBA and we want all our images to be NRGBA
	if cfg.ColorModel == color.RGBAModel && img.ColorModel() != color.NRGBAModel {
		t.Error("Image is not in the correct color format")
	}

}
