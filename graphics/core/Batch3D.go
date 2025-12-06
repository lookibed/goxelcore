package core

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/Zyko0/go-sdl3/gl"
	goxelcore // For Vec2f, Vec3f, Vec4f, size_t
	"goxelcore/maths" // For UVRegion
)

// Batch3DVertex corresponds to C++ Batch3DVertex struct.
type Batch3DVertex struct {
	Position goxelcore.Vec3f
	UV       goxelcore.Vec2f
	Color    goxelcore.Vec4f
}

// ATTRIBUTES defines the vertex attributes for Batch3DVertex.
// Corresponds to C++ Batch3DVertex::ATTRIBUTES.
var Batch3DVertexAttributes = []VertexAttribute{
	{Type: VertexAttributeTypeFloat, Normalized: false, Count: 3}, // Position (x, y, z)
	{Type: VertexAttributeTypeFloat, Normalized: false, Count: 2}, // UV (u, v)
	{Type: VertexAttributeTypeFloat, Normalized: false, Count: 4}, // Color (r, g, b, a)
}

// Batch3D corresponds to C++ Batch3D class.
type Batch3D struct {
	buffer         []Batch3DVertex
	capacity       goxelcore.size_t
	mesh           *Mesh // Uses the Mesh struct directly
	blank          *Texture // Blank texture for untextured drawing
	index          goxelcore.size_t // Current vertex index in buffer
	tint           goxelcore.Vec4f // Current tint color (default to white)
	currentTexture *Texture
	region         *maths.UVRegion // Current UV region (default to full)
}

// NewBatch3D creates a new Batch3D instance.
// Corresponds to C++ Batch3D constructor.
func NewBatch3D(capacity goxelcore.size_t) (*Batch3D, error) {
	b := &Batch3D{
		capacity: capacity,
		buffer:   make([]Batch3DVertex, capacity),
		index:    0,
		tint:     goxelcore.Vec4f{X: 1.0, Y: 1.0, Z: 1.0, W: 1.0}, // Default to white
		region:   maths.DefaultUVRegion(),
	}

	// Create a 1x1 blank white texture for untextured drawing
	blankImageData := NewImageData(ImageFormatRGBA8888, 1, 1)
	blankImageData.Data[0] = 255
	blankImageData.Data[1] = 255
	blankImageData.Data[2] = 255
	blankImageData.Data[3] = 255
	blankTexture, err := TextureFromImageData(blankImageData)
	if err != nil {
		return nil, fmt.Errorf("failed to create blank texture for Batch3D: %w", err)
	}
	b.blank = blankTexture
	b.currentTexture = b.blank

	// Create a Mesh to handle rendering of the batched vertices
	meshData := NewMeshData(
		make([]byte, int(capacity)*int(unsafe.Sizeof(Batch3DVertex{}))),
		nil, // No initial indices
		Batch3DVertexAttributes,
	)
	mesh, err := NewMesh(meshData)
	if err != nil {
		return nil, fmt.Errorf("failed to create mesh for Batch3D: %w", err)
	}
	b.mesh = mesh

	return b, nil
}

// Destructor for Batch3D.
func (b *Batch3D) Delete() {
	if b.mesh != nil {
		b.mesh.Delete()
	}
	if b.blank != nil {
		b.blank.Delete()
	}
}

// Begin prepares the batch for drawing.
// Corresponds to C++ Batch3D::begin()
func (b *Batch3D) Begin() {
	b.index = 0
}

// Texture sets the current texture for subsequent draws.
// Corresponds to C++ Batch3D::texture()
func (b *Batch3D) Texture(texture *Texture) {
	if texture == nil {
		b.currentTexture = b.blank
	} else {
		b.currentTexture = texture
	}
}

// SetRegion sets the UV region for subsequent draws.
// Corresponds to C++ Batch3D::setRegion()
func (b *Batch3D) SetRegion(region *maths.UVRegion) {
	b.region = region
}

// vertex adds a single vertex to the buffer.
// Corresponds to C++ Batch3D::vertex() overloads.
func (b *Batch3D) vertex(position goxelcore.Vec3f, uv goxelcore.Vec2f, color goxelcore.Vec4f) {
	if b.index >= b.capacity {
		b.Flush() // Flush if buffer is full
	}
	b.buffer[b.index] = Batch3DVertex{
		Position: position,
		UV:       uv,
		Color:    color,
	}
	b.index++
}

// Face helper adds a quad face.
// Corresponds to C++ Batch3D::face()
func (b *Batch3D) face(
	coord goxelcore.Vec3f,
	w, h float32,
	axisX, axisY goxelcore.Vec3f,
	region *maths.UVRegion,
	tint goxelcore.Vec4f,
) {
	b.Texture(b.currentTexture) // Ensure texture is set

	u1, v1 := region.U1, region.V1
	u2, v2 := region.U2, region.V2

	// Calculate vertices for the quad
	p0 := coord
	p1 := goxelcore.Vec3f{X: coord.X + axisX.X*w, Y: coord.Y + axisX.Y*w, Z: coord.Z + axisX.Z*w}
	p2 := goxelcore.Vec3f{X: coord.X + axisY.X*h, Y: coord.Y + axisY.Y*h, Z: coord.Z + axisY.Z*h}
	p3 := goxelcore.Vec3f{X: p1.X + axisY.X*h, Y: p1.Y + axisY.Y*h, Z: p1.Z + axisY.Z*h}

	// Top-left triangle
	b.vertex(p0, goxelcore.Vec2f{X: u1, Y: v2}, tint)
	b.vertex(p2, goxelcore.Vec2f{X: u1, Y: v1}, tint)
	b.vertex(p1, goxelcore.Vec2f{X: u2, Y: v2}, tint)

	// Bottom-right triangle
	b.vertex(p1, goxelcore.Vec2f{X: u2, Y: v2}, tint)
	b.vertex(p2, goxelcore.Vec2f{X: u1, Y: v1}, tint)
	b.vertex(p3, goxelcore.Vec2f{X: u2, Y: v1}, tint)
}

// Sprite draws a 3D sprite billboarded. (Simplified)
// Corresponds to C++ Batch3D::sprite()
func (b *Batch3D) Sprite(
	pos, up, right goxelcore.Vec3f,
	w, h float32,
	uv *maths.UVRegion,
	tint goxelcore.Vec4f,
) {
	b.Texture(b.currentTexture)

	// Vertices for the billboard
	// This is a simplified implementation, actual billboard would involve camera alignment
	// For now, assume a quad in XY plane relative to pos
	
	halfW := w / 2.0
	halfH := h / 2.0

	// Calculate points relative to center and axes
	p0 := goxelcore.Vec3f{X: pos.X - right.X*halfW - up.X*halfH, Y: pos.Y - right.Y*halfW - up.Y*halfH, Z: pos.Z - right.Z*halfW - up.Z*halfH} // bottom-left
	p1 := goxelcore.Vec3f{X: pos.X + right.X*halfW - up.X*halfH, Y: pos.Y + right.Y*halfW - up.Y*halfH, Z: pos.Z + right.Z*halfW - up.Z*halfH} // bottom-right
	p2 := goxelcore.Vec3f{X: pos.X - right.X*halfW + up.X*halfH, Y: pos.Y - right.Y*halfW + up.Y*halfH, Z: pos.Z - right.Z*halfW + up.Z*halfH} // top-left
	p3 := goxelcore.Vec3f{X: pos.X + right.X*halfW + up.X*halfH, Y: pos.Y + right.Y*halfW + up.Y*halfH, Z: pos.Z + right.Z*halfW + up.Z*halfH} // top-right

	u1, v1 := uv.U1, uv.V1
	u2, v2 := uv.U2, uv.V2

	// Two triangles forming a quad
	b.vertex(p0, goxelcore.Vec2f{X: u1, Y: v1}, tint) // Bottom-left
	b.vertex(p1, goxelcore.Vec2f{X: u2, Y: v1}, tint) // Bottom-right
	b.vertex(p2, goxelcore.Vec2f{X: u1, Y: v2}, tint) // Top-left

	b.vertex(p1, goxelcore.Vec2f{X: u2, Y: v1}, tint) // Bottom-right
	b.vertex(p3, goxelcore.Vec2f{X: u2, Y: v2}, tint) // Top-right
	b.vertex(p2, goxelcore.Vec2f{X: u1, Y: v2}, tint) // Top-left
}

// XSprite (Stub)
// Corresponds to C++ Batch3D::xSprite()
func (b *Batch3D) XSprite(w, h float32, uv *maths.UVRegion, tint goxelcore.Vec4f, shading bool) {
	log.Printf("Batch3D.XSprite: Stub - w:%.2f h:%.2f uv:%v tint:%v shading:%t\n", w, h, uv, tint, shading)
}

// Cube (Stub)
// Corresponds to C++ Batch3D::cube()
func (b *Batch3D) Cube(coords, size goxelcore.Vec3f, texfaces [6]*maths.UVRegion, tint goxelcore.Vec4f, shading bool) {
	log.Printf("Batch3D.Cube: Stub - coords:%v size:%v tint:%v shading:%t\n", coords, size, tint, shading)
}

// BlockCube (Stub)
// Corresponds to C++ Batch3D::blockCube()
func (b *Batch3D) BlockCube(size goxelcore.Vec3f, texfaces [6]*maths.UVRegion, tint goxelcore.Vec4f, shading bool) {
	log.Printf("Batch3D.BlockCube: Stub - size:%v tint:%v shading:%t\n", size, tint, shading)
}

// Point draws a single 3D point.
// Corresponds to C++ Batch3D::point()
func (b *Batch3D) Point(pos goxelcore.Vec3f, tint goxelcore.Vec4f) {
	// For 3D points, a different primitive type might be needed, or GL_POINTS.
	// For now, this adds a vertex for a point.
	b.vertex(pos, goxelcore.Vec2f{X: 0, Y: 0}, tint)
}

// Flush renders all batched vertices and clears the buffer.
// Corresponds to C++ Batch3D::flush()
func (b *Batch3D) Flush() {
	if b.index == 0 {
		return // Nothing to draw
	}

	b.currentTexture.Bind()
	
	// Convert Batch3DVertex slice to byte slice for MeshData
	vertexBytes := unsafe.Slice((*byte)(unsafe.Pointer(&b.buffer[0])), int(b.index)*int(unsafe.Sizeof(Batch3DVertex{})))

	// Reload mesh with current buffer data
	b.mesh.Reload(vertexBytes, b.index, nil, Batch3DVertexAttributes) // No indices yet

	// Draw the mesh
	gl.BindVertexArray(b.mesh.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(b.index)) // Assuming we draw triangles
	gl.BindVertexArray(0)
	
b.currentTexture.Unbind()

	b.index = 0 // Reset buffer index
}

// FlushPoints renders only point primitives (Stub for now, needs distinct point handling).
// Corresponds to C++ Batch3D::flushPoints()
func (b *Batch3D) FlushPoints() {
	log.Println("Batch3D.FlushPoints: Stub")
}

// SetColor sets the current tint color for subsequent draws.
// Corresponds to C++ Batch3D::setColor()
func (b *Batch3D) SetColor(color goxelcore.Vec4f) {
	b.tint = color
}

// GetColor returns the current tint color.
// Corresponds to C++ Batch3D::getColor()
func (b *Batch3D) GetColor() goxelcore.Vec4f {
	return b.tint
}
