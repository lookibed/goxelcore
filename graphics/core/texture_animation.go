package core

import (
	"log"
	"goxelcore" // For Vec2i
)

const (
	// DEFAULT_FRAME_DURATION corresponds to C++ constexpr float DEFAULT_FRAME_DURATION.
	DEFAULT_FRAME_DURATION = 0.150
)

// Frame corresponds to C++ Frame struct in voxelcore/src/graphics/core/TextureAnimation.hpp
type Frame struct {
	SrcPos   goxelcore.Vec2i
	DstPos   goxelcore.Vec2i
	Size     goxelcore.Vec2i
	Duration float32
}

// NewFrame creates a new Frame instance.
func NewFrame(srcPos, dstPos, size goxelcore.Vec2i, duration float32) Frame {
	return Frame{
		SrcPos:   srcPos,
		DstPos:   dstPos,
		Size:     size,
		Duration: duration,
	}
}

// TextureAnimation corresponds to C++ TextureAnimation class in voxelcore/src/graphics/core/TextureAnimation.hpp
type TextureAnimation struct {
	CurrentFrame goxelcore.size_t
	Timer        float32
	SrcTexture   *Texture
	DstTexture   *Texture
	Frames       []Frame
}

// NewTextureAnimation creates a new TextureAnimation instance.
// Corresponds to C++ TextureAnimation(Texture* srcTex, Texture* dstTex) constructor.
func NewTextureAnimation(srcTex, dstTex *Texture) *TextureAnimation {
	return &TextureAnimation{
		SrcTexture: srcTex,
		DstTexture: dstTex,
		CurrentFrame: 0,
		Timer: 0.0,
		Frames: make([]Frame, 0),
	}
}

// AddFrame adds a frame to the animation.
// Corresponds to C++ TextureAnimation::addFrame().
func (ta *TextureAnimation) AddFrame(frame Frame) {
	ta.Frames = append(ta.Frames, frame)
}

// TextureAnimator corresponds to C++ TextureAnimator class in voxelcore/src/graphics/core/TextureAnimation.hpp
type TextureAnimator struct {
	fboR       uint32 // Renderbuffer for color, assumed to be part of FBO
	fboD       uint32 // Renderbuffer for depth
	Animations []TextureAnimation
}

// NewTextureAnimator creates a new TextureAnimator instance.
// Corresponds to C++ TextureAnimator() constructor.
func NewTextureAnimator() *TextureAnimator {
	return &TextureAnimator{
		Animations: make([]TextureAnimation, 0),
	}
}

// AddAnimation adds a single animation to the animator.
// Corresponds to C++ TextureAnimator::addAnimation().
func (ta *TextureAnimator) AddAnimation(animation TextureAnimation) {
	ta.Animations = append(ta.Animations, animation)
}

// AddAnimations adds multiple animations to the animator.
// Corresponds to C++ TextureAnimator::addAnimations().
func (ta *TextureAnimator) AddAnimations(animations []TextureAnimation) {
	ta.Animations = append(ta.Animations, animations...)
}

// Update updates all animations. (Stub for now)
// Corresponds to C++ TextureAnimator::update().
func (ta *TextureAnimator) Update(delta float32) {
	// Actual implementation would iterate through animations, update their timers,
	// and potentially render frames to framebuffers.
	log.Printf("TextureAnimator.Update: Stub - delta: %f\n", delta)
	for i := range ta.Animations {
		anim := &ta.Animations[i]
		anim.Timer += delta
		if anim.CurrentFrame < goxelcore.size_t(len(anim.Frames)) {
			currentFrame := anim.Frames[anim.CurrentFrame]
			if anim.Timer >= currentFrame.Duration {
				anim.Timer -= currentFrame.Duration
				anim.CurrentFrame++
				if anim.CurrentFrame >= goxelcore.size_t(len(anim.Frames)) {
					anim.CurrentFrame = 0 // Loop animation
				}
				// TODO: Trigger rendering or texture update based on new frame
			}
		}
	}
}

// Delete cleans up resources. (Stub for now)
// Corresponds to C++ ~TextureAnimator() destructor.
func (ta *TextureAnimator) Delete() {
	log.Println("TextureAnimator.Delete: Stub")
	// Actual cleanup would involve deleting OpenGL FBO/Renderbuffers if they were created here.
}
