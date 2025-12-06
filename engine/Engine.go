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
	"goxelcore" // For Vec2u
	"goxelcore/assets" // Import the assets package
	"goxelcore/content" // Import the content package
	"goxelcore/debug" // Import our new debug package
	"goxelcore/devtools" // Import devtools
	"goxelcore/graphics/core" // For DrawContext
	"goxelcore/io" // For io.Path (needed for logger filename)
	"goxelcore/logic" // Import the logic package for EngineController
	"goxelcore/window" // Assuming a window package will be created
	"goxelcore/graphics/ui" // Import ui package
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
	logger       *debug.Logger  // Add Logger component
	controller   *logic.EngineController // Add EngineController component
	assets       *assets.Assets // Add Assets component
	assetsLoader *assets.AssetsLoader // Add AssetsLoader component
	content      *content.ContentControl // Add ContentControl component
	editor       *devtools.Editor // Add editor field
	gui          *ui.GUI          // Add GUI component
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
			logger:   debug.NewLogger("engine"), // Initialize Logger
			assets:   assets.NewAssets(), // Initialize Assets
			// paths, controller, assetsLoader, content, editor, gui will be initialized in Initialize method as they depend on the engine itself
		}
	})
	return engineInstance
}

// Initialize initializes the engine with the given core parameters.
// Corresponds to C++ Engine::initialize(CoreParameters coreParameters)
func (e *Engine) Initialize(coreParameters *CoreParameters) error {
	e.params = coreParameters // Update with provided parameters

	e.logger.Info("Engine: Initializing...")

	// Initialize EnginePaths
	e.paths = NewEnginePaths(e.params) // Initialize paths
	
	// Initialize debug logger to file
	debug.Init(e.paths.GetUserFilesFolder() + "/latest.log")
	e.logger.Info("Logger initialized to file: %s/latest.log", e.paths.GetUserFilesFolder())


	// Initialize SDL Video for window creation
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		e.logger.Error("Failed to initialize SDL: %v", err)
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
		e.logger.Error("Failed to initialize window and input: %v", err)
		return fmt.Errorf("failed to initialize window and input: %w", err)
	}
	e.window = win
	e.input = inp
	e.glContext = win.GetGLContext() // Get GLContext from the initialized window

	// Initialize Time
	e.time.Set(float64(time.Now().UnixNano()) / float64(time.Second))

	// Initialize EngineController
	e.controller = logic.NewEngineController(e)

	// Initialize AssetsLoader
	e.assetsLoader = assets.NewAssetsLoader(e, e.assets, &e.paths.ResPaths)

	// Initialize Editor
	e.editor = devtools.NewEditor(e) // Initialize editor

	// Initialize GUI
	e.gui = ui.NewGUI(e)

	// Define postContent callback
	postContentCallback := func() {
		e.logger.Info("Engine: Post content load callback triggered (stub).")
		// In C++, Assets::setup() is called here
		e.assets.Setup()
		e.editor.LoadTools() // Call editor's LoadTools
		e.gui.OnAssetsLoad(e.assets) // Call GUI's OnAssetsLoad
	}

	// Initialize ContentControl
	e.content = content.NewContentControl(e.project, e.paths, e.input.(*window.SDLInput), postContentCallback)


	e.logger.Info("Engine: Initialization complete.")
	return nil
}

// Run starts the main engine loop.
// Corresponds to C++ Engine::run()
func (e *Engine) Run() {
	e.logger.Info("Engine: Starting run loop...")

	// Setup signal handler for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) // SIGTERM corresponds to C++ std::signal(SIGTERM, ...)

	go func() {
		<-sigChan
		e.logger.Info("Engine: Received termination signal. Quitting...")
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
			e.logger.Info("Escape key pressed. Quitting...")
			e.quit()
		}

		// In C++, the quit event is often handled in the main event loop
		// The sdl.PollEvent above in SDLInput.PollEvents() should capture sdl.EVENT_QUIT
		// and it will set the quitSignal. So this check should be enough.

		// GUI Act (Update logic)
		w, h := e.window.GetSize()
		e.gui.Act(float32(e.time.GetDelta()), goxelcore.Vec2u{X: uint32(w), Y: uint32(h)})

		// Placeholder for applicationTick(), updateFrontend(), renderFrame()
		// For now, just present the renderer
		gl.ClearColor(0.2, 0.3, 0.3, 1.0) // Example clear color
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// GUI Draw
		// C++ uses DrawContext pctx(nullptr, *window, nullptr);
		// For now, let's create a dummy DrawContext
		// We pass e.gui.GetContainer().batch2D. This requires Batch2D to be public.
		dummyDrawContext := core.NewDrawContext(nil, e.window, e.gui.GetContainer().GetBatch2D()) 
		e.gui.Draw(*dummyDrawContext, e.assets)

		// GUI PostAct
		e.gui.PostAct()


		e.window.GetWindow().GLSwapWindow() // Swap buffers
		
		return nil
	})

	e.logger.Info("Engine: Run loop ended.")
}

// quit sets the internal quit signal.
// Corresponds to C++ Engine::quit()
func (e *Engine) quit() {
	e.quitSignal = true
}

// Terminate performs cleanup before the engine exits.
// Corresponds to C++ Engine::terminate()
func Terminate() {
	e := GetInstance() // Get the instance to access the logger
	e.logger.Info("Engine: Terminating...")
	debug.Flush() // Flush logs before closing

	if e != nil {
		if e.assets != nil {
			e.assets.Delete() // Clean up assets
		}
		if e.window != nil {
			// Destroy GL context
			if e.glContext != nil {
				sdl.GLDeleteContext(e.glContext)
			}
			// Destroy window
			e.window.GetWindow().Destroy()
		}
		// Additional cleanup for input or other components if necessary
	}
	sdl.Quit()
	e.logger.Info("Engine: Termination complete.")
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

// GetController returns the EngineController instance.
// Corresponds to C++ Engine::getController()
func (e *Engine) GetController() *logic.EngineController {
	return e.controller
}

// GetAssets returns the Assets instance.
// Corresponds to C++ Engine::getAssets()
func (e *Engine) GetAssets() *assets.Assets {
	return e.assets
}

// GetAssetsLoader returns the AssetsLoader instance.
func (e *Engine) GetAssetsLoader() *assets.AssetsLoader {
	return e.assetsLoader
}

// GetContentControl returns the ContentControl instance.
// Corresponds to C++ Engine::getContentControl()
func (e *Engine) GetContentControl() *content.ContentControl {
	return e.content
}

// GetEditor returns the Editor instance.
// Corresponds to C++ GUI::getEditor() and Engine::getEditor()
func (e *Engine) GetEditor() *devtools.Editor {
	return e.editor
}

// GetGUI returns the GUI instance.
// Corresponds to C++ Engine::getGUI()
func (e *Engine) GetGUI() *ui.GUI {
	return e.gui
}