package core

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// Entity represents a unique game object identifier
type Entity uint32

// Component interface that all components must implement
type Component interface{}

// System interface that all systems must implement
type System interface {
	Update(dt float64, entities []Entity, manager *ECSManager)
	GetRequiredComponents() []reflect.Type
}

// RenderSystem interface for systems that need rendering
type RenderSystem interface {
	System
	Render(renderer *sdl.Renderer, interpolation float64, entities []Entity, manager *ECSManager)
}

// ECSManager manages entities, components, and systems
type ECSManager struct {
	mutex        sync.RWMutex
	nextEntityID Entity
	entities     map[Entity]bool
	components   map[Entity]map[reflect.Type]Component
	// Systems by category
	physicsSystems  []System       // runs in fixed timestep
	gameplaySystems []System       // runs in variable timestep
	renderSystems   []RenderSystem // runs every frame with interpolation
}

// NewECSManager creates a new ECS manager
func NewECSManager() *ECSManager {
	return &ECSManager{
		nextEntityID:    1,
		entities:        make(map[Entity]bool),
		components:      make(map[Entity]map[reflect.Type]Component),
		physicsSystems:  make([]System, 0),
		gameplaySystems: make([]System, 0),
		renderSystems:   make([]RenderSystem, 0),
	}
}

// CreateEntity creates a new entity
func (ecs *ECSManager) CreateEntity() Entity {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()

	entity := ecs.nextEntityID
	ecs.nextEntityID++
	ecs.entities[entity] = true
	ecs.components[entity] = make(map[reflect.Type]Component)

	return entity
}

// DestroyEntity marks an entity for destruction
func (ecs *ECSManager) DestroyEntity(entity Entity) {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()

	delete(ecs.entities, entity)
	delete(ecs.components, entity)
}

// AddComponent adds a component to an entity
func (ecs *ECSManager) AddComponent(entity Entity, component Component) error {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()

	if _, ok := ecs.entities[entity]; !ok {
		return fmt.Errorf("entity %d does not exist", entity)
	}
	componentType := reflect.TypeOf(component)
	ecs.components[entity][componentType] = component

	return nil
}

// GetComponent retrieves a component from an entity
func (ecs *ECSManager) GetComponent(entity Entity, componentType reflect.Type) (Component, bool) {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()

	if !ecs.entities[entity] {
		return nil, false
	}

	component, exists := ecs.components[entity][componentType]
	return component, exists
}

// RemoveComponent removes a component from an entity
func (ecs *ECSManager) RemoveComponent(entity Entity, componentType reflect.Type) {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()

	if !ecs.entities[entity] {
		return
	}

	delete(ecs.components[entity], componentType)
}

// HasComponent checks if an entity has a specific component
func (ecs *ECSManager) HasComponent(entity Entity, componentType reflect.Type) bool {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()

	if !ecs.entities[entity] {
		return false
	}

	_, exists := ecs.components[entity][componentType]
	return exists
}

// GetEntitiesWithComponents returns entities that have all specified components
func (ecs *ECSManager) GetEntitiesWithComponents(requiredTypes []reflect.Type) []Entity {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()

	var entities []Entity
	for entity := range ecs.entities {
		hasAll := true
		for _, componentType := range requiredTypes {
			if _, exists := ecs.components[entity][componentType]; !exists {
				hasAll = false
				break
			}
		}
		if hasAll {
			entities = append(entities, entity)
		}
	}
	return entities
}

// UpdateGameplay runs all gameplay loop systems with the given delta time
func (ecs *ECSManager) UpdateGameplay(dt float64) {
	for _, system := range ecs.gameplaySystems {
		entities := ecs.GetEntitiesWithComponents(system.GetRequiredComponents())
		system.Update(dt, entities, ecs)
	}
}

// UpdatePhysics runs all physics at determine time
func (ecs *ECSManager) UpdatePhysics(dt float64) {
	for _, system := range ecs.physicsSystems {
		entities := ecs.GetEntitiesWithComponents(system.GetRequiredComponents())
		system.Update(dt, entities, ecs)
	}
}

// RenderSystems runs all render systems
func (ecs *ECSManager) RenderSystems(renderer *sdl.Renderer, interpolation float64) {
	for _, system := range ecs.renderSystems {
		entities := ecs.GetEntitiesWithComponents(system.GetRequiredComponents())
		system.Render(renderer, interpolation, entities, ecs)
	}
}

func (ecs *ECSManager) AddPhysicsSystem(system System) {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()
	ecs.physicsSystems = append(ecs.physicsSystems, system)
}

func (ecs *ECSManager) AddGameplaySystem(system System) {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()
	ecs.gameplaySystems = append(ecs.gameplaySystems, system)
}

func (ecs *ECSManager) AddRenderSystem(system RenderSystem) {
	ecs.mutex.Lock()
	defer ecs.mutex.Unlock()
	ecs.renderSystems = append(ecs.renderSystems, system)
}

// CopyEntity creates a new entity and copies all components from the source entity
func (ecs *ECSManager) CopyEntity(source Entity) (Entity, error) {
	ecs.mutex.RLock()
	sourceComponents, exists := ecs.components[source]
	ecs.mutex.RUnlock()

	if !exists {
		return 0, fmt.Errorf("source entity %d does not exist", source)
	}

	newEntity := ecs.CreateEntity()

	for compType, comp := range sourceComponents {
		// Attempt to deep copy the component
		copiedComp, err := deepCopyComponent(comp)
		if err != nil {
			return 0, fmt.Errorf("failed to copy component %v: %v", compType, err)
		}
		if err := ecs.AddComponent(newEntity, copiedComp); err != nil {
			return 0, err
		}
	}

	return newEntity, nil
}

func deepCopyComponent(comp Component) (Component, error) {
	origVal := reflect.ValueOf(comp)

	// Handle pointer types
	if origVal.Kind() == reflect.Ptr {
		origVal = origVal.Elem()
	}

	// Create a new instance of the same type
	copyVal := reflect.New(origVal.Type()).Elem()

	// Copy the value
	copyVal.Set(origVal)

	// Return as interface
	if compVal := copyVal.Interface(); compVal != nil {
		return compVal, nil
	}
	return nil, fmt.Errorf("failed to copy component of type %v", origVal.Type())
}

// GetEntityCount returns the number of active entities
func (ecs *ECSManager) GetEntityCount() int {
	ecs.mutex.RLock()
	defer ecs.mutex.RUnlock()
	return len(ecs.entities)
}
