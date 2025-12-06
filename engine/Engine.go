package engine

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	// For now, using the display settings directly from EngineSettings
	// C++ uses DisplaySettings* settings and a title.
	displaySettings := e.settings.Display
	width := displaySettings.Width.Value
	height := displaySettings.Height.Value
	title := "GoxelCore" // Default title, later can be derived from CoreParameters or other settings

	// Initialize the Window and Input components
	// The C++ version returns a tuple of Window and Input
	// For now, let's create a basic SDL window directly. This will be refactored
	// to use a dedicated window package later.
	sdlWindow, sdlRenderer, err := sdl.CreateWindowAndRenderer(title, width, height, 0)
	if err != nil {
		return fmt.Errorf("failed to create SDL window and renderer: %w", err)
	}
	// Defer destroy actions until Terminate is called
	// defer sdlRenderer.Destroy()
	// defer sdlWindow.Destroy()

	e.window = window.NewSDLWindow(sdlWindow, sdlRenderer) // Placeholder for a Go window abstraction

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

		var event sdl.Event
		for sdl.PollEvent(&event) {
			if event.Type == sdl.EVENT_QUIT {
				e.quit() // Set quit signal
			}
			// Handle other events here
		}

		// Placeholder for applicationTick(), updateFrontend(), renderFrame()
		// For now, just present the renderer
		e.window.GetRenderer().SetDrawColor(255, 255, 255, 255) // White background
		e.window.GetRenderer().Clear()
		e.window.GetRenderer().DebugText(50, 50, "GoxelCore Go Port")
		e.window.GetRenderer().Present()

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
	if engineInstance != nil && engineInstance.window != nil {
		engineInstance.window.GetRenderer().Destroy()
		engineInstance.window.GetWindow().Destroy()
	}
	// For now, directly calling sdl.Quit() and binsdl.Load().Unload()
	// In a more structured port, these might be managed differently.
	sdl.Quit()
	// binsdl.Load().Unload() // This seems to be for loading/unloading the SDL shared library, not a direct C++ equivalent.
	// It's already deferred in goxelcore/main.go in the original example, so I'll keep it there for now or rethink its placement.
	log.Println("Engine: Termination complete.")
}
