package window

import "github.com/Zyko0/go-sdl3/sdl"

// Window represents an abstract window.
type Window struct {
	sdlWindow   *sdl.Window
	sdlRenderer *sdl.Renderer
	// Add other window-related fields here
}

// NewSDLWindow creates a new Window abstraction around SDL window and renderer.
func NewSDLWindow(sdlWindow *sdl.Window, sdlRenderer *sdl.Renderer) *Window {
	return &Window{
		sdlWindow:   sdlWindow,
		sdlRenderer: sdlRenderer,
	}
}

// GetWindow returns the underlying SDL window.
func (w *Window) GetWindow() *sdl.Window {
	return w.sdlWindow
}

// GetRenderer returns the underlying SDL renderer.
func (w *Window) GetRenderer() *sdl.Renderer {
	return w.sdlRenderer
}

// Additional methods will be added as C++ Window virtual methods are ported.