package components

import "github.com/veandco/go-sdl2/sdl"

type Sprite struct {
	Texture *sdl.Texture
	Width   int32
	Height  int32
	Color   sdl.Color
}

type Renderable struct {
	Color sdl.Color
}
