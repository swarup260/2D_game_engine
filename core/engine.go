package core

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type GameEngine struct {
	window     *sdl.Window
	renderer   *sdl.Renderer
	running    bool
	fixedDelta float64
	// Timing - Fixed timestep for physics, variable for rendering
	targetFPS     int
	frameTime     time.Duration
	lastFrameTime time.Time

	// Fixed timestep for physics
	physicsTimestep float64 // Fixed timestep (e.g., 1/60 = 0.0166...)
	accumulator     float64 // Accumulated time for physics updates
	maxFrameTime    float64 // Cap to prevent spiral of death

	// Variable timestep for rendering
	renderDeltaTime float64 // Time since last render

	// Frame rate tracking
	frameCount int
	fpsTimer   time.Time
	currentFPS int

	// Core systems (to be implemented)
	// Subsystems
	Input *InputManager
	ECS   *ECSManager
	Render  *Renderer
	Scenes  *SceneManager
	// Physics *PhysicsSystem
	// Audio   *AudioManager
}

func NewGameEngine(title string, width, height int32, targetFPS int) (*GameEngine, error) {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		return nil, fmt.Errorf("failed to initialize SDL: %v", err)
	}

	if err := img.Init(img.INIT_PNG); err != nil {
		return nil, fmt.Errorf("failed to initialize SDL Image: %v", err)
	}

	// Create window
	window, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		width,
		height,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %v", err)
	}

	// Create renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		window.Destroy()
		return nil, fmt.Errorf("failed to create renderer: %v", err)
	}

	// Fixed physics timestep (60 FPS physics regardless of render FPS)
	physicsTimestep := 1.0 / 60.0

	engine := &GameEngine{
		window:          window,
		renderer:        renderer,
		running:         false,
		targetFPS:       targetFPS,
		frameTime:       time.Duration(1000/targetFPS) * time.Millisecond,
		lastFrameTime:   time.Now(),
		physicsTimestep: physicsTimestep,
		accumulator:     0.0,
		maxFrameTime:    0.25,              // Cap at 250ms to prevent spiral of death
		Input:           NewInputManager(), // Initialize the inputManager
		ECS:             NewECSManager(),   // Initialize the inputManager
		Render:          NewRenderer(renderer,sdl.Color{R: 0, G: 0, B: 0, A: 255}),   // Black background
		Scenes:          NewSceneManager(),
	}

	return engine, nil

}

func (ge *GameEngine) Run() error {
	defer ge.cleanup()

	ge.running = true
	ge.lastFrameTime = time.Now()
	ge.fpsTimer = time.Now()

	fmt.Printf("Starting game engine:\n")
	fmt.Printf("- Target render FPS: %d\n", ge.targetFPS)
	fmt.Printf("- Fixed physics timestep: %.4fs (%.0f FPS)\n", ge.physicsTimestep, 1.0/ge.physicsTimestep)

	for ge.running {
		frameStart := time.Now()

		// Calculate frame time
		currentTime := time.Now()
		frameTime := currentTime.Sub(ge.lastFrameTime).Seconds()
		ge.lastFrameTime = currentTime

		// Cap frame time to prevent spiral of death
		if frameTime > ge.maxFrameTime {
			frameTime = ge.maxFrameTime
		}

		// Store render delta time (variable timestep for rendering)
		ge.renderDeltaTime = frameTime

		// Add frame time to accumulator for physics
		ge.accumulator += frameTime

		// Input processing
		ge.handleEvents()

		// Fixed timestep physics updates
		// Run physics multiple times if we've accumulated enough time
		for ge.accumulator >= ge.physicsTimestep {
			ge.updatePhysics(ge.physicsTimestep)
			ge.accumulator -= ge.physicsTimestep
		}

		// Calculate interpolation factor for smooth rendering
		// This allows rendering between physics steps
		interpolation := ge.accumulator / ge.physicsTimestep

		// Variable timestep updates (gameplay logic, animations, etc.)
		ge.updateGameplay(ge.renderDeltaTime)

		// Render with interpolation for smooth movement
		ge.render(interpolation)

		// Frame rate limiting for rendering
		ge.limitFrameRate(frameStart)

		// Update FPS counter
		ge.updateFPS()

	}

	return nil

}

// handleEvents processes SDL events and input
func (ge *GameEngine) handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.WindowEvent:
			if e.Event == sdl.WINDOWEVENT_RESIZED {
				// Handle window resize
				ge.handleWindowResize(e.Data1, e.Data2)
			}
		}
		ge.Input.Update()
		if ge.Input.ShouldQuit() {
			ge.Stop()
		}
	}
}

// updatePhysics handles fixed timestep physics updates
func (ge *GameEngine) updatePhysics(fixedDeltaTime float64) {
	// Update physics systems with fixed timestep
	// This ensures consistent physics regardless of frame rate

	// Update physics in scene manager
	ge.Scenes.UpdatePhysics(fixedDeltaTime)

	// Example physics operations:
	// - Collision detection and response
	// - Rigid body dynamics
	// - Particle physics
	// - Any time-critical simulations
}

// updateGameplay handles variable timestep gameplay updates
func (ge *GameEngine) updateGameplay(deltaTime float64) {
	// Update input manager
	ge.Input.Update()

	// Update gameplay logic with variable timestep
	// This allows for smooth animations and non-critical updates
	ge.Scenes.Update(deltaTime)

	// Update audio system
	// TODO AUDIO HANDLER

	// ECS System
	ge.ECS.UpdateSystems(deltaTime)

	// Example gameplay operations:
	// - UI animations
	// - Particle effects (non-physics)
	// - Audio synchronization
	// - Camera movement
	// - Visual effects
}

// render handles all rendering with interpolation
func (ge *GameEngine) render(interpolation float64) {
	// Clear screen with background color
	ge.Render.BeginFrame()
	// Render current scene with interpolation for smooth movement
	// Interpolation allows rendering positions between physics steps
	ge.Scenes.Render(interpolation)

	// Present the frame
	ge.Render.EndFrame()
}

// limitFrameRate ensures consistent render timing
func (ge *GameEngine) limitFrameRate(frameStart time.Time) {
	frameTime := time.Since(frameStart)

	if frameTime < ge.frameTime {
		sleepTime := ge.frameTime - frameTime
		time.Sleep(sleepTime)
	}
}

// updateFPS tracks and displays current FPS
func (ge *GameEngine) updateFPS() {
	ge.frameCount++

	if time.Since(ge.fpsTimer) >= time.Second {
		ge.currentFPS = ge.frameCount
		ge.frameCount = 0
		ge.fpsTimer = time.Now()

		// Optional: Print FPS and timing info
		fmt.Printf("Render FPS: %d, Physics: %.1f FPS, Accumulator: %.4f\n",
			ge.currentFPS, 1.0/ge.physicsTimestep, ge.accumulator)
	}
}

// handleWindowResize handles window resize events
func (ge *GameEngine) handleWindowResize(width, height int32) {
	fmt.Printf("Window resized to: %dx%d\n", width, height)
	// ge.sceneManager.HandleResize(width, height)
}

// SetPhysicsTimestep allows changing the fixed physics timestep
func (ge *GameEngine) SetPhysicsTimestep(timestep float64) {
	ge.physicsTimestep = timestep
}

// GetPhysicsTimestep returns the current physics timestep
func (ge *GameEngine) GetPhysicsTimestep() float64 {
	return ge.physicsTimestep
}

// Stop gracefully stops the engine
func (ge *GameEngine) Stop() {
	ge.running = false
}

// IsRunning returns whether the engine is currently running
func (ge *GameEngine) IsRunning() bool {
	return ge.running
}

// GetRenderDeltaTime returns the variable timestep for rendering/gameplay
func (ge *GameEngine) GetRenderDeltaTime() float64 {
	return ge.renderDeltaTime
}

// GetFPS returns current frames per second
func (ge *GameEngine) GetFPS() int {
	return ge.currentFPS
}

// GetRenderer returns the SDL renderer for direct access if needed
func (ge *GameEngine) GetRenderer() *sdl.Renderer {
	return ge.renderer
}

// Destroy cleans up engine resources
func (ge *GameEngine) cleanup() {
	if ge.renderer != nil {
		ge.renderer.Destroy()
	}
	if ge.window != nil {
		ge.window.Destroy()
	}
	img.Quit()
	sdl.Quit()
	ge.Input.Cleanup()
}
