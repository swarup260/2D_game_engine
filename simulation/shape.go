package simulation

import (
	"fmt"

	"game_engine/core"
	"game_engine/ecs"
	"game_engine/physics"

	"github.com/veandco/go-sdl2/sdl"
)

type ShapeScene struct {
	ECS      *core.ECSManager
	renderer *core.Renderer
	entities []core.Entity
}

func NewShapeScene(manager *core.ECSManager, renderer *core.Renderer) *ShapeScene {
	return &ShapeScene{ECS: manager, renderer: renderer}
}

func (gs *ShapeScene) Init() error {

	gs.ECS.AddSystem(&ecs.ShapeRenderSystem{})
	gs.ECS.AddSystem(&ecs.MovementSystem{})

	circle := gs.ECS.CreateEntity()
	s := ecs.Circle{Center: physics.Vector2D{X: 100, Y: 200}, Radius: 30}
	gs.ECS.AddComponent(circle, &ecs.Shape{
		Fill:  false,
		Shape: s,
		Color: sdl.Color{R: 0, G: 255, B: 0, A: 255}, // Green player
	})

	gs.ECS.AddComponent(circle, &ecs.Transform{
		X: 400, Y: 300, Rotation: 0, ScaleX: 1, ScaleY: 1,
	})
	gs.ECS.AddComponent(circle, &ecs.Velocity{X: 0, Y: 0})

	return nil
}

func (gs *ShapeScene) HandleInput(im core.InputManager) {}

func (gs *ShapeScene) Update(dt float64) {
	// Fixed updates are handled in GameEngine loop
	gs.ECS.UpdateSystems(dt)
}

func (gs *ShapeScene) UpdatePhysics(dt float64) {
	// Fixed updates are handled in GameEngine loop
}

func (gs *ShapeScene) Render(alpha float64) {
	// Render all entities through the ECS render systems
	gs.ECS.RenderSystems(gs.renderer.GetRenderer())
}

func (gs *ShapeScene) Cleanup() {
	fmt.Println("Cleaning up Game Scene...")

	// Destroy all entities
	for _, entity := range gs.entities {
		gs.ECS.DestroyEntity(entity)
	}

	// Clear collections
	gs.entities = gs.entities[:0]

	fmt.Println("Game Scene cleaned up")
}
