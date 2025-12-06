package core

import (
	"fmt"
	"log"
	"goxelcore/maths" // For UVRegion

	"github.com/Zyko0/go-sdl3/gl"
)

// Texture corresponds to C++ Texture class in voxelcore/src/graphics/core/Texture.hpp
type Texture struct {
	id     uint32
	width  uint
	height uint
}

// MAX_RESOLUTION corresponds to C++ static uint MAX_RESOLUTION.
const MAX_RESOLUTION uint = 16384 // Example value, check actual C++ value if possible.

// NewTexture creates a new Texture with a given OpenGL ID, width, and height.
// Corresponds to C++ Texture(uint id, uint width, uint height)
func NewTexture(id uint32, width, height uint) *Texture {
	return &Texture{
		id:     id,
		width:  width,
		height: height,
	}
}

// NewTextureFromData creates a new Texture from raw pixel data.
// Corresponds to C++ Texture(const ubyte* data, uint width, uint height, ImageFormat format)
func NewTextureFromData(data []byte, width, height uint, format ImageFormat) (*Texture, error) {
	var glFormat int32
	var glType int32
	var glInternalFormat int32

	switch format {
	case ImageFormatRGB888:
		glFormat = gl.RGB
		glType = gl.UNSIGNED_BYTE
		glInternalFormat = gl.RGB
	case ImageFormatRGBA8888:
		glFormat = gl.RGBA
		glType = gl.UNSIGNED_BYTE
		glInternalFormat = gl.RGBA
	default:
		return nil, fmt.Errorf("unsupported image format: %v", format)
	}

	var id uint32
	gl.GenTextures(1, &id)
	if id == 0 {
		return nil, fmt.Errorf("failed to generate texture ID")
	}

	texture := NewTexture(id, width, height)
	texture.Bind()
	gl.TexImage2D(gl.TEXTURE_2D, 0, glInternalFormat, int32(width), int32(height), 0, glFormat, glType, gl.Ptr(data))
	texture.Unbind()

	return texture, nil
}

// TextureFromImageData creates a new Texture from an ImageData object.
// Corresponds to C++ static std::unique_ptr<Texture> from(const ImageData* image)
func TextureFromImageData(image *ImageData) (*Texture, error) {
	if image == nil {
		return nil, fmt.Errorf("ImageData is nil")
	}
	return NewTextureFromData(image.Data, image.Width, image.Height, image.Format)
}


// Bind activates the texture unit.
// Corresponds to C++ Texture::bind()
func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

// Unbind deactivates the texture unit.
// Corresponds to C++ Texture::unbind()
func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Reload updates the entire texture data.
// Corresponds to C++ Texture::reload(const ubyte* data, uint w, uint h)
func (t *Texture) Reload(data []byte, w, h uint) error {
	// Assuming format remains the same as when texture was created.
	// For now, this is a simplified reload. A more robust implementation
	// would need the format or re-create the texture.
	t.Bind()
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, int32(w), int32(h), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
	t.Unbind()
	t.width = w
	t.height = h
	return nil
}

// ReloadPartial updates a sub-region of the texture data.
// Corresponds to C++ Texture::reloadPartial(const ImageData& image, uint x, uint y, uint w, uint h)
func (t *Texture) ReloadPartial(image *ImageData, x, y, w, h uint) error {
	t.Bind()
	var glFormat int32
	var glType int32
	switch image.Format {
	case ImageFormatRGB888:
		glFormat = gl.RGB
		glType = gl.UNSIGNED_BYTE
	case ImageFormatRGBA8888:
		glFormat = gl.RGBA
		glType = gl.UNSIGNED_BYTE
	default:
		return fmt.Errorf("unsupported image format for partial reload: %v", image.Format)
	}
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, int32(x), int32(y), int32(w), int32(h), glFormat, glType, gl.Ptr(image.Data))
	t.Unbind()
	return nil
}

// SetNearestFilter sets the texture min/mag filter to GL_NEAREST.
// Corresponds to C++ Texture::setNearestFilter()
func (t *Texture) SetNearestFilter() {
	t.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	t.Unbind()
}

// ReloadFromImageData updates the entire texture from an ImageData object.
// Corresponds to C++ Texture::reload(const ImageData& image)
func (t *Texture) ReloadFromImageData(image *ImageData) error {
	return t.Reload(image.Data, image.Width, image.Height)
}

// SetMipMapping configures mipmapping for the texture.
// Corresponds to C++ Texture::setMipMapping()
func (t *Texture) SetMipMapping(flag, pixelated bool) {
	t.Bind()
	if flag {
		gl.GenerateMipmap(gl.TEXTURE_2D)
		if pixelated {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		}
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // Mag filter is usually linear for mipmaps
	} else {
		// Disable mipmapping (use base level texture)
		if pixelated {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		}
	}
	t.Unbind()
}

// ReadData reads pixel data back from the texture. (Stub)
// Corresponds to C++ std::unique_ptr<ImageData> Texture::readData()
func (t *Texture) ReadData() *ImageData {
	log.Println("Texture.ReadData: Not implemented yet")
	return nil
}

// GetId returns the OpenGL texture ID.
// Corresponds to C++ Texture::getId()
func (t *Texture) GetId() uint32 {
	return t.id
}

// GetUVRegion returns a UVRegion covering the entire texture.
// Corresponds to C++ Texture::getUVRegion()
func (t *Texture) GetUVRegion() *maths.UVRegion {
	return maths.DefaultUVRegion()
}

// GetWidth returns the texture width.
// Corresponds to C++ Texture::getWidth()
func (t *Texture) GetWidth() uint {
	return t.width
}

// GetHeight returns the texture height.
// Corresponds to C++ Texture::getHeight()
func (t *Texture) GetHeight() uint {
	return t.height
}

// Delete cleans up the OpenGL texture.
// Corresponds to C++ Texture destructor.
func (t *Texture) Delete() {
	gl.DeleteTextures(1, &t.id)
}
