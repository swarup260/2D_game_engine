package geometry

import (
	"math"
)

// -----------------------------------------------------------------------------
// The Shape interface is the core of this unified approach.
// -----------------------------------------------------------------------------
type Shape interface {
	// Support returns the point on the shape that is furthest in a given direction.
	// This is the only function needed by the GJK algorithm.
	Support(direction Vector2D) Vector2D
}

// -----------------------------------------------------------------------------
// Circle implements the Shape interface.
// -----------------------------------------------------------------------------
type Circle struct {
	Center Vector2D
	Radius float64
}

// Support for a Circle is the center point plus the radius in the given direction.
func (c *Circle) Support(direction Vector2D) Vector2D {
	// Normalize the direction vector and scale it by the radius.
	return c.Center.Add(direction.Normalize().Multiply(c.Radius))
}

type Polygon struct {
	Vertices []Vector2D
	Position Vector2D
	Rotation float64 // Angle in degrees 📐
}

// Support for a Polygon is the vertex that has the maximum dot product with the direction.
func (p *Polygon) Support(direction Vector2D) Vector2D {
    maxDot := math.Inf(-1)
    var supportPoint Vector2D

    // Convert rotation from degrees to radians.
    rotationRad := p.Rotation * math.Pi / 180
    cosTheta := math.Cos(rotationRad)
    sinTheta := math.Sin(rotationRad)

    // Iterate through all vertices to find the one with the max dot product.
    for _, v := range p.Vertices {
        // Apply rotation to the vertex.
        rotatedX := v.X*cosTheta - v.Y*sinTheta
        rotatedY := v.X*sinTheta + v.Y*cosTheta
        rotatedV := Vector2D{X: rotatedX, Y: rotatedY}

        // Calculate the dot product of the rotated vertex with the direction.
        dotProduct := rotatedV.Dot(direction)

        // If this is the furthest point so far, store it.
        if dotProduct > maxDot {
            maxDot = dotProduct
            supportPoint = rotatedV
        }
    }

    // Add the polygon's position to get the absolute world coordinate.
    return Vector2D{
        X: supportPoint.X + p.Position.X,
        Y: supportPoint.Y + p.Position.Y,
    }
}

type Triangle struct {
	Vertices []Vector2D
	Position Vector2D
	Rotation float64 // Angle in degrees 📐
}


func (t *Triangle) Support(direction Vector2D) Vector2D {
	maxDot := -math.MaxFloat64
	var supportPoint Vector2D 

	// Convert rotation from degrees to radians.
	rotationRad := t.Rotation * math.Pi / 180
	cosTheta := math.Cos(rotationRad)
	sinTheta := math.Sin(rotationRad)

	// Iterate through all vertices to find the one with the max dot product.
	for _, v := range t.Vertices {
		// Apply rotation to the vertex first.
		rotatedX := v.X*cosTheta - v.Y*sinTheta
		rotatedY := v.X*sinTheta + v.Y*cosTheta
		rotatedV := Vector2D{X: rotatedX, Y: rotatedY}

		// Calculate the dot product of the rotated vertex with the direction.
		dotProduct := rotatedV.Dot(direction)

		// If this is the furthest point so far, store it.
		if dotProduct > maxDot {
			maxDot = dotProduct
			supportPoint = rotatedV
		}
	}

	// Add the triangle's position to get the absolute world coordinate.
	return Vector2D{
		X: supportPoint.X + t.Position.X,
		Y: supportPoint.Y + t.Position.Y,
	}
}

// GJK is a simplified version for demonstration. A full GJK would handle
// 3-point simplexes (triangles) to find if the origin is contained.
// For this example, we'll assume a collision is found after two points.
// A full implementation requires more complex geometry logic.
