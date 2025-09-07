package core

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Renderer struct {
    renderer *sdl.Renderer
    clearCol sdl.Color
}

// NewRendererEngine creates a new renderer wrapper
func NewRenderer(r *sdl.Renderer, clearColor sdl.Color) *Renderer {
    return &Renderer{
        renderer: r,
        clearCol: clearColor,
    }
}

// BeginFrame clears the screen
func (re *Renderer) BeginFrame() {
    re.renderer.SetDrawColor(re.clearCol.R, re.clearCol.G, re.clearCol.B, re.clearCol.A)
    re.renderer.Clear()
}

// EndFrame presents the rendered frame
func (re *Renderer) EndFrame() {
    re.renderer.Present()
}

func (re *Renderer) GetRenderer() *sdl.Renderer {
    return re.renderer
}