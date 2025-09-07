package physics

type ShapeType int

const (
	ShapeCircle ShapeType = iota
	ShapeAABB
	// ShapePolygon // extend later with SAT
)

type Circle struct{ Radius float64 }
type AABB struct{ HalfW, HalfH float64 }

type Collider struct {
	Type   ShapeType
	Circle Circle
	AABB   AABB
}
