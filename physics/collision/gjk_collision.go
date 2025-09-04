package collision

import (
	"2d_game_engine/physics/geometry"
)


type Vector2D = geometry.Vector2D
type Shape = geometry.Shape

// SupportMinkowskiDifference finds the support point for the Minkowski difference
// of two shapes in a given direction.
func SupportMinkowskiDifference(shapeA, shapeB Shape, direction Vector2D) Vector2D {
	// Support(A-B) = Support(A) - Support(-B)
	supportA := shapeA.Support(direction)
	supportB := shapeB.Support(direction.Multiply(-1))
	return supportA.Subtract(supportB)
}

// GJKResult stores the result of the GJK collision check.
type GJKResult struct {
	Collision bool
	Simplex   []geometry.Vector2D
}


// GJKDetectCollision checks for collision and returns the simplex if one is found.
func GJKDetectCollision(shapeA, shapeB Shape) GJKResult {
	// Initial search direction: vector from B to A's centers.
	direction := shapeA.Support(Vector2D{X: 1, Y: 0}).Subtract(shapeB.Support(Vector2D{X: -1, Y: 0}))
	if direction.Length() == 0 {
		direction.X = 1
	}

	simplex := make([]Vector2D, 0, 3)
	simplex = append(simplex, SupportMinkowskiDifference(shapeA, shapeB, direction))
	direction = simplex[0].Multiply(-1)

	for {
		newPoint := SupportMinkowskiDifference(shapeA, shapeB, direction)
		if newPoint.Dot(direction) <= 0 {
			// Origin is not in the Minkowski difference. No collision.
			return GJKResult{Collision: false}
		}

		simplex = append(simplex, newPoint)

		if containsOrigin(simplex, &direction) {
			// Origin is inside the simplex. Collision detected.
			return GJKResult{Collision: true, Simplex: simplex}
		}
	}
}


// containsOrigin updates the simplex and returns true if the origin is contained.
// It also updates the search direction for the next iteration.
func containsOrigin(simplex []Vector2D, direction *Vector2D) bool {
	// Get the last added point (A) and the point before it (B).
	A := simplex[len(simplex)-1]
	B := simplex[len(simplex)-2]

	// The vector from B to A.
	AB := B.Subtract(A)
	// The vector from A to the origin.
	AO := A.Multiply(-1)

	// Determine if the origin is on the same side as the perp vector.
	// This finds the closest edge to the origin.
	if AB.Perp().Dot(AO) < 0 {
		// Origin is not on the side of the perp vector.
		// The new search direction is towards the origin from the last added point.
		*direction = AB.Perp().Multiply(-1)
	} else {
		*direction = AB.Perp()
	}

	return false
}