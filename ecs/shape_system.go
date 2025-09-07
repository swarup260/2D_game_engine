package ecs

import (
	"math"
	"reflect"

	"game_engine/core"
	"game_engine/physics"

	"github.com/veandco/go-sdl2/sdl"
)

type DrawShape interface {
	Draw(fill bool, renderer *sdl.Renderer)
}

type Circle struct {
	Center physics.Vector2D
	Radius float64
}

func (c Circle) Draw(fill bool, renderer *sdl.Renderer) {
	if fill {
		drawFilledCircle(renderer, c.Center, c.Radius)
	} else {
		drawCircle(renderer, c.Center, c.Radius)

	}
}

func drawCircle(renderer *sdl.Renderer, center physics.Vector2D, radius float64) {
	x := int32(radius)
	y := int32(0)
	err := int32(0)

	centerX := int32(center.X)
	centerY := int32(center.Y)

	for x >= y {
		renderer.DrawPoint(centerX+x, centerY+y)
		renderer.DrawPoint(centerX+y, centerY+x)
		renderer.DrawPoint(centerX-y, centerY+x)
		renderer.DrawPoint(centerX-x, centerY+y)
		renderer.DrawPoint(centerX-x, centerY-y)
		renderer.DrawPoint(centerX-y, centerY-x)
		renderer.DrawPoint(centerX+y, centerY-x)
		renderer.DrawPoint(centerX+x, centerY-y)

		if err <= 0 {
			y += 1
			err += 2*y + 1
		}
		if err > 0 {
			x -= 1
			err -= 2*x + 1
		}
	}
}

func drawFilledCircle(renderer *sdl.Renderer, center physics.Vector2D, radius float64) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				renderer.DrawPoint(int32(center.X+x), int32(center.Y+y))
			}
		}
	}
}

type Polygon struct {
	Vertices []physics.Vector2D
	Position physics.Vector2D
}

func (p Polygon) Draw(fill bool, renderer *sdl.Renderer) {
	if fill {
		drawFilledPolygon(renderer, &p)
	} else {
		drawPolygon(renderer, &p)

	}
}

// drawPolygon renders the polygon by drawing lines between its vertices.
// It translates the vertex coordinates based on the polygon's position.
func drawPolygon(renderer *sdl.Renderer, p *Polygon) {
	// Check for a valid polygon with at least 3 vertices.
	if len(p.Vertices) < 3 {
		return
	}

	// Iterate through the vertices to draw the edges.
	for i := 0; i < len(p.Vertices); i++ {
		p1 := p.Vertices[i]
		// The next vertex, wrapping around to the first for the last edge.
		p2 := p.Vertices[(i+1)%len(p.Vertices)]

		// Translate the vertex coordinates by the polygon's position.
		x1 := int32(p1.X + p.Position.X)
		y1 := int32(p1.Y + p.Position.Y)
		x2 := int32(p2.X + p.Position.X)
		y2 := int32(p2.Y + p.Position.Y)

		// Draw a line between the two vertices.
		renderer.DrawLine(x1, y1, x2, y2)
	}
}

// DrawFilledPolygon draws a filled polygon using a scanline fill algorithm.
func drawFilledPolygon(renderer *sdl.Renderer, p *Polygon) {
	if len(p.Vertices) < 3 {
		return
	}

	// 1. Find the top and bottom bounds of the polygon.
	minY := math.Inf(1)
	maxY := math.Inf(-1)

	// Create a slice of absolute vertices.
	absoluteVertices := make([]physics.Vector2D, len(p.Vertices))
	for i, v := range p.Vertices {
		absoluteVertices[i] = physics.Vector2D{X: v.X + p.Position.X, Y: v.Y + p.Position.Y}
		if absoluteVertices[i].Y < minY {
			minY = absoluteVertices[i].Y
		}
		if absoluteVertices[i].Y > maxY {
			maxY = absoluteVertices[i].Y
		}
	}

	// 2. Iterate through each horizontal scanline from top to bottom.
	for y := int32(minY); y <= int32(maxY); y++ {
		var nodes []float64

		// 3. Find the intersection points of the scanline with the polygon's edges.
		for i := 0; i < len(absoluteVertices); i++ {
			p1 := absoluteVertices[i]
			p2 := absoluteVertices[(i+1)%len(absoluteVertices)]

			// Check if the scanline intersects the edge.
			if (p1.Y < float64(y) && p2.Y >= float64(y)) || (p2.Y < float64(y) && p1.Y >= float64(y)) {
				// Calculate the x-coordinate of the intersection.
				intersectX := (float64(y)-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y) + p1.X
				nodes = append(nodes, intersectX)
			}
		}

		// Sort the intersection points.
		if len(nodes) > 1 {
			// A simple bubble sort for a small number of nodes.
			// A more efficient sort could be used for complex polygons.
			for i := 0; i < len(nodes)-1; i++ {
				for j := i + 1; j < len(nodes); j++ {
					if nodes[i] > nodes[j] {
						nodes[i], nodes[j] = nodes[j], nodes[i]
					}
				}
			}

			// 4. Draw a horizontal line between each pair of intersection points.
			for i := 0; i < len(nodes); i += 2 {
				if i+1 < len(nodes) {
					renderer.DrawLine(int32(nodes[i]), y, int32(nodes[i+1]), y)
				}
			}
		}
	}
}

type Triangle struct {
	Vertices []physics.Vector2D
	Position physics.Vector2D
}

func (t Triangle) Draw(fill bool, renderer *sdl.Renderer) {
	if fill {
		drawFilledTriangle(renderer, &t)
	} else {
		drawTriangle(renderer, &t)

	}
}

func drawTriangle(renderer *sdl.Renderer, t *Triangle) {
	if len(t.Vertices) != 3 {
		return
	}

	points := make([]sdl.Point, 4)
	for i := 0; i < 3; i++ {
		// Calculate the absolute position of each vertex by adding the triangle's position.
		absoluteX := t.Vertices[i].X + t.Position.X
		absoluteY := t.Vertices[i].Y + t.Position.Y
		points[i] = sdl.Point{X: int32(absoluteX), Y: int32(absoluteY)}
	}
	points[3] = points[0] // Close the triangle.

	renderer.DrawLines(points)
}

// DrawFilledTriangle draws a filled triangle using a scanline fill algorithm.
func drawFilledTriangle(renderer *sdl.Renderer, t *Triangle) {
	if len(t.Vertices) != 3 {
		return
	}

	// 1. Get absolute vertices and sort them by Y-coordinate.
	absVertices := make([]physics.Vector2D, 3)
	for i, v := range t.Vertices {
		absVertices[i] = physics.Vector2D{X: v.X + t.Position.X, Y: v.Y + t.Position.Y}
	}

	// Sort vertices by Y-coordinate (top to bottom).
	for i := 0; i < 2; i++ {
		for j := i + 1; j < 3; j++ {
			if absVertices[i].Y > absVertices[j].Y {
				absVertices[i], absVertices[j] = absVertices[j], absVertices[i]
			}
		}
	}

	v1 := absVertices[0]
	v2 := absVertices[1]
	v3 := absVertices[2]

	// 2. Iterate through scanlines to fill the top and bottom parts of the triangle.

	// Top half of the triangle (from v1 to v2).
	for y := int32(v1.Y); y <= int32(v2.Y); y++ {
		if v2.Y != v1.Y && v3.Y != v1.Y {
			// Find x-coordinates on the two main edges.
			x1 := v1.X + (float64(y)-v1.Y)*(v2.X-v1.X)/(v2.Y-v1.Y)
			x2 := v1.X + (float64(y)-v1.Y)*(v3.X-v1.X)/(v3.Y-v1.Y)

			// Draw a horizontal line between them.
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			renderer.DrawLine(int32(x1), y, int32(x2), y)
		}
	}

	// Bottom half of the triangle (from v2 to v3).
	for y := int32(v2.Y); y <= int32(v3.Y); y++ {
		if v3.Y != v2.Y && v3.Y != v1.Y {
			// Find x-coordinates on the two main edges.
			x1 := v2.X + (float64(y)-v2.Y)*(v3.X-v2.X)/(v3.Y-v2.Y)
			x2 := v1.X + (float64(y)-v1.Y)*(v3.X-v1.X)/(v3.Y-v1.Y)

			// Draw a horizontal line between them.
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			renderer.DrawLine(int32(x1), y, int32(x2), y)
		}
	}
}

type Shape struct {
	Fill  bool
	Shape DrawShape
	Color sdl.Color
}

type ShapeRenderSystem struct{}

func (srs *ShapeRenderSystem) GetRequiredComponents() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&Shape{}),
	}
}

func (srs *ShapeRenderSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
	// This system doesn't need to update anything in the Update phase
}

func (srs *ShapeRenderSystem) Render(renderer *sdl.Renderer, entities []core.Entity, manager *core.ECSManager) {
	for _, entity := range entities {
		shape, _ := manager.GetComponent(entity, reflect.TypeOf(&Shape{}))

		s := shape.(*Shape)
		renderer.SetDrawColor(s.Color.R, s.Color.G, s.Color.B, s.Color.A)
		s.Shape.Draw(s.Fill, renderer)
	}
}
