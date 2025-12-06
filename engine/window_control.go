package engine

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Zyko0/go-sdl3/sdl"
	"goxelcore/window"
	// "goxelcore/io" // Assuming io package exists for file operations
)

// VC_BUILD_NAME and ENGINE_DEBUG_BUILD would be Go build tags or constants.
// For now, let's define them as constants for demonstration.
const (
	VC_BUILD_NAME = "" // Example: "dev"
	ENGINE_VERSION_MAJOR = 0
	ENGINE_VERSION_MINOR = 1
	ENGINE_DEBUG_BUILD = true // Set to false for release build
)

// WindowControl corresponds to C++ WindowControl class in voxelcore/src/engine/WindowControl.hpp
type WindowControl struct {
	engine *Engine // Use a pointer to the Engine instance
}

// NewWindowControl creates a new WindowControl instance.
func NewWindowControl(eng *Engine) *WindowControl {
	return &WindowControl{
		engine: eng,
	}
}

// Initialize creates and initializes the Window and Input instances.
// Corresponds to C++ WindowControl::initialize()
func (wc *WindowControl) Initialize() (*window.Window, window.Input, error) {
	project := wc.engine.GetProject() // Assuming GetProject() exists
	settings := wc.engine.GetSettings()

	title := project.Title // Assuming Project struct has a Title field
	if title != "" {
		title += " - "
	}

	buildName := VC_BUILD_NAME
	if buildName == "" {
		title += "VoxelCore v" + strconv.Itoa(ENGINE_VERSION_MAJOR) + "." + strconv.Itoa(ENGINE_VERSION_MINOR)
	} else {
		title += buildName
	}
	if ENGINE_DEBUG_BUILD {
		title += " [debug]"
	}
	// TODO: Check if debuggingServer exists
	// if wc.engine.GetDebuggingServer() != nil {
	// 	title = "[debugging] " + title
	// }

	// For now, using the display settings directly from EngineSettings
	// C++ uses DisplaySettings* settings and a title.
	displaySettings := settings.Display
	width := displaySettings.Width.Value
	height := displaySettings.Height.Value

	// Create SDL window with OpenGL flag
	sdlWindow, err := sdl.CreateWindow(title, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SDL window: %w", err)
	}
	// Create OpenGL context
	glContext, err := sdlWindow.GLCreateContext()
	if err != nil {
		sdlWindow.Destroy()
		return nil, nil, fmt.Errorf("failed to create OpenGL context: %w", err)
	}
	// Make the context current
	if err := sdlWindow.GLMakeCurrent(glContext); err != nil {
		sdl.GLDeleteContext(glContext)
		sdlWindow.Destroy()
		return nil, nil, fmt.Errorf("failed to make OpenGL context current: %w", err)
	}
	// Initialize GL bindings - This would typically be done by SDL automatically
	// In SDL3, OpenGL functions are available after creating the context


	w := window.NewSDLWindowWithGL(sdlWindow, glContext)
	inp := window.NewSDLInput()

	w.SetFramerate(settings.Display.Framerate.Value)
	
	// Load icon (stub for now)
	iconData := wc.loadIcon()
	if iconData != nil {
		w.SetIcon(iconData)
	}

	return w, inp, nil
}

// loadIcon corresponds to C++ load_icon() function.
// This is a stub for now.
func (wc *WindowControl) loadIcon() interface{} { // Returns ImageData*
	log.Println("WindowControl.loadIcon: Not implemented yet (requires ImageData and imageio porting)")
	// In a real implementation:
	// file := "res:textures/misc/icon.png"
	// if io.Exists(file) {
	// 	if imageData, err := imageio.Read(file); err == nil {
	// 		imageData.FlipY()
	// 		return imageData
	// 	}
	// }
	return nil
}

// NextFrame performs per-frame updates for the window and input.
// Corresponds to C++ WindowControl::nextFrame()
func (wc *WindowControl) NextFrame(waitForRefresh bool) {
	settings := wc.engine.GetSettings()
	win := wc.engine.GetWindow()
	inp := wc.engine.GetInput()

	targetFramerate := settings.Display.Framerate.Value
	if win.IsIconified() && settings.Display.LimitFpsIconified.Value {
		targetFramerate = 20
	}
	win.SetFramerate(targetFramerate)
	
	win.SwapBuffers() // This calls GLSwapWindow now
	inp.PollEvents(waitForRefresh && !win.CheckShouldRefresh())
}

// SaveScreenshot captures a screenshot of the window.
// Corresponds to C++ WindowControl::saveScreenshot()
func (wc *WindowControl) SaveScreenshot() {
	win := wc.engine.GetWindow()
	//paths := wc.engine.GetPaths() // Assuming GetPaths() exists

	log.Println("WindowControl.SaveScreenshot: Not implemented yet (requires ImageData and imageio porting)")
	// In a real implementation:
	// imageData := win.TakeScreenshot()
	// if imageData != nil {
	// 	imageData.FlipY()
	// 	filename := paths.GetNewScreenshotFile("png") // Assuming this method exists
	// 	if err := imageio.Write(filename, imageData); err == nil {
	// 		log.Printf("saved screenshot as %s\n", filename)
	// 	} else {
	// 		log.Printf("failed to save screenshot: %v\n", err)
	// 	}
	// }
}

// ToggleFullscreen toggles the window's fullscreen mode.
// Corresponds to C++ WindowControl::toggleFullscreen()
func (wc *WindowControl) ToggleFullscreen() {
	settings := wc.engine.GetSettings()
	windowModeSetting := settings.Display.WindowMode
	
	if windowModeSetting.Value != int(window.WindowModeFullscreen) {
		windowModeSetting.Value = int(window.WindowModeFullscreen)
	} else {
		windowModeSetting.Value = int(window.WindowModeWindowed)
	}
	wc.engine.GetWindow().SetMode(window.WindowMode(windowModeSetting.Value))
}
