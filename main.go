package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"strings"

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
	rays    *Rays

	colorBuffer        *ColorBuffer
	colorBufferTexture *sdl.Texture

	textures map[string]*image.NRGBA

	showFPS = flag.Bool("showFPS", false, "Show current FPS and on exit display the average FPS.")
)

func castAllRays() {
	// initial ray angle
	angle := player.rotationAngle - (FOV / 2)

	for column := 0; column < NumRays; column++ {
		ray := rays[column]
		ray.cast(angle)
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

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE) // -1 is the default driver (the graphics driver)
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

func loadTextures() {
	imageDir := "./images/"
	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		log.Fatal(err)
	}

	textures = make(map[string]*image.NRGBA, len(files))
	for _, file := range files {
		filename := file.Name()
		f, err := os.Open(imageDir + filename)
		if err != nil {
			log.Fatalf("Could not open file: %s", err)
		}
		defer f.Close()

		imgNRGBA, err := decodeImage(f)
		if err != nil {
			log.Fatalf("Could not decode image from file: %s. Error: %s", err, imageDir+filename)
		}

		textures[strings.TrimSuffix(filename, path.Ext(filename))] = imgNRGBA
	}
}

func decodeImage(r io.Reader) (*image.NRGBA, error) {
	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	var imgNRGBA *image.NRGBA
	var ok bool
	if imgNRGBA, ok = img.(*image.NRGBA); !ok {
		switch img.ColorModel() {
		case color.RGBAModel:
			imgNRGBA = image.NewNRGBA(img.Bounds())
			draw.Draw(imgNRGBA, img.Bounds(), img, image.Point{}, draw.Src)
		case color.GrayModel, color.Gray16Model, color.AlphaModel, color.Alpha16Model:
			fallthrough
		default:
			log.Fatal("Unsupported image format.")
		}
	}
	return imgNRGBA, nil
}

func setup() {
	// Load textures from images directory
	loadTextures()

	// initialize map
	gameMap = NewGameMap(level1)

	// initialize rays
	rays = NewRays()

	// initialize the player
	player = &Player{
		x:             WindowWidth / 2,
		y:             WindowHeight / 2,
		width:         1,
		height:        1,
		turnDirection: 0,
		walkDirection: 0,
		rotationAngle: 2 * PI / 2,
		walkSpeed:     100,
		turnSpeed:     70 * (PI / 180),
	}

	// initialize the color buffer
	colorBuffer = NewColorBuffer(WindowWidth, WindowHeight)

	// create color buffer texture
	var err error
	colorBufferTexture, err = renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888, // endianess https://forums.libsdl.org/viewtopic.php?p=39284
		sdl.TEXTUREACCESS_STREAMING,
		WindowWidth,
		WindowHeight,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating the texture: %s", err)
		panic(err)
	}

}

func update(elapsedMS float64) {
	// deltaTime := float64(sdl.GetTicks()-ticksLastFrame) / 1000.0
	// ticksLastFrame = sdl.GetTicks()

	player.update(elapsedMS*1000.0, gameMap)

	castAllRays()

	/* stop and waste some time until we reach the target frame time length we want
	 * timeout = SDL_GetTicks() + frameTimeLength
	 * !SDL_TICKS_PASSED(SDL.GetTicks(), timeout)
	 */
	// sdl.Delay(uint32(FrameTimeLength - deltaTime))
}

func renderColorBuffer() {
	// update the sdl texture
	colorBufferTexture.Update(nil, colorBuffer.Pixels, colorBuffer.GetPitch())

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
			colorBuffer.Set(i, y, 0x333333FF)
		}

		// same for all the columns of X
		var textureOffsetX int
		if ray.wasHitVertical { // use Y to get the offset instead
			textureOffsetX = int(ray.wallHitY) % TextureHeight
		} else {
			textureOffsetX = int(ray.wallHitX) % TextureWidth
		}

		// render the wall from top to bottom - cols
		for y := wallTopPixel; y < wallBottomPixel; y++ {
			distanceFromTop := y + (wallStripHeight / 2) - (WindowHeight / 2)
			textureOffsetY := float64(distanceFromTop) * float64(TextureHeight) / float64(wallStripHeight)

			texel := textures["redbrick"].NRGBAAt(int(textureOffsetX), int(textureOffsetY))
			var c uint32 = uint32(texel.R)<<24 | uint32(texel.G)<<16 | uint32(texel.B)<<8 | uint32(texel.A)
			colorBuffer.Set(i, y, c)
		}

		// set color for the floor
		for y := wallBottomPixel; y < WindowHeight; y++ {
			colorBuffer.Set(i, y, 0x777777FF)
		}
	}
}

func render() {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear() // clear back buffer

	project3d()
	renderColorBuffer()
	// colorBuffer.Clear(0x00000000) // clear the color buffer

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
	flag.Parse()

	if err := run(); err != nil {
		destroy()
		os.Exit(1)
	}

	var (
		counter           = 0
		elapsedMS, sumFPS float64
	)

	setup()

	running = true
	for running {
		start := sdl.GetPerformanceCounter()
		processInput()
		update(elapsedMS)
		render()
		end := sdl.GetPerformanceCounter()

		elapsedMS = float64(end-start) / float64(sdl.GetPerformanceFrequency()*1000.0)

		sdl.Delay(uint32(math.Floor(16.666 - elapsedMS)))
		elapsed := float64(end-start) / float64(sdl.GetPerformanceFrequency())

		if *showFPS {
			counter++
			currentFPS := 1.0 / elapsed
			sumFPS += currentFPS

			fmt.Printf("FPS: %f\n", 1.0/elapsed)
		}
	}

	destroy()

	if *showFPS {
		fmt.Printf("Average FPS: %f\n", sumFPS/float64(counter))
	}

	os.Exit(0)
}
