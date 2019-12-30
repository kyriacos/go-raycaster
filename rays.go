package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Rays - many rays
type Rays [NumRays]Ray

// Ray - struct
type Ray struct {
	angle, wallHitX, wallHitY, distance float64
	wasHitVertical                      bool

	isRayFacingUp, isRayFacingDown,
	isRayFacingLeft, isRayFacingRight bool

	wallHitContent int // store the actual content of the wall once we find a hit
}

func (r *Ray) render(renderer *sdl.Renderer, x, y float64) {
	renderer.SetDrawColor(255, 0, 0, 30)
	renderer.DrawLine(
		int32(MinimapScaleFactor*x),
		int32(MinimapScaleFactor*y),
		int32(MinimapScaleFactor*r.wallHitX),
		int32(MinimapScaleFactor*r.wallHitY),
	)
}

// NewRay - constructor
func NewRay(angle float64) Ray {
	nAngle := normalizeAngle(angle)

	// we have to figure out which way the ray is facing when trying to calculate the intersects.
	// Math.PI = 180 Degrees
	isRayFacingDown := false
	if nAngle > 0 && nAngle < PI {
		isRayFacingDown = true
	}

	// 0.5 * Math.PI = 90 Degrees
	// 1.5 * Math.PI = 270 Degrees
	isRayFacingRight := false
	if nAngle < 0.5*PI || nAngle > 1.5*PI {
		isRayFacingRight = true
	}

	return Ray{
		angle:            nAngle,
		wallHitX:         0.0,
		wallHitY:         0.0,
		distance:         0.0,
		wasHitVertical:   false,
		isRayFacingDown:  isRayFacingDown,
		isRayFacingUp:    !isRayFacingDown,
		isRayFacingRight: isRayFacingRight,
		isRayFacingLeft:  !isRayFacingRight,

		wallHitContent: -1, // set to -1 for debugging in case
	}
}
