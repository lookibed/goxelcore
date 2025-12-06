package window

import (
	"log"

	"github.com/Zyko0/go-sdl3/sdl"
)

// SDLInput is a concrete implementation of the Input interface using SDL.
type SDLInput struct {
	bindings    *Bindings
	cursorState CursorState

	// State for keyboard
	keyboardState []uint8
	prevKeyboardState []uint8
	pressedKeys []Keycode
	keyCallbacks map[Keycode][]func() // Simplified: no HandlerList

	// State for mouse
	mouseState uint32
	prevMouseState uint32
	scrollValue int

	// Codepoints (for text input)
	codepoints []rune
}

// NewSDLInput creates a new SDLInput instance.
func NewSDLInput() *SDLInput {
	return &SDLInput{
		bindings:     NewBindings(),
		keyCallbacks: make(map[Keycode][]func()),
		keyboardState: sdl.GetKeyboardState(),
		prevKeyboardState: make([]uint8, sdl.NUM_SCANCODES),
	}
}

// PollEvents polls for SDL events and updates input states.
func (s *SDLInput) PollEvents(waitForRefresh bool) {
	// Update previous states
	copy(s.prevKeyboardState, s.keyboardState)
	s.prevMouseState = s.mouseState
	s.scrollValue = 0 // Reset scroll for current frame
	s.codepoints = []rune{} // Clear codepoints for current frame
	
	// Reset JustChanged for all bindings
	for _, b := range s.bindings.GetAll() {
		b.JustChanged = false
	}

	// Process all pending SDL events
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			// Handled by engine
		case *sdl.KeyboardEvent:
			keycode := Keycode(e.Keysym.Sym)
			isPressed := (e.State == sdl.PRESSED)

			// Update internal state
			if isPressed {
				// Prevent double counting if key is held down and only event is a repeat
				if e.Repeat == 0 { 
					s.keyboardState[e.Keysym.Scancode] = 1 // Use scancode for internal state tracking
				}
			} else {
				s.keyboardState[e.Keysym.Scancode] = 0
			}
			
			// Update bindings
			for _, binding := range s.bindings.GetAll() {
				if binding.InputType == InputTypeKeyboard && Keycode(binding.Code) == keycode {
					if isPressed != binding.State {
						binding.JustChanged = true
						binding.State = isPressed
						if isPressed {
							// Trigger callback if binding is active
							// if binding.onactived != nil {
							//     binding.onactived.Call() // Simplified
							// }
						}
					}
				}
			}

			// Trigger key callbacks
			if isPressed && e.Repeat == 0 {
				if callbacks, ok := s.keyCallbacks[keycode]; ok {
					for _, cb := range callbacks {
						cb()
					}
				}
			}

		case *sdl.MouseButtonEvent:
			mousecode := Mousecode(e.Button)
			isPressed := (e.State == sdl.PRESSED)

			// Update internal state
			if isPressed {
				s.mouseState |= sdl.BUTTON(mousecode)
			} else {
				s.mouseState &= ^sdl.BUTTON(mousecode)
			}

			// Update bindings
			for _, binding := range s.bindings.GetAll() {
				if binding.InputType == InputTypeMouse && Mousecode(binding.Code) == mousecode {
					if isPressed != binding.State {
						binding.JustChanged = true
						binding.State = isPressed
						if isPressed {
							// Trigger callback if binding is active
							// if binding.onactived != nil {
							//     binding.onactived.Call() // Simplified
							// }
						}
					}
				}
			}

		case *sdl.MouseMotionEvent:
			s.cursorState.Pos.X = float32(e.X)
			s.cursorState.Pos.Y = float32(e.Y)
			s.cursorState.Delta.X = float32(e.XRel)
			s.cursorState.Delta.Y = float32(e.YRel)

		case *sdl.MouseWheelEvent:
			s.scrollValue = int(e.Y) // e.Y is vertical scroll
			
		case *sdl.TextEditingEvent:
			// Handle text editing if needed
		case *sdl.TextInputEvent:
			// Collect codepoints for text input
			for _, r := range e.Text {
				s.codepoints = append(s.codepoints, r)
			}
		}
	}
	
	// Update pressedKeys for GetPressedKeys()
	s.pressedKeys = s.pressedKeys[:0] // Clear slice
	for scancode := sdl.SCANCODE_UNKNOWN + 1; scancode < sdl.NUM_SCANCODES; scancode++ {
		if s.keyboardState[scancode] == 1 {
			s.pressedKeys = append(s.pressedKeys, sdl.GetKeyFromScancode(scancode))
		}
	}
}

// GetClipboardText returns the current text on the clipboard.
func (s *SDLInput) GetClipboardText() string {
	text, err := sdl.GetClipboardText()
	if err != nil {
		log.Printf("Error getting clipboard text: %v", err)
		return ""
	}
	return text
}

// SetClipboardText sets the clipboard text.
func (s *SDLInput) SetClipboardText(text string) {
	err := sdl.SetClipboardText(text)
	if err != nil {
		log.Printf("Error setting clipboard text: %v", err)
	}
}

// GetScroll returns the scroll value for the current frame.
func (s *SDLInput) GetScroll() int {
	return s.scrollValue
}

// Pressed checks if a key is currently pressed.
func (s *SDLInput) Pressed(keycode Keycode) bool {
	scancode := sdl.GetScancodeFromKey(keycode)
	return s.keyboardState[scancode] == 1
}

// JustPressed checks if a key was pressed in the current frame.
func (s *SDLInput) JustPressed(keycode Keycode) bool {
	scancode := sdl.GetScancodeFromKey(keycode)
	return s.keyboardState[scancode] == 1 && s.prevKeyboardState[scancode] == 0
}

// Clicked checks if a mouse button is currently pressed.
func (s *SDLInput) Clicked(mousecode Mousecode) bool {
	return (s.mouseState & sdl.BUTTON(mousecode)) != 0
}

// JustClicked checks if a mouse button was clicked in the current frame.
func (s *SDLInput) JustClicked(mousecode Mousecode) bool {
	return (s.mouseState&sdl.BUTTON(mousecode)) != 0 && (s.prevMouseState&sdl.BUTTON(mousecode)) == 0
}

// GetCursor returns the current cursor state.
func (s *SDLInput) GetCursor() CursorState {
	return s.cursorState
}

// IsCursorLocked returns true if the cursor is locked.
func (s *SDLInput) IsCursorLocked() bool {
	return s.cursorState.Locked
}

// ToggleCursor toggles the cursor lock state.
func (s *SDLInput) ToggleCursor() {
	s.cursorState.Locked = !s.cursorState.Locked
	if s.cursorState.Locked {
		sdl.SetRelativeMouseMode(true)
	} else {
		sdl.SetRelativeMouseMode(false)
	}
}

// GetBindings returns the Bindings manager.
func (s *SDLInput) GetBindings() *Bindings {
	return s.bindings
}

// AddKeyCallback adds a callback for a specific key.
func (s *SDLInput) AddKeyCallback(key Keycode, callback func()) {
	s.keyCallbacks[key] = append(s.keyCallbacks[key], callback)
}

// GetPressedKeys returns a slice of currently pressed keys.
func (s *SDLInput) GetPressedKeys() []Keycode {
	return s.pressedKeys
}

// GetCodepoints returns a slice of codepoints (text input) from the current frame.
func (s *SDLInput) GetCodepoints() []rune {
	return s.codepoints
}
