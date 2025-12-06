package core

// DrawPrimitive corresponds to C++ DrawPrimitive enum in voxelcore/src/graphics/core/commons.hpp
type DrawPrimitive int

const (
	DrawPrimitivePoint DrawPrimitive = iota
	DrawPrimitiveLine
	DrawPrimitiveTriangle
)

// BlendMode corresponds to C++ BlendMode enum in voxelcore/src/graphics/core/commons.hpp
type BlendMode int

const (
	BlendModeNormal BlendMode = iota
	BlendModeAddition
	BlendModeInversion
)

// Flushable interface corresponds to C++ Flushable class.
type Flushable interface {
	Flush()
}

// Bindable interface corresponds to C++ Bindable class.
type Bindable interface {
	Bind()
	Unbind()
}
