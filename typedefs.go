package goxelcore // This will be the main package for common types

// Vec2u corresponds to glm::uvec2
type Vec2u struct {
	X, Y uint32
}

// Vec2f corresponds to glm::vec2
type Vec2f struct {
	X, Y float32
}

// Vec3f corresponds to glm::vec3
type Vec3f struct {
	X, Y, Z float32
}

// Vec4f corresponds to glm::vec4
type Vec4f struct {
	X, Y, Z, W float32
}

// Add other typedefs as needed from C++ typedefs.hpp
// For now, these are the only ones directly needed by DrawContext.