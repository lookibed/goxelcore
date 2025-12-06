package ui

import "goxelcore" // For Vec2f, Vec4f

// UINode is the base interface/struct for all UI elements.
// Corresponds to C++ gui::UINode.
type UINode struct {
	// Common properties of a UI node
	Visible  bool
	Position goxelcore.Vec2f
	Size     goxelcore.Vec2f
	// Add more common properties as discovered from C++ UINode (if it exists)
	// or from inspecting other UI element headers.
}

// Placeholder for basic UINode methods
func (n *UINode) IsVisible() bool { return n.Visible }
func (n *UINode) GetPosition() goxelcore.Vec2f { return n.Position }
func (n *UINode) GetSize() goxelcore.Vec2f { return n.Size }
