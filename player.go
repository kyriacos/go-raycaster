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

func (p *Player) Render() {
	Renderer.SetDrawColor(255, 255, 255, 255)
	rect := &sdl.Rect{
		X: int32(MinimapScaleFactor * p.x),
		Y: int32(MinimapScaleFactor * p.y),
		W: int32(MinimapScaleFactor * p.width),
		H: int32(MinimapScaleFactor * p.height),
	}
	Renderer.FillRect(rect)

	/*
	 * Add a line to see which angle my player is turning
	 *
	 *    |\
	 * 30 | \  y
	 *    |  \
	 *    |a  \
	 *    -----
	 *      x
	 *
	 */
	length := 30 * MinimapScaleFactor
	Renderer.DrawLine(
		int32(MinimapScaleFactor*p.x),
		int32(MinimapScaleFactor*p.y),
		int32(MinimapScaleFactor*p.x+math.Cos(p.rotationAngle)*length),
		int32(MinimapScaleFactor*p.y+math.Sin(p.rotationAngle)*length),
	)
}

func (p *Player) Update(deltaTime float64) {
	p.move(deltaTime)
}

func (p *Player) move(deltaTime float64) {
	// Turning: its the turn direction -1/+1/0 multiplied by the rotation speed
	p.rotationAngle += float64(p.turnDirection) * p.turnSpeed * deltaTime

	// Moving:  direction -1/+1/0 multiplied by how fast (speed)
	//          the player should move to calculate how much of a jump/step we make
	moveStep := float64(p.walkDirection) * p.walkSpeed * deltaTime

	newX := p.x + math.Cos(p.rotationAngle)*moveStep
	newY := p.y + math.Sin(p.rotationAngle)*moveStep

	// perform wall collision check
	if !G.GameMap.HasWallAt(newX, newY) {
		p.x = newX
		p.y = newY
	}
}
