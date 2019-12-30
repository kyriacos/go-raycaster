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

// func minimapScale()
