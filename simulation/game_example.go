package simulation

import (
	"fmt"
	"reflect"

	"game_engine/core"
	"game_engine/ecs"

	"github.com/veandco/go-sdl2/sdl"
)

type GameScene struct {
	ecsManager   *core.ECSManager
	inputManager *core.InputManager
	renderer     *sdl.Renderer

	// Entities
	player  core.Entity
	enemies []core.Entity

	// Previous state for interpolation
	prevPositions map[core.Entity]ecs.Transform
}

func NewGameScene(ecsManager *core.ECSManager, inputManager *core.InputManager, renderer *sdl.Renderer) *GameScene {
	return &GameScene{
		ecsManager:    ecsManager,
		inputManager:  inputManager,
		renderer:      renderer,
		enemies:       make([]core.Entity, 0),
		prevPositions: make(map[core.Entity]ecs.Transform),
	}
}

func (gs *GameScene) Init() error {
	fmt.Println("Initializing Game Scene...")

	// Add systems to ECS manager
	gs.ecsManager.AddSystem(ecs.NewInputSystem(gs.inputManager))
	gs.ecsManager.AddSystem(&ecs.MovementSystem{})
	gs.ecsManager.AddSystem(&ecs.SpriteRenderSystem{})

	// Create player entity
	gs.player = gs.ecsManager.CreateEntity()

	// Add components to player
	gs.ecsManager.AddComponent(gs.player, &ecs.Transform{
		X: 400, Y: 300, Rotation: 0, ScaleX: 1, ScaleY: 1,
	})
	gs.ecsManager.AddComponent(gs.player, &ecs.Velocity{X: 0, Y: 0})
	gs.ecsManager.AddComponent(gs.player, &ecs.Sprite{
		Texture: nil, // Using colored rectangle
		Width:   50,
		Height:  50,
		Color:   sdl.Color{R: 0, G: 255, B: 0, A: 255}, // Green player
	})
	gs.ecsManager.AddComponent(gs.player, &ecs.Health{Current: 100, Max: 100})
	gs.ecsManager.AddComponent(gs.player, &ecs.Input{})

	// Create some enemy entities
	for i := 0; i < 3; i++ {
		enemy := gs.ecsManager.CreateEntity()
		gs.enemies = append(gs.enemies, enemy)

		gs.ecsManager.AddComponent(enemy, &ecs.Transform{
			X:        float64(100 + i*200),
			Y:        float64(100 + i*50),
			Rotation: 0,
			ScaleX:   1,
			ScaleY:   1,
		})
		gs.ecsManager.AddComponent(enemy, &ecs.Velocity{
			X: float64(50 + i*25),
			Y: float64(25 + i*10),
		})
		gs.ecsManager.AddComponent(enemy, &ecs.Sprite{
			Texture: nil,
			Width:   30,
			Height:  30,
			Color:   sdl.Color{R: 255, G: 0, B: 0, A: 255}, // Red enemies
		})
		gs.ecsManager.AddComponent(enemy, &ecs.Health{Current: 50, Max: 50})
	}

	fmt.Printf("Scene initialized with %d entities\n", gs.ecsManager.GetEntityCount())
	return nil
}

func (gs *GameScene) HandleInput(im core.InputManager) {
	// Update input manager
	gs.inputManager.Update()

	// Update player input component
	if playerInputComp, ok := gs.ecsManager.GetComponent(gs.player, reflect.TypeOf(&ecs.Input{})); ok {
		input := playerInputComp.(*ecs.Input)
		input.MoveUp = gs.inputManager.IsKeyPressed(sdl.SCANCODE_W) || gs.inputManager.IsKeyPressed(sdl.SCANCODE_UP)
		input.MoveDown = gs.inputManager.IsKeyPressed(sdl.SCANCODE_S) || gs.inputManager.IsKeyPressed(sdl.SCANCODE_DOWN)
		input.MoveLeft = gs.inputManager.IsKeyPressed(sdl.SCANCODE_A) || gs.inputManager.IsKeyPressed(sdl.SCANCODE_LEFT)
		input.MoveRight = gs.inputManager.IsKeyPressed(sdl.SCANCODE_D) || gs.inputManager.IsKeyPressed(sdl.SCANCODE_RIGHT)
	}
}

func (gs *GameScene) Update(dt float64) {
	// Store previous positions for interpolation
	gs.storePreviousPositions()

	// Update all systems
	gs.ecsManager.UpdateSystems(dt)

	// Simple enemy AI - bounce off walls
	for _, enemy := range gs.enemies {
		if transformComp, ok := gs.ecsManager.GetComponent(enemy, reflect.TypeOf(&ecs.Transform{})); ok {
			if velocityComp, ok := gs.ecsManager.GetComponent(enemy, reflect.TypeOf(&ecs.Velocity{})); ok {
				transform := transformComp.(*ecs.Transform)
				velocity := velocityComp.(*ecs.Velocity)

				// Bounce off walls
				if transform.X <= 0 || transform.X >= 770 {
					velocity.X = -velocity.X
				}
				if transform.Y <= 0 || transform.Y >= 570 {
					velocity.Y = -velocity.Y
				}
			}
		}
	}
}

func (gs *GameScene) UpdatePhysics(dt float64) {
	// Physics updates would go here (collision detection, physics simulation, etc.)
	// For this example, movement is handled in the regular Update
}

func (gs *GameScene) Render(alpha float64) {
	// Render all entities through the ECS render systems
	gs.ecsManager.RenderSystems(gs.renderer)

	// Render UI elements (could be separate UI system)
	gs.renderUI()
}

func (gs *GameScene) renderUI() {
	// Simple UI rendering - health bar for player
	if healthComp, ok := gs.ecsManager.GetComponent(gs.player, reflect.TypeOf(&ecs.Health{})); ok {
		health := healthComp.(*ecs.Health)

		// Health bar background
		bgRect := &sdl.Rect{X: 10, Y: 10, W: 200, H: 20}
		gs.renderer.SetDrawColor(100, 100, 100, 255)
		gs.renderer.FillRect(bgRect)

		// Health bar foreground
		healthWidth := int32(float64(health.Current) / float64(health.Max) * 200)
		healthRect := &sdl.Rect{X: 10, Y: 10, W: healthWidth, H: 20}
		gs.renderer.SetDrawColor(255, 0, 0, 255)
		gs.renderer.FillRect(healthRect)
	}

	// Entity count display (simple debug info)
	// In a real game, you'd use a proper font rendering system
}

func (gs *GameScene) storePreviousPositions() {
	// Store previous positions for smooth interpolation
	allEntities := append([]core.Entity{gs.player}, gs.enemies...)
	for _, entity := range allEntities {
		if transformComp, ok := gs.ecsManager.GetComponent(entity, reflect.TypeOf(&ecs.Transform{})); ok {
			transform := transformComp.(*ecs.Transform)
			gs.prevPositions[entity] = *transform // Copy the transform
		}
	}
}

func (gs *GameScene) Cleanup() {
	fmt.Println("Cleaning up Game Scene...")

	// Destroy all entities
	gs.ecsManager.DestroyEntity(gs.player)
	for _, enemy := range gs.enemies {
		gs.ecsManager.DestroyEntity(enemy)
	}

	// Clear collections
	gs.enemies = gs.enemies[:0]
	gs.prevPositions = make(map[core.Entity]ecs.Transform)

	fmt.Println("Game Scene cleaned up")
}
