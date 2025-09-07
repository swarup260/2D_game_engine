package core

import (
	"reflect"
)

type Transform struct {
	X, Y           float64
	ScaleX, ScaleY float64
}

type Velocity struct {
	X, Y float64
}

type RigidBody struct {
	Mass     float64
	Friction float64
	Restitution float64
}

// MovementSystem handles entity movement based on velocity

type MovementSystem struct { }


func (s *MovementSystem) Update(dt float64, entities []Entity, manager *ECSManager) {
	for _, entity := range entities {
		transform, _ := manager.GetComponent(entity, reflect.TypeOf(&Transform{}))
		velocity, _ := manager.GetComponent(entity, reflect.TypeOf(&Velocity{}))

		t := transform.(*Transform)
		v := velocity.(*Velocity)

		t.X += v.X * dt
		t.Y += v.Y * dt
	}
}

func (s *MovementSystem) GetRequiredComponents() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(&Transform{}), reflect.TypeOf(&Velocity{})}
}


// PhysicsSystem handles physics simulation
type PhysicsSystem struct {
	Gravity float64
}

func (s *PhysicsSystem) Update(dt float64, entities []Entity, manager *ECSManager) {
	for _, entity := range entities {
		velocity, _ := manager.GetComponent(entity, reflect.TypeOf(&Velocity{}))
		rigidbody, _ := manager.GetComponent(entity, reflect.TypeOf(&RigidBody{}))
		
		v := velocity.(*Velocity)
		rb := rigidbody.(*RigidBody)
		
		// Apply gravity
		v.Y += s.Gravity * dt
		
		// Apply friction
		v.X *= (1.0 - rb.Friction*dt)
		v.Y *= (1.0 - rb.Friction*dt)
	}
}

func (s *PhysicsSystem) GetRequiredComponents() []reflect.Type{
	return []reflect.Type{reflect.TypeOf(&Velocity{}), reflect.TypeOf(&RigidBody{})}
}