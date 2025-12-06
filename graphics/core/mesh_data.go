package core

import (
	"fmt"
	"log"
	"unsafe" // For Sizeof, to calculate attribute size
)

// VertexAttributeType corresponds to C++ VertexAttribute::Type enum.
type VertexAttributeType int

const (
	VertexAttributeTypeFloat VertexAttributeType = iota
	VertexAttributeTypeInt
	VertexAttributeTypeUnsignedInt
	VertexAttributeTypeShort
	VertexAttributeTypeUnsignedShort
	VertexAttributeTypeByte
	VertexAttributeTypeUnsignedByte
)

// VertexAttribute describes a single vertex attribute.
// Corresponds to C++ VertexAttribute struct in voxelcore/src/graphics/core/MeshData.hpp
type VertexAttribute struct {
	Type       VertexAttributeType
	Normalized bool
	Count      uint8 // Corresponds to ubyte count
}

// Size calculates the size of the attribute in bytes.
// Corresponds to C++ VertexAttribute::size()
func (attr *VertexAttribute) Size() uint32 {
	var typeSize uint32
	switch attr.Type {
	case VertexAttributeTypeFloat:
		typeSize = uint32(unsafe.Sizeof(float32(0)))
	case VertexAttributeTypeInt, VertexAttributeTypeUnsignedInt:
		typeSize = uint32(unsafe.Sizeof(int32(0))) // Assuming 32-bit int/uint for GL
	case VertexAttributeTypeShort, VertexAttributeTypeUnsignedShort:
		typeSize = uint32(unsafe.Sizeof(int16(0)))
	case VertexAttributeTypeByte, VertexAttributeTypeUnsignedByte:
		typeSize = uint32(unsafe.Sizeof(int8(0)))
	default:
		log.Printf("Warning: Unknown VertexAttributeType: %v\n", attr.Type)
		return 0
	}
	return uint32(attr.Count) * typeSize
}

// MeshData holds raw mesh data.
// Corresponds to C++ MeshData<VertexStructure> struct in voxelcore/src/graphics/core/MeshData.hpp
// In Go, we use a slice of bytes for raw vertex data and slices of uint32 for indices
// to make it generic across different vertex structures.
type MeshData struct {
	Vertices []byte // Raw vertex data, e.g., []float32 cast to []byte
	Indices  [][]uint32 // Vector of index buffers
	Attrs    []VertexAttribute // Vertex attributes description
	Stride   uint32 // Calculated vertex stride
}

// NewMeshData creates a new MeshData instance.
// Corresponds to C++ MeshData constructor.
func NewMeshData(vertices []byte, indices [][]uint32, attrs []VertexAttribute) *MeshData {
	md := &MeshData{
		Vertices: vertices,
		Indices:  indices,
		Attrs:    attrs,
	}
	md.CalculateStride()
	return md
}

// CalculateStride computes the total stride for a single vertex.
func (md *MeshData) CalculateStride() {
	md.Stride = 0
	for _, attr := range md.Attrs {
		md.Stride += attr.Size()
	}
}

// GetVertexCount calculates the number of vertices based on raw data and stride.
func (md *MeshData) GetVertexCount() int {
	if md.Stride == 0 {
		return 0
	}
	return len(md.Vertices) / int(md.Stride)
}

// For now, no `util/Buffer.hpp` Go equivalent is explicitly created.
// Go slices handle the buffer functionality.
