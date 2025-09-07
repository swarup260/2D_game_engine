package ecs

import (

	"game_engine/physics"

	"github.com/veandco/go-sdl2/sdl"
)

// Transform component for position, rotation, and scale
type Transform struct {
	X, Y     float64
	Rotation float64
	ScaleX   float64
	ScaleY   float64
}

// Velocity component for movement
type Velocity struct {
	X, Y float64
}

// Sprite component for rendering
type Sprite struct {
	Texture *sdl.Texture
	Width   int32
	Height  int32
	Color   sdl.Color
}

// Health component for game objects
type Health struct {
	Current int
	Max     int
}

// Input component for player-controlled entities
type Input struct {
	MoveUp    bool
	MoveDown  bool
	MoveLeft  bool
	MoveRight bool
}

// RigidBody wraps a physics.Body
type RigidBody struct {
    Body *physics.Body
}