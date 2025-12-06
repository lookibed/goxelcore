package coders

import (
	"log"
	"strings"

	"goxelcore/graphics/core" // For core.Param
	// "goxelcore/io" // For io.path (stubbed for now)
	// "goxelcore/engine" // For ResPaths (stubbed for now)
)

// ResPathsStub is a placeholder for ResPaths class.
// Corresponds to C++ ResPaths class.
type ResPathsStub struct {
	// Add fields as needed if ResPaths needs to be ported later
}

// ProcessingResult corresponds to C++ GLSLExtension::ProcessingResult struct.
type ProcessingResult struct {
	Code string
	Params map[string]core.Param // Assuming core.Param is defined in graphics/core package
}


// GLSLExtension corresponds to C++ GLSLExtension class in voxelcore/src/coders/GLSLExtension.hpp
type GLSLExtension struct {
	headers map[string]ProcessingResult
	defines map[string]string

	paths *ResPathsStub // Placeholder for ResPaths
	traceOutput bool
}


// NewGLSLExtension creates a new GLSLExtension instance.
func NewGLSLExtension() *GLSLExtension {
	return &GLSLExtension{
		headers: make(map[string]ProcessingResult),
		defines: make(map[string]string),
	}
}

// SetPaths sets the resource paths.
// Corresponds to C++ GLSLExtension::setPaths()
func (ext *GLSLExtension) SetPaths(paths *ResPathsStub) {
	ext.paths = paths
}

// SetTraceOutput enables or disables trace output.
// Corresponds to C++ GLSLExtension::setTraceOutput()
func (ext *GLSLExtension) SetTraceOutput(enabled bool) {
	ext.traceOutput = enabled
}

// Define adds or updates a GLSL define.
// Corresponds to C++ GLSLExtension::define()
func (ext *GLSLExtension) Define(name, value string) {
	ext.defines[name] = value
}

// Undefine removes a GLSL define.
// Corresponds to C++ GLSLExtension::undefine()
func (ext *GLSLExtension) Undefine(name string) {
	delete(ext.defines, name)
}

// SetDefined sets whether a GLSL define is present.
// Corresponds to C++ GLSLExtension::setDefined()
func (ext *GLSLExtension) SetDefined(name string, defined bool) {
	if defined {
		ext.defines[name] = "" // Value might not matter for simple defines
	} else {
		delete(ext.defines, name)
	}
}

// AddHeader adds a processed header.
// Corresponds to C++ GLSLExtension::addHeader()
func (ext *GLSLExtension) AddHeader(name string, header ProcessingResult) {
	ext.headers[name] = header
}

// GetHeader returns a processed header.
// Corresponds to C++ GLSLExtension::getHeader()
func (ext *GLSLExtension) GetHeader(name string) (ProcessingResult, bool) {
	result, ok := ext.headers[name]
	return result, ok
}

// GetDefine returns the value of a GLSL define.
// Corresponds to C++ GLSLExtension::getDefine()
func (ext *GLSLExtension) GetDefine(name string) string {
	return ext.defines[name]
}

// GetDefines returns all GLSL defines.
// Corresponds to C++ GLSLExtension::getDefines()
func (ext *GLSLExtension) GetDefines() map[string]string {
	return ext.defines
}

// HasHeader checks if a header exists.
// Corresponds to C++ GLSLExtension::hasHeader()
func (ext *GLSLExtension) HasHeader(name string) bool {
	_, ok := ext.headers[name]
	return ok
}

// HasDefine checks if a define exists.
// Corresponds to C++ GLSLExtension::hasDefine()
func (ext *GLSLExtension) HasDefine(name string) bool {
	_, ok := ext.defines[name]
	return ok
}

// LoadHeader (stub for now).
// Corresponds to C++ GLSLExtension::loadHeader()
func (ext *GLSLExtension) LoadHeader(name string) {
	log.Printf("GLSLExtension.LoadHeader: Stub - %s\n", name)
}

// Process processes GLSL source code (stub for now).
// Corresponds to C++ GLSLExtension::process()
func (ext *GLSLExtension) Process(file string, source string, isHeader bool, defines []string) (ProcessingResult, error) {
	// For a stub, just return the source as is, and an empty param map.
	// Actual preprocessing logic would go here.
	if ext.traceOutput {
		log.Printf("GLSLExtension.Process (stub) for file %s, isHeader: %t, defines: %v\n", file, isHeader, defines)
	}
	// Apply explicit defines to source for basic replacement
	processedCode := source
	for _, def := range defines {
		parts := strings.SplitN(def, "=", 2)
		if len(parts) == 2 {
			processedCode = strings.ReplaceAll(processedCode, "#define "+parts[0], "#define "+def)
		} else {
			processedCode = strings.ReplaceAll(processedCode, "#define "+def, "#define "+def+"\n") // Simple define
		}
	}


	return ProcessingResult{
		Code: processedCode,
		Params: make(map[string]core.Param), // Use the correct Params map type
	}, nil
}

// GLSL_VERSION constant.
// Corresponds to C++ static inline std::string VERSION = "330 core";
const GLSL_VERSION = "330 core"
