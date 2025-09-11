package ecs

import (
	"reflect"

	"game_engine/core"
	"game_engine/physics"

	"github.com/veandco/go-sdl2/sdl"
)

// MovementSystem handles entity movement based on velocity
type MovementSystem struct{}

func (ms *MovementSystem) GetRequiredComponents() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&Transform{}),
		reflect.TypeOf(&Velocity{}),
	}
}

func (ms *MovementSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
	for _, entity := range entities {
		transform, _ := manager.GetComponent(entity, reflect.TypeOf(&Transform{}))
		velocity, _ := manager.GetComponent(entity, reflect.TypeOf(&Velocity{}))

		t := transform.(*Transform)
		v := velocity.(*Velocity)

		// Update position based on velocity
		t.X += v.X * dt
		t.Y += v.Y * dt

		// Keep entities within screen bounds (800x600)
		if t.X < 0 {
			t.X = 0
			v.X = 0
		} else if t.X > 750 {
			t.X = 750
			v.X = 0
		}

		if t.Y < 0 {
			t.Y = 0
			v.Y = 0
		} else if t.Y > 550 {
			t.Y = 550
			v.Y = 0
		}
	}
}



// SpriteRenderSystem renders sprites
type SpriteRenderSystem struct{}

func (srs *SpriteRenderSystem) GetRequiredComponents() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&Transform{}),
		reflect.TypeOf(&Sprite{}),
	}
}

func (srs *SpriteRenderSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
	// This system doesn't need to update anything in the Update phase
}

func (srs *SpriteRenderSystem) Render(renderer *sdl.Renderer, entities []core.Entity, manager *core.ECSManager) {
	for _, entity := range entities {
		transform, _ := manager.GetComponent(entity, reflect.TypeOf(&Transform{}))
		sprite, _ := manager.GetComponent(entity, reflect.TypeOf(&Sprite{}))

		t := transform.(*Transform)
		s := sprite.(*Sprite)

		// Create destination rectangle
		dstRect := &sdl.Rect{
			X: int32(t.X),
			Y: int32(t.Y),
			W: s.Width,
			H: s.Height,
		}

		if s.Texture != nil {
			renderer.Copy(s.Texture, nil, dstRect)
		} else {
			// Render as colored rectangle if no texture
			renderer.SetDrawColor(s.Color.R, s.Color.G, s.Color.B, s.Color.A)
			renderer.FillRect(dstRect)
		}
	}
}

// Physics System handles entity physics
type PhysicsSystem struct {
	World *physics.World
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{
		World: physics.NewWorld(),
	}
}

func (ps *PhysicsSystem) AddRigidBody(entity core.Entity, body *physics.Body, manager *core.ECSManager) {
	manager.AddComponent(entity, &RigidBody{Body: body})
	ps.World.AddBody(body)
}

func (ps *PhysicsSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
	for _, entity := range entities {
		rigidBody, _ := manager.GetComponent(entity, reflect.TypeOf(&RigidBody{}))
		transform, _ := manager.GetComponent(entity, reflect.TypeOf(&Transform{}))

		body := rigidBody.(*RigidBody)
		t := transform.(*Transform)

		t.X, t.Y = body.Body.Position.X, body.Body.Position.Y

	}
}

func (ps *PhysicsSystem) UpdatePhysics(dt float64) {
	ps.World.Step(dt)
}

func (ps *PhysicsSystem) GetRequiredComponents() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(&RigidBody{})}
}
