package devtools

import "goxelcore/engine"

// Editor is a placeholder for the C++ devtools::Editor class.
type Editor struct {
	engine *engine.Engine
	// Add fields/methods as needed
}

// NewEditor creates a new Editor instance.
func NewEditor(eng *engine.Engine) *Editor {
	return &Editor{engine: eng}
}

// LoadTools is a stub.
func (e *Editor) LoadTools() {
	e.engine.GetLogger().Info("Editor.LoadTools: Stub")
}
