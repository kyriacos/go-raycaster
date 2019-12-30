package main

import "math"

const (
	FPS             = 30
	FrameTimeLength = 1000 / FPS

	PI    = math.Pi
	TwoPI = 2.0 * PI

	TileSize   = 64
	MapNumRows = 13
	MapNumCols = 20

	MinimapScaleFactor = 1.0

	WindowWidth  = MapNumCols * TileSize
	WindowHeight = MapNumRows * TileSize

	FOV     = 60 * (math.Pi / 180)
	NumRays = WindowWidth
)
