package window

import (
	"fmt"
	"strings"

	"github.com/Zyko0/go-sdl3/sdl"
)

// InputType corresponds to C++ InputType enum in voxelcore/src/window/input.hpp
type InputType int

const (
	InputTypeKeyboard InputType = iota
	InputTypeMouse
)

// Keycode corresponds to C++ Keycode enum in voxelcore/src/window/input.hpp
// We map them to SDL keycodes.
type Keycode = sdl.Keycode

// Mousecode corresponds to C++ Mousecode enum in voxelcore/src/window/input.hpp
// We map them to SDL mouse button codes.
type Mousecode = uint8

// CursorState corresponds to C++ CursorState struct in voxelcore/src/window/input.hpp
type CursorState struct {
	Locked bool
	Pos    struct{ X, Y float32 } // glm::vec2
	Delta  struct{ X, Y float32 } // glm::vec2
}

// Binding corresponds to C++ Binding struct in voxelcore/src/window/input.hpp
type Binding struct {
	// onactived util::HandlersList<> // Simplified for now
	InputType InputType
	Code      int // Can be Keycode (sdl.Keycode) or Mousecode (sdl.MouseButtonID)
	State     bool
	JustChanged bool
	Enabled   bool
}

// NewBinding creates a new Binding.
func NewBinding(inputType InputType, code int) *Binding {
	return &Binding{
		InputType: inputType,
		Code:      code,
		State:     false,
		JustChanged: false,
		Enabled:   true,
	}
}

// Active returns true if the binding is currently active.
func (b *Binding) Active() bool {
	return b.State
}

// JActive returns true if the binding became active in the current frame.
func (b *Binding) JActive() bool {
	return b.State && b.JustChanged
}

// Reset resets the binding.
func (b *Binding) Reset(inputType InputType, code int) {
	b.InputType = inputType
	b.Code = code
	b.State = false
	b.JustChanged = false
}

// Text returns a string representation of the binding.
func (b *Binding) Text() string {
	switch b.InputType {
	case InputTypeKeyboard:
		return sdl.GetKeyName(Keycode(b.Code))
	case InputTypeMouse:
		return sdl.GetMouseButtonName(Mousecode(b.Code))
	default:
		return "<unknown input type>"
	}
}

// Bindings corresponds to C++ Bindings class in voxelcore/src/window/input.hpp
type Bindings struct {
	bindings map[string]*Binding
}

// NewBindings creates a new Bindings instance.
func NewBindings() *Bindings {
	return &Bindings{
		bindings: make(map[string]*Binding),
	}
}

// Active checks if a binding by name is active.
func (bs *Bindings) Active(name string) bool {
	if b, ok := bs.bindings[name]; ok {
		return b.Active()
	}
	return false
}

// JActive checks if a binding by name just became active.
func (bs *Bindings) JActive(name string) bool {
	if b, ok := bs.bindings[name]; ok {
		return b.JActive()
	}
	return false
}

// Get returns a binding by name.
func (bs *Bindings) Get(name string) *Binding {
	return bs.bindings[name]
}

// Require returns a binding by name, creating it if it doesn't exist.
// C++ version threw an exception if not found, here we'll panic for now or return error
func (bs *Bindings) Require(name string, inputType InputType, code int) *Binding {
	if b, ok := bs.bindings[name]; ok {
		return b
	}
	b := NewBinding(inputType, code)
	bs.bindings[name] = b
	return b
}


// Bind adds a new binding.
func (bs *Bindings) Bind(name string, inputType InputType, code int) {
	if _, exists := bs.bindings[name]; !exists {
		bs.bindings[name] = NewBinding(inputType, code)
	}
}

// Rebind changes an existing binding.
func (bs *Bindings) Rebind(name string, inputType InputType, code int) {
	if b, ok := bs.bindings[name]; ok {
		b.Reset(inputType, code)
	} else {
		// C++ `require` behavior: create if not exists
		bs.bindings[name] = NewBinding(inputType, code)
	}
}

// GetAll returns all bindings.
func (bs *Bindings) GetAll() map[string]*Binding {
	return bs.bindings
}

// EnableAll enables all bindings.
func (bs *Bindings) EnableAll() {
	for _, b := range bs.bindings {
		b.Enabled = true
	}
}

// Input represents the input handling system.
// Corresponds to C++ Input class in voxelcore/src/window/input.hpp
// This will be an interface with a concrete SDLInput implementation.
type Input interface {
	PollEvents(waitForRefresh bool)
	GetClipboardText() string
	SetClipboardText(text string)
	GetScroll() int
	Pressed(keycode Keycode) bool
	JustPressed(keycode Keycode) bool
	Clicked(mousecode Mousecode) bool
	JustClicked(mousecode Mousecode) bool
	GetCursor() CursorState
	IsCursorLocked() bool
	ToggleCursor()
	GetBindings() *Bindings
	AddKeyCallback(key Keycode, callback func()) // Simplified for now
	GetPressedKeys() []Keycode
	GetCodepoints() []rune // C++ uses uint, Go uses rune for unicode codepoints
}

// input_util namespace in C++ would be functions in Go
// keycode_from, mousecode_from are not directly ported as we use SDL's directly
// to_string functions are implemented in Binding.Text()
