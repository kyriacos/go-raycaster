package main

type Game struct {
	Running        bool   //= false
	TicksLastFrame uint32 // = 0

	Player  *Player
	GameMap *GameMap
	Rays    *Rays
}
