package engine

// Project represents basic project information.
// Corresponds to C++ Project struct in voxelcore/src/devtools/Project.hpp (and used by Engine)
type Project struct {
	Title string
	// Add other relevant fields if they become necessary for Engine initialization.
}

// NewProject creates a new Project struct.
func NewProject() *Project {
	return &Project{
		Title: "VoxelCore Project", // Default title
	}
}
