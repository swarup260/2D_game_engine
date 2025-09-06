package main

import (
	"2d_game_engine/core"
)

func main() {
	engine, err := core.NewGameEngine("My 2D Game", 1280, 720, 144)
	if err != nil {
		panic(err)
	}

	// You can also customize physics timestep
	engine.SetPhysicsTimestep(1.0 / 120.0) // 120 FPS physics

	// // debug console
	// if err := engine.InitializeDebugSystems(); err != nil {
	// 	panic(err)
	// }

	if err := engine.Run(); err != nil {
		panic(err)
	}

}
