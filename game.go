package main

import (
	"image"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer

	Running        bool   //= false
	TicksLastFrame uint32 // = 0

	Player  *Player
	GameMap *GameMap
	Rays    *Rays

	ColorBuffer        *ColorBuffer
	ColorBufferTexture *sdl.Texture

	Textures map[string]*image.NRGBA
}
