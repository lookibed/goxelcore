package core

import (
	"fmt"
	"goxelcore" // For typedefs like Vec2u, Vec4f
	"goxelcore/window" // For Window type
)

// DrawContext manages rendering state.
// Corresponds to C++ DrawContext class in voxelcore/src/graphics/core/DrawContext.hpp
type DrawContext struct {
	window   *window.Window // C++ uses reference, Go uses pointer
	parent   *DrawContext   // C++ uses const DrawContext*
	viewport goxelcore.Vec2u
	g2d      Batch2D // C++ uses Batch2D*, Go uses Batch2D interface
	flushable Flushable
	fbo      Bindable // C++ uses Bindable* for Framebuffer, Go uses Bindable interface
	depthMask bool
	depthTest bool
	cullFace  bool
	blendMode BlendMode
	scissorsCount int
	lineWidth float32
}

// NewDrawContext creates a new DrawContext.
// Corresponds to C++ DrawContext constructor.
func NewDrawContext(parent *DrawContext, win *window.Window, g2d Batch2D) *DrawContext {
	// Default values from C++ constructor (implicitly set or from defaults)
	// blendMode is normal by default
	return &DrawContext{
		parent:        parent,
		window:        win,
		g2d:           g2d,
		depthMask:     true,
		depthTest:     false,
		cullFace:      false,
		blendMode:     BlendModeNormal,
		scissorsCount: 0,
		lineWidth:     1.0,
		viewport:      goxelcore.Vec2u{X: 0, Y: 0}, // Default to 0,0 for now
	}
}

// GetBatch2D returns the 2D batch renderer.
// Corresponds to C++ DrawContext::getBatch2D()
func (dc *DrawContext) GetBatch2D() Batch2D {
	return dc.g2d
}

// GetViewport returns the current viewport dimensions.
// Corresponds to C++ DrawContext::getViewport()
func (dc *DrawContext) GetViewport() goxelcore.Vec2u {
	return dc.viewport
}

// Sub creates a sub-context.
// Corresponds to C++ DrawContext::sub()
func (dc *DrawContext) Sub(flushable Flushable) DrawContext {
	subCtx := *dc // Copy the current context
	subCtx.parent = dc
	subCtx.flushable = flushable
	return subCtx
}

// SetViewport sets the drawing viewport.
// Corresponds to C++ DrawContext::setViewport()
func (dc *DrawContext) SetViewport(viewport goxelcore.Vec2u) {
	dc.viewport = viewport
	// gl.Viewport(int32(viewport.X), int32(viewport.Y), int32(viewport.X), int32(viewport.Y))
	// Actual OpenGL call will be made when we integrate with GL
	fmt.Printf("DrawContext.SetViewport: Not implemented yet (GL call needed: %v)\n", viewport)
}

// SetFramebuffer sets the active framebuffer object.
// Corresponds to C++ DrawContext::setFramebuffer()
func (dc *DrawContext) SetFramebuffer(fbo Bindable) {
	if dc.fbo == fbo {
		return
	}
	if dc.fbo != nil {
		dc.fbo.Unbind()
	}
	dc.fbo = fbo
	if dc.fbo != nil {
		dc.fbo.Bind()
	}
	fmt.Println("DrawContext.SetFramebuffer: Not implemented yet (GL call needed)")
}

// SetDepthMask enables or disables writing to the depth buffer.
// Corresponds to C++ DrawContext::setDepthMask()
func (dc *DrawContext) SetDepthMask(flag bool) {
	if dc.depthMask == flag {
		return
	}
	dc.depthMask = flag
	// gl.DepthMask(flag)
	fmt.Printf("DrawContext.SetDepthMask: Not implemented yet (GL call needed: %v)\n", flag)
}

// SetDepthTest enables or disables depth testing.
// Corresponds to C++ DrawContext::setDepthTest()
func (dc *DrawContext) SetDepthTest(flag bool) {
	if dc.depthTest == flag {
		return
	}
	dc.depthTest = flag
	// if flag { gl.Enable(gl.DEPTH_TEST) } else { gl.Disable(gl.DEPTH_TEST) }
	fmt.Printf("DrawContext.SetDepthTest: Not implemented yet (GL call needed: %v)\n", flag)
}

// SetCullFace enables or disables face culling.
// Corresponds to C++ DrawContext::setCullFace()
func (dc *DrawContext) SetCullFace(flag bool) {
	if dc.cullFace == flag {
		return
	}
	dc.cullFace = flag
	// if flag { gl.Enable(gl.CULL_FACE) } else { gl.Disable(gl.CULL_FACE) }
	fmt.Printf("DrawContext.SetCullFace: Not implemented yet (GL call needed: %v)\n", flag)
}

// SetBlendMode sets the blending mode.
// Corresponds to C++ DrawContext::setBlendMode()
func (dc *DrawContext) SetBlendMode(mode BlendMode) {
	if dc.blendMode == mode {
		return
	}
	dc.blendMode = mode
	// Actual GL blend function calls needed here
	fmt.Printf("DrawContext.SetBlendMode: Not implemented yet (GL call needed: %v)\n", mode)
}

// SetScissors sets the scissor rectangle.
// Corresponds to C++ DrawContext::setScissors()
func (dc *DrawContext) SetScissors(area goxelcore.Vec4f) {
	// Requires GL scissor test calls
	fmt.Printf("DrawContext.SetScissors: Not implemented yet (GL call needed: %v)\n", area)
}

// SetLineWidth sets the line width for drawing.
// Corresponds to C++ DrawContext::setLineWidth()
func (dc *DrawContext) SetLineWidth(width float32) {
	if dc.lineWidth == width {
		return
	}
	dc.lineWidth = width
	// gl.LineWidth(width)
	fmt.Printf("DrawContext.SetLineWidth: Not implemented yet (GL call needed: %v)\n", width)
}

// Batch2D is a placeholder interface for Batch2D.
// A concrete implementation will be ported later.
type Batch2D interface {
	Flushable // Batch2D is Flushable
	// Add other Batch2D methods here
}
