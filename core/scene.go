package core

type Scene interface {
	Init() error
	HandleInput(im InputManager)
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
func (sm *SceneManager) HandleInput(im InputManager) {
	if scene := sm.Current(); scene != nil {
		scene.HandleInput(im)
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

/* package core

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Scene interface {
	Init() error    // Lightweight setup
	Load() error    // Heavy asset loading (can run async)
	IsLoaded() bool // True when ready to switch
	HandleInput(ev sdl.Event)
	Update(dt float64)        // Variable timestep updates
	UpdatePhysics(dt float64) // Fixed timestep updates
	Render(alpha float64)
	Cleanup()
}

type SceneManager struct {
	scenes       []Scene
	loadingScene Scene
	loading      bool
	loadDone     chan struct{}
}

func NewSceneManager() *SceneManager {
	return &SceneManager{
		scenes:   make([]Scene, 0),
		loadDone: make(chan struct{}, 1),
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

// Async scene loading with a loading screen
func (sm *SceneManager) LoadSceneAsync(newScene Scene, loadingScene Scene) {
	sm.loadingScene = loadingScene
	sm.loading = true
	sm.Push(loadingScene)

	go func() {
		newScene.Init()
		newScene.Load()
		sm.loadDone <- struct{}{}
	}()
}

func (sm *SceneManager) Update(dt float64) {
	if sm.loading {
		select {
		case <-sm.loadDone:
			// Replace loading scene with new scene
			sm.Pop()
			sm.Push(sm.loadingScene) // cleanup loading scene
			sm.Pop()
			sm.Push(sm.loadingScene) // cleanup loading scene
			sm.Pop()
			sm.Push(sm.loadingScene) // cleanup loading scene
		default:
			// Still loading
		}
	}
	if scene := sm.Current(); scene != nil {
		scene.Update(dt)
	}
}

func (sm *SceneManager) UpdatePhysics(dt float64) {
	if scene := sm.Current(); scene != nil {
		scene.UpdatePhysics(dt)
	}
}

func (sm *SceneManager) Render(alpha float64) {
	if scene := sm.Current(); scene != nil {
		scene.Render(alpha)
	}
}

func (sm *SceneManager) HandleInput(ev sdl.Event) {
	if scene := sm.Current(); scene != nil {
		scene.HandleInput(ev)
	}
}

/* future coee


type LoadingScene struct {
    progress float32
}

func (ls *LoadingScene) Init() error { ls.progress = 0; return nil }
func (ls *LoadingScene) Load() error { return nil } // No heavy load here
func (ls *LoadingScene) IsLoaded() bool { return true }
func (ls *LoadingScene) HandleInput(ev sdl.Event) {}
func (ls *LoadingScene) Update(dt float32) {
    // Fake progress animation
    if ls.progress < 1 {
        ls.progress += dt * 0.5
    }
}
func (ls *LoadingScene) UpdatePhysics(dt float32) {}
func (ls *LoadingScene) Render(alpha float32) {
    // Draw loading bar or spinner
}
func (ls *LoadingScene) Cleanup() {}


type HeavyScene struct {
    loaded bool
}

func (hs *HeavyScene) Init() error { return nil }
func (hs *HeavyScene) Load() error {
    // Simulate heavy asset loading
    time.Sleep(3 * time.Second)
    hs.loaded = true
    return nil
}
func (hs *HeavyScene) IsLoaded() bool { return hs.loaded }
func (hs *HeavyScene) HandleInput(ev sdl.Event) {}
func (hs *HeavyScene) Update(dt float32) {}
func (hs *HeavyScene) UpdatePhysics(dt float32) {}
func (hs *HeavyScene) Render(alpha float32) {}
func (hs *HeavyScene) Cleanup() {}

// Using It in GameEngine

loading := &LoadingScene{}
heavy := &HeavyScene{}
engine.Scenes.LoadSceneAsync(heavy, loading)


The loop will:

Push the loading scene immediately

Load the heavy scene in a goroutine

Swap to the heavy scene when loading finishes

✅ Benefits

No frame stalls — loading happens in the background

Loading screen can animate, play music, or show tips

Works with your fixed-update / variable-render loop

Easy to extend for streaming assets in chunks


*/
