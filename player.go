package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Player - stuff
type Player struct {
	x, y          float64
	width, height float64

	turnDirection int // -1 left, +1 right
	walkDirection int // -1 left, +1 right

	rotationAngle float64

	walkSpeed float64
	turnSpeed float64
}

func (p *Player) render(renderer *sdl.Renderer) {
	renderer.SetDrawColor(255, 255, 255, 255)
	rect := &sdl.Rect{
		X: int32(minimapScaleFactor * player.x),
		Y: int32(minimapScaleFactor * player.y),
		W: int32(minimapScaleFactor * player.width),
		H: int32(minimapScaleFactor * player.height),
	}
	renderer.FillRect(rect)

	// draw line to show direction the player is looking
	renderer.DrawLine(
		int32(minimapScaleFactor*player.x),
		int32(minimapScaleFactor*player.y),
		int32(minimapScaleFactor*player.x+math.Cos(player.rotationAngle)*40),
		int32(minimapScaleFactor*player.y+math.Sin(player.rotationAngle)*40),
	)
}

// tick
// note: maybe i don't need to pass the gamemap and just have the collision detection inside my 'main' update
func (p *Player) update(deltaTime float64, gm *GameMap) {
	move(deltaTime, gm)
}

func move(deltaTime float64, gm *GameMap) {
	// Turning: its the turn direction -1/+1/0 multiplied by the rotation speed
	player.rotationAngle += float64(player.turnDirection) * player.turnSpeed * deltaTime

	// Moving:  direction -1/+1/0 multiplied by how fast (speed)
	//          the player should move to calculate how much of a jump/step we make
	moveStep := float64(player.walkDirection) * player.walkSpeed * deltaTime

	newX := player.x + math.Cos(player.rotationAngle)*moveStep
	newY := player.y + math.Sin(player.rotationAngle)*moveStep

	// perform wall collision check
	if !gameMap.hasWallAt(newX, newY) {
		player.x = newX
		player.y = newY
	}
}
