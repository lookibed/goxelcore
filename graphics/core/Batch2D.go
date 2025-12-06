package core

import (
	"fmt"
	"log"
	"unsafe"

	gl "github.com/go-gl/gl/v3.3-core/gl"
	"goxelcore" // For Vec2f, Vec4f
	"goxelcore/maths" // For UVRegion
)

// Batch2DVertex corresponds to C++ Batch2DVertex struct.
type Batch2DVertex struct {
	Position goxelcore.Vec2f
	UV       goxelcore.Vec2f
	Color    goxelcore.Vec4f
}

// ATTRIBUTES defines the vertex attributes for Batch2DVertex.
// Corresponds to C++ Batch2DVertex::ATTRIBUTES.
var Batch2DVertexAttributes = []VertexAttribute{
	{Type: VertexAttributeTypeFloat, Normalized: false, Count: 2}, // Position
	{Type: VertexAttributeTypeFloat, Normalized: false, Count: 2}, // UV
	{Type: VertexAttributeTypeFloat, Normalized: false, Count: 4}, // Color
}

// Batch2D corresponds to C++ Batch2D class.
type Batch2D struct {
	buffer     []Batch2DVertex
	capacity   goxelcore.size_t
	mesh       *Mesh // Uses the Mesh struct directly
	blank      *Texture // Blank texture for untextured drawing
	index      goxelcore.size_t // Current vertex index in buffer
	color      goxelcore.Vec4f // Current tint color
	currentTexture *Texture
	primitive  DrawPrimitive
	region     *maths.UVRegion
}

// NewBatch2D creates a new Batch2D instance.
// Corresponds to C++ Batch2D constructor.
func NewBatch2D(capacity goxelcore.size_t) (*Batch2D, error) {
	b := &Batch2D{
		capacity: capacity,
		buffer:   make([]Batch2DVertex, capacity),
		index:    0,
		color:    goxelcore.Vec4f{X: 1.0, Y: 1.0, Z: 1.0, W: 1.0}, // Default to white
		primitive: DrawPrimitiveTriangle,
		region:    maths.DefaultUVRegion(),
	}

	// Create a 1x1 blank white texture for untextured drawing
	blankImageData := NewImageData(ImageFormatRGBA8888, 1, 1)
	blankImageData.Data[0] = 255
	blankImageData.Data[1] = 255
	blankImageData.Data[2] = 255
	blankImageData.Data[3] = 255
	blankTexture, err := TextureFromImageData(blankImageData)
	if err != nil {
		return nil, fmt.Errorf("failed to create blank texture for Batch2D: %w", err)
	}
	b.blank = blankTexture
	b.currentTexture = b.blank

	// Create a Mesh to handle rendering of the batched vertices
	// For now, create a MeshData with dummy vertex data, will be reloaded in Flush
	meshData := NewMeshData(
		make([]byte, int(capacity)*int(unsafe.Sizeof(Batch2DVertex{}))), // Buffer for max capacity
		nil, // No initial indices
		Batch2DVertexAttributes,
	)
	mesh, err := NewMesh(meshData)
	if err != nil {
		return nil, fmt.Errorf("failed to create mesh for Batch2D: %w", err)
	}
	b.mesh = mesh

	return b, nil
}

// Destructor for Batch2D.
func (b *Batch2D) Delete() {
	if b.mesh != nil {
		b.mesh.Delete()
	}
	if b.blank != nil {
		b.blank.Delete()
	}
}


// Begin prepares the batch for drawing.
// Corresponds to C++ Batch2D::begin()
func (b *Batch2D) Begin() {
	b.index = 0
}

// Texture sets the current texture for subsequent draws.
// Corresponds to C++ Batch2D::texture()
func (b *Batch2D) Texture(texture *Texture) {
	if texture == nil {
		b.currentTexture = b.blank
	} else {
		b.currentTexture = texture
	}
}

// Untexture unsets the current texture.
// Corresponds to C++ Batch2D::untexture()
func (b *Batch2D) Untexture() {
	b.currentTexture = b.blank
}

// SetRegion sets the UV region for subsequent draws.
// Corresponds to C++ Batch2D::setRegion()
func (b *Batch2D) SetRegion(region *maths.UVRegion) {
	b.region = region
}

// vertex adds a single vertex to the buffer.
// Corresponds to C++ Batch2D::vertex()
func (b *Batch2D) vertex(position goxelcore.Vec2f, uv goxelcore.Vec2f, color goxelcore.Vec4f) {
	if b.index >= b.capacity {
		b.Flush() // Flush if buffer is full
	}
	b.buffer[b.index] = Batch2DVertex{
		Position: position,
		UV:       uv,
		Color:    color,
	}
	b.index++
}

// setPrimitive sets the drawing primitive type.
// Corresponds to C++ Batch2D::setPrimitive()
func (b *Batch2D) setPrimitive(primitive DrawPrimitive) {
	if b.primitive != primitive {
		b.Flush() // Flush before changing primitive
		b.primitive = primitive
	}
}


// Sprite draws a textured sprite.
// Corresponds to C++ Batch2D::sprite()
func (b *Batch2D) Sprite(x, y, w, h float32, region *maths.UVRegion, tint goxelcore.Vec4f) {
	b.setPrimitive(DrawPrimitiveTriangle)
	b.Texture(b.currentTexture) // Ensure texture is set

	u1, v1 := region.U1, region.V1
	u2, v2 := region.U2, region.V2

	// Vertices for a rectangle
	b.vertex(goxelcore.Vec2f{X: x, Y: y + h}, goxelcore.Vec2f{X: u1, Y: v2}, tint)
	b.vertex(goxelcore.Vec2f{X: x, Y: y}, goxelcore.Vec2f{X: u1, Y: v1}, tint)
	b.vertex(goxelcore.Vec2f{X: x + w, Y: y}, goxelcore.Vec2f{X: u2, Y: v1}, tint)

	b.vertex(goxelcore.Vec2f{X: x, Y: y + h}, goxelcore.Vec2f{X: u1, Y: v2}, tint)
	b.vertex(goxelcore.Vec2f{X: x + w, Y: y}, goxelcore.Vec2f{X: u2, Y: v1}, tint)
	b.vertex(goxelcore.Vec2f{X: x + w, Y: y + h}, goxelcore.Vec2f{X: u2, Y: v2}, tint)
}

// Other sprite/rect/line drawing functions will be added as needed.
// For now, these are stubs or simplified.

// Point draws a single point.
// Corresponds to C++ Batch2D::point()
func (b *Batch2D) Point(x, y, r, g, b, a float32) {
	b.setPrimitive(DrawPrimitivePoint)
	b.Texture(b.blank) // Points are untextured
	col := goxelcore.Vec4f{X: r, Y: g, Z: b, W: a}
	b.vertex(goxelcore.Vec2f{X: x, Y: y}, goxelcore.Vec2f{X: 0, Y: 0}, col)
}

// SetColor sets the current tint color for subsequent draws.
// Corresponds to C++ Batch2D::setColor(const glm::vec4& color)
func (b *Batch2D) SetColor(color goxelcore.Vec4f) {
	b.color = color
}

// SetColorRGBA sets the current tint color using RGBA ints (0-255).
// Corresponds to C++ Batch2D::setColor(int r, int g, int b, int a=255)
func (b *Batch2D) SetColorRGBA(r, g, b, a int) {
	b.color = goxelcore.Vec4f{X: float32(r) / 255.0, Y: float32(g) / 255.0, Z: float32(b) / 255.0, W: float32(a) / 255.0}
}

// ResetColor resets the current tint color to white.
// Corresponds to C++ Batch2D::resetColor()
func (b *Batch2D) ResetColor() {
	b.color = goxelcore.Vec4f{X: 1.0, Y: 1.0, Z: 1.0, W: 1.0}
}

// GetColor returns the current tint color.
// Corresponds to C++ Batch2D::getColor()
func (b *Batch2D) GetColor() goxelcore.Vec4f {
	return b.color
}

// Line draws a line.
// Corresponds to C++ Batch2D::line()
func (b *Batch2D) Line(x1, y1, x2, y2, r, g, b, a float32) {
	b.setPrimitive(DrawPrimitiveLine)
	b.Texture(b.blank) // Lines are untextured
	col := goxelcore.Vec4f{X: r, Y: g, Z: b, W: a}
	b.vertex(goxelcore.Vec2f{X: x1, Y: y1}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x2, Y: y2}, goxelcore.Vec2f{X: 0, Y: 0}, col)
}

// LineRect draws a rectangle outline.
// Corresponds to C++ Batch2D::lineRect()
func (b *Batch2D) LineRect(x, y, w, h float32) {
	b.setPrimitive(DrawPrimitiveLine)
	b.Texture(b.blank) // Lines are untextured
	col := b.color
	b.vertex(goxelcore.Vec2f{X: x, Y: y}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x + w, Y: y}, goxelcore.Vec2f{X: 0, Y: 0}, col)

	b.vertex(goxelcore.Vec2f{X: x + w, Y: y}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x + w, Y: y + h}, goxelcore.Vec2f{X: 0, Y: 0}, col)

	b.vertex(goxelcore.Vec2f{X: x + w, Y: y + h}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x, Y: y + h}, goxelcore.Vec2f{X: 0, Y: 0}, col)

	b.vertex(goxelcore.Vec2f{X: x, Y: y + h}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x, Y: y}, goxelcore.Vec2f{X: 0, Y: 0}, col)
}

// Rect draws a filled rectangle. (Simplified)
// Corresponds to C++ Batch2D::rect(float x, float y, float w, float h)
func (b *Batch2D) Rect(x, y, w, h float32) {
	b.setPrimitive(DrawPrimitiveTriangle)
	b.Texture(b.blank) // Filled rect is untextured
	col := b.color

	b.vertex(goxelcore.Vec2f{X: x, Y: y + h}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x, Y: y}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x + w, Y: y}, goxelcore.Vec2f{X: 0, Y: 0}, col)

	b.vertex(goxelcore.Vec2f{X: x, Y: y + h}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x + w, Y: y}, goxelcore.Vec2f{X: 0, Y: 0}, col)
	b.vertex(goxelcore.Vec2f{X: x + w, Y: y + h}, goxelcore.Vec2f{X: 0, Y: 0}, col)
}


// Flush renders all batched vertices and clears the buffer.
// Corresponds to C++ Batch2D::flush()
func (b *Batch2D) Flush() {
	if b.index == 0 {
		return // Nothing to draw
	}

	b.currentTexture.Bind()
	
	// Convert Batch2DVertex slice to byte slice for MeshData
	// unsafe.Slice is used here for zero-copy conversion,
	// but needs care to ensure underlying memory is not GC'd too early
	// before GL operations are complete. For this project context (no compile/run), it's okay.
	vertexBytes := unsafe.Slice((*byte)(unsafe.Pointer(&b.buffer[0])), int(b.index)*int(unsafe.Sizeof(Batch2DVertex{})))

	// Reload mesh with current buffer data
	b.mesh.Reload(vertexBytes, b.index, nil, Batch2DVertexAttributes) // No indices yet

	// Draw the mesh
	var primitiveGL uint32
	switch b.primitive {
	case DrawPrimitivePoint:
		primitiveGL = gl.POINTS
	case DrawPrimitiveLine:
		primitiveGL = gl.LINES
	case DrawPrimitiveTriangle:
		primitiveGL = gl.TRIANGLES
	default:
		log.Printf("Warning: Unsupported primitive type for Batch2D: %v\n", b.primitive)
		return
	}
	
	// For Batch2D, assume we always draw with vertices directly (no IBO)
	// or create indices on the fly if needed for specific shapes.
	// For now, let's assume gl.DrawArrays for simplicity.
	gl.BindVertexArray(b.mesh.vao)
	gl.DrawArrays(primitiveGL, 0, int32(b.index)) // Draw directly from VBO
	gl.BindVertexArray(0)
	
	b.currentTexture.Unbind()

	b.index = 0 // Reset buffer index
}

// LineWidth sets the width of lines.
// Corresponds to C++ Batch2D::lineWidth()
func (b *Batch2D) LineWidth(width float32) {
	// gl.LineWidth(width) // This is a GL state change, can be applied before Flush
	log.Printf("Batch2D.LineWidth: Stub - Setting line width to %f (GL call needed).\n", width)
}
