package main

import (
	"fmt"
	"image/color"
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

	colorBuffer        *ColorBuffer
	colorBufferTexture *sdl.Texture
)

func castAllRays() {
	// initial ray angle
	angle := player.rotationAngle - (FOV / 2)

	for column := 0; column < NumRays; column++ {
		ray := NewRay(angle)
		rays[column] = ray.cast()
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
	defer colorBufferTexture.Destroy()
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

	// initialize the color buffer
	colorBuffer = (new(ColorBuffer)).Init()

	// create color buffer texture
	var err error
	colorBufferTexture, err = renderer.CreateTexture(
		sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_STREAMING,
		WindowWidth,
		WindowHeight,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating the texture: %s", err)
		panic(err)
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

func renderColorBuffer() {
	// update the sdl texture
	colorBufferTexture.Update(nil, colorBuffer.Pix, calculatePitch())

	// copy the texture to the renderer
	renderer.Copy(colorBufferTexture, nil, nil) // nil and nil since we want to use the entire texture (src and dest used if you want to get a subset of the texture)
}

func project3d() {
	for i := 0; i < NumRays; i++ {
		ray := rays[i]
		// calculate perpendicular distance to remove the fisheye effect
		perpendicularDistance := ray.distance * math.Cos(ray.angle-player.rotationAngle)
		distanceToProjPlane := (WindowWidth / 2) / math.Tan(FOV/2)
		projectedWallHeight := (TileSize / perpendicularDistance) * distanceToProjPlane

		wallStripHeight := int(projectedWallHeight)

		// where the wall starts - starts right after our ceiling
		wallTopPixel := (WindowHeight / 2) - (wallStripHeight / 2) // middle of the screen and half the wall height
		if wallTopPixel < 0 {
			wallTopPixel = 0
		}
		// ends where the floor starts rendering
		wallBottomPixel := (WindowHeight / 2) + (wallStripHeight / 2)
		if wallBottomPixel > WindowHeight {
			wallBottomPixel = WindowHeight
		}

		// set color for the ceiling
		for y := 0; y < wallTopPixel; y++ {
			colorBuffer.Set(i, y, uint32ToColorRGBA(0x333333FF))
		}

		// render the wall from top to bottom - cols
		for y := wallTopPixel; y < wallBottomPixel; y++ {
			var c color.RGBA
			if ray.wasHitVertical {
				c = uint32ToColorRGBA(0xFFFFFFFF)
			} else {
				c = uint32ToColorRGBA(0xCCCCCCFF) // grey
			}
			colorBuffer.Set(i, y, c)
		}

		// set color for the floor
		for y := wallBottomPixel; y < WindowHeight; y++ {
			colorBuffer.Set(i, y, uint32ToColorRGBA(0x777777FF))
		}
	}
}

func render() {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear() // clear back buffer

	project3d()
	renderColorBuffer()
	colorBuffer.Clear(0x00000000) // clear the color buffer

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
