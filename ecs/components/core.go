package components

import (
	EngineMath "game_engine/math"
)

type Transform struct {
	Position     EngineMath.Vector
	PrevPosition EngineMath.Vector
	Rotation     float64
}

// Input component for player-controlled entities
type Input struct {
	MoveUp    bool
	MoveDown  bool
	MoveLeft  bool
	MoveRight bool
}
