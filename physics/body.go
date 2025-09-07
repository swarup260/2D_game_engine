package physics

import (
	"math"
)

type BodyType int

const (
	BodyStatic BodyType = iota
	BodyKinematic
	BodyDynamic
)

type Body struct {
	// Pose
	Position Vector2D
	Velocity Vector2D
	Angle    float64
	AngVel   float64

	// Physical
	Type            BodyType
	Mass            float64
	InvMass         float64
	Restitution     float64 // bounciness [0..1]
	StaticFriction  float64
	DynamicFriction float64
	LinearDamping   float64
	// (No rotation support in AABB; keep AngVel 0 unless polygons are added)

	// Forces
	Force Vector2D

	// Collision
	Collider   Collider
	IsSleeping bool
	UserData   any
}

func NewDynamic(pos Vector2D, col Collider, density float64) *Body {
	area := float64(1.0)
	switch col.Type {
	case ShapeCircle:
		area = math.Pi * col.Circle.Radius * col.Circle.Radius
	case ShapeAABB:
		area = (col.AABB.HalfW * 2) * (col.AABB.HalfH * 2)
	}
	m := density * area
	if m <= 0 {
		m = 1
	}
	return &Body{
		Position: pos, Type: BodyDynamic, Mass: m, InvMass: 1.0 / m,
		Restitution: 0.2, StaticFriction: 0.5, DynamicFriction: 0.3,
		LinearDamping: 0.01, Collider: col,
	}
}

func NewStatic(pos Vector2D, col Collider) *Body {
	return &Body{
		Position: pos, Type: BodyStatic, Mass: 0, InvMass: 0,
		Restitution: 0.2, StaticFriction: 0.6, DynamicFriction: 0.4,
		Collider: col,
	}
}

type Contact struct {
	A, B        *Body
	Normal      Vector2D // from A to B
	Penetration float64
	Point       Vector2D // contact point (approx; not used heavily for AABB)
}
