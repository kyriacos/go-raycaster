package main

import (
	"fmt"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	fps             = 30
	frameTimeLength = 1000 / fps

	pi    = math.Pi
	twoPI = 2 * pi

	tileSize   = 64
	mapNumRows = 13
	mapNumCols = 20

	minimapScaleFactor = 0.2

	windowWidth  = mapNumCols * tileSize
	windowHeight = mapNumRows * tileSize

	fov     = 60 * (math.Pi / 180)
	numRays = windowWidth
)

var grid = [mapNumRows][mapNumCols]int{
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

func run() (err error) {
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing SDL: %s\n", err)
		return
	}

	window, err = sdl.CreateWindow(
		"RayCaster",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		windowWidth,
		windowHeight,
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
	// initialize the player
	// player := &Player{
	// 	x:             windowWidth / 2,
	// 	y:             windowHeight / 2,
	// 	width:         5,
	// 	height:        5,
	// 	turnDirection: 0,
	// 	walkDirection: 0,
	// 	rotationAngle: pi / 2,
	// 	walkSpeed:     100,
	// 	turnSpeed:     45 * (pi / 180),
	// }
}

func update() {
	/* stop and waste some time until we reach the target frame time length we want
	 * timeout = SDL_GetTicks() + frameTimeLength
	 * !SDL_TICKS_PASSED(SDL.GetTicks(), timeout)
	 */
	sdl.Delay(frameTimeLength)

	// deltaTime := float64(sdl.GetTicks()-ticksLastFrame) / 1000.0
	// ticksLastFrame = sdl.GetTicks()

}

func render() {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear() // clear back buffer

	// render all game objects for current frame
	renderMap()
	renderPlayer()

	// swap current buffer with back buffer
	renderer.Present()
}

func renderMap() {
	for i := 0; i < mapNumRows; i++ {
		for j := 0; j < mapNumCols; j++ {
			tileX := j * tileSize // column
			tileY := i * tileSize // row

			var tileColor uint8 = 0
			if grid[i][j] != 0 {
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

func renderPlayer() {}

func processInput() {
	if event := sdl.PollEvent(); event != nil {
		switch k := event.(type) {
		case *sdl.QuitEvent: // sdl.QUIT
			// println("Quit")
			running = false
		case *sdl.KeyboardEvent:
			if k.Keysym.Sym == sdl.K_ESCAPE {
				running = false
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
