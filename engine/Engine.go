package engine

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Zyko0/go-sdl3/gl" // Import go-sdl3/gl
	"github.com/Zyko0/go-sdl3/sdl"
	"goxelcore/window" // Assuming a window package will be created
)

// Engine is the main engine struct, implementing a singleton pattern.
// Corresponds to C++ Engine class in voxelcore/src/engine/Engine.hpp
type Engine struct {
	params       *CoreParameters
	settings     *EngineSettings
	quitSignal   bool
	window       *window.Window // Placeholder for the Window component
	input        window.Input   // Add input field
	glContext    sdl.GLContext  // Add OpenGL context
	// Add other fields as needed, mirroring C++ Engine.hpp
}

var (
	engineInstance *Engine
	engineOnce     sync.Once
)

// GetInstance returns the singleton instance of the Engine.
// Corresponds to C++ Engine::getInstance()
func GetInstance() *Engine {
	engineOnce.Do(func() {
		engineInstance = &Engine{
			params:   NewCoreParameters(),
			settings: NewEngineSettings(),
		}
	})
	return engineInstance
}

// Initialize initializes the engine with the given core parameters.
// Corresponds to C++ Engine::initialize(CoreParameters coreParameters)
func (e *Engine) Initialize(coreParameters *CoreParameters) error {
	e.params = coreParameters // Update with provided parameters

	log.Println("Engine: Initializing...")

	// Initialize SDL Video for window creation
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fmt.Errorf("failed to initialize SDL: %w", err)
	}

	// Set OpenGL attributes
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)

	// For now, using the display settings directly from EngineSettings
	// C++ uses DisplaySettings* settings and a title.
	displaySettings := e.settings.Display
	width := displaySettings.Width.Value
	height := displaySettings.Height.Value
	title := "GoxelCore" // Default title, later can be derived from CoreParameters or other settings

	// Create SDL window with OpenGL flag
	sdlWindow, err := sdl.CreateWindow(title, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		return fmt.Errorf("failed to create SDL window: %w", err)
	}
	// No sdlRenderer for OpenGL context, use sdlWindow directly for GLContext

	// Create OpenGL context
	glContext, err := sdlWindow.GLCreateContext()
	if err != nil {
		return fmt.Errorf("failed to create OpenGL context: %w", err)
	}
	e.glContext = glContext

	// Make the context current
	if err := sdlWindow.GLMakeCurrent(glContext); err != nil {
		return fmt.Errorf("failed to make OpenGL context current: %w", err)
	}

	// Initialize GL bindings
	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize GL bindings: %w", err)
	}

	// Create our window abstraction
	e.window = window.NewSDLWindowWithGL(sdlWindow, glContext) // New constructor for window.Window
	
	// Initialize Input component
	e.input = window.NewSDLInput()

	log.Println("Engine: Initialization complete.")
	return nil
}

// Run starts the main engine loop.
// Corresponds to C++ Engine::run()
func (e *Engine) Run() {
	log.Println("Engine: Starting run loop...")

	// Setup signal handler for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) // SIGTERM corresponds to C++ std::signal(SIGTERM, ...)

	go func() {
		<-sigChan
		log.Println("Engine: Received termination signal. Quitting...")
		e.quit() // Call the engine's quit method
	}()

	sdl.RunLoop(func() error {
		if e.quitSignal {
			return sdl.EndLoop
		}

		// Poll for events using our Input component
		e.input.PollEvents(false) // Assuming waitForRefresh = false for now
		
		// Check for quit event (e.g., window close button)
		// For example, if SDLK_ESCAPE is pressed or a binding named "quit" is active
		if e.input.JustPressed(window.Keycode(sdl.K_ESCAPE)) {
			log.Println("Escape key pressed. Quitting...")
			e.quit()
		}

		// In C++, the quit event is often handled in the main event loop
		// The sdl.PollEvent above in SDLInput.PollEvents() should capture sdl.EVENT_QUIT
		// and it will set the quitSignal. So this check should be enough.

		// Placeholder for applicationTick(), updateFrontend(), renderFrame()
		// For now, just present the renderer
		gl.ClearColor(0.2, 0.3, 0.3, 1.0) // Example clear color
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		e.window.GetWindow().GLSwapWindow() // Swap buffers
		
		return nil
	})

	log.Println("Engine: Run loop ended.")
}

// quit sets the internal quit signal.
// Corresponds to C++ Engine::quit()
func (e *Engine) quit() {
	e.quitSignal = true
}

// Terminate performs cleanup before the engine exits.
// Corresponds to C++ Engine::terminate()
func Terminate() {
	log.Println("Engine: Terminating...")
	if engineInstance != nil {
		if engineInstance.window != nil {
			// Destroy GL context
			if engineInstance.glContext != nil {
				sdl.GLDeleteContext(engineInstance.glContext)
			}
			// Destroy window
			engineInstance.window.GetWindow().Destroy()
		}
		// Additional cleanup for input or other components if necessary
	}
	sdl.Quit()
	log.Println("Engine: Termination complete.")
}

// GetInput returns the Input system instance.
// Corresponds to C++ Engine::getInput()
func (e *Engine) GetInput() window.Input {
    return e.input
}

// GetWindow returns the Window instance.
// Corresponds to C++ Engine::getWindow()
func (e *Engine) GetWindow() *window.Window {
    return e.window
}

// GetGLContext returns the OpenGL context.
func (e *Engine) GetGLContext() sdl.GLContext {
	return e.glContext
}