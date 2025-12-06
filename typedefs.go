package goxelcore // This will be the main package for common types

// Vec2u corresponds to glm::uvec2
type Vec2u struct {
	X, Y uint32
}

// Vec2f corresponds to glm::vec2
type Vec2f struct {
	X, Y float32
}

// Vec2i corresponds to glm::ivec2
type Vec2i struct {
	X, Y int
}

// Vec3f corresponds to glm::vec3
type Vec3f struct {
	X, Y, Z float32
}

// Vec4f corresponds to glm::vec4
type Vec4f struct {
	X, Y, Z, W float32
}

// size_t corresponds to C++ size_t
type SizeT uint