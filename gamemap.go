package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Level - type
type Level [MapNumRows][MapNumCols]int

// GameMap - comment
type GameMap struct {
	level Level
}

func (gm *GameMap) hasWallAt(x float64, y float64) bool {
	if x < 0 || x > WindowWidth || y < 0 || y > WindowHeight {
		return true
	}

	mapGridIndexX := int(math.Floor(x / TileSize))
	mapGridIndexY := int(math.Floor(y / TileSize))

	return gm.level[mapGridIndexY][mapGridIndexX] != 0
}

func (gm *GameMap) render(renderer *sdl.Renderer) {
	for i := 0; i < MapNumRows; i++ {
		for j := 0; j < MapNumCols; j++ {
			tileX := j * TileSize // column
			tileY := i * TileSize // row

			var tileColor uint8 = 0
			if gameMap.level[i][j] != 0 {
				tileColor = 255
			}

			renderer.SetDrawColor(tileColor, tileColor, tileColor, 255)
			rect := &sdl.Rect{
				X: int32(MinimapScaleFactor * float64(tileX)),
				Y: int32(MinimapScaleFactor * float64(tileY)),
				W: int32(math.Floor(MinimapScaleFactor * TileSize)),
				H: int32(math.Floor(MinimapScaleFactor * TileSize)),
			}
			renderer.FillRect(rect)
		}
	}
}
