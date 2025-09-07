package simulation

import (
	"fmt"
	"game_engine/core"
	"game_engine/ecs"
	"game_engine/physics"

	"github.com/veandco/go-sdl2/sdl"
)

type GravityBallScene struct {
	ECS           *core.ECSManager
	renderer      *core.Renderer
	physicsSystem *ecs.PhysicsSystem
	entities      []core.Entity
}

func NewGravityBallScene(manager *core.ECSManager, renderer *core.Renderer) *GravityBallScene {
	return &GravityBallScene{ECS: manager, renderer: renderer, physicsSystem: ecs.NewPhysicsSystem()}
}


func (gs *GravityBallScene) Init() error {

	gs.ECS.AddSystem(gs.physicsSystem)
	gs.ECS.AddSystem(&ecs.SpriteRenderSystem{})

	for i := 0; i < 10; i++ {
		// Ball (dynamic circle)
		ballEntity := gs.ECS.CreateEntity()
		ball := physics.NewDynamic(physics.Vector2D{X: 400, Y: 100},
			physics.Collider{Type: physics.ShapeCircle, Circle: physics.Circle{Radius: 16}},
			1.0)
		gs.physicsSystem.AddRigidBody(ballEntity, ball, gs.ECS)
		gs.ECS.AddComponent(ballEntity, &ecs.Transform{
			X: float64(400 + (10 * i)), Y: float64(100 + (10 * i)), Rotation: 0, ScaleX: 1, ScaleY: 1,
		})

		gs.ECS.AddComponent(ballEntity, &ecs.Sprite{
			Texture: nil, // Using colored rectangle
			Width:   50,
			Height:  50,
			Color:   sdl.Color{R: 255, G: 0, B: 0, A: 255}, // Green player
		})

		gs.entities = append(gs.entities, ballEntity)
	}

	// Ground (static AABB)
	groundEntity := gs.ECS.CreateEntity()
	ground := physics.NewStatic(physics.Vector2D{X: 0, Y: 670},
		physics.Collider{Type: physics.ShapeAABB, AABB: physics.AABB{HalfW: 1280, HalfH: 50}},
	)
	gs.physicsSystem.AddRigidBody(groundEntity, ground, gs.ECS)
	gs.ECS.AddComponent(groundEntity, &ecs.Transform{})
	gs.ECS.AddComponent(groundEntity, &ecs.Sprite{
		Texture: nil, // Using colored rectangle
		Width:   1280,
		Height:  50,
		Color:   sdl.Color{R: 0, G: 255, B: 0, A: 255}, // Green player
	})

	gs.entities = append(gs.entities, groundEntity)

	fmt.Printf("Scene initialized with %d entities\n", gs.ECS.GetEntityCount())

	return nil
}

func (gs *GravityBallScene) HandleInput(im core.InputManager) {}

func (gs *GravityBallScene) Update(dt float64) {
	// Fixed updates are handled in GameEngine loop
	gs.ECS.UpdateSystems(dt)
}

func (gs *GravityBallScene) UpdatePhysics(dt float64) {
	// Fixed updates are handled in GameEngine loop
	gs.physicsSystem.UpdatePhysics(dt)
}

func (gs *GravityBallScene) Render(alpha float64) {
	// Render all entities through the ECS render systems
	gs.ECS.RenderSystems(gs.renderer.GetRenderer())
}

func (gs *GravityBallScene) Cleanup() {
	fmt.Println("Cleaning up Game Scene...")

	// Destroy all entities
	for _, entity := range gs.entities {
		gs.ECS.DestroyEntity(entity)
	}

	// Clear collections
	gs.entities = gs.entities[:0]

	fmt.Println("Game Scene cleaned up")
}
