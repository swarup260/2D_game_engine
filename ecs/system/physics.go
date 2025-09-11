package system

import (
	"math"
	"reflect"

	EngineMath "game_engine/math"

	"game_engine/core" // Your ECS manager package (e.g., with World, Query2)
	"game_engine/ecs/components"
)

type PhysicsSystem struct{}

func (s *PhysicsSystem) GetRequiredComponents() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&components.Transform{}),
		reflect.TypeOf(&components.RigidBody{}),
	}
}

func (s *PhysicsSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
	for _, e := range entities {
		rb, _ := manager.GetComponent(e, reflect.TypeOf(&components.RigidBody{}))
		t, _ := manager.GetComponent(e, reflect.TypeOf(&components.Transform{}))

		body := rb.(*components.RigidBody)
		transform := t.(*components.Transform)

		// Save previous state for interpolation
		transform.PrevPosition = transform.Position

		if body.IsStatic {
			continue
		}

		// --- Dynamic bodies ---
		if !body.IsStatic {
			// Acceleration = Force / Mass
			acceleration := body.Force.Mul(body.InverseMass)

			// Update velocity
			body.Velocity = body.Velocity.Add(acceleration.Mul(dt))

			// Apply damping
			body.Velocity = body.Velocity.Mul(1.0 - body.Damping*dt)

			// Integrate position
			transform.Position = transform.Position.Add(body.Velocity.Mul(dt))
		}

		// Clear forces after integration
		body.Force = EngineMath.Vector{}
		body.Torque = 0

		// --- Kinematic bodies ---
		// Forces are ignored, but velocity (set externally) is applied
		if math.IsInf(body.Mass, 1) {
			transform.Position = transform.Position.Add(body.Velocity.Mul(dt))
		}
	}
}

// var Gravity = gameEngineMath.Vector{0, 5}

// // ForceApplicationSystem: Applies gravity, damping, and accumulated forces.
// type ForceApplicationSystem struct{}

// func (s *ForceApplicationSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
// 	for _, entity := range entities {

// 		rigidBody, _ := manager.GetComponent(entity, reflect.TypeOf(&components.RigidBody{}))
// 		rb := rigidBody.(*components.RigidBody)

// 		if rb.IsStatic {
// 			return
// 		}

// 		forceAccumulator, _ := manager.GetComponent(entity, reflect.TypeOf(&components.ForceAccumulator{}))

// 		acc := forceAccumulator.(*components.ForceAccumulator)

// 		// Apply gravity
// 		gravityForce := Gravity.Mul(rb.Mass * rb.GravityScale)
// 		rb.Velocity = rb.Velocity.Add(gravityForce.Mul(dt * rb.InverseMass))

// 		// Apply accumulated forces
// 		for _, f := range acc.Forces {
// 			rb.Velocity = rb.Velocity.Add(f.Mul(dt * rb.InverseMass))
// 		}
// 		for _, t := range acc.Torques {
// 			rb.AngularVelocity += t * dt // Simplify, assume inertia=1 or compute properly
// 		}
// 		acc.Forces = acc.Forces[:0] // Clear
// 		acc.Torques = acc.Torques[:0]

// 		// Damping
// 		rb.Velocity = rb.Velocity.Mul(math.Pow(1-rb.Damping, dt))
// 		rb.AngularVelocity *= math.Pow(1-rb.AngularDamping, dt)

// 	}
// }

// func (s *ForceApplicationSystem) GetRequiredComponents() []reflect.Type {
// 	return []reflect.Type{
// 		reflect.TypeOf(&components.RigidBody{}),
// 		reflect.TypeOf(&components.ForceAccumulator{}),
// 	}
// }

// // IntegrationSystem: Updates positions/rotations from velocities.
// type IntegrationSystem struct{}

// func (s *IntegrationSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
// 	for _, entity := range entities {
// 		rigidBody, _ := manager.GetComponent(entity, reflect.TypeOf(&components.RigidBody{}))
// 		rb := rigidBody.(*components.RigidBody)

// 		if rb.IsStatic {
// 			return
// 		}

// 		transform, _ := manager.GetComponent(entity, reflect.TypeOf(&components.Transform{}))
// 		t := transform.(*components.Transform)

// 		t.Position = t.Position.Add(rb.Velocity.Mul(dt))
// 		t.Rotation += rb.AngularVelocity * dt
// 	}
// }

// func (s *IntegrationSystem) GetRequiredComponents() []reflect.Type {
// 	return []reflect.Type{
// 		reflect.TypeOf(&components.RigidBody{}),
// 		reflect.TypeOf(&components.Transform{}),
// 	}
// }

// type SyncPreviousSystem struct{}

// func (s *SyncPreviousSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
// 	for _, entity := range entities {
// 		transform, _ := manager.GetComponent(entity, reflect.TypeOf(&components.Transform{}))
// 		t := transform.(*components.Transform)
// 		// Assume entities with Transform also get PreviousTransform on spawn
// 		prevTransform, _ := manager.GetComponent(entity, reflect.TypeOf(&components.PreviousTransform{}))
// 		prev := prevTransform.(*components.PreviousTransform)

// 		prev.Position = t.Position
// 		prev.Rotation = t.Rotation
// 	}
// }

// func (a *SyncPreviousSystem) GetRequiredComponents() []reflect.Type {
// 	return []reflect.Type{
// 		reflect.TypeOf(&components.Transform{}),
// 		reflect.TypeOf(&components.PreviousTransform{}),
// 	}
// }

// // CollisionDetectionSystem: Detects contacts, stores as temporary data or events.
// type CollisionDetectionSystem struct {
// 	Contacts []Contact // Buffer for contacts, clear each frame
// }

// type Contact struct {
// 	EntityA, EntityB core.Entity
// 	Penetration      float64
// 	Normal           gameEngineMath.Vector
// 	Point            gameEngineMath.Vector
// }

// func (s *CollisionDetectionSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
// 	s.Contacts = s.Contacts[:0] // Clear

// 	entities = manager.GetEntitiesWithComponents([]reflect.Type{
// 		reflect.TypeOf(&components.Transform{}),
// 		reflect.TypeOf(&components.RigidBody{}),
// 		reflect.TypeOf(&components.Collider{}),
// 	})

// 	for i := 0; i < len(entities); i++ {
// 		for j := i + 1; j < len(entities); j++ {
// 			e1 := entities[i]
// 			e2 := entities[j]

// 			trans1Comp, _ := manager.GetComponent(e1, reflect.TypeOf(&components.Transform{}))
// 			rb1Comp, _ := manager.GetComponent(e1, reflect.TypeOf(&components.RigidBody{}))
// 			col1Comp, _ := manager.GetComponent(e1, reflect.TypeOf(&components.Collider{}))

// 			trans2Comp, _ := manager.GetComponent(e2, reflect.TypeOf(&components.Transform{}))
// 			rb2Comp, _ := manager.GetComponent(e2, reflect.TypeOf(&components.RigidBody{}))
// 			col2Comp, _ := manager.GetComponent(e2, reflect.TypeOf(&components.Collider{}))

// 			trans1 := trans1Comp.(*components.Transform)
// 			rb1 := rb1Comp.(*components.RigidBody)
// 			col1 := col1Comp.(*components.Collider)

// 			trans2 := trans2Comp.(*components.Transform)
// 			rb2 := rb2Comp.(*components.RigidBody)
// 			col2 := col2Comp.(*components.Collider)

// 			if rb1.IsStatic && rb2.IsStatic {
// 				continue
// 			}
// 			if col1.IsSensor || col2.IsSensor {
// 				// Handle triggers separately
// 				continue
// 			}

// 			// Compute world positions
// 			pos1 := trans1.Position.Add(col1.Offset)
// 			pos2 := trans2.Position.Add(col2.Offset)

// 			// Circle-Circle collision
// 			if col1.Shape == "circle" && col2.Shape == "circle" {
// 				dist := pos1.Sub(pos2).Length()
// 				sumRadius := col1.Radius + col2.Radius
// 				if dist < sumRadius {
// 					penetration := sumRadius - dist
// 					normal := pos2.Sub(pos1).Normalize()
// 					s.Contacts = append(s.Contacts, Contact{e1, e2, penetration, normal, pos1.Add(normal.Mul(col1.Radius))})
// 				}
// 				continue
// 			}

// 			// Box-Box (AABB, ignoring rotation for simplicity)
// 			if col1.Shape == "box" && col2.Shape == "box" {
// 				half1 := col1.HalfSize
// 				half2 := col2.HalfSize
// 				if math.Abs(pos1.X-pos2.X) < half1.X+half2.X && math.Abs(pos1.Y-pos2.Y) < half1.Y+half2.Y {
// 					// Overlap, compute penetration and normal (simplified)
// 					overlapX := half1.X + half2.X - math.Abs(pos1.X-pos2.X)
// 					overlapY := half1.Y + half2.Y - math.Abs(pos1.Y-pos2.Y)
// 					var penetration float64
// 					var normal gameEngineMath.Vector
// 					if overlapX < overlapY {
// 						penetration = overlapX
// 						normal = gameEngineMath.Vector{X: 1, Y: 0}
// 						if pos1.X > pos2.X {
// 							normal.X = -1
// 						}
// 					} else {
// 						penetration = overlapY
// 						normal = gameEngineMath.Vector{X: 0, Y: 1}
// 						if pos1.Y > pos2.Y {
// 							normal.Y = -1
// 						}
// 					}
// 					s.Contacts = append(s.Contacts, Contact{e1, e2, penetration, normal, gameEngineMath.Vector{}})
// 				}
// 			}
// 			// Add more shapes (e.g., circle-box) as needed
// 		}
// 	}

// }

// func (s *CollisionDetectionSystem) GetRequiredComponents() []reflect.Type {
// 	return []reflect.Type{
// 		reflect.TypeOf(&components.Transform{}),
// 		reflect.TypeOf(&components.RigidBody{}),
// 		reflect.TypeOf(&components.Collider{}),
// 	}
// }

// // CollisionResolutionSystem: Resolves contacts with impulses.
// type CollisionResolutionSystem struct {
// 	Detection *CollisionDetectionSystem // Reference to get contacts
// }

// func (s *CollisionResolutionSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
// 	for _, contact := range s.Detection.Contacts {

// 		// Get components
// 		trans1Comp, _ := manager.GetComponent(contact.EntityA, reflect.TypeOf(&components.Transform{}))
// 		rb1Comp, _ := manager.GetComponent(contact.EntityA, reflect.TypeOf(&components.RigidBody{}))
// 		col1Comp, _ := manager.GetComponent(contact.EntityA, reflect.TypeOf(&components.Collider{}))

// 		trans2Comp, _ := manager.GetComponent(contact.EntityB, reflect.TypeOf(&components.Transform{}))
// 		rb2Comp, _ := manager.GetComponent(contact.EntityB, reflect.TypeOf(&components.RigidBody{}))
// 		col2Comp, _ := manager.GetComponent(contact.EntityB, reflect.TypeOf(&components.Collider{}))

// 		trans1 := trans1Comp.(*components.Transform)
// 		rb1 := rb1Comp.(*components.RigidBody)
// 		col1 := col1Comp.(*components.Collider)

// 		trans2 := trans2Comp.(*components.Transform)
// 		rb2 := rb2Comp.(*components.RigidBody)
// 		col2 := col2Comp.(*components.Collider)

// 		// Relative velocity along normal
// 		relVel := rb1.Velocity.Sub(rb2.Velocity)
// 		velAlongNormal := relVel.X*contact.Normal.X + relVel.Y*contact.Normal.Y

// 		if velAlongNormal > 0 {
// 			continue // Separating
// 		}

// 		// Restitution (average)
// 		e := math.Min(col1.Restitution, col2.Restitution)

// 		// Impulse scalar
// 		j := -(1 + e) * velAlongNormal
// 		invMassSum := rb1.InverseMass + rb2.InverseMass
// 		if invMassSum == 0 {
// 			continue // Both infinite mass
// 		}
// 		j /= invMassSum

// 		impulse := contact.Normal.Mul(j)

// 		// Apply impulses
// 		rb1.Velocity = rb1.Velocity.Sub(impulse.Mul(rb1.InverseMass))
// 		rb2.Velocity = rb2.Velocity.Add(impulse.Mul(rb2.InverseMass))

// 		// Position correction (simple penetration resolve)
// 		correction := contact.Normal.Mul(contact.Penetration / invMassSum * 0.8) // Positional correction factor
// 		trans1.Position = trans1.Position.Sub(correction.Mul(rb1.InverseMass))
// 		trans2.Position = trans2.Position.Add(correction.Mul(rb2.InverseMass))

// 		// Friction (simplified static)
// 		// ... Add tangential impulse if needed
// 	}
// }

// func (s *CollisionResolutionSystem) GetRequiredComponents() []reflect.Type {
// 	return []reflect.Type{
// 		reflect.TypeOf(&components.Transform{}),
// 		reflect.TypeOf(&components.RigidBody{}),
// 		reflect.TypeOf(&components.Collider{}),
// 	}
// }
