package window

import (
	"fmt" // For potential error formatting
	"github.com/Zyko0/go-sdl3/sdl"
)

// WindowMode corresponds to C++ WindowMode enum in voxelcore/src/window/Window.hpp
type WindowMode int

const (
	WindowModeWindowed WindowMode = iota
	WindowModeFullscreen
	WindowModeBorderless
)

// CursorShape corresponds to C++ CursorShape enum in voxelcore/src/graphics/core/commons.hpp
// Map to SDL cursors where possible.
type CursorShape int

const (
	CursorShapeArrow CursorShape = iota // sdl.SYSTEM_CURSOR_ARROW
	CursorShapeText                     // sdl.SYSTEM_CURSOR_IBEAM
	CursorShapeCrosshair                // sdl.SYSTEM_CURSOR_CROSSHAIR
	CursorShapePointer                  // sdl.SYSTEM_CURSOR_HAND
	CursorShapeEWResize                 // sdl.SYSTEM_CURSOR_SIZEWE
	CursorShapeNSResize                 // sdl.SYSTEM_CURSOR_SIZENS
	CursorShapeNWSEResize               // sdl.SYSTEM_CURSOR_SIZENWSE
	CursorShapeNESWResize               // sdl.SYSTEM_CURSOR_SIZENESW
	CursorShapeAllResize                // sdl.SYSTEM_CURSOR_SIZEALL
	CursorShapeNotAllowed               // sdl.SYSTEM_CURSOR_NO
)

// Window represents an abstract window.
type Window struct {
	sdlWindow   *sdl.Window
	sdlRenderer *sdl.Renderer // Will be nil if using OpenGL
	glContext   sdl.GLContext // Add GL context
	size        struct{ X, Y int } // Corresponds to glm::ivec2 size
	mode        WindowMode
	shouldRefresh bool
}

// NewSDLWindow creates a new Window abstraction around SDL window and renderer (for 2D rendering).
// Keeping this for compatibility or if there are 2D parts.
func NewSDLWindow(sdlWindow *sdl.Window, sdlRenderer *sdl.Renderer) *Window {
	w, h := sdlWindow.Size()
	return &Window{
		sdlWindow:   sdlWindow,
		sdlRenderer: sdlRenderer,
		size:        struct{ X, Y int }{X: int(w), Y: int(h)},
		mode:        WindowModeWindowed, // Default for now
		shouldRefresh: false,
	}
}

// NewSDLWindowWithGL creates a new Window abstraction with an OpenGL context.
func NewSDLWindowWithGL(sdlWindow *sdl.Window, glContext sdl.GLContext) *Window {
	w, h := sdlWindow.Size()
	return &Window{
		sdlWindow:   sdlWindow,
		glContext:   glContext,
		size:        struct{ X, Y int }{X: int(w), Y: int(h)},
		mode:        WindowModeWindowed, // Default for now
		shouldRefresh: false,
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

// GetGLContext returns the OpenGL context.
func (w *Window) GetGLContext() sdl.GLContext {
	return w.glContext
}

// GetSize returns the window size.
func (w *Window) GetSize() (int, int) {
	return w.size.X, w.size.Y
}

// SwapBuffers presents the renderer's buffer to the screen.
// Corresponds to C++ Window::swapBuffers()
func (w *Window) SwapBuffers() {
	// For OpenGL, this means swapping the window's buffers.
	w.sdlWindow.GLSwap()
}

// IsMaximized returns true if the window is maximized.
// Corresponds to C++ Window::isMaximized()
func (w *Window) IsMaximized() bool {
	return (w.sdlWindow.GetFlags() & sdl.WINDOW_MAXIMIZED) != 0
}

// IsFocused returns true if the window has keyboard focus.
// Corresponds to C++ Window::isFocused()
func (w *Window) IsFocused() bool {
	return (w.sdlWindow.GetFlags() & sdl.WINDOW_INPUT_FOCUS) != 0
}

// IsIconified returns true if the window is minimized.
// Corresponds to C++ Window::isIconified()
func (w *Window) IsIconified() bool {
	return (w.sdlWindow.GetFlags() & sdl.WINDOW_MINIMIZED) != 0
}

// IsShouldClose returns true if the window should close.
// Corresponds to C++ Window::isShouldClose()
func (w *Window) IsShouldClose() bool {
	// This is typically handled by polling SDL_QUIT event.
	// The Engine's Run loop handles the quit signal, so this can return that signal.
	return false // Engine's quitSignal manages this
}

// SetShouldClose requests the window to close.
// Corresponds to C++ Window::setShouldClose()
func (w *Window) SetShouldClose(flag bool) {
	// In SDL, we usually don't explicitly set "should close" on the window.
	// Instead, we emit an SDL_QUIT event or set an internal engine quit flag.
	// For now, it can be a no-op or set an internal flag.
}

// SetCursor sets the cursor shape.
// Corresponds to C++ Window::setCursor()
func (w *Window) SetCursor(shape CursorShape) {
	var sdlCursor *sdl.Cursor
	var err error // Declare error variable

	switch shape {
	case CursorShapeArrow:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)
	case CursorShapeText:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM)
	case CursorShapeCrosshair:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_CROSSHAIR)
	case CursorShapePointer:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND)
	case CursorShapeEWResize:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZEWE)
	case CursorShapeNSResize:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENS)
	case CursorShapeNWSEResize:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENWSE)
	case CursorShapeNESWResize:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENESW)
	case CursorShapeAllResize:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZEALL)
	case CursorShapeNotAllowed:
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_NO)
	default:
		// Default to arrow or log error
		sdlCursor, err = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW) // Handle error here too
	}

	if err != nil { // Check for error after all assignments
		fmt.Printf("Failed to create system cursor %v: %v\n", shape, err)
		return // Exit if cursor creation failed
	}

	if sdlCursor != nil {
		sdl.SetCursor(sdlCursor)
		// sdlCursor.Destroy() // Cursors created with CreateSystemCursor do not need to be freed manually.
	} else {
		fmt.Printf("Failed to set cursor shape %v (cursor was nil but no error reported)\n", shape)
	}
}

// SetMode sets the window mode (windowed, fullscreen, borderless).
// Corresponds to C++ Window::setMode()
func (w *Window) SetMode(mode WindowMode) {
	var sdlFullscreenFlag uint32
	switch mode {
	case WindowModeFullscreen:
		sdlFullscreenFlag = sdl.WINDOW_FULLSCREEN
	case WindowModeBorderless:
		sdlFullscreenFlag = sdl.WINDOW_FULLSCREEN_DESKTOP // This is borderless fullscreen
	case WindowModeWindowed:
		sdlFullscreenFlag = 0 // Not fullscreen
	}
	err := w.sdlWindow.SetFullscreen(sdlFullscreenFlag)
	if err != nil {
		fmt.Printf("Failed to set window mode to %v: %v\n", mode, err)
	} else {
		w.mode = mode
	}
}

// GetMode returns the current window mode.
// Corresponds to C++ Window::getMode()
func (w *Window) GetMode() WindowMode {
	return w.mode
}

// Focus brings the window to the foreground and gives it input focus.
// Corresponds to C++ Window::focus()
func (w *Window) Focus() {
	err := w.sdlWindow.Raise()
	if err != nil {
		fmt.Printf("Failed to raise window: %v\n", err)
	}
}

// SetTitle sets the window title.
// Corresponds to C++ Window::setTitle()
func (w *Window) SetTitle(title string) {
	w.sdlWindow.SetTitle(title)
}

// SetIcon sets the window icon.
// Corresponds to C++ Window::setIcon(const ImageData* image)
func (w *Window) SetIcon(imageData interface{}) { // ImageData is a complex type, use interface{} for now
	// This requires porting ImageData struct and loading image data.
	// For now, this is a stub.
	_ = imageData // Suppress unused warning
	fmt.Println("Window.SetIcon: Not implemented yet")
}

// PushScissor pushes a new scissor rectangle.
// Corresponds to C++ Window::pushScissor()
func (w *Window) PushScissor(area interface{}) { // glm::vec4 area, use interface{} for now
	// Requires graphics context implementation. Stub for now.
	_ = area
	fmt.Println("Window.PushScissor: Not implemented yet")
}

// PopScissor pops the last scissor rectangle.
// Corresponds to C++ Window::popScissor()
func (w *Window) PopScissor() {
	// Requires graphics context implementation. Stub for now.
	fmt.Println("Window.PopScissor: Not implemented yet")
}

// ResetScissor resets the scissor rectangle to cover the entire window.
// Corresponds to C++ Window::resetScissor()
func (w *Window) ResetScissor() {
	// Requires graphics context implementation. Stub for now.
	fmt.Println("Window.ResetScissor: Not implemented yet")
}

// SetShouldRefresh sets a flag indicating the window content should be redrawn.
// Corresponds to C++ Window::setShouldRefresh()
func (w *Window) SetShouldRefresh() {
	w.shouldRefresh = true
}

// CheckShouldRefresh checks if the refresh flag is set and clears it.
// Corresponds to C++ Window::checkShouldRefresh()
func (w *Window) CheckShouldRefresh() bool {
	refresh := w.shouldRefresh
	w.shouldRefresh = false
	return refresh
}

// Time returns the current time in seconds (usually since window creation).
// Corresponds to C++ Window::time()
func (w *Window) Time() float64 {
	// SDL_GetPerformanceCounter() / SDL_GetPerformanceFrequency() can provide high-res time.
	// Or just time.Since(startTime).Seconds()
	return float64(sdl.GetPerformanceCounter()) / float64(sdl.GetPerformanceFrequency())
}

// SetFramerate sets the target framerate.
// Corresponds to C++ Window::setFramerate()
func (w *Window) SetFramerate(framerate int) {
	// This usually involves VSync or manual frame pacing.
	// SDL_GL_SetSwapInterval(1) for VSync.
	// For now, just a stub.
	fmt.Printf("Window.SetFramerate(%d): Not fully implemented (VSync control needed).\n", framerate)
	if framerate > 0 {
		// Can set a swap interval if using GL. For now, SDL renderer does not directly control framerate this way.
		// sdl.GLSetSwapInterval(1) for VSync (often used for framerate limiting).
	}
}

// TakeScreenshot captures the current window content.
// Corresponds to C++ Window::takeScreenshot()
func (w *Window) TakeScreenshot() interface{} { // Returns std::unique_ptr<ImageData>
	// Requires ImageData porting and implementation. Stub for now.
	fmt.Println("Window.TakeScreenshot: Not implemented yet")
	return nil
}

// display namespace functions.
// These are not methods of the Window class, but global display utilities.
// We can place them as package-level functions in the window package.

// Clear clears the current rendering target.
// Corresponds to C++ display::clear()
func Clear() {
	// Assuming this clears the currently bound renderer.
	// This might be more appropriate in a graphics package.
	// For now, if called, it might clear the default engine renderer.
	fmt.Println("display.Clear: Not implemented directly on default renderer.")
}

// ClearDepth clears the depth buffer.
// Corresponds to C++ display::clearDepth()
func ClearDepth() {
	// Requires graphics context. Stub for now.
	fmt.Println("display.ClearDepth: Not implemented yet")
}

// SetBgColor sets the background color (vec3).
// Corresponds to C++ display::setBgColor(glm::vec3 color)
func SetBgColor(r, g, b float32) { // Using float32 for R, G, B
	// Requires graphics context. Stub for now.
	fmt.Printf("display.SetBgColor(%f, %f, %f): Not implemented yet\n", r, g, b)
}

// SetBgColorRGBA sets the background color (vec4).
// Corresponds to C++ display::setBgColor(glm::vec4 color)
func SetBgColorRGBA(r, g, b, a float32) { // Using float32 for R, G, B, A
	// Requires graphics context. Stub for now.
	fmt.Printf("display.SetBgColorRGBA(%f, %f, %f, %f): Not implemented yet\n", r, g, b, a)
}