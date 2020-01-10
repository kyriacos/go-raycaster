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

// All the game globals
var (
	Window   *sdl.Window   // The main window that we render to
	Renderer *sdl.Renderer // The SDL renderer

	CB        *ColorBuffer // The instance of ColorBuffer we use to update every tick
	CBTexture *sdl.Texture

	Textures map[string]*image.NRGBA // Stores all the texture images

	G *Game // The game instance

	showFPS = flag.Bool("showFPS", false, "Show current FPS and on exit display the average FPS.")
)

func castAllRays() {
	// initial ray angle
	angle := G.Player.rotationAngle - (FOV / 2)

	for column := 0; column < NumRays; column++ {
		ray := G.Rays[column]
		ray.Cast(angle)
		angle += FOV / NumRays
	}
}

func run() (err error) {
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing SDL: %s\n", err)
		return
	}

	Window, err = sdl.CreateWindow(
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

	Renderer, err = sdl.CreateRenderer(Window, -1, sdl.RENDERER_SOFTWARE) // -1 is the default driver (the graphics driver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating SDL renderer: %s\n", err)
		return
	}

	if err = Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set blend mode: %s", err)
		return
	}

	return nil
}

func destroy() {
	// defer order?
	defer sdl.Quit()
	defer Window.Destroy()
	defer Renderer.Destroy()
	defer CBTexture.Destroy()
}

func loadTextures() {
	imageDir := "./images/"
	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		log.Fatal(err)
	}

	Textures = make(map[string]*image.NRGBA, len(files))
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

		Textures[strings.TrimSuffix(filename, path.Ext(filename))] = imgNRGBA
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
	level := LoadLevel("./levels/level1.json")
	G.GameMap = NewGameMap(level)

	// initialize rays
	G.Rays = NewRays()

	// initialize the player
	G.Player = &Player{
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
	CB = NewColorBuffer(WindowWidth, WindowHeight)

	// create color buffer texture
	var err error
	CBTexture, err = Renderer.CreateTexture(
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
	G.Player.Update(elapsedMS * 1000.0)

	castAllRays()
}

func renderColorBuffer() {
	// update the sdl texture
	CBTexture.Update(nil, CB.Pixels, CB.GetPitch())

	// copy the texture to the renderer
	Renderer.Copy(CBTexture, nil, nil) // nil and nil since we want to use the entire texture (src and dest used if you want to get a subset of the texture)
}

func project3d() {
	for i := 0; i < NumRays; i++ {
		ray := G.Rays[i]
		// calculate perpendicular distance to remove the fisheye effect
		perpendicularDistance := ray.distance * math.Cos(ray.angle-G.Player.rotationAngle)
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
			CB.Set(i, y, 0x333333FF)
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

			texel := Textures["redbrick"].NRGBAAt(int(textureOffsetX), int(textureOffsetY))
			var c uint32 = uint32(texel.R)<<24 | uint32(texel.G)<<16 | uint32(texel.B)<<8 | uint32(texel.A)
			CB.Set(i, y, c)
		}

		// set color for the floor
		for y := wallBottomPixel; y < WindowHeight; y++ {
			CB.Set(i, y, 0x777777FF)
		}
	}
}

func render() {
	Renderer.SetDrawColor(0, 0, 0, 255)
	Renderer.Clear() // clear back buffer

	project3d()
	renderColorBuffer()

	// render all game objects for current frame
	G.GameMap.Render()
	G.Player.Render()

	for _, ray := range G.Rays {
		ray.Render(Renderer, G.Player.x, G.Player.y)
	}

	// swap current buffer with back buffer
	Renderer.Present()
}

func processInput() {
	if event := sdl.PollEvent(); event != nil {
		switch t := event.(type) {
		case *sdl.QuitEvent: // sdl.QUIT
			// println("Quit")
			G.Running = false
		case *sdl.KeyboardEvent:
			key := t.Keysym.Sym
			if t.Type == sdl.KEYDOWN {
				switch key {
				case sdl.K_ESCAPE:
					G.Running = false
				case sdl.K_UP:
					G.Player.walkDirection = 1
				case sdl.K_DOWN:
					G.Player.walkDirection = -1
				case sdl.K_RIGHT:
					G.Player.turnDirection = 1
				case sdl.K_LEFT:
					G.Player.turnDirection = -1
				}
			}
			if t.Type == sdl.KEYUP {
				switch key {
				case sdl.K_UP, sdl.K_DOWN:
					G.Player.walkDirection = 0
				case sdl.K_RIGHT, sdl.K_LEFT:
					G.Player.turnDirection = 0
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	G = &Game{
		Running:        false,
		TicksLastFrame: 0,
	}

	if err := run(); err != nil {
		destroy()
		os.Exit(1)
	}

	var (
		counter           = 0
		elapsedMS, sumFPS float64
	)

	setup()

	G.Running = true
	for G.Running {
		start := sdl.GetPerformanceCounter()

		processInput()
		update(elapsedMS)
		render()

		end := sdl.GetPerformanceCounter()

		elapsedMS = float64(end-start) / float64(sdl.GetPerformanceFrequency()*1000.0)

		sdl.Delay(uint32(math.Floor(FrameTimeLength - elapsedMS))) // pause until we reach the target frames

		if *showFPS {
			elapsed := float64(end-start) / float64(sdl.GetPerformanceFrequency())
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
