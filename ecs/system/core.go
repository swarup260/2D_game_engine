package system

import (
	"reflect"

	"game_engine/core"
	"game_engine/ecs/components"

	"github.com/veandco/go-sdl2/sdl"
)

// InputSystem handles player input
type InputSystem struct {
	inputManager *core.InputManager
}

func NewInputSystem(im *core.InputManager) *InputSystem {
	return &InputSystem{inputManager: im}
}

func (is *InputSystem) GetRequiredComponents() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&components.Input{}),
		reflect.TypeOf(&components.RigidBody{}),
	}
}

func (is *InputSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
	// Get current keyboard state
	keyState := sdl.GetKeyboardState()

	for _, entity := range entities {
		// SAFE WAY 1: Check exists before type assertion
		inputComp, _ := manager.GetComponent(entity, reflect.TypeOf(&components.Input{}))
		rb, _ := manager.GetComponent(entity, reflect.TypeOf(&components.RigidBody{}))

		if inputComp == nil {
			continue // Skip this entity if components are missing
		}

		// Now safe to type assert because we checked for nil
		input := inputComp.(*components.Input)
		body := rb.(*components.RigidBody)

		if body.IsStatic {
			continue // Skip if type assertion failed
		}

		// Update input state
		input.MoveUp = keyState[sdl.SCANCODE_W] != 0 || keyState[sdl.SCANCODE_UP] != 0
		input.MoveDown = keyState[sdl.SCANCODE_S] != 0 || keyState[sdl.SCANCODE_DOWN] != 0
		input.MoveLeft = keyState[sdl.SCANCODE_A] != 0 || keyState[sdl.SCANCODE_LEFT] != 0
		input.MoveRight = keyState[sdl.SCANCODE_D] != 0 || keyState[sdl.SCANCODE_RIGHT] != 0

		// Reset velocity
		body.Velocity.X = 0
		body.Velocity.Y = 0

		// Apply movement based on input
		speed := 300.0
		if input.MoveLeft {
			body.Velocity.X -= speed
		}
		if input.MoveRight {
			body.Velocity.X += speed
		}
		if input.MoveUp {
			body.Velocity.Y -= speed
		}
		if input.MoveDown {
			body.Velocity.Y += speed
		}

		// Normalize diagonal movement
		if body.Velocity.X != 0 && body.Velocity.Y != 0 {
			body.Velocity.X *= 0.707
			body.Velocity.Y *= 0.707
		}
	}
}
