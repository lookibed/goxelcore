package core

import (
	"fmt"
	"log"
	"unsafe"

	gl "github.com/go-gl/gl/v3.3-core/gl"
)

// MeshStats corresponds to C++ MeshStats struct in voxelcore/src/graphics/core/Mesh.hpp
var MeshStats = struct {
	MeshesCount int
	DrawCalls   int
}{}

// IndexBufferData corresponds to C++ IndexBufferData struct in voxelcore/src/graphics/core/Mesh.hpp
type IndexBufferData struct {
	Indices      []uint32 // Pointer in C++, slice in Go
	IndicesCount size_t   // size_t from C++ typedefs.hpp
}

// Mesh corresponds to C++ template class Mesh<VertexStructure> in voxelcore/src/graphics/core/Mesh.hpp
type Mesh struct {
	vao         uint32
	vbo         uint32
	ibos        []IndexBuffer // For multiple index buffers
	vertexCount size_t
}

// IndexBuffer internal struct
type IndexBuffer struct {
	ibo        uint32
	indexCount size_t
}

// NewMesh creates a new Mesh from MeshData.
// Corresponds to C++ explicit Mesh(const MeshData<VertexStructure>& data)
func NewMesh(data *MeshData) (*Mesh, error) {
	if data.Stride == 0 {
		return nil, fmt.Errorf("mesh data stride is zero, cannot create mesh")
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data.Vertices), gl.Ptr(data.Vertices), gl.STATIC_DRAW)

	var offset uintptr
	for i, attr := range data.Attrs {
		gl.EnableVertexAttribArray(uint32(i))
		var glType uint32
		var typeSize uint32
		
		switch attr.Type {
		case VertexAttributeTypeFloat:
			glType = gl.FLOAT
			typeSize = uint32(unsafe.Sizeof(float32(0)))
		case VertexAttributeTypeInt:
			glType = gl.INT
			typeSize = uint32(unsafe.Sizeof(int32(0)))
		case VertexAttributeTypeUnsignedInt:
			glType = gl.UNSIGNED_INT
			typeSize = uint32(unsafe.Sizeof(uint32(0)))
		case VertexAttributeTypeShort:
			glType = gl.SHORT
			typeSize = uint32(unsafe.Sizeof(int16(0)))
		case VertexAttributeTypeUnsignedShort:
			glType = gl.UNSIGNED_SHORT
			typeSize = uint32(unsafe.Sizeof(uint16(0)))
		case VertexAttributeTypeByte:
			glType = gl.BYTE
			typeSize = uint32(unsafe.Sizeof(int8(0)))
		case VertexAttributeTypeUnsignedByte:
			glType = gl.UNSIGNED_BYTE
			typeSize = uint32(unsafe.Sizeof(uint8(0)))
		default:
			return nil, fmt.Errorf("unsupported vertex attribute type: %v", attr.Type)
		}

		if attr.Type == VertexAttributeTypeFloat {
			gl.VertexAttribPointer(uint32(i), int32(attr.Count), glType, attr.Normalized, int32(data.Stride), gl.PtrOffset(int(offset)))
		} else {
			// For integer attributes, use VertexAttribIPointer
			gl.VertexAttribIPointer(uint32(i), int32(attr.Count), glType, int32(data.Stride), gl.PtrOffset(int(offset)))
		}

		offset += uintptr(uint32(attr.Count) * typeSize)
	}

	ibos := make([]IndexBuffer, len(data.Indices))
	for i, indicesData := range data.Indices {
		var ibo uint32
		gl.GenBuffers(1, &ibo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicesData)*int(unsafe.Sizeof(uint32(0))), gl.Ptr(indicesData), gl.STATIC_DRAW)
		ibos[i] = IndexBuffer{ibo: ibo, indexCount: size_t(len(indicesData))}
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	MeshStats.MeshesCount++
	return &Mesh{
		vao:         vao,
		vbo:         vbo,
		ibos:        ibos,
		vertexCount: size_t(data.GetVertexCount()),
	}, nil
}

// NewMeshFromBuffers creates a new Mesh from raw buffers.
// Corresponds to C++ Mesh(const VertexStructure* vertexBuffer, size_t vertices, std::vector<IndexBufferData> indices)
// and Mesh(const VertexStructure* vertexBuffer, size_t vertices)
func NewMeshFromBuffers(vertexBuffer []byte, vertexCount size_t, indexBufferData []IndexBufferData, attrs []VertexAttribute) (*Mesh, error) {
	md := NewMeshData(vertexBuffer, make([][]uint32, 0), attrs) // Create MeshData temporarily
	md.Indices = make([][]uint32, len(indexBufferData))
	for i, ibd := range indexBufferData {
		md.Indices[i] = ibd.Indices
	}
	md.CalculateStride() // Ensure stride is calculated for the attrs

	return NewMesh(md)
}

// Reload updates GL vertex and index buffers data.
// Corresponds to C++ Mesh::reload()
func (m *Mesh) Reload(vertexBuffer []byte, vertexCount size_t, indexBufferData []IndexBufferData, attrs []VertexAttribute) error {
	gl.BindVertexArray(m.vao)

	// Reload VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBuffer), gl.Ptr(vertexBuffer), gl.STATIC_DRAW)
	m.vertexCount = vertexCount

	// Reload IBOs
	for _, ibo := range m.ibos {
		gl.DeleteBuffers(1, &ibo.ibo)
	}
	m.ibos = make([]IndexBuffer, len(indexBufferData))
	for i, indicesData := range indexBufferData {
		var iboId uint32
		gl.GenBuffers(1, &iboId)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, iboId)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicesData.Indices)*int(unsafe.Sizeof(uint32(0))), gl.Ptr(indicesData.Indices), gl.STATIC_DRAW)
		m.ibos[i] = IndexBuffer{ibo: iboId, indexCount: indicesData.IndicesCount}
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	return nil
}

// Draw renders the mesh with specified primitive type and index buffer.
// Corresponds to C++ Mesh::draw(unsigned int primitive, int iboIndex = 0)
func (m *Mesh) Draw(primitive uint32, iboIndex int) {
	if iboIndex >= len(m.ibos) {
		log.Printf("Warning: IBO index %d out of bounds for mesh with %d IBOs\n", iboIndex, len(m.ibos))
		return
	}

	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ibos[iboIndex].ibo)
	gl.DrawElements(primitive, int32(m.ibos[iboIndex].indexCount), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	MeshStats.DrawCalls++
}

// DrawTriangles renders the mesh as triangles.
// Corresponds to C++ Mesh::draw()
func (m *Mesh) DrawTriangles() {
	if len(m.ibos) > 0 {
		m.Draw(gl.TRIANGLES, 0)
	} else {
		// If no index buffer, draw arrays directly
		gl.BindVertexArray(m.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(m.vertexCount))
		gl.BindVertexArray(0)
		MeshStats.DrawCalls++
	}
}

// Delete cleans up OpenGL resources.
// Corresponds to C++ ~Mesh() destructor.
func (m *Mesh) Delete() {
	gl.DeleteVertexArrays(1, &m.vao)
	gl.DeleteBuffers(1, &m.vbo)
	for _, ibo := range m.ibos {
		gl.DeleteBuffers(1, &ibo.ibo)
	}
	MeshStats.MeshesCount--
}
