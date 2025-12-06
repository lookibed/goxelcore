package maths

import (
	"goxelcore" // For Vec2f, Vec4f
	"math"
)

// UVRegion corresponds to C++ UVRegion struct in voxelcore/src/maths/UVRegion.hpp
type UVRegion struct {
	U1 float32
	V1 float32
	U2 float32
	V2 float32
}

// NewUVRegion creates a new UVRegion with specified coordinates.
// Corresponds to C++ UVRegion(float u1, float v1, float u2, float v2)
func NewUVRegion(u1, v1, u2, v2 float32) *UVRegion {
	return &UVRegion{U1: u1, V1: v1, U2: u2, V2: v2}
}

// DefaultUVRegion creates a new UVRegion with default values (0,0,1,1).
// Corresponds to C++ UVRegion() constructor.
func DefaultUVRegion() *UVRegion {
	return &UVRegion{U1: 0.0, V1: 0.0, U2: 1.0, V2: 1.0}
}

// GetWidth returns the width of the UV region.
// Corresponds to C++ UVRegion::getWidth()
func (uv *UVRegion) GetWidth() float32 {
	return float32(math.Abs(float64(uv.U2 - uv.U1)))
}

// GetHeight returns the height of the UV region.
// Corresponds to C++ UVRegion::getHeight()
func (uv *UVRegion) GetHeight() float32 {
	return float32(math.Abs(float64(uv.V2 - uv.V1)))
}

// AutoSub adjusts the UV region based on sub-coordinates.
// Corresponds to C++ UVRegion::autoSub()
func (uv *UVRegion) AutoSub(w, h, x, y float32) {
	x *= 1.0 - w
	y *= 1.0 - h
	uvw := uv.GetWidth()
	uvh := uv.GetHeight()
	uv.U1 = uv.U1 + uvw*x
	uv.V1 = uv.V1 + uvh*y
	uv.U2 = uv.U1 + uvw*w
	uv.V2 = uv.V1 + uvh*h
}

// Apply transforms a given UV coordinate by the region.
// Corresponds to C++ UVRegion::apply()
func (uv *UVRegion) Apply(vec goxelcore.Vec2f) goxelcore.Vec2f {
	w := uv.GetWidth()
	h := uv.GetHeight()
	return goxelcore.Vec2f{X: uv.U1 + vec.X*w, Y: uv.V1 + vec.Y*h}
}

// Scale scales the UV region.
// Corresponds to C++ UVRegion::scale(float x, float y)
func (uv *UVRegion) Scale(x, y float32) {
	w := uv.U2 - uv.U1
	h := uv.V2 - uv.V1
	cx := (uv.U1 + uv.U2) * 0.5
	cy := (uv.V1 + uv.V2) * 0.5
	uv.U1 = cx - w*0.5*x
	uv.V1 = cy - h*0.5*y
	uv.U2 = cx + w*0.5*x
	uv.V2 = cy + h*0.5*y
}

// ScaleVec scales the UV region using a Vec2f.
// Corresponds to C++ UVRegion::scale(const glm::vec2& vec)
func (uv *UVRegion) ScaleVec(vec goxelcore.Vec2f) {
	uv.Scale(vec.X, vec.Y)
}

// Set sets the UV region coordinates from a Vec4f.
// Corresponds to C++ UVRegion::set(const glm::vec4& vec)
func (uv *UVRegion) Set(vec goxelcore.Vec4f) {
	uv.U1 = vec.X
	uv.V1 = vec.Y
	uv.U2 = vec.Z
	uv.V2 = vec.W
}

// Multiply scales the UVRegion by a Vec2f and returns a new UVRegion.
// Corresponds to C++ UVRegion::operator*(const glm::vec2& scale)
func (uv *UVRegion) Multiply(scale goxelcore.Vec2f) *UVRegion {
	copyUV := *uv
	copyUV.Scale(scale.X, scale.Y)
	return &copyUV
}

// IsFull checks if the UV region covers the full texture (0,0,1,1).
// Corresponds to C++ UVRegion::isFull()
func (uv *UVRegion) IsFull() bool {
	e := float32(1e-7)
	return math.Abs(float64(uv.U1-0.0)) < float64(e) && math.Abs(float64(uv.V1-0.0)) < float64(e) &&
		math.Abs(float64(uv.U2-1.0)) < float64(e) && math.Abs(float64(uv.V2-1.0)) < float64(e)
}
