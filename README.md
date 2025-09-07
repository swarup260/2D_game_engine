# 2D_game_engine

<!-- https://gameprogrammingpatterns.com/game-loop.html -->



## Core Architecture Overview
1. Game Loop: Fixed timestep for updates, variable rendering.
    * Establish the core game loop
        > Initialize -> While(Running) { ProcessInput -> Update -> Render } -> Shutdown.
2. Rendering: Use SDL's renderer for hardware-accelerated drawing.
3. Input: Poll SDL events for keys and mouse.
4. ECS: Entities as IDs, components as structs, systems as functions that process them.
5. Physics: Simple velocity-based movement with AABB collision.
6. Assets: Load textures with SDL_image.


Main Game Loop

* Fixed timestep for physics
* Variable timestep for rendering
* Input processing
* Update/Render separation

=========================================================================================

1. Fixed Timestep for Physics:

Physics runs at consistent 60 FPS regardless of render frame rate
Uses accumulator pattern to handle multiple physics steps per frame
Prevents physics instability and ensures deterministic behavior
Includes "spiral of death" protection with maxFrameTime

2. Variable Timestep for Rendering:

Rendering uses actual frame delta time for smooth animations
Allows for high refresh rate displays (120Hz, 144Hz, etc.)
Gameplay logic (UI, effects) uses variable timestep

3. Input Processing:

Handled once per frame in handleEvents()
Proper event delegation to input manager

4. Update/Render Separation:

updatePhysics() - Fixed timestep physics
updateGameplay() - Variable timestep game logic
render() - Rendering with interpolation

Key Benefits:
Deterministic Physics:

Physics behaves identically regardless of frame rate
Essential for multiplayer games and replays

Smooth Rendering:

Interpolation between physics steps for buttery smooth movement
High refresh rate support (144Hz, 240Hz displays)

Performance Adaptive:

Automatically handles frame rate drops gracefully
Can run physics slower than rendering or vice versa

Professional Architecture:

Follows industry-standard game loop patterns
Used by engines like Unity, Unreal, and custom engines

The main loop now properly separates concerns and provides the foundation for a robust 2D game engine!



```go

// Getting window size
var w, h int32
window.GetSize(&w, &h)

// Getting mouse position  
var x, y int32
sdl.GetMouseState(&x, &y)

// Getting renderer output size
var w, h int32
renderer.GetOutputSize(&w, &h)

```



```

my-game-engine/
├── go.mod
├── go.sum
├── main.go
├── engine/
│   ├── core.go
│   ├── renderer.go
│   ├── input.go
│   ├── ecs.go
│   ├── physics.go
│   ├── audio.go
│   ├── scenes.go
│   ├── assets.go
│   ├── console.go    // New: Debug console
│   └── commands.go   // New: Command parsing
└── assets/
    ├── player.png
    ├── enemy.png
    ├── bgm.ogg
    ├── collision.wav
    └── arial.ttf     // Font for console

```





<!-- Future arch to faster the ECS -->

<!-- https://www.youtube.com/watch?v=71RSWVyOMEY -->