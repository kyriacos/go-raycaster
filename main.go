package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600

	fps             = 30
	frameTimeLength = 1000 / fps
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer

	running               = false
	ticksLastFrame uint32 = 0

	playerX, playerY = 0, 0
)

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
		sdl.WINDOW_BORDERLESS)
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

}

func update() {
	// stop and waste some time until we reach the target frame time length we want
	// timeout = SDL_GetTicks() + frameTimeLength
	// !SDL_TICKS_PASSED(SDL.GetTicks(), timeout)
	sdl.Delay(frameTimeLength)

	deltaTime := float64(sdl.GetTicks()-ticksLastFrame) / 1000.0
	ticksLastFrame = sdl.GetTicks()

	playerX += int(50.0 * deltaTime)
	playerY += int(50.0 * deltaTime)
}

func render() {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear() // clear back buffer

	// render all game objects for current frame
	renderer.SetDrawColor(255, 255, 0, 255)
	rect := sdl.Rect{X: int32(playerX), Y: int32(playerY), W: 20, H: 20}
	renderer.FillRect(&rect)

	// swap current buffer with back buffer
	renderer.Present()
}

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
