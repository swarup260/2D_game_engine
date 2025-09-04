package geometry
// 2d_game_engine/physics/geometry

import (
	"math"
)

//-----------------------------------------------------------------------------
// Vector2D represents a 2D point or vector.
//-----------------------------------------------------------------------------
type Vector2D struct {
	X float64
	Y float64
}

// Add two vectors.
func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{X: v.X + other.X, Y: v.Y + other.Y}
}

// Subtracts one vector from another.
func (v Vector2D) Subtract(other Vector2D) Vector2D {
	return Vector2D{X: v.X - other.X, Y: v.Y - other.Y}
}

// Multiply a vector by a scalar.
func (v Vector2D) Multiply(scalar float64) Vector2D {
	return Vector2D{X: v.X * scalar, Y: v.Y * scalar}
}

// Divide a vector by a scalar.
func (v Vector2D) Divide(scalar float64) Vector2D {
	return Vector2D{X: v.X / scalar, Y: v.Y / scalar}
}

// Dot returns the dot product of two vectors.
func (v Vector2D) Dot(other Vector2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Cross returns the cross product of two vectors (in 2D, this is a scalar).
func (v Vector2D) Cross(other Vector2D) float64 {
	return v.X*other.Y - v.Y*other.X
}
// Cross3 returns the cross product of three vectors (in 2D, this is a scalar).
func (v Vector2D) Cross3(vectorA,vectorB,vectorC Vector2D) float64 {
	return (vectorB.X - vectorA.X) * (vectorC.Y - vectorA.Y) - (vectorB.Y - vectorA.Y) * (vectorC.X - vectorA.X)
}
// Length calculates the magnitude of the vector.
func (v Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// LengthSaqured calculates the squared magnitude of the vector.
func (v Vector2D) LengthSaqured() float64 {
	return (v.X*v.X) + (v.Y*v.Y)
}

// Rotate rotates the vector by a given angle in radians.
func (v Vector2D) Rotate (angle float64) Vector2D {
	cosTheta := math.Cos(angle)
	sinTheta := math.Sin(angle)
	return Vector2D{
		X: v.X*cosTheta - v.Y*sinTheta,
		Y: v.X*sinTheta + v.Y*cosTheta,
	}
}

// Normalize returns a unit vector (length 1) in the same direction.
func (v Vector2D) Normalize() Vector2D {
	length := v.Length()
	if length == 0 {
		return Vector2D{}
	}
	return Vector2D{X: v.X / length, Y: v.Y / length}
}

// Perp returns the perpendicular vector (90 degree rotation clockwise).
func (v Vector2D) Perp() Vector2D {
	return Vector2D{X: -v.Y, Y: v.X}
}


// Negate returns the negated vector.
func (v Vector2D) Negate() Vector2D {
	return Vector2D{X: -v.X, Y: -v.Y}
}

// Angle returns the angle in radians between two vectors.
func (v Vector2D) Angle(vectorA,vectorB  Vector2D) float64 {
	return math.Atan2(vectorB.Y-vectorA.Y, vectorB.X-vectorA.X)
}