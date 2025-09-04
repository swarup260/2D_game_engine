package collision

import (
	"math"
)

//-----------------------------------------------------------------------------
// EPA Algorithm: Collision Resolution
//-----------------------------------------------------------------------------

// Edge represents an edge of the simplex with a distance to the origin.
type Edge struct {
	Distance float64
	Normal   Vector2D
	Index    int
}

// EPA finds the penetration vector and depth.
func EPA(simplex []Vector2D, shapeA, shapeB Shape) (Vector2D, float64) {
	const maxIterations = 50
	for i := 0; i < maxIterations; i++ {
		// Find the edge closest to the origin.
		closestEdge := getClosestEdge(simplex)
		normal := closestEdge.Normal
		distance := closestEdge.Distance

		// Get a new support point in the direction of the closest edge's normal.
		support := SupportMinkowskiDifference(shapeA, shapeB, normal)
		supportDistance := support.Dot(normal)

		if math.Abs(supportDistance-distance) < 1e-6 {
			// If the new support point is not further from the origin than
			// the closest edge, we have found the minimum penetration.
			return normal, distance
		} else {
			// Otherwise, add the new point to the simplex to refine the shape.
			simplex = insertPoint(simplex, support, closestEdge.Index+1)
		}
	}

	// Fallback in case the algorithm doesn't converge.
	return Vector2D{}, 0.0
}

// getClosestEdge finds the edge of the simplex closest to the origin.
func getClosestEdge(simplex []Vector2D) Edge {
	closest := Edge{
		Distance: math.Inf(1),
		Normal:   Vector2D{},
		Index:    -1,
	}

	for i := 0; i < len(simplex); i++ {
		p1 := simplex[i]
		p2 := simplex[(i+1)%len(simplex)]

		edge := p2.Subtract(p1)
		// Vector from p1 to the origin.
		p1ToOrigin := p1.Multiply(-1)

		// The vector to the closest point on the line segment from the origin.
		// Projection of p1ToOrigin onto the edge vector.
		t := p1ToOrigin.Dot(edge) / edge.Dot(edge)

		if t < 0 {
			t = 0
		} else if t > 1 {
			t = 1
		}

		closestPoint := p1.Add(edge.Multiply(t))
		normal := closestPoint.Normalize()
		distance := closestPoint.Length()

		if distance < closest.Distance {
			closest.Distance = distance
			closest.Normal = normal
			closest.Index = i
		}
	}

	return closest
}

// insertPoint inserts a new point into the simplex at a specified index.
func insertPoint(simplex []Vector2D, point Vector2D, index int) []Vector2D {
	newSimplex := make([]Vector2D, len(simplex)+1)
	copy(newSimplex[:index], simplex[:index])
	newSimplex[index] = point
	copy(newSimplex[index+1:], simplex[index:])
	return newSimplex
}