package components

import (
	EngineMath "game_engine/math"
	"math"
)

type RigidBody struct {
	Velocity        EngineMath.Vector
	AngularVelocity float64
	Mass            float64 // >0 for dynamic, 0 for static
	InverseMass     float64 // Computed as 1/Mass or 0
	Damping         float64 // Linear damping (0-1)
	AngularDamping  float64 // (0-1)
	GravityScale    float64 // Multiplier for global gravity
	IsStatic        bool    // If true, immovable

	// Force accumulator
	Force  EngineMath.Vector
	Torque float64
}

// NewDynamicBody creates a movable rigid body affected by forces & gravity
func NewDynamicBody(mass float64) *RigidBody {
	if mass <= 0 {
		panic("Dynamic body must have mass > 0")
	}

	return &RigidBody{
		Velocity:        EngineMath.Vector{0, 0},
		AngularVelocity: 0,
		Mass:            mass,
		InverseMass:     1.0 / mass,
		Damping:         0.98,  // default linear damping
		AngularDamping:  0.98,  // default angular damping
		GravityScale:    1.0,   // affected by gravity
		IsStatic:        false, // movable
	}
}

// NewStaticBody creates an immovable rigid body
func NewStaticBody() *RigidBody {
	return &RigidBody{
		Velocity:        EngineMath.Vector{0, 0},
		AngularVelocity: 0,
		Mass:            0,
		InverseMass:     0,
		Damping:         1.0,  // no decay
		AngularDamping:  1.0,  // no decay
		GravityScale:    0.0,  // not affected by gravity
		IsStatic:        true, // immovable
	}
}

// NewKinematicBody creates a body that is moved manually (not by forces)
func NewKinematicBody() *RigidBody {
	return &RigidBody{
		Velocity:        EngineMath.Vector{0, 0},
		AngularVelocity: 0,
		Mass:            math.Inf(1), // treat as infinite mass
		InverseMass:     0,           // no response to forces
		Damping:         1.0,         // doesn’t auto-damp
		AngularDamping:  1.0,
		GravityScale:    0.0,   // not affected by gravity
		IsStatic:        false, // it moves, but only by code
	}
}

type Collider struct {
	Shape       string            // "circle" or "box"
	Radius      float64           // For circle
	HalfSize    EngineMath.Vector // For box (half-width, half-height)
	Offset      EngineMath.Vector // Local offset from transform
	Friction    float64           // 0-1
	Restitution float64           // Bounciness 0-1
	IsSensor    bool              // Triggers only, no response
}

// Optional: For forces applied this frame
type ForceAccumulator struct {
	Forces  []EngineMath.Vector
	Torques []float64
}
