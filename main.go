package main

import (
	"fmt"
	"os"

	"2d_game_engine/physics/geometry"
	render "2d_game_engine/renderer"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Initialize SDL2.
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}
	defer sdl.Quit()

	// Create a window and a renderer.
	window, renderer, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	defer window.Destroy()
	defer renderer.Destroy()

	fmt.Println("GJK + EPA Collision System Demonstration")
	fmt.Println("---------------------------------------")

	// Define two colliding polygons.

	polygonA := &geometry.Polygon{
		Vertices: []geometry.Vector2D{
			{X: 0, Y: 0},
			{X: 100, Y: 0},
			{X: 100, Y: 100},
			{X: 0, Y: 100},
		},
		Position: geometry.Vector2D{X: 350, Y: 250},
	}

	triangle := &geometry.Triangle{
		Vertices: []geometry.Vector2D{
			{X: 0, Y: -100},   // Top point
			{X: -100, Y: 100}, // Bottom-left point
			{X: 100, Y: 100},  // Bottom-right point
		},
		Position: geometry.Vector2D{X: 400, Y: 500},
	}

	running := true
	for running {
		// Handle events.
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Clear the renderer with a solid background color.
		renderer.SetDrawColor(30, 30, 30, 255) // Dark gray
		renderer.Clear()

		// Set the color for drawing the polygon.
		renderer.SetDrawColor(255, 255, 255, 255) // White
		render.DrawFilledPolygon(renderer, polygonA)

		renderer.SetDrawColor(255, 255, 255, 255) // White
		render.DrawFilledTriangle(renderer, triangle)

		renderer.SetDrawColor(255, 255, 255, 255) // White
		render.DrawCircle(renderer, geometry.Vector2D{400, 300}, 50)

		// Present the changes to the window.
		renderer.Present()

		// Short delay to avoid excessive CPU usage.
		sdl.Delay(16) // ~60 FPS
	}

	// polygonB := &geometry.Polygon{
	// 	Vertices: []geometry.Vector2D{
	// 		{X: 3, Y: 3},
	// 		{X: 8, Y: 3},
	// 		{X: 8, Y: 8},
	// 		{X: 3, Y: 8},
	// 	},
	// }

	// fmt.Println("Polygon A vertices:", polygonA.Vertices)
	// fmt.Println("Polygon A vertices support:", polygonA.Support(geometry.Vector2D{X: 1, Y: 0}))
	// fmt.Println("Polygon B vertices:", polygonB.Vertices)
	// fmt.Println("Polygon B vertices support:", polygonB.Support(geometry.Vector2D{X: 1, Y: 0}))

	// // Perform the GJK collision check.
	// result := collision.GJKDetectCollision(polygonA, polygonB)

	// if result.Collision {
	// 	fmt.Println("Collision detected by GJK!")
	// 	fmt.Printf("Initial simplex has %d points.\n", len(result.Simplex))

	// 	// If GJK detects a collision, use EPA to find the penetration vector.
	// 	// penetrationNormal, penetrationDepth := collision.EPA(result.Simplex, polygonA, polygonB)

	// 	// fmt.Printf("EPA calculated penetration:\n")
	// 	// fmt.Printf("  Normal: (%.2f, %.2f)\n", penetrationNormal.X, penetrationNormal.Y)
	// 	// fmt.Printf("  Depth: %.2f\n", penetrationDepth)
	// } else {
	// 	fmt.Println("No collision detected.")
	// }
	return nil
}
