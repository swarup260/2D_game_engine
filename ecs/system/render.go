package system

import (
	"reflect"

	"game_engine/core"
	"game_engine/ecs/components"

	"github.com/veandco/go-sdl2/sdl"
)

type SpriteRenderSystem struct{}

func (s *SpriteRenderSystem) GetRequiredComponents() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&components.Transform{}),
		reflect.TypeOf(&components.Sprite{}),
	}
}

func (s *SpriteRenderSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {
	// no-op, this is only a render system
}

func (s *SpriteRenderSystem) Render(renderer *sdl.Renderer, interpolation float64, entities []core.Entity, manager *core.ECSManager) {
	for _, e := range entities {

		t, _ := manager.GetComponent(e, reflect.TypeOf(&components.Transform{}))
		s, _ := manager.GetComponent(e, reflect.TypeOf(&components.Sprite{}))

		transform := t.(*components.Transform)
		sprite := s.(*components.Sprite)

		// Interpolated position
		x := transform.PrevPosition.X + (transform.Position.X-transform.PrevPosition.X)*interpolation
		y := transform.Position.Y + (transform.Position.Y-transform.Position.Y)*interpolation

		dst := &sdl.Rect{X: int32(x), Y: int32(y), W: sprite.Width, H: sprite.Height}
		if sprite.Texture != nil {
			renderer.Copy(sprite.Texture, nil, dst)
		} else {
			// Render as colored rectangle if no texture
			renderer.SetDrawColor(sprite.Color.R, sprite.Color.G, sprite.Color.B, sprite.Color.A)
			renderer.FillRect(dst)
		}
	}
}

// type RenderingSystem struct {
// 	interpPos gameEngineMath.Vector
// }

// func (s *RenderingSystem) Update(dt float64, entities []core.Entity, manager *core.ECSManager) {

// 	for _, entity := range entities {
// 		transform, _ := manager.GetComponent(entity, reflect.TypeOf(&components.Transform{}))
// 		prevTransform, _ := manager.GetComponent(entity, reflect.TypeOf(&components.PreviousTransform{}))

// 		trans := transform.(*components.Transform)
// 		prev := prevTransform.(*components.PreviousTransform)

// 		// Interpolate position
// 		s.interpPos = prev.Position.Add(trans.Position.Sub(prev.Position).Mul(dt))

// 	}

// }

// func (s *RenderingSystem) Render(renderer *sdl.Renderer, entities []core.Entity, manager *core.ECSManager) {
// 	fmt.Println("RENDERING", len(entities))
// 	for _, entity := range entities {
// 		collider, _ := manager.GetComponent(entity, reflect.TypeOf(&components.Collider{}))
// 		render, _ := manager.GetComponent(entity, reflect.TypeOf(&components.Renderable{}))

// 		col := collider.(*components.Collider)
// 		rend := render.(*components.Renderable)

// 		// Interp rotation if needed: interpRot := prev.Rotation + (trans.Rotation - prev.Rotation) * alpha

// 		renderer.SetDrawColor(rend.Color.R, rend.Color.G, rend.Color.B, rend.Color.A)

// 		// Draw based on collider shape (for physics viz)
// 		worldPos := s.interpPos.Add(col.Offset) // Assuming no rotation for simplicity

// 		switch col.Shape {
// 		case "circle":
// 			s.drawCircle(renderer, worldPos, col.Radius, rend.DrawMode == "outline")
// 		case "box":
// 			half := col.HalfSize
// 			rect := &sdl.Rect{
// 				X: int32(worldPos.X - half.X),
// 				Y: int32(worldPos.Y - half.Y),
// 				W: int32(half.X * 2),
// 				H: int32(half.Y * 2),
// 			}
// 			if rend.DrawMode == "outline" {
// 				renderer.DrawRect(rect)
// 			} else {
// 				renderer.FillRect(rect)
// 			}
// 			// Add more shapes
// 		}

// 	}
// }

// func (s *RenderingSystem) GetRequiredComponents() []reflect.Type {
// 	return []reflect.Type{
// 		reflect.TypeOf(&components.Transform{}),
// 		reflect.TypeOf(&components.PreviousTransform{}),
// 		reflect.TypeOf(&components.Collider{}),
// 		reflect.TypeOf(&components.Renderable{}),
// 	}
// }

// // Helper: Approximate circle with lines (SDL2 has no built-in circle; use this or sdl_gfx)
// // In systems/render.go
// func (s *RenderingSystem) drawCircle(renderer *sdl.Renderer, center gameEngineMath.Vector, radius float64, outline bool) {
// 	if !outline {
// 		// Simple "fill" approximation: draw lines from center to edges
// 		steps := 32
// 		angleStep := 2 * math.Pi / float64(steps)
// 		cx, cy := int32(center.X), int32(center.Y)
// 		for i := 0; i < steps; i++ {
// 			angle := float64(i) * angleStep
// 			x := int32(center.X + radius*math.Cos(angle))
// 			y := int32(center.Y + radius*math.Sin(angle))
// 			renderer.DrawLine(cx, cy, x, y)
// 		}
// 		// Note: This is not true fill; for proper fill, use sdl2_gfx or render to texture.
// 	} else {
// 		// Outline as before
// 		steps := 32
// 		angleStep := 2 * math.Pi / float64(steps)
// 		prevX := int32(center.X + radius)
// 		prevY := int32(center.Y)
// 		for i := 0; i <= steps; i++ {
// 			angle := float64(i) * angleStep
// 			x := int32(center.X + radius*math.Cos(angle))
// 			y := int32(center.Y + radius*math.Sin(angle))
// 			renderer.DrawLine(prevX, prevY, x, y)
// 			prevX, prevY = x, y
// 		}
// 	}
// }
