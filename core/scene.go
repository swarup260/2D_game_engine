package core

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Scene interface {
	Init() error
	HandleInput(ev sdl.Event)
	Update(dt float64)    // variable or fixed depending on call site
	UpdatePhysics(dt float64)     // fixed timestep physics updates
	Render(alpha float64) // interpolation factor for smooth rendering
	Cleanup()
}

type SceneManager struct {
	scenes []Scene // stack of scenes
}

func NewSceneManager() *SceneManager {
	return &SceneManager{
		scenes: make([]Scene, 0),
	}
}

func (sm *SceneManager) Push(scene Scene) error {
	if err := scene.Init(); err != nil {
		return err
	}
	sm.scenes = append(sm.scenes, scene)
	return nil
}

func (sm *SceneManager) Pop() {
	if len(sm.scenes) == 0 {
		return
	}
	top := sm.scenes[len(sm.scenes)-1]
	top.Cleanup()
	sm.scenes = sm.scenes[:len(sm.scenes)-1]
}

func (sm *SceneManager) Replace(scene Scene) error {
	sm.Pop()
	return sm.Push(scene)
}

func (sm *SceneManager) Current() Scene {
	if len(sm.scenes) == 0 {
		return nil
	}
	return sm.scenes[len(sm.scenes)-1]
}

// Delegation helpers
func (sm *SceneManager) HandleInput(ev sdl.Event) {
	if scene := sm.Current(); scene != nil {
		scene.HandleInput(ev)
	}
}

// NEW: Physics update delegation
func (sm *SceneManager) UpdatePhysics(dt float64) {
    if scene := sm.Current(); scene != nil {
        scene.UpdatePhysics(dt)
    }
}

func (sm *SceneManager) Update(dt float64) {
	if scene := sm.Current(); scene != nil {
		scene.Update(dt)
	}
}

func (sm *SceneManager) Render(alpha float64) {
	if scene := sm.Current(); scene != nil {
		scene.Render(alpha)
	}
}
