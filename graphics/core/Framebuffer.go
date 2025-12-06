package core

import (
	"fmt"
	"log"

	"github.com/Zyko0/go-sdl3/gl"
)

// Framebuffer corresponds to C++ Framebuffer class in voxelcore/src/graphics/core/Framebuffer.hpp
type Framebuffer struct {
	fbo     uint32
	depth   uint32 // Renderbuffer ID for depth/stencil
	width   uint
	height  uint
	format  uint32 // GL texture internal format, e.g., GL_RGBA8, GL_RGB8
	texture *Texture
}

// NewFramebuffer creates a new Framebuffer with a color texture and optional depth buffer.
// Corresponds to C++ Framebuffer(uint width, uint height, bool alpha=false)
func NewFramebuffer(width, height uint, alpha bool) (*Framebuffer, error) {
	var fboID uint32
	gl.GenFramebuffers(1, &fboID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fboID)

	var texFormat uint32 = gl.RGB8
	var glFormat int32 = gl.RGB
	if alpha {
		texFormat = gl.RGBA8
		glFormat = gl.RGBA
	}

	// Create a texture to attach to the framebuffer
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, int32(texFormat), int32(width), int32(height), 0, glFormat, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, textureID, 0)

	texture := NewTexture(textureID, width, height)

	// Create a depth and stencil renderbuffer
	var depthID uint32
	gl.GenRenderbuffers(1, &depthID)
	gl.BindRenderbuffer(gl.RENDERBUFFER, depthID)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, depthID)

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		gl.DeleteFramebuffers(1, &fboID)
		gl.DeleteRenderbuffers(1, &depthID)
		texture.Delete() // Delete the attached texture
		return nil, fmt.Errorf("framebuffer incomplete: %x", status)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	return &Framebuffer{
		fbo:     fboID,
		depth:   depthID,
		width:   width,
		height:  height,
		format:  texFormat,
		texture: texture,
	}, nil
}

// NewFramebufferFromExisting creates a new Framebuffer using existing GL IDs and a texture.
// Corresponds to C++ Framebuffer(uint fbo, uint depth, std::unique_ptr<Texture> texture)
func NewFramebufferFromExisting(fboID, depthID uint32, texture *Texture) *Framebuffer {
	if texture == nil {
		log.Println("NewFramebufferFromExisting: texture is nil")
		return nil
	}
	return &Framebuffer{
		fbo:     fboID,
		depth:   depthID,
		width:   texture.Width,
		height:  texture.Height,
		format:  0, // Format unknown from existing texture
		texture: texture,
	}
}


// Bind activates the framebuffer.
// Corresponds to C++ Framebuffer::bind()
func (fb *Framebuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo)
	gl.Viewport(0, 0, int32(fb.width), int32(fb.height))
}

// Unbind deactivates the framebuffer, binding the default framebuffer.
// Corresponds to C++ Framebuffer::unbind()
func (fb *Framebuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	// Reset viewport to default window size or a known size.
	// For now, assume engine will set the main viewport.
}

// Resize updates the framebuffer's texture and depth buffer size.
// Corresponds to C++ Framebuffer::resize()
func (fb *Framebuffer) Resize(width, height uint) error {
	if fb.width == width && fb.height == height {
		return nil
	}

	fb.width = width
	fb.height = height

	fb.Bind()

	// Resize texture
	gl.BindTexture(gl.TEXTURE_2D, fb.texture.id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, int32(fb.format), int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil) // Assume RGBA for resize
	// Re-attach texture
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fb.texture.id, 0)

	// Resize depth renderbuffer
	gl.BindRenderbuffer(gl.RENDERBUFFER, fb.depth)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int32(width), int32(height))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, fb.depth)

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		fb.Unbind()
		return fmt.Errorf("framebuffer incomplete after resize: %x", status)
	}

	fb.Unbind()
	return nil
}

// GetTexture returns the color attachment texture of the framebuffer.
// Corresponds to C++ Framebuffer::getTexture()
func (fb *Framebuffer) GetTexture() *Texture {
	return fb.texture
}

// GetWidth returns the framebuffer width.
// Corresponds to C++ Framebuffer::getWidth()
func (fb *Framebuffer) GetWidth() uint {
	return fb.width
}

// GetHeight returns the framebuffer height.
// Corresponds to C++ Framebuffer::getHeight()
func (fb *Framebuffer) GetHeight() uint {
	return fb.height
}

// GetFBO returns the OpenGL Framebuffer Object ID.
// Corresponds to C++ Framebuffer::getFBO()
func (fb *Framebuffer) GetFBO() uint32 {
	return fb.fbo
}

// Delete cleans up OpenGL resources.
// Corresponds to C++ ~Framebuffer() destructor.
func (fb *Framebuffer) Delete() {
	if fb.fbo != 0 {
		gl.DeleteFramebuffers(1, &fb.fbo)
	}
	if fb.depth != 0 {
		gl.DeleteRenderbuffers(1, &fb.depth)
	}
	if fb.texture != nil {
		fb.texture.Delete()
	}
}
