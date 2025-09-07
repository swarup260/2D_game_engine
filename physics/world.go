package physics

import (
	"math"
)

// ---------------- World ----------------

type World struct {
	Bodies     []*Body
	Gravity    Vector2D
	Iterations int     // solver iterations
	cellSize   float64 // Broadphase spatial hash
	grid       map[int]map[int][]*Body
}

func NewWorld() *World {
	return &World{
		Gravity: Vector2D{0, 800}, Iterations: 8,
		cellSize: 128, grid: make(map[int]map[int][]*Body),
	}
}

func (w *World) AddBody(b *Body) { w.Bodies = append(w.Bodies, b) }

func (w *World) ClearForces() {
	for _, b := range w.Bodies {
		b.Force = Vector2D{}
	}
}

// ---------------- Integration ----------------

func (w *World) IntegrateForces(dt float64) {
	for _, b := range w.Bodies {
		if b.Type != BodyDynamic || b.IsSleeping {
			continue
		}
		// Apply gravity
		b.Force = b.Force.Add(w.Gravity.Multiply(b.Mass))
		// Semi-implicit Euler
		acc := b.Force.Multiply(b.InvMass)
		b.Velocity = b.Velocity.Add(acc.Multiply(dt))
		// Damping (simple)
		b.Velocity = b.Velocity.Multiply(1.0 / (1.0 + b.LinearDamping*dt))
	}
}

func (w *World) IntegrateVelocities(dt float64) {
	for _, b := range w.Bodies {
		if b.Type == BodyStatic || b.IsSleeping {
			continue
		}
		b.Position = b.Position.Add(b.Velocity.Multiply(dt))
	}
}

// ---------------- Broadphase (spatial hash) ----------------

func (w *World) aabbForBody(b *Body) (minX, minY, maxX, maxY float64) {
	switch b.Collider.Type {
	case ShapeCircle:
		r := b.Collider.Circle.Radius
		return b.Position.X - r, b.Position.Y - r, b.Position.X + r, b.Position.Y + r
	case ShapeAABB:
		hw, hh := b.Collider.AABB.HalfW, b.Collider.AABB.HalfH
		return b.Position.X - hw, b.Position.Y - hh, b.Position.X + hw, b.Position.Y + hh
	}
	return 0, 0, 0, 0
}

func (w *World) hashClear() {
	for k := range w.grid {
		delete(w.grid, k)
	}
}

func (w *World) cell(p float64) int { return int(math.Floor(float64(p / w.cellSize))) }

func (w *World) hashInsert(b *Body) {
	minX, minY, maxX, maxY := w.aabbForBody(b)
	cx0, cy0 := w.cell(minX), w.cell(minY)
	cx1, cy1 := w.cell(maxX), w.cell(maxY)
	for cx := cx0; cx <= cx1; cx++ {
		row, ok := w.grid[cx]
		if !ok {
			row = make(map[int][]*Body)
			w.grid[cx] = row
		}
		for cy := cy0; cy <= cy1; cy++ {
			row[cy] = append(row[cy], b)
		}
	}
}

func (w *World) BroadphasePairs() [][2]*Body {
	w.hashClear()
	for _, b := range w.Bodies {
		if b.IsSleeping {
			continue
		}
		w.hashInsert(b)
	}
	out := make([][2]*Body, 0)
	seen := make(map[*Body]map[*Body]bool)
	for _, col := range w.grid {
		for _, list := range col {
			n := len(list)
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					a, b := list[i], list[j]
					if a.Type == BodyStatic && b.Type == BodyStatic {
						continue
					}
					if seen[a] == nil {
						seen[a] = make(map[*Body]bool)
					}
					if seen[a][b] {
						continue
					}
					seen[a][b] = true
					out = append(out, [2]*Body{a, b})
				}
			}
		}
	}
	return out
}

// ---------------- Narrowphase ----------------

func collideCircleCircle(a, b *Body) (bool, Contact) {
	ra := a.Collider.Circle.Radius
	rb := b.Collider.Circle.Radius
	d := b.Position.Subtract(a.Position)
	dist2 := d.Len2()
	r := ra + rb
	if dist2 > r*r {
		return false, Contact{}
	}
	dist := float64(math.Sqrt(float64(dist2)))
	normal := Vector2D{}
	pen := r - dist
	if dist != 0 {
		normal = d.Multiply(1.0 / dist)
	} else {
		normal = Vector2D{1, 0}
	}
	point := a.Position.Add(normal.Multiply(ra))
	return true, Contact{A: a, B: b, Normal: normal, Penetration: pen, Point: point}
}

func collideAABBAABB(a, b *Body) (bool, Contact) {
	aw, ah := a.Collider.AABB.HalfW, a.Collider.AABB.HalfH
	bw, bh := b.Collider.AABB.HalfW, b.Collider.AABB.HalfH
	dx := b.Position.X - a.Position.X
	px := (aw + bw) - float64(math.Abs(float64(dx)))
	if px <= 0 {
		return false, Contact{}
	}
	dy := b.Position.Y - a.Position.Y
	py := (ah + bh) - float64(math.Abs(float64(dy)))
	if py <= 0 {
		return false, Contact{}
	}
	// choose axis of minimum penetration
	if px < py {
		normal := Vector2D{1, 0}
		if dx < 0 {
			normal = Vector2D{-1, 0}
		}
		return true, Contact{A: a, B: b, Normal: normal, Penetration: px}
	}
	normal := Vector2D{0, 1}
	if dy < 0 {
		normal = Vector2D{0, -1}
	}
	return true, Contact{A: a, B: b, Normal: normal, Penetration: py}
}

func collideCircleAABB(circ, box *Body) (bool, Contact) {
	// clamp circle center to box
	hw, hh := box.Collider.AABB.HalfW, box.Collider.AABB.HalfH
	local := circ.Position.Subtract(box.Position)
	clamped := Vector2D{
		X: float64(math.Max(float64(-hw), math.Min(float64(local.X), float64(hw)))),
		Y: float64(math.Max(float64(-hh), math.Min(float64(local.Y), float64(hh)))),
	}
	closest := box.Position.Add(clamped)
	d := circ.Position.Subtract(closest)
	r := circ.Collider.Circle.Radius
	d2 := d.Len2()
	if d2 > r*r {
		return false, Contact{}
	}
	dist := float64(math.Sqrt(float64(d2)))
	normal := Vector2D{}
	pen := r - dist
	if dist != 0 {
		normal = d.Multiply(1.0 / dist)
	} else {
		// circle center is inside box
		// pick axis to push out
		if math.Abs(float64(local.X)) > math.Abs(float64(local.Y)) {
			if local.X > 0 {
				normal = Vector2D{1, 0}
			} else {
				normal = Vector2D{-1, 0}
			}
		} else {
			if local.Y > 0 {
				normal = Vector2D{0, 1}
			} else {
				normal = Vector2D{0, -1}
			}
		}
	}
	return true, Contact{A: circ, B: box, Normal: normal, Penetration: pen, Point: closest}
}

func Narrowphase(a, b *Body) (bool, Contact) {
	switch a.Collider.Type {
	case ShapeCircle:
		switch b.Collider.Type {
		case ShapeCircle:
			return collideCircleCircle(a, b)
		case ShapeAABB:
			hit, m := collideCircleAABB(a, b)
			return hit, m
		}
	case ShapeAABB:
		switch b.Collider.Type {
		case ShapeCircle:
			hit, m := collideCircleAABB(b, a)
			// flip normal from A to B
			m.A, m.B = a, b
			m.Normal = m.Normal.Multiply(-1)
			return hit, m
		case ShapeAABB:
			return collideAABBAABB(a, b)
		}
	}
	return false, Contact{}
}

// ---------------- Solver ----------------

func (w *World) Solve(contacts []Contact) {
	for i := 0; i < w.Iterations; i++ {
		for idx := range contacts {
			c := &contacts[idx]
			a, b := c.A, c.B
			// Relative velocity along normal
			rv := b.Velocity.Subtract(a.Velocity)
			velAlongNormal := rv.Dot(c.Normal)
			// Do not resolve if separating
			if velAlongNormal > 0 {
				continue
			}

			// Restitution
			e := float64(math.Max(float64(a.Restitution), float64(b.Restitution)))

			// Impulse scalar
			jn := -(1 + e) * velAlongNormal
			invMassSum := a.InvMass + b.InvMass
			if invMassSum == 0 {
				continue
			}
			jn /= invMassSum

			impulse := c.Normal.Multiply(jn)
			if a.Type == BodyDynamic {
				a.Velocity = a.Velocity.Subtract(impulse.Multiply(a.InvMass))
			}
			if b.Type != BodyStatic {
				b.Velocity = b.Velocity.Add(impulse.Multiply(b.InvMass))
			}

			// Friction
			rv = b.Velocity.Subtract(a.Velocity)
			tangent := rv.Subtract(c.Normal.Multiply(rv.Dot(c.Normal))).Normalize()
			jt := -rv.Dot(tangent)
			jt /= invMassSum

			mu_s := float64(math.Sqrt(float64(a.StaticFriction * b.StaticFriction)))
			mu_d := float64(math.Sqrt(float64(a.DynamicFriction * b.DynamicFriction)))

			var frictionImpulse Vector2D
			if float64(math.Abs(float64(jt))) < jn*mu_s {
				frictionImpulse = tangent.Multiply(jt)
			} else {
				frictionImpulse = tangent.Multiply(-jn * mu_d)
			}
			if a.Type == BodyDynamic {
				a.Velocity = a.Velocity.Subtract(frictionImpulse.Multiply(a.InvMass))
			}
			if b.Type != BodyStatic {
				b.Velocity = b.Velocity.Add(frictionImpulse.Multiply(b.InvMass))
			}

			// Positional correction (baumgarte)
			const percent = 0.4
			const slop = 0.01
			corr := c.Normal.Multiply(percent * float64(math.Max(0, float64(c.Penetration-slop))) / invMassSum)
			if a.Type == BodyDynamic {
				a.Position = a.Position.Subtract(corr.Multiply(a.InvMass))
			}
			if b.Type != BodyStatic {
				b.Position = b.Position.Add(corr.Multiply(b.InvMass))
			}
		}
	}
}

// ---------------- Step ----------------

func (w *World) Step(dt float64) {
	// 1) Integrate forces (gravity etc.)
	w.IntegrateForces(dt)

	// 2) Generate contacts
	pairs := w.BroadphasePairs()
	contacts := make([]Contact, 0, len(pairs))
	for _, p := range pairs {
		if ok, m := Narrowphase(p[0], p[1]); ok {
			contacts = append(contacts, m)
		}
	}

	// 3) Solve contacts (iterative impulses)
	w.Solve(contacts)

	// 4) Integrate velocities (update positions)
	w.IntegrateVelocities(dt)

	// 5) Clear forces for next step
	w.ClearForces()
}
