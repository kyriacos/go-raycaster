package main

import "math"

const (
	fps             = 30
	frameTimeLength = 1000 / fps

	pi    = math.Pi
	twoPI = 2 * pi

	tileSize   = 64
	mapNumRows = 13
	mapNumCols = 20

	minimapScaleFactor = 1.0

	windowWidth  = mapNumCols * tileSize
	windowHeight = mapNumRows * tileSize

	fov     = 60 * (math.Pi / 180)
	numRays = windowWidth
)
