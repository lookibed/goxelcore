package ui

import "goxelcore" // For Vec2f, Vec4f

// Container is a UI element that can hold other UI elements.
// Corresponds to C++ gui::Container class.
type Container struct {
	UINode // Embed UINode
	Children []interface{} // Children elements, interface{} for now
	// Add more container-specific properties as discovered
	padding   goxelcore.Vec4f
	interval  float32
	orientation Orientation
}

// NewContainer creates a new Container.
func NewContainer(gui *GUI, size goxelcore.Vec2f, padding goxelcore.Vec4f, interval float32, orientation Orientation) *Container {
	return &Container{
		UINode: UINode{Visible: true, Size: size},
		Children: make([]interface{}, 0),
		padding: padding,
		interval: interval,
		orientation: orientation,
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
