package body

import (
	"2d_game_engine/physics/geometry"
)

type Vector2D = geometry.Vector2D

type Body struct {
	id                int
	angle             float64 // Angle in degrees 📐
	vertices          []Vector2D
	position          Vector2D
	velocity          Vector2D
	acceleration      Vector2D
	force             Vector2D
	torque            float64
	positionImpulse   Vector2D
	constraintImpulse struct {
		linear  Vector2D
		angular float64
	}
	speed  float64
	motion float64
	render struct {
		visible bool
		opacity float64
		color   struct{ r, g, b, a uint8 }
		sprite  struct {
			xScale,
			yScale,
			xOffset,
			yOffset float64
		}
	}
	restitution     int32
	friction        float64
	frictionStatic  float64
	frictionAir     float64
	slop            float64
	totalContacts   int
	angularSpeed    float64
	angularVelocity float64
	isSensor        bool
	isStatic        bool
	isSleeping      bool
	mass            float64
	inertia         float64
	density         float64
	deltaTime       float64
	area            float64
	inverseInertia  float64
	inverseMass     float64
}

func NewBody() *Body {
	return &Body{
		id:              0,
		angle:           0,
		isSensor:        false,
		isStatic:        false,
		isSleeping:      false,
		area:            0,
		mass:            0,
		inertia:         0,
		deltaTime:       1000 / 60,
		friction:        0.1,
		frictionStatic:  0.5,
		frictionAir:     0.01,
		restitution:     0,
		slop:            0.05,
		totalContacts:   0,
		angularSpeed:    0,
		angularVelocity: 0,
		density:         0.001,
		vertices:        []Vector2D{},
		position:        Vector2D{X: 0, Y: 0},
		velocity:        Vector2D{X: 0, Y: 0},
		acceleration:    Vector2D{X: 0, Y: 0},
		force:           Vector2D{X: 0, Y: 0},
		torque:          0,
		positionImpulse: Vector2D{X: 0, Y: 0},
		speed:           0,
		constraintImpulse: struct {
			linear  Vector2D
			angular float64
		}{
			linear:  Vector2D{X: 0, Y: 0},
			angular: 0,
		},
		motion: 0,
		render: struct {
			visible bool
			opacity float64
			color   struct{ r, g, b, a uint8 }
			sprite  struct {
				xScale,
				yScale,
				xOffset,
				yOffset float64
			}
		}{
			visible: true,
			opacity: 1,
			color:   struct{ r, g, b, a uint8 }{r: 255, g: 255, b: 255, a: 255},
			sprite: struct {
				xScale,
				yScale,
				xOffset,
				yOffset float64
			}{
				xScale:  1,
				yScale:  1,
				xOffset: 0,
				yOffset: 0,
			},
		},
	}
}

// Getters

func (b *Body) GetID() int {
	return b.id
}

func (b *Body) GetAngle() float64 {
	return b.angle
}

func (b *Body) GetVertices() []Vector2D {
	return b.vertices
}

func (b *Body) GetPosition() Vector2D {
	return b.position
}

func (b *Body) GetVelocity() Vector2D {
	return b.velocity
}

func (b *Body) GetAcceleration() Vector2D {
	return b.acceleration
}

func (b *Body) GetForce() Vector2D {
	return b.force
}

func (b *Body) GetTorque() float64 {
	return b.torque
}

func (b *Body) GetPositionImpulse() Vector2D {
	return b.positionImpulse
}

func (b *Body) GetConstraintImpulse() struct {
	linear  Vector2D
	angular float64
} {
	return b.constraintImpulse
}

func (b *Body) GetSpeed() float64 {
	return b.speed
}

func (b *Body) GetMotion() float64 {
	return b.motion
}

func (b *Body) GetRender() struct {
	visible bool
	opacity float64
	color   struct{ r, g, b, a uint8 }
	sprite  struct {
		xScale,
		yScale,
		xOffset,
		yOffset float64
	}
} {
	return b.render
}

func (b *Body) GetRestitution() int32 {
	return b.restitution
}

func (b *Body) GetFriction() float64 {
	return b.friction
}

func (b *Body) GetFrictionStatic() float64 {
	return b.frictionStatic
}

func (b *Body) GetFrictionAir() float64 {
	return b.frictionAir
}

func (b *Body) GetSlop() float64 {
	return b.slop
}

func (b *Body) GetTotalContacts() int {
	return b.totalContacts
}

func (b *Body) GetAngularSpeed() float64 {
	return b.angularSpeed
}

func (b *Body) GetAngularVelocity() float64 {
	return b.angularVelocity
}

func (b *Body) GetIsSensor() bool {
	return b.isSensor
}

func (b *Body) GetIsStatic() bool {
	return b.isStatic
}

func (b *Body) GetIsSleeping() bool {
	return b.isSleeping
}

func (b *Body) GetMass() float64 {
	return b.mass
}

func (b *Body) GetInertia() float64 {
	return b.inertia
}

func (b *Body) GetDensity() float64 {
	return b.density
}

func (b *Body) GetDeltaTime() float64 {
	return b.deltaTime
}

func (b *Body) GetArea() float64 {
	return b.area
}

// Setters

func (b *Body) SetID(id int) {
	b.id = id
}

func (b *Body) SetAngle(angle float64) {
	b.angle = angle
}

func (b *Body) SetVertices(vertices []Vector2D) {
	b.vertices = vertices
}

func (b *Body) SetPosition(position Vector2D) {
	b.position = position
}

func (b *Body) SetVelocity(velocity Vector2D) {
	b.velocity = velocity
}

func (b *Body) SetAcceleration(acceleration Vector2D) {
	b.acceleration = acceleration
}

func (b *Body) SetForce(force Vector2D) {
	b.force = force
}

func (b *Body) SetTorque(torque float64) {
	b.torque = torque
}

func (b *Body) SetPositionImpulse(positionImpulse Vector2D) {
	b.positionImpulse = positionImpulse
}

func (b *Body) SetConstraintImpulse(linear Vector2D, angular float64) {
	b.constraintImpulse.linear = linear
	b.constraintImpulse.angular = angular
}

func (b *Body) SetSpeed(speed float64) {
	b.speed = speed
}

func (b *Body) SetMotion(motion float64) {
	b.motion = motion
}

func (b *Body) SetVisible(visible bool) {
	b.render.visible = visible
}

func (b *Body) SetOpacity(opacity float64) {
	b.render.opacity = opacity
}

func (b *Body) SetColor(red, green, blue, alpha uint8) {
	b.render.color.r = red
	b.render.color.g = green
	b.render.color.b = blue
	b.render.color.a = alpha
}

func (b *Body) SetSprite(xScale, yScale, xOffset, yOffset float64) {
	b.render.sprite.xScale = xScale
	b.render.sprite.yScale = yScale
	b.render.sprite.xOffset = xOffset
	b.render.sprite.yOffset = yOffset
}

func (b *Body) SetRestitution(restitution int32) {
	b.restitution = restitution
}

func (b *Body) SetFriction(friction float64) {
	b.friction = friction
}

func (b *Body) SetFrictionStatic(frictionStatic float64) {
	b.frictionStatic = frictionStatic
}

func (b *Body) SetFrictionAir(frictionAir float64) {
	b.frictionAir = frictionAir
}

func (b *Body) SetSlop(slop float64) {
	b.slop = slop
}

func (b *Body) SetTotalContacts(totalContacts int) {
	b.totalContacts = totalContacts
}

func (b *Body) SetAngularSpeed(angularSpeed float64) {
	b.angularSpeed = angularSpeed
}

func (b *Body) SetAngularVelocity(angularVelocity float64) {
	b.angularVelocity = angularVelocity
}

func (b *Body) SetIsSensor(isSensor bool) {
	b.isSensor = isSensor
}

func (b *Body) SetIsStatic(isStatic bool) {
	b.isStatic = isStatic
}

func (b *Body) SetIsSleeping(isSleeping bool) {
	b.isSleeping = isSleeping
}

func (b *Body) SetMass(mass float64) {
	var moment = b.inertia / (b.mass / 6)
	b.inertia = moment * (mass / 6)
	b.inverseInertia = 1 / b.inertia

	b.mass = mass
	b.inverseMass = 1 / b.mass
	b.density = b.mass / b.area
}

func (b *Body) SetInertia(inertia float64) {
	b.inertia = inertia
	b.inverseInertia = 1 / b.inertia
}

func (b *Body) SetDensity(density float64) {
	b.SetMass(density * b.area)
	b.density = density
}

func (b *Body) SetDeltaTime(deltaTime float64) {
	b.deltaTime = deltaTime
}

func (b *Body) SetArea(area float64) {
	b.area = area
}
