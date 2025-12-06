package ui

import (
	"fmt"
	"log"
	"time"

	"goxelcore"           // For Vec2u
	"goxelcore/assets"     // For Assets
	"goxelcore/engine"     // For Engine
	"goxelcore/graphics"   // For Camera (stub)
	"goxelcore/graphics/core" // For Batch2D, DrawContext
	"goxelcore/window"     // For Input, CursorState
	"goxelcore/devtools"   // For Editor (stub)
)

// PageLoaderFunc corresponds to C++ gui::PageLoaderFunc.
type PageLoaderFunc func(name string) *UINode

// GUI is the main UI controller.
// Corresponds to C++ gui::GUI class.
type GUI struct {
	engine *engine.Engine // C++ uses reference
	input  window.Input   // C++ uses reference

	batch2D     *core.Batch2D
	container   *Container
	hover       *UINode
	pressed     *UINode
	focus       *UINode
	tooltip     *UINode
	rootDocument *UiDocument // C++ uses std::shared_ptr<UiDocument>
	syntaxColorScheme *FontStylesScheme // C++ uses std::unique_ptr<FontStylesScheme>
	storage     map[string]*UINode // C++ uses std::unordered_map<std::string, std::shared_ptr<UINode>>

	uiCamera *graphics.Camera // C++ uses std::unique_ptr<Camera>
	menu     *Menu            // C++ uses std::shared_ptr<Menu>
	postRunnables []func()     // C++ uses std::queue<runnable>

	pagesLoader PageLoaderFunc

	tooltipTimer   float32
	doubleClickTimer float32
	doubleClickDelay float32
	doubleClicked  bool
	debug          bool
}

// NewGUI creates a new GUI instance.
// Corresponds to C++ GUI(Engine& engine) constructor.
func NewGUI(eng *engine.Engine) *GUI {
	gui := &GUI{
		engine: eng,
		input:  eng.GetInput(), // Get input from engine
		storage: make(map[string]*UINode),
		postRunnables: make([]func(), 0),
		doubleClickDelay: 0.5, // Default from C++
		debug: false,
	}

	// Initialize Batch2D

batch2D, err := core.NewBatch2D(10000) // Example capacity
	if err != nil {
		log.Printf("GUI: Failed to create Batch2D: %v\n", err)
	}
	gui.batch2D = batch2D

	// Initialize container, menu, camera, etc.
	gui.container = NewContainer(gui, goxelcore.Vec2f{X: 0, Y: 0}, goxelcore.Vec4f{W:0,X:0,Y:0,Z:0}, 0, OrientationVertical) // Placeholder default size/padding
	gui.menu = NewMenu(gui)
	gui.uiCamera = &graphics.Camera{} // Stub Camera
	gui.rootDocument = &UiDocument{} // Stub UiDocument
	gui.syntaxColorScheme = &FontStylesScheme{} // Stub FontStylesScheme

	return gui
}

// SetPageLoader sets the function to load UI pages.
// Corresponds to C++ GUI::setPageLoader().
func (g *GUI) SetPageLoader(pageLoader PageLoaderFunc) {
	g.pagesLoader = pageLoader
}

// GetPagesLoader returns the page loader function.
// Corresponds to C++ GUI::getPagesLoader().
func (g *GUI) GetPagesLoader() PageLoaderFunc {
	return g.pagesLoader
}

// GetMenu returns the main menu node.
// Corresponds to C++ GUI::getMenu().
func (g *GUI) GetMenu() *Menu {
	return g.menu
}

// GetFocused returns the currently focused node.
// Corresponds to C++ GUI::getFocused().
func (g *GUI) GetFocused() *UINode {
	return g.focus
}

// IsFocusCaught checks if user input is caught by an element.
// Corresponds to C++ GUI::isFocusCaught().
func (g *GUI) IsFocusCaught() bool {
	// Stub implementation
	return false
}

// Act handles input and updates UI logic.
// Corresponds to C++ GUI::act(float delta, const glm::uvec2& viewport).
func (g *GUI) Act(delta float32, viewport goxelcore.Vec2u) {
	// Stub
	g.actMouse(delta, g.input.GetCursor())
	g.actFocused()
	g.updateTooltip(delta)
	log.Printf("GUI.Act: Stub - delta: %f, viewport: %v\n", delta, viewport)
}

// Draw renders all visible elements.
// Corresponds to C++ GUI::draw(const DrawContext& pctx, const Assets& assets).
func (g *GUI) Draw(pctx core.DrawContext, assets *assets.Assets) {
	// Stub
	log.Printf("GUI.Draw: Stub - Drawing UI elements.\n")
	// Actual drawing would involve iterating through container children and drawing them.
	// g.batch2D.Begin()
	// for _, child := range g.container.Children {
	//    // Draw child
	// }
	// g.batch2D.Flush()
}

// PostAct performs post-update actions.
// Corresponds to C++ GUI::postAct().
func (g *GUI) PostAct() {
	// Stub
	log.Println("GUI.PostAct: Stub")
	for _, r := range g.postRunnables {
		r()
	}
	g.postRunnables = g.postRunnables[:0] // Clear queue
}

// Add adds an element to the main container.
// Corresponds to C++ GUI::add().
func (g *GUI) Add(node *UINode) {
	// Stub
	log.Printf("GUI.Add: Stub - Adding UINode: %v\n", node)
	// g.container.AddChild(node)
}

// Remove removes an element from the main container.
// Corresponds to C++ GUI::remove().
func (g *GUI) Remove(node *UINode) {
	// Stub
	log.Printf("GUI.Remove: Stub - Removing UINode: %v\n", node)
}

// Store stores a node in the GUI's dictionary.
// Corresponds to C++ GUI::store().
func (g *GUI) Store(name string, node *UINode) {
	g.storage[name] = node
	log.Printf("GUI.Store: Stored UINode '%s'\n", name)
}

// Get retrieves a stored node.
// Corresponds to C++ GUI::get().
func (g *GUI) Get(name string) *UINode {
	return g.storage[name]
}

// RemoveByName removes a node from the GUI's dictionary by name.
// Corresponds to C++ GUI::remove(const std::string& name).
func (g *GUI) RemoveByName(name string) {
	delete(g.storage, name)
	log.Printf("GUI.RemoveByName: Removed UINode '%s'\n", name)
}

// SetFocus sets the currently focused node.
// Corresponds to C++ GUI::setFocus().
func (g *GUI) SetFocus(node *UINode) {
	g.focus = node
	log.Printf("GUI.SetFocus: Set focus to node: %v\n", node)
}

// GetContainer returns the main container.
// Corresponds to C++ GUI::getContainer().
func (g *GUI) GetContainer() *Container {
	return g.container
}

// SetSyntaxColorScheme sets the syntax color scheme.
// Corresponds to C++ GUI::setSyntaxColorScheme().
func (g *GUI) SetSyntaxColorScheme(scheme *FontStylesScheme) {
	g.syntaxColorScheme = scheme
}

// GetSyntaxColorScheme returns the syntax color scheme.
// Corresponds to C++ GUI::getSyntaxColorScheme().
func (g *GUI) GetSyntaxColorScheme() *FontStylesScheme {
	return g.syntaxColorScheme
}

// OnAssetsLoad handles actions after assets are loaded.
// Corresponds to C++ GUI::onAssetsLoad().
func (g *GUI) OnAssetsLoad(assets *assets.Assets) {
	log.Println("GUI.OnAssetsLoad: Stub")
}

// PostRunnable enqueues a function to be run after the current frame.
// Corresponds to C++ GUI::postRunnable().
func (g *GUI) PostRunnable(callback func()) {
	g.postRunnables = append(g.postRunnables, callback)
}

// SetDoubleClickDelay sets the delay for double-click detection.
// Corresponds to C++ GUI::setDoubleClickDelay().
func (g *GUI) SetDoubleClickDelay(delay float32) {
	g.doubleClickDelay = delay
}

// GetDoubleClickDelay returns the delay for double-click detection.
// Corresponds to C++ GUI::getDoubleClickDelay().
func (g *GUI) GetDoubleClickDelay() float32 {
	return g.doubleClickDelay
}

// ToggleDebug toggles debug mode.
// Corresponds to C++ GUI::toggleDebug().
func (g *GUI) ToggleDebug() {
	g.debug = !g.debug
	log.Printf("GUI.ToggleDebug: Debug mode is now %t\n", g.debug)
}

// GetInput returns the Input system.
// Corresponds to C++ GUI::getInput().
func (g *GUI) GetInput() window.Input {
	return g.input
}

// GetWindow returns the Window.
// Corresponds to C++ GUI::getWindow().
func (g *GUI) GetWindow() *window.Window {
	return g.engine.GetWindow()
}

// GetEditor returns the Editor.
// Corresponds to C++ GUI::getEditor().
func (g *GUI) GetEditor() *devtools.Editor {
	return g.engine.GetEditor() // Assuming engine has a GetEditor()
}

// GetEngine returns the Engine instance.
// Corresponds to C++ GUI::getEngine().
func (g *GUI) GetEngine() *engine.Engine {
	return g.engine
}

// Private helper methods
func (g *GUI) actMouse(delta float32, cursor window.CursorState) {
	// Stub for mouse interaction logic
	log.Printf("GUI.actMouse: Stub - delta: %f, cursor: %v\n", delta, cursor)
}

func (g *GUI) actFocused() {
	// Stub for focused element logic
	log.Println("GUI.actFocused: Stub")
}

func (g *GUI) updateTooltip(delta float32) {
	// Stub for tooltip logic
	log.Printf("GUI.updateTooltip: Stub - delta: %f\n", delta)
}

func (g *GUI) resetTooltip() {
	// Stub for resetting tooltip
	log.Println("GUI.resetTooltip: Stub")
}