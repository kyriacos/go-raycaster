package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Level - type
type Level [mapNumRows][mapNumCols]int

// GameMap - comment
type GameMap struct {
	level Level
}

func (gm *GameMap) hasWallAt(x float64, y float64) bool {
	if x < 0 || x > windowWidth || y < 0 || y > windowHeight {
		return true
	}

	mapGridIndexX := int(math.Floor(x / tileSize))
	mapGridIndexY := int(math.Floor(y / tileSize))

	return gm.level[mapGridIndexY][mapGridIndexX] != 0
}

func (gm *GameMap) render(renderer *sdl.Renderer) {
	for i := 0; i < mapNumRows; i++ {
		for j := 0; j < mapNumCols; j++ {
			tileX := j * tileSize // column
			tileY := i * tileSize // row

			var tileColor uint8 = 0
			if gameMap.level[i][j] != 0 {
				tileColor = 255
			}

			renderer.SetDrawColor(tileColor, tileColor, tileColor, 255)
			rect := &sdl.Rect{
				X: int32(minimapScaleFactor * float64(tileX)),
				Y: int32(minimapScaleFactor * float64(tileY)),
				W: int32(math.Floor(minimapScaleFactor * tileSize)),
				H: int32(math.Floor(minimapScaleFactor * tileSize)),
			}
			renderer.FillRect(rect)
		}
	}
}
