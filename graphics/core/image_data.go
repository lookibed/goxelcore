package core

import (
	"fmt"
	"goxelcore" // For Vec4f
	// "image/color" // For potential color conversions
)

// ImageFormat corresponds to C++ ImageFormat enum in voxelcore/src/graphics/core/ImageData.hpp
type ImageFormat int

const (
	ImageFormatRGB888 ImageFormat = iota
	ImageFormatRGBA8888
)

// ImageData represents raw image pixel data.
// Corresponds to C++ ImageData class in voxelcore/src/graphics/core/ImageData.hpp
type ImageData struct {
	Format ImageFormat
	Width  uint
	Height uint
	Data   []uint8 // Corresponds to std::unique_ptr<ubyte[]> data
}

// NewImageData creates a new ImageData instance with uninitialized data.
// Corresponds to C++ ImageData(ImageFormat format, uint width, uint height)
func NewImageData(format ImageFormat, width, height uint) *ImageData {
	channels := uint(3)
	if format == ImageFormatRGBA8888 {
		channels = 4
	}
	dataSize := width * height * channels
	return &ImageData{
		Format: format,
		Width:  width,
		Height: height,
		Data:   make([]uint8, dataSize),
	}
}

// NewImageDataFromBytes creates a new ImageData instance with provided data.
// Corresponds to C++ ImageData(ImageFormat format, uint width, uint height, std::unique_ptr<ubyte[]> data)
// and ImageData(ImageFormat format, uint width, uint height, const ubyte* data)
func NewImageDataFromBytes(format ImageFormat, width, height uint, data []uint8) *ImageData {
	channels := uint(3)
	if format == ImageFormatRGBA8888 {
		channels = 4
	}
	expectedSize := width * height * channels
	if uint(len(data)) != expectedSize {
		// Log an error or return an error, for now just create new data
		fmt.Printf("Warning: Provided data size %d does not match expected size %d for format %v, width %d, height %d. Resizing data.\n",
			len(data), expectedSize, format, width, height)
		newData := make([]uint8, expectedSize)
		copy(newData, data)
		data = newData
	}
	return &ImageData{
		Format: format,
		Width:  width,
		Height: height,
		Data:   data,
	}
}

// GetPixel returns the color of the pixel at (x, y).
// Not directly in C++ API, but useful for implementing other methods.
func (img *ImageData) GetPixel(x, y uint) (r, g, b, a uint8) {
	if x >= img.Width || y >= img.Height {
		return 0, 0, 0, 0 // Out of bounds
	}
	channels := uint(3)
	if img.Format == ImageFormatRGBA8888 {
		channels = 4
	}
	idx := (y*img.Width + x) * channels
	r = img.Data[idx]
	g = img.Data[idx+1]
	b = img.Data[idx+2]
	if channels == 4 {
		a = img.Data[idx+3]
	} else {
		a = 255 // Opaque for RGB
	}
	return r, g, b, a
}

// SetPixel sets the color of the pixel at (x, y).
// Not directly in C++ API, but useful for implementing other methods.
func (img *ImageData) SetPixel(x, y uint, r, g, b, a uint8) {
	if x >= img.Width || y >= img.Height {
		return // Out of bounds
	}
	channels := uint(3)
	if img.Format == ImageFormatRGBA8888 {
		channels = 4
	}
	idx := (y*img.Width + x) * channels
	img.Data[idx] = r
	img.Data[idx+1] = g
	img.Data[idx+2] = b
	if channels == 4 {
		img.Data[idx+3] = a
	}
}

// FlipX flips the image horizontally.
// Corresponds to C++ ImageData::flipX()
func (img *ImageData) FlipX() {
	if img.Width == 0 || img.Height == 0 {
		return
	}
	channels := uint(3)
	if img.Format == ImageFormatRGBA8888 {
		channels = 4
	}
	rowSize := img.Width * channels
	for y := uint(0); y < img.Height; y++ {
		for x := uint(0); x < img.Width/2; x++ {
			leftIdx := (y*img.Width + x) * channels
			rightIdx := (y*img.Width + (img.Width - 1 - x)) * channels
			for c := uint(0); c < channels; c++ {
				img.Data[leftIdx+c], img.Data[rightIdx+c] = img.Data[rightIdx+c], img.Data[leftIdx+c]
			}
		}
	}
}

// FlipY flips the image vertically.
// Corresponds to C++ ImageData::flipY()
func (img *ImageData) FlipY() {
	if img.Width == 0 || img.Height == 0 {
		return
	}
	channels := uint(3)
	if img.Format == ImageFormatRGBA8888 {
		channels = 4
	}
	rowSize := img.Width * channels
	for y := uint(0); y < img.Height/2; y++ {
		topRowIdx := y * rowSize
		bottomRowIdx := (img.Height - 1 - y) * rowSize
		for i := uint(0); i < rowSize; i++ {
			img.Data[topRowIdx+i], img.Data[bottomRowIdx+i] = img.Data[bottomRowIdx+i], img.Data[topRowIdx+i]
		}
	}
}

// DrawLine draws a line between two points. (Stub)
// Corresponds to C++ ImageData::drawLine()
func (img *ImageData) DrawLine(x1, y1, x2, y2 int, col goxelcore.Vec4f) {
	// Implementation would involve a line drawing algorithm like Bresenham's.
	// For now, it's a stub.
	fmt.Printf("ImageData.DrawLine: Stub - (%d,%d) to (%d,%d) with color %v\n", x1, y1, x2, y2, col)
}

// DrawRect draws a rectangle. (Stub)
// Corresponds to C++ ImageData::drawRect()
func (img *ImageData) DrawRect(x, y, w, h int, col goxelcore.Vec4f) {
	// Implementation would involve drawing four lines or filling pixels.
	// For now, it's a stub.
	fmt.Printf("ImageData.DrawRect: Stub - (%d,%d) %dx%d with color %v\n", x, y, w, h, col)
}

// Blit copies one image onto another. (Stub)
// Corresponds to C++ ImageData::blit()
func (img *ImageData) Blit(other *ImageData, x, y int) {
	// This would involve pixel-by-pixel copying, handling different formats.
	// For now, it's a stub.
	fmt.Printf("ImageData.Blit: Stub - Blitting image of size %dx%d at (%d,%d)\n", other.Width, other.Height, x, y)
}

// Extrude (Stub)
// Corresponds to C++ ImageData::extrude()
func (img *ImageData) Extrude(x, y, w, h int) {
	fmt.Printf("ImageData.Extrude: Stub - Extruding area (%d,%d) %dx%d\n", x, y, w, h)
}

// FixAlphaColor (Stub)
// Corresponds to C++ ImageData::fixAlphaColor()
func (img *ImageData) FixAlphaColor() {
	fmt.Println("ImageData.FixAlphaColor: Stub")
}

// MulColor multiplies image colors by a given color. (Stub)
// Corresponds to C++ ImageData::mulColor()
func (img *ImageData) MulColor(col goxelcore.Vec4f) {
	fmt.Printf("ImageData.MulColor: Stub - Multiplying with color %v\n", col)
}

// MulColorImage multiplies image colors by another image's colors. (Stub)
// Corresponds to C++ ImageData::mulColor(const ImageData& other)
func (img *ImageData) MulColorImage(other *ImageData) {
	fmt.Printf("ImageData.MulColorImage: Stub - Multiplying with image of size %dx%d\n", other.Width, other.Height)
}

// AddColor adds color to the image. (Stub)
// Corresponds to C++ ImageData::addColor()
func (img *ImageData) AddColor(col goxelcore.Vec4f, multiplier int) {
	fmt.Printf("ImageData.AddColor: Stub - Adding color %v with multiplier %d\n", col, multiplier)
}

// AddColorImage adds colors from another image. (Stub)
// Corresponds to C++ ImageData::addColor(const ImageData& other, int multiplier)
func (img *ImageData) AddColorImage(other *ImageData, multiplier int) {
	fmt.Printf("ImageData.AddColorImage: Stub - Adding image of size %dx%d with multiplier %d\n", other.Width, other.Height, multiplier)
}

// Extend extends the image to a new size. (Stub)
// Corresponds to C++ ImageData::extend()
func (img *ImageData) Extend(newWidth, newHeight uint) {
	fmt.Printf("ImageData.Extend: Stub - Extending to %dx%d\n", newWidth, newHeight)
}

// Cropped returns a new ImageData instance representing a cropped section. (Stub)
// Corresponds to C++ ImageData::cropped()
func (img *ImageData) Cropped(x, y, width, height int) *ImageData {
	fmt.Printf("ImageData.Cropped: Stub - Cropping area (%d,%d) %dx%d\n", x, y, width, height)
	return nil // Return nil for stub
}

// GetData returns the raw pixel data.
// Corresponds to C++ ImageData::getData()
func (img *ImageData) GetData() []uint8 {
	return img.Data
}

// GetFormat returns the image format.
// Corresponds to C++ ImageData::getFormat()
func (img *ImageData) GetFormat() ImageFormat {
	return img.Format
}

// GetWidth returns the image width.
// Corresponds to C++ ImageData::getWidth()
func (img *ImageData) GetWidth() uint {
	return img.Width
}

// GetHeight returns the image height.
// Corresponds to C++ ImageData::getHeight()
func (img *ImageData) GetHeight() uint {
	return img.Height
}

// GetDataSize returns the size of the raw pixel data in bytes.
// Corresponds to C++ ImageData::getDataSize()
func (img *ImageData) GetDataSize() uint {
	channels := uint(3)
	if img.Format == ImageFormatRGBA8888 {
		channels = 4
	}
	return img.Width * img.Height * channels
}

// AddAtlasMargins (Stub)
// Corresponds to add_atlas_margins() global function.
func AddAtlasMargins(img *ImageData, gridSize int) *ImageData {
	fmt.Printf("AddAtlasMargins: Stub - Adding margins to %dx%d image with grid size %d\n", img.Width, img.Height, gridSize)
	return img // Return original for stub
}
