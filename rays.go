package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Rays - many rays
type Rays [NumRays]*Ray

func NewRays() *Rays {
	r := new(Rays)
	for i := range r {
		r[i] = NewRay()
	}
	return r
}

// Ray - struct
type Ray struct {
	angle, wallHitX, wallHitY, distance float64
	wasHitVertical                      bool

	isRayFacingUp, isRayFacingDown,
	isRayFacingLeft, isRayFacingRight bool

	wallHitContent int // store the actual content of the wall once we find a hit
}

// NewRay - constructor
func NewRay() *Ray {
	return &Ray{
		angle:            0.0,
		wallHitX:         -1,
		wallHitY:         -1,
		distance:         -1,
		wasHitVertical:   false,
		isRayFacingDown:  false,
		isRayFacingUp:    false,
		isRayFacingRight: false,
		isRayFacingLeft:  false,

		wallHitContent: -1, // set to -1 for debugging in case
	}
}

func (r *Ray) Render(renderer *sdl.Renderer, x, y float64) {
	renderer.SetDrawColor(255, 0, 0, 30)
	renderer.DrawLine(
		int32(MinimapScaleFactor*x),
		int32(MinimapScaleFactor*y),
		int32(MinimapScaleFactor*r.wallHitX),
		int32(MinimapScaleFactor*r.wallHitY),
	)
}

func (r *Ray) Cast(angle float64) *Ray {
	var xIntercept, yIntercept, xStep, yStep float64

	r.angle = normalizeAngle(angle)

	// we have to figure out which way the ray is facing when trying to calculate the intersects.
	// Math.PI = 180 Degrees
	r.isRayFacingDown = false
	if r.angle > 0 && r.angle < PI {
		r.isRayFacingDown = true
	}
	r.isRayFacingUp = !r.isRayFacingDown

	// 0.5 * Math.PI = 90 Degrees
	// 1.5 * Math.PI = 270 Degrees
	r.isRayFacingRight = false
	if r.angle < 0.5*PI || r.angle > 1.5*PI {
		r.isRayFacingRight = true
	}
	r.isRayFacingLeft = !r.isRayFacingRight

	/*
	 * ================================
	 * HORIZONTAL ray-grid intersection
	 * ================================
	 *
	 */

	foundHorizontalWallHit := false
	horzWallHitX := 0.0
	horzWallHitY := 0.0
	horzWallContent := 0

	/* Find the y-coordinate of the closest horizontal grid intersection
	 * =================================================================
	 *
	 * Ray facing up. (standard - base assumption)
	 * -----*--  => Ay (yIntercept) when ray is facing up. That's what we always calculate
	 * |      |     No need to add TILE_SIZE
	 * |      |
	 * |      |
	 * |      |
	 * |      |
	 * --------
	 *
	 * Ray facing down
	 * --------
	 * |      |
	 * |      |
	 * |      |
	 * |      |
	 * |      |
	 * -----*--  => Ay (yIntercept) when ray is facing down.
	 *              We take yIntercept = Math.floor(player.y / TILE_SIZE) * TILE_SIZE;
	 *              And we add TILE_SIZE
	 *
	 * So the yIntercept value above relates to if the ray is actually pointing up.
	 * We actually calculate that initially using:
	 * yIntercept = Math.floor(player.y / TILE_SIZE) * TILE_SIZE
	 *
	 * Think of the box (tile) - yIntercept calculated initially with the above formula will
	 * give the Ay coordinate sitting at the 'top' of the tile i.e. when the ray is facing up.
	 *
	 * If however it's actually facing down then we add TILE_SIZE to it
	 * yIntercept += this.isRayFacingDown ?  TILE_SIZE : 0;
	 *
	 */
	yIntercept = math.Floor(G.Player.y/TileSize) * TileSize // this always gets the 'top' Ay coordinate i.e. ray facing up
	if r.isRayFacingDown {                                  // else += 0
		yIntercept += TileSize
	}

	// Find the x-coordinate of the closest horizontal grid interception
	// xIntercept = G.Player.x + ((G.Player.y - yIntercept) / Math.tan(this.angle));
	xIntercept = G.Player.x + ((yIntercept - G.Player.y) / math.Tan(r.angle))

	// Calculate the increment xstep and ystep
	yStep = TileSize
	if r.isRayFacingUp { // depending on where the ray is facing we invert or not
		yStep *= -1
	}

	xStep = TileSize / math.Tan(r.angle)
	if r.isRayFacingLeft && xStep > 0 { // if the xstep is positive but the ray is facing left we invert
		xStep *= -1
	}
	if r.isRayFacingRight && xStep < 0 { // if the xstep is negative but the ray is facing right we invert
		xStep *= -1
	}
	nextHorzTouchX := xIntercept
	nextHorzTouchY := yIntercept

	// increment xstep and ystep until we find a wall
	for nextHorzTouchX >= 0 &&
		nextHorzTouchX <= WindowWidth &&
		nextHorzTouchY >= 0 &&
		nextHorzTouchY <= WindowHeight {

		testTouchX := nextHorzTouchX
		testTouchY := nextHorzTouchY
		// Since right now the Y position is sitting right on the border of the tile to check
		// if we are actually in the wall i.e. not skip over it we can force it by
		// just pushing it by 1 (i.e. subtract 1 if its facing up) to get it in
		// only used for checking we don't want to change the value of nextHorzTouchY as otherwise we will be skipping a value each time we loop
		if r.isRayFacingUp { // force one pixel up
			testTouchY = nextHorzTouchY - 1
		}

		// Found a wall hit
		if G.GameMap.HasWallAt(testTouchX, testTouchY) {
			horzWallHitX = nextHorzTouchX
			horzWallHitY = nextHorzTouchY
			horzWallContent = G.GameMap.Level.At(int(math.Floor(testTouchY/TileSize)), int(math.Floor(testTouchX/TileSize)))
			foundHorizontalWallHit = true
			break
		} else {
			nextHorzTouchX += xStep
			nextHorzTouchY += yStep
		}
	}

	/*
	 *
	 * ================================
	 * VERTICAL ray-grid intersection
	 * ================================
	 *
	 */

	foundVerticalWallHit := false
	vertWallHitX := 0.0
	vertWallHitY := 0.0
	vertWallContent := 0

	// Find the x-coordinate of the closest vertical grid interception
	xIntercept = math.Floor(G.Player.x/TileSize) * TileSize
	if r.isRayFacingRight { // add 32 (tile_size) if facing right
		xIntercept += TileSize
	}

	// Find the y-coordinate of the closest vertical grid interception
	// yIntercept = G.Player.y + ((G.Player.x - xIntercept) * Math.tan(this.angle));
	yIntercept = G.Player.y + ((xIntercept - G.Player.x) * math.Tan(r.angle))

	// Calculate the increment xstep and ystep
	xStep = TileSize
	if r.isRayFacingLeft { // depending on where the ray is facing we invert or not
		xStep *= -1
	}

	yStep = TileSize * math.Tan(r.angle)
	if r.isRayFacingUp && yStep > 0 {
		yStep *= -1
	}
	if r.isRayFacingDown && yStep < 0 {
		yStep *= -1
	}

	nextVertTouchX := xIntercept
	nextVertTouchY := yIntercept

	// increment xstep and ystep until we find a wall
	for nextVertTouchX >= 0 &&
		nextVertTouchX <= WindowWidth &&
		nextVertTouchY >= 0 &&
		nextVertTouchY <= WindowHeight {

		testTouchX := nextVertTouchX
		testTouchY := nextVertTouchY
		// since right now the X position is sitting right on the border of the tile
		// to check if we are actually in the wall i.e. not skip over it
		// we can force it by just pushing it by 1 (i.e. subtract 1 if its facing left) to get it in the next tile/wall/pixel
		// only used for checking we don't want to change the value of nextVertTouchX as otherwise we will be skipping a value each time we loop
		if r.isRayFacingLeft { // force one pixel left
			testTouchX = nextVertTouchX - 1
		}

		if G.GameMap.HasWallAt(testTouchX, nextVertTouchY) {
			vertWallHitX = nextVertTouchX
			vertWallHitY = nextVertTouchY
			vertWallContent = G.GameMap.Level.At(int(math.Floor(testTouchY/TileSize)), int(math.Floor(testTouchX/TileSize)))

			foundVerticalWallHit = true
			break
		} else {
			nextVertTouchX += xStep
			nextVertTouchY += yStep
		}
	}

	// Calculate both horizontal and vertical distances and choose the smallest value
	horzHitDistance := math.MaxFloat64 // if we didn't get a hit then we basically just set it to a really large value
	if foundHorizontalWallHit {
		horzHitDistance = distanceBetweenPoints(G.Player.x, G.Player.y, horzWallHitX, horzWallHitY)
	}
	vertHitDistance := math.MaxFloat64 // if we didn't get a hit then we basically just set it to a really large value
	if foundVerticalWallHit {
		vertHitDistance = distanceBetweenPoints(G.Player.x, G.Player.y, vertWallHitX, vertWallHitY)
	}

	// Compare the two and store the smallest one
	if horzHitDistance < vertHitDistance {
		r.wallHitX = horzWallHitX
		r.wallHitY = horzWallHitY
		r.distance = horzHitDistance
		r.wallHitContent = horzWallContent
		r.wasHitVertical = false
	} else {
		r.wallHitX = vertWallHitX
		r.wallHitY = vertWallHitY
		r.distance = vertHitDistance
		r.wallHitContent = vertWallContent
		r.wasHitVertical = true
	}

	return r

}
