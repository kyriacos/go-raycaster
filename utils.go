package main

import "math"

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
	case int8, int, int32, int64, float32:
		v = val.(float64)
	default:
		v = val.(float64)
	}

	return int32(MinimapScaleFactor * v)
}
