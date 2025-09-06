package core

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// InputState represents the state of an input (pressed, held, released)
type InputState int

const (
	InputStateUp InputState = iota
	InputStatePressed
	InputStateHeld
	InputStateReleased
)

// MouseButton represents mouse button indices
type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonMiddle
	MouseButtonRight
)

// InputManager handles all input events and states
type InputManager struct {
	sync.Mutex // Mutex to protect the maps from concurrent access
	// Keyboard state
	keyboardState     []uint8
	prevKeyboardState []uint8

	// Mouse state
	mouseX, mouseY   int32
	mouseButtons     [3]InputState
	prevMouseButtons [3]bool
	mouseWheel       int32

	// Controller state
	controllers           map[sdl.JoystickID]*sdl.GameController
	controllerAxes        map[sdl.JoystickID][]int16
	controllerButtons     map[sdl.JoystickID][]bool
	prevControllerButtons map[sdl.JoystickID][]bool

	// Event queue for custom handling
	events []sdl.Event

	// Input bindings
	keyBindings        map[string]sdl.Scancode
	mouseBindings      map[string]MouseButton
	controllerBindings map[string]uint8

	quit bool
}

// NewInputManager creates and initializes a new InputManager
func NewInputManager() *InputManager {
	im := &InputManager{
		keyboardState:         make([]uint8, sdl.NUM_SCANCODES),
		prevKeyboardState:     make([]uint8, sdl.NUM_SCANCODES),
		controllers:           make(map[sdl.JoystickID]*sdl.GameController),
		controllerAxes:        make(map[sdl.JoystickID][]int16),
		controllerButtons:     make(map[sdl.JoystickID][]bool),
		prevControllerButtons: make(map[sdl.JoystickID][]bool),
		events:                make([]sdl.Event, 0),
		keyBindings:           make(map[string]sdl.Scancode),
		mouseBindings:         make(map[string]MouseButton),
		controllerBindings:    make(map[string]uint8),
		quit:                  false,
	}

	// Initialize SDL subsystems
	if err := sdl.Init(sdl.INIT_GAMECONTROLLER); err != nil {
		panic(err)
	}

	// Initialize default key bindings
	im.initDefaultBindings()

	return im
}

// initDefaultBindings sets up common key bindings
func (im *InputManager) initDefaultBindings() {
	// Movement
	im.keyBindings["up"] = sdl.SCANCODE_W
	im.keyBindings["down"] = sdl.SCANCODE_S
	im.keyBindings["left"] = sdl.SCANCODE_A
	im.keyBindings["right"] = sdl.SCANCODE_D
	im.keyBindings["jump"] = sdl.SCANCODE_SPACE

	// Actions
	im.keyBindings["action"] = sdl.SCANCODE_E
	im.keyBindings["attack"] = sdl.SCANCODE_F
	im.keyBindings["menu"] = sdl.SCANCODE_ESCAPE

	// Mouse bindings
	im.mouseBindings["primary"] = MouseButtonLeft
	im.mouseBindings["secondary"] = MouseButtonRight
	im.mouseBindings["tertiary"] = MouseButtonMiddle

	// Controller bindings (Xbox layout)
	im.controllerBindings["jump"] = 0   // A button
	im.controllerBindings["action"] = 1 // B button
	im.controllerBindings["attack"] = 2 // X button
	im.controllerBindings["menu"] = 6   // Back button
}

// Update processes all pending SDL events and updates input states
func (im *InputManager) Update() {
	// Store previous states
	copy(im.prevKeyboardState, im.keyboardState)
	for i := range im.prevMouseButtons {
		im.prevMouseButtons[i] = im.mouseButtons[i] == InputStateHeld || im.mouseButtons[i] == InputStatePressed
	}

	// Copy previous controller button states
	for id := range im.controllerButtons {
		if im.prevControllerButtons[id] == nil {
			im.prevControllerButtons[id] = make([]bool, len(im.controllerButtons[id]))
		}
		copy(im.prevControllerButtons[id], im.controllerButtons[id])
	}

	// Reset mouse wheel
	im.mouseWheel = 0

	// Clear events
	im.events = im.events[:0]

	// Process SDL events
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		im.events = append(im.events, event)

		switch e := event.(type) {
		case *sdl.QuitEvent:
			im.quit = true

		case *sdl.MouseWheelEvent:
			im.mouseWheel = e.Y

		case *sdl.ControllerDeviceEvent:
			im.handleControllerEvent(e)
		}
	}

	// Update keyboard state
	keyState := sdl.GetKeyboardState()
	copy(im.keyboardState, keyState)

	// Update mouse state
	x, y, _ := sdl.GetMouseState()
	im.mouseX = x
	im.mouseX = y

	// Update mouse button states
	im.updateMouseButtonStates()

	// Update controller states
	im.updateControllerStates()
}

// updateMouseButtonStates updates the state of mouse buttons
func (im *InputManager) updateMouseButtonStates() {
	_, _, mouseState := sdl.GetMouseState()
	buttons := []uint32{sdl.BUTTON_LEFT, sdl.BUTTON_MIDDLE, sdl.BUTTON_RIGHT}

	for i, button := range buttons {
		pressed := mouseState&sdl.Button(button) != 0
		prevPressed := im.prevMouseButtons[i]

		if pressed && !prevPressed {
			im.mouseButtons[i] = InputStatePressed
		} else if pressed && prevPressed {
			im.mouseButtons[i] = InputStateHeld
		} else if !pressed && prevPressed {
			im.mouseButtons[i] = InputStateReleased
		} else {
			im.mouseButtons[i] = InputStateUp
		}
	}
}

// updateControllerStates updates controller button and axis states
func (im *InputManager) updateControllerStates() {
	for id, controller := range im.controllers {
		if controller == nil {
			continue
		}

		// Update button states
        for i := 0; i < len(im.controllerButtons[id]); i++ {
            im.controllerButtons[id][i] = controller.Button(sdl.GameControllerButton(i)) == 1
        }

        // Update axis states
        for i := 0; i < len(im.controllerAxes[id]); i++ {
            im.controllerAxes[id][i] = controller.Axis(sdl.GameControllerAxis(i))
        }
	}
}

// handleControllerEvent handles controller connection and disconnection events.
// It is thread-safe due to the use of a mutex.
func (im *InputManager) handleControllerEvent(event *sdl.ControllerDeviceEvent) {
	im.Lock()
	defer im.Unlock()

	switch event.Type {
	case sdl.CONTROLLERDEVICEADDED:
		// Attempt to open the newly added controller.
		controller := sdl.GameControllerOpen(int(event.Which))
		if controller == nil {
			fmt.Printf("Failed to open controller with ID %d: %s", event.Which, sdl.GetError())
			return
		}

		// Check if it's already in the map to prevent duplicates.
		if _, exists := im.controllers[event.Which]; exists {
			fmt.Printf("Controller with ID %d already exists. Skipping.", event.Which)
			controller.Close()
			return
		}

		// Store the new controller and initialize its state arrays.
		fmt.Printf("Controller added: %s (ID: %d)", controller.Name(), event.Which)
		im.controllers[event.Which] = controller
		im.controllerAxes[event.Which] = make([]int16, sdl.CONTROLLER_AXIS_MAX)
		im.controllerButtons[event.Which] = make([]bool, sdl.CONTROLLER_BUTTON_MAX)
		im.prevControllerButtons[event.Which] = make([]bool, sdl.CONTROLLER_BUTTON_MAX)

	case sdl.CONTROLLERDEVICEREMOVED:
		// Look up the controller by its ID and remove it.
		if controller, exists := im.controllers[event.Which]; exists {
			fmt.Printf("Controller removed: %s (ID: %d)", controller.Name(), event.Which)
			controller.Close()
			// The `sdl.JoystickID` type is correct for deleting from the map.
			// No type conversion to `int32` is necessary here.
			delete(im.controllers, event.Which)
			delete(im.controllerAxes, event.Which)
			delete(im.controllerButtons, event.Which)
			delete(im.prevControllerButtons, event.Which)
		} else {
			fmt.Printf("Attempted to remove a non-existent controller with ID %d", event.Which)
		}
	}
}

// GetButtonState returns the current state of a specific button for a given controller.
func (im *InputManager) GetButtonState(id sdl.JoystickID, button int) (bool, bool) {
	im.Lock()
	defer im.Unlock()

	if buttons, exists := im.controllerButtons[id]; exists && int(button) < len(buttons) {
		prevButtons := im.prevControllerButtons[id]
		return buttons[button], buttons[button] && !prevButtons[button]
	}
	return false, false
}

// Keyboard input methods
func (im *InputManager) IsKeyDown(scancode sdl.Scancode) bool {
	return im.keyboardState[scancode] == 1
}

func (im *InputManager) IsKeyPressed(scancode sdl.Scancode) bool {
	return im.keyboardState[scancode] == 1 && im.prevKeyboardState[scancode] == 0
}

func (im *InputManager) IsKeyReleased(scancode sdl.Scancode) bool {
	return im.keyboardState[scancode] == 0 && im.prevKeyboardState[scancode] == 1
}

// Action-based input methods (using bindings)
func (im *InputManager) IsActionDown(action string) bool {
	if scancode, exists := im.keyBindings[action]; exists {
		return im.IsKeyDown(scancode)
	}
	if mouseBtn, exists := im.mouseBindings[action]; exists {
		return im.mouseButtons[mouseBtn] == InputStateHeld || im.mouseButtons[mouseBtn] == InputStatePressed
	}
	if controllerBtn, exists := im.controllerBindings[action]; exists {
		return im.IsControllerButtonDown(controllerBtn)
	}
	return false
}

func (im *InputManager) IsActionPressed(action string) bool {
	if scancode, exists := im.keyBindings[action]; exists {
		return im.IsKeyPressed(scancode)
	}
	if mouseBtn, exists := im.mouseBindings[action]; exists {
		return im.mouseButtons[mouseBtn] == InputStatePressed
	}
	if controllerBtn, exists := im.controllerBindings[action]; exists {
		return im.IsControllerButtonPressed(controllerBtn)
	}
	return false
}

func (im *InputManager) IsActionReleased(action string) bool {
	if scancode, exists := im.keyBindings[action]; exists {
		return im.IsKeyReleased(scancode)
	}
	if mouseBtn, exists := im.mouseBindings[action]; exists {
		return im.mouseButtons[mouseBtn] == InputStateReleased
	}
	if controllerBtn, exists := im.controllerBindings[action]; exists {
		return im.IsControllerButtonReleased(controllerBtn)
	}
	return false
}

// Mouse input methods
func (im *InputManager) GetMousePosition() (int32, int32) {
	return im.mouseX, im.mouseY
}

func (im *InputManager) IsMouseButtonDown(button MouseButton) bool {
	return im.mouseButtons[button] == InputStateHeld || im.mouseButtons[button] == InputStatePressed
}

func (im *InputManager) IsMouseButtonPressed(button MouseButton) bool {
	return im.mouseButtons[button] == InputStatePressed
}

func (im *InputManager) IsMouseButtonReleased(button MouseButton) bool {
	return im.mouseButtons[button] == InputStateReleased
}

func (im *InputManager) GetMouseWheel() int32 {
	return im.mouseWheel
}

// Controller input methods
func (im *InputManager) IsControllerButtonDown(button uint8) bool {
	for _, buttons := range im.controllerButtons {
		if int(button) < len(buttons) && buttons[button] {
			return true
		}
	}
	return false
}

func (im *InputManager) IsControllerButtonPressed(button uint8) bool {
	for id, buttons := range im.controllerButtons {
		if int(button) < len(buttons) && buttons[button] {
			prevButtons := im.prevControllerButtons[id]
			if int(button) < len(prevButtons) && !prevButtons[button] {
				return true
			}
		}
	}
	return false
}

func (im *InputManager) IsControllerButtonReleased(button uint8) bool {
	for id, buttons := range im.controllerButtons {
		if int(button) < len(buttons) && !buttons[button] {
			prevButtons := im.prevControllerButtons[id]
			if int(button) < len(prevButtons) && prevButtons[button] {
				return true
			}
		}
	}
	return false
}

func (im *InputManager) GetControllerAxis(axis uint8) int16 {
	for _, axes := range im.controllerAxes {
		if int(axis) < len(axes) {
			return axes[axis]
		}
	}
	return 0
}

// Binding management
func (im *InputManager) BindKey(action string, scancode sdl.Scancode) {
	im.keyBindings[action] = scancode
}

func (im *InputManager) BindMouse(action string, button MouseButton) {
	im.mouseBindings[action] = button
}

func (im *InputManager) BindController(action string, button uint8) {
	im.controllerBindings[action] = button
}

// Utility methods
func (im *InputManager) GetEvents() []sdl.Event {
	return im.events
}

func (im *InputManager) ShouldQuit() bool {
	return im.quit
}

func (im *InputManager) SetQuit(quit bool) {
	im.quit = quit
}

// Cleanup releases resources
func (im *InputManager) Cleanup() {
	for _, controller := range im.controllers {
		if controller != nil {
			controller.Close()
		}
	}
	im.controllers = make(map[sdl.JoystickID]*sdl.GameController)
}
