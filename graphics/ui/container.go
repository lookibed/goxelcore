package ui

import (
	"goxelcore" // For Vec2f, Vec4f
	"goxelcore/graphics/core" // For Batch2D
)

// Container is a UI element that can hold other UI elements.
// Corresponds to C++ gui::Container class.
type Container struct {
	UINode // Embed UINode
	Children []interface{} // Children elements, interface{} for now
	// Add more container-specific properties as discovered
	padding   goxelcore.Vec4f
	interval  float32
	orientation Orientation
	batch2D   *core.Batch2D // Added for direct access from Container
}

// NewContainer creates a new Container.
func NewContainer(gui *GUI, size goxelcore.Vec2f, padding goxelcore.Vec4f, interval float32, orientation Orientation) *Container {
	return &Container{
		UINode: UINode{Visible: true, Size: size},
		Children: make([]interface{}, 0),
		padding: padding,
		interval: interval,
		orientation: orientation,
		batch2D: gui.batch2D, // Get Batch2D from GUI
	}
}

// Orientation corresponds to C++ Orientation enum/type (implied by BasePanel)
type Orientation int

const (
	OrientationVertical Orientation = iota
	OrientationHorizontal
)

// AddChild adds a child UI element to the container.
func (c *Container) AddChild(child interface{}) { // Use interface{} for now
	c.Children = append(c.Children, child)
}

// GetBatch2D returns the Batch2D instance used by the container.
func (c *Container) GetBatch2D() *core.Batch2D {
	return c.batch2D
}