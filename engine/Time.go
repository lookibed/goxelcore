package engine

import (
	"time"
)

// Time corresponds to C++ Time class in voxelcore/src/engine/Time.hpp
type Time struct {
	frame    uint64
	lastTime float64 // In seconds
	delta    float64 // In seconds
}

// NewTime creates a new Time instance.
// Corresponds to C++ Time() constructor.
func NewTime() *Time {
	return &Time{}
}

// Update calculates the delta time and updates the last time.
// Corresponds to C++ Time::update(double currentTime)
func (t *Time) Update(currentTime float64) {
	t.frame++
	t.delta = currentTime - t.lastTime
	t.lastTime = currentTime
}

// Step advances time by a given delta.
// Corresponds to C++ Time::step(double delta)
func (t *Time) Step(delta float64) {
	t.frame++
	t.lastTime += delta
	t.delta = delta
}

// Set sets the last recorded time.
// Corresponds to C++ Time::set(double currentTime)
func (t *Time) Set(currentTime float64) {
	t.lastTime = currentTime
}

// GetDelta returns the time elapsed since the last update/step, in seconds.
// Corresponds to C++ Time::getDelta()
func (t *Time) GetDelta() float64 {
	return t.delta
}

// GetTime returns the total time elapsed since the timer was set/started, in seconds.
// Corresponds to C++ Time::getTime()
func (t *Time) GetTime() float64 {
	return t.lastTime
}
