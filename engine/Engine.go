package engine

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time" // Import Go's time package

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
	project      *Project       // Add Project field
	windowControl *WindowControl // Add WindowControl field
	paths        *EnginePaths   // Add EnginePaths field
	time         Time           // Add Time component
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
			project:  NewProject(), // Initialize Project
			time:     Time{},       // Initialize Time
			// paths will be initialized in Initialize method as it depends on CoreParameters
		}
	})
	return engineInstance
}

// Initialize initializes the engine with the given core parameters.
// Corresponds to C++ Engine::initialize(CoreParameters coreParameters)
func (e *Engine) Initialize(coreParameters *CoreParameters) error {
	e.params = coreParameters // Update with provided parameters

	log.Println("Engine: Initializing...")

	// Initialize EnginePaths
	e.paths = NewEnginePaths(e.params) // Initialize paths

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

	// Initialize WindowControl
	e.windowControl = NewWindowControl(e)

	// Use WindowControl to initialize Window and Input
	win, inp, err := e.windowControl.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize window and input: %w", err)
	}
	e.window = win
	e.input = inp
	e.glContext = win.GetGLContext() // Get GLContext from the initialized window

	// Initialize time
	e.time.Set(float64(time.Now().UnixNano()) / float64(time.Second))

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

		// Update time
		e.time.Update(float64(time.Now().UnixNano()) / float64(time.Second))

		// Poll for events using our Input component
		e.windowControl.NextFrame(false) // Use WindowControl's nextFrame
		
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

// GetProject returns the Project instance.
// Corresponds to C++ Engine::getProject()
func (e *Engine) GetProject() *Project {
	return e.project
}

// GetPaths returns the EnginePaths instance.
// Corresponds to C++ Engine::getPaths()
func (e *Engine) GetPaths() *EnginePaths {
	return e.paths
}

// GetTime returns the Time instance.
// Corresponds to C++ Time::getTime()
func (e *Engine) GetTime() *Time {
	return &e.time
}