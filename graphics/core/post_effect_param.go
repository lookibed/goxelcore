package core

import "goxelcore" // For typedefs like Vec2f, Vec3f, Vec4f

// ParamType corresponds to C++ PostEffect::Param::Type enum.
type ParamType int

const (
	ParamTypeInt ParamType = iota
	ParamTypeFloat
	ParamTypeVec2 // Corresponding to glm::vec2 (will use goxelcore.Vec2f)
	ParamTypeVec3 // Corresponding to glm::vec3 (will use goxelcore.Vec3f)
	ParamTypeVec4 // Corresponding to glm::vec4 (will use goxelcore.Vec4f)
)

// Param represents a shader uniform parameter.
// Corresponds to C++ PostEffect::Param struct.
type Param struct {
	Type ParamType
	// Value is a variant in C++. In Go, we'll use an interface{} for flexibility
	// or define a custom type that holds one of the possible values.
	// For now, let's keep it simple with explicit fields or interface.
	DefValue interface{} // Default value
	Value    interface{} // Current value
	Array    bool
	Dirty    bool
}

// NewParam creates a new PostEffect Param.
func NewParam(paramType ParamType, defValue interface{}, array bool) Param {
	return Param{
		Type:     paramType,
		DefValue: defValue,
		Value:    defValue,
		Array:    array,
		Dirty:    true,
	}
}