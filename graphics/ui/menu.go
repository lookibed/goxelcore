package ui

import "goxelcore"

// Menu is a specialized container for menu items.
// Corresponds to C++ gui::Menu class.
type Menu struct {
	Container // Embed Container
	// Add menu-specific properties
}

// NewMenu creates a new Menu.
func NewMenu(gui *GUI) *Menu {
	return &Menu{
		Container: *NewContainer(gui, goxelcore.Vec2f{X: 0, Y: 0}, goxelcore.Vec4f{W:0,X:0,Y:0,Z:0}, 0, OrientationVertical),
	}
}

// HasOpenPage checks if there is an open page in the menu.
func (m *Menu) HasOpenPage() bool {
	// Stub
	return false
}

// Back navigates back in the menu.
func (m *Menu) Back() bool {
	// Stub
	return false
}
