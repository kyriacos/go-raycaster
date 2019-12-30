package main

import (
	"fmt"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var level1 = Level{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

var (
	window   *sdl.Window
	renderer *sdl.Renderer

	running               = false
	ticksLastFrame uint32 = 0

	player  *Player
	gameMap *GameMap

	rays Rays = Rays{}
)

func castRay(ray Ray) Ray {
	var xIntercept, yIntercept, xStep, yStep float64

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
	yIntercept = math.Floor(player.y/TileSize) * TileSize // this always gets the 'top' Ay coordinate i.e. ray facing up
	if ray.isRayFacingDown {                              // else += 0
		yIntercept += TileSize
	}

	// Find the x-coordinate of the closest horizontal grid interception
	// xIntercept = player.x + ((player.y - yIntercept) / Math.tan(this.angle));
	xIntercept = player.x + ((yIntercept - player.y) / math.Tan(ray.angle))

	// Calculate the increment xstep and ystep
	yStep = TileSize
	if ray.isRayFacingUp { // depending on where the ray is facing we invert or not
		yStep *= -1
	}

	xStep = TileSize / math.Tan(ray.angle)
	if ray.isRayFacingLeft && xStep > 0 { // if the xstep is positive but the ray is facing left we invert
		xStep *= -1
	}
	if ray.isRayFacingRight && xStep < 0 { // if the xstep is negative but the ray is facing right we invert
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
		if ray.isRayFacingUp { // force one pixel up
			testTouchY = nextHorzTouchY - 1
		}

		// Found a wall hit
		if gameMap.hasWallAt(testTouchX, testTouchY) {
			horzWallHitX = nextHorzTouchX
			horzWallHitY = nextHorzTouchY
			horzWallContent = gameMap.level[int(math.Floor(testTouchY/TileSize))][int(math.Floor(testTouchX/TileSize))]
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
	xIntercept = math.Floor(player.x/TileSize) * TileSize
	if ray.isRayFacingRight { // add 32 (tile_size) if facing right
		xIntercept += TileSize
	}

	// Find the y-coordinate of the closest vertical grid interception
	// yIntercept = player.y + ((player.x - xIntercept) * Math.tan(this.angle));
	yIntercept = player.y + ((xIntercept - player.x) * math.Tan(ray.angle))

	// Calculate the increment xstep and ystep
	xStep = TileSize
	if ray.isRayFacingLeft { // depending on where the ray is facing we invert or not
		xStep *= -1
	}

	yStep = TileSize * math.Tan(ray.angle)
	if ray.isRayFacingUp && yStep > 0 {
		yStep *= -1
	}
	if ray.isRayFacingDown && yStep < 0 {
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
		if ray.isRayFacingLeft { // force one pixel left
			testTouchX = nextVertTouchX - 1
		}

		if gameMap.hasWallAt(testTouchX, nextVertTouchY) {
			vertWallHitX = nextVertTouchX
			vertWallHitY = nextVertTouchY
			vertWallContent = gameMap.level[int(math.Floor(testTouchY/TileSize))][int(math.Floor(testTouchX/TileSize))]

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
		horzHitDistance = distanceBetweenPoints(player.x, player.y, horzWallHitX, horzWallHitY)
	}
	vertHitDistance := math.MaxFloat64 // if we didn't get a hit then we basically just set it to a really large value
	if foundVerticalWallHit {
		vertHitDistance = distanceBetweenPoints(player.x, player.y, vertWallHitX, vertWallHitY)
	}

	// Compare the two and store the smallest one
	if horzHitDistance < vertHitDistance {
		ray.wallHitX = horzWallHitX
		ray.wallHitY = horzWallHitY
		ray.distance = horzHitDistance
		ray.wallHitContent = horzWallContent
		ray.wasHitVertical = false
	} else {
		ray.wallHitX = vertWallHitX
		ray.wallHitY = vertWallHitY
		ray.distance = vertHitDistance
		ray.wallHitContent = vertWallContent
		ray.wasHitVertical = true
	}

	return ray

}

func castAllRays() {
	// initial ray angle
	angle := player.rotationAngle - (FOV / 2)

	for column := 0; column < NumRays; column++ {
		ray := NewRay(angle)
		rays[column] = castRay(ray)
		angle += FOV / NumRays
	}
}

func run() (err error) {
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing SDL: %s\n", err)
		return
	}

	window, err = sdl.CreateWindow(
		"RayCaster",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		WindowWidth,
		WindowHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating SDL window: %s\n", err)
		return
	}

	renderer, err = sdl.CreateRenderer(window, -1, 0) // -1 is the default driver (the graphics driver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating SDL renderer: %s\n", err)
		return
	}

	if err = renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set blend mode: %s", err)
		return
	}

	return nil
}

func destroy() {
	// defer order?
	defer sdl.Quit()
	defer window.Destroy()
	defer renderer.Destroy()
}

func setup() {
	// initialize map
	gameMap = &GameMap{
		level: level1,
	}

	// initialize the player
	player = &Player{
		x:             WindowWidth / 2,
		y:             WindowHeight / 2,
		width:         1,
		height:        1,
		turnDirection: 0,
		walkDirection: 0,
		rotationAngle: PI / 2,
		walkSpeed:     100,
		turnSpeed:     70 * (PI / 180),
	}
}

func update() {
	/* stop and waste some time until we reach the target frame time length we want
	 * timeout = SDL_GetTicks() + frameTimeLength
	 * !SDL_TICKS_PASSED(SDL.GetTicks(), timeout)
	 */
	sdl.Delay(FrameTimeLength)

	deltaTime := float64(sdl.GetTicks()-ticksLastFrame) / 1000.0
	ticksLastFrame = sdl.GetTicks()

	player.update(deltaTime, gameMap)

	castAllRays()
}

func render() {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear() // clear back buffer

	// render all game objects for current frame
	gameMap.render(renderer)
	player.render(renderer)

	for _, ray := range rays {
		ray.render(renderer, player.x, player.y)
	}

	// swap current buffer with back buffer
	renderer.Present()
}

func processInput() {
	if event := sdl.PollEvent(); event != nil {
		switch t := event.(type) {
		case *sdl.QuitEvent: // sdl.QUIT
			// println("Quit")
			running = false
		case *sdl.KeyboardEvent:
			key := t.Keysym.Sym
			if t.Type == sdl.KEYDOWN {
				switch key {
				case sdl.K_ESCAPE:
					running = false
				case sdl.K_UP:
					player.walkDirection = 1
				case sdl.K_DOWN:
					player.walkDirection = -1
				case sdl.K_RIGHT:
					player.turnDirection = 1
				case sdl.K_LEFT:
					player.turnDirection = -1
				}
			}
			if t.Type == sdl.KEYUP {
				switch key {
				case sdl.K_UP, sdl.K_DOWN:
					player.walkDirection = 0
				case sdl.K_RIGHT, sdl.K_LEFT:
					player.turnDirection = 0
				}
			}
		}
	}
}

func main() {
	if err := run(); err != nil {
		destroy()
		os.Exit(1)
	}

	setup()

	running = true
	for running {
		processInput()
		update()
		render()
	}

	destroy()
	os.Exit(0)
}
