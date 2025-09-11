package simulation

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"game_engine/core"
	"game_engine/ecs/components"
	"game_engine/ecs/system"
	EngineMath "game_engine/math"

	"github.com/veandco/go-sdl2/sdl"
)

type GameEngineScene struct {
	GameEngine *core.GameEngine
	entities   []core.Entity
	player     core.Entity
}

func (gs *GameEngineScene) Init() error {

	// Register  systems
	gs.GameEngine.ECS.AddPhysicsSystem(&system.PhysicsSystem{})
	gs.GameEngine.ECS.AddGameplaySystem(system.NewInputSystem(gs.GameEngine.Input))
	gs.GameEngine.ECS.AddRenderSystem(&system.SpriteRenderSystem{})

	gs.player = gs.GameEngine.ECS.CreateEntity()
	gs.GameEngine.ECS.AddComponent(gs.player, &components.Transform{Position: EngineMath.Vector{X: 100, Y: 100}})
	gs.GameEngine.ECS.AddComponent(gs.player, components.NewDynamicBody(1.0)) // 1kg player
	gs.GameEngine.ECS.AddComponent(gs.player, &components.Input{})
	gs.GameEngine.ECS.AddComponent(gs.player, &components.Sprite{Texture: nil, Width: 32, Height: 32, Color: sdl.Color{R: 255, G: 0, B: 0, A: 255}})

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := 0; i < 5; i++ {

		min := 50
		max := 300

		X := float64(r.Intn(int(gs.GameEngine.WINDOW_WIDTH)))
		Y := float64(r.Intn(int(gs.GameEngine.WINDOW_HEIGHT)))

		// c := uint8(r.Intn(225))
		Width := r.Int31n(int32(max-min)) + int32(min)

		ground := gs.GameEngine.ECS.CreateEntity()
		gs.GameEngine.ECS.AddComponent(ground, &components.Transform{Position: EngineMath.Vector{X, Y}})
		gs.GameEngine.ECS.AddComponent(ground, components.NewStaticBody()) // immovable floor
		gs.GameEngine.ECS.AddComponent(ground, &components.Sprite{Texture: nil, Width: Width, Height: 30, Color: sdl.Color{R: 0, G: 225, B: 0, A: 255}})

		gs.entities = append(gs.entities, ground)
	}

	fmt.Printf("Scene initialized with %d entities\n", gs.GameEngine.ECS.GetEntityCount())

	return nil
}

func (gs *GameEngineScene) HandleInput(im core.InputManager) {
	gs.GameEngine.Input.Update()

	if playerInputComp, ok := gs.GameEngine.ECS.GetComponent(gs.player, reflect.TypeOf(&components.Input{})); ok {
		input := playerInputComp.(*components.Input)
		input.MoveUp = gs.GameEngine.Input.IsKeyPressed(sdl.SCANCODE_W) || gs.GameEngine.Input.IsKeyPressed(sdl.SCANCODE_UP)
		input.MoveDown = gs.GameEngine.Input.IsKeyPressed(sdl.SCANCODE_S) || gs.GameEngine.Input.IsKeyPressed(sdl.SCANCODE_DOWN)
		input.MoveLeft = gs.GameEngine.Input.IsKeyPressed(sdl.SCANCODE_A) || gs.GameEngine.Input.IsKeyPressed(sdl.SCANCODE_LEFT)
		input.MoveRight = gs.GameEngine.Input.IsKeyPressed(sdl.SCANCODE_D) || gs.GameEngine.Input.IsKeyPressed(sdl.SCANCODE_RIGHT)
	}

}

func (gs *GameEngineScene) Update(dt float64) {
	gs.GameEngine.ECS.UpdateGameplay(dt)
}

func (gs *GameEngineScene) UpdatePhysics(dt float64) {
	gs.GameEngine.ECS.UpdatePhysics(dt)
}

func (gs *GameEngineScene) Render(alpha float64) {
	renderer := gs.GameEngine.Render.GetRenderer()
	gs.GameEngine.ECS.RenderSystems(renderer, alpha)
}

func (gs *GameEngineScene) Cleanup() {
	fmt.Println("Cleaning up Game Scene...")

	// Destroy all entities
	gs.GameEngine.ECS.DestroyEntity(gs.player)
	for _, entity := range gs.entities {
		gs.GameEngine.ECS.DestroyEntity(entity)
	}

	// Clear collections
	gs.entities = gs.entities[:0]

	fmt.Println("Game Scene cleaned up")
}
