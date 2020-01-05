package main

import "math"

// Various Constants
const (
	FPS             = 30
	FrameTimeLength = 1000 / FPS

	PI    = math.Pi
	TwoPI = 2.0 * PI

	TileSize   = 64
	MapNumRows = 13
	MapNumCols = 20

	MinimapScaleFactor = 0.2

	WindowWidth  = MapNumCols * TileSize
	WindowHeight = MapNumRows * TileSize

	FOV     = 60 * (math.Pi / 180)
	NumRays = WindowWidth

	TextureWidth  = 64
	TextureHeight = 64
)

// Some base colors
var (
	ColorCeiling = uint32ToColorNRGBA(0x333333FF)
	ColorFloor   = uint32ToColorNRGBA(0x777777FF)
	ColorBlack   = uint32ToColorNRGBA(0x000000FF)
)
