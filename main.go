package main

import (
	"game_engine/core"
	"game_engine/simulation"
)

func main() {
	engine, err := core.NewGameEngine("My 2D Game", 1280, 720, 75)
	if err != nil {
		panic(err)
	}

	// engine.Scenes.Push(simulation.NewGameScene(engine.ECS,engine.Input ,engine.Render.GetRenderer()))
	// engine.Scenes.Push(simulation.NewGravityBallScene(engine.ECS, engine.Render))
	engine.Scenes.Push(simulation.NewShapeScene(engine.ECS, engine.Render))

	if err := engine.Run(); err != nil {
		panic(err)
	}

}
