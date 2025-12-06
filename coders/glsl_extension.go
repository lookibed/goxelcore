package coders

import (
	"fmt"
	"log"
	"strings"

	"goxelcore/engine"          // For engine.ResPaths
	"goxelcore/graphics/core" // For core.Param
	"goxelcore/io"            // For io.Path, io.ReadString
)

// ProcessingResult corresponds to C++ GLSLExtension::ProcessingResult struct.
type ProcessingResult struct {
	Code   string
	Params map[string]core.Param // Assuming core.Param is defined in graphics/core package
}

// GLSLExtension corresponds to C++ GLSLExtension class in voxelcore/src/coders/GLSLExtension.hpp
type GLSLExtension struct {
	headers map[string]ProcessingResult
	defines map[string]string

	paths       *engine.ResPaths // Changed from *ResPathsStub to *engine.ResPaths
	traceOutput bool
}

// GLSL_VERSION constant.
// Corresponds to C++ static inline std::string VERSION = "330 core";
const GLSL_VERSION = "330 core"

// NewGLSLExtension creates a new GLSLExtension instance.
func NewGLSLExtension() *GLSLExtension {
	return &GLSLExtension{
	headers: make(map[string]ProcessingResult),
	defines: make(map[string]string),
	}
}

// SetPaths sets the resource paths.
// Corresponds to C++ GLSLExtension::setPaths()
func (ext *GLSLExtension) SetPaths(paths *engine.ResPaths) {
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
		ext.define(name, "TRUE")
	} else {
		ext.undefine(name)
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

// LoadHeader loads and processes a GLSL header file.
// Corresponds to C++ GLSLExtension::loadHeader()
func (ext *GLSLExtension) LoadHeader(name string) error {
	if ext.paths == nil {
		return fmt.Errorf("ResPaths not set for GLSLExtension")
	}
	
	// Check if already loaded
	if _, ok := ext.headers[name]; ok {
		return nil
	}

	file := ext.paths.Find("shaders/lib/" + name + ".glsl")
	if file.IsEmpty() {
		return fmt.Errorf("GLSL header '%s' not found", name)
	}

	source, err := io.ReadString(file)
	if err != nil {
		return fmt.Errorf("failed to read GLSL header '%s': %w", file.String(), err)
	}
	
	// Add a placeholder to prevent infinite recursion
	ext.addHeader(name, ProcessingResult{Code: "// #include recursion guard for " + name + "\n"})

	result, err := ext.Process(file, source, true, nil) // Process as header, no extra defines
	if err != nil {
		return fmt.Errorf("failed to process GLSL header '%s': %w", file.String(), err)
	}
	ext.addHeader(name, result)
	return nil
}

// GLSLParser is an internal parser for GLSL extension directives.
type GLSLParser struct {
	*BasicParser
	glsl *GLSLExtension
	// C++ GLSLParser has `std::unordered_map<std::string, PostEffect::Param> params;`
	// but for now, we only care about #include.
	processedCode strings.Builder
	extraDefines []string
	isHeader bool
}

func NewGLSLParser(glsl *GLSLExtension, file io.Path, source string, isHeader bool, defines []string) *GLSLParser {
	p := &GLSLParser{
		BasicParser: NewBasicParser(file.String(), source),
		glsl: glsl,
		extraDefines: defines,
		isHeader: isHeader,
	}
	p.clikeComment = true // Enable C-like comments
	return p
}

func (p *GLSLParser) processIncludeDirective() error {
	p.SkipWhitespace(false)
	if p.PeekNoJump() != '<' {
		return p.error("'<' expected")
	}
	p.Skip(1)
	p.SkipWhitespace(false)
	headerName := p.ParseName()
	p.SkipWhitespace(false)
	if p.PeekNoJump() != '>' {
		return p.error("'>' expected")
	}
	p.Skip(1)
	p.SkipWhitespace(false)
	p.SkipLine()

	// Load and process the included header
	if err := p.glsl.LoadHeader(headerName); err != nil {
		return err
	}
	headerResult, _ := p.glsl.GetHeader(headerName)
	p.processedCode.WriteString(headerResult.Code)
	p.processedCode.WriteString(fmt.Sprintf("#line %d\n", p.line)) // Restore line number
	return nil
}

func (p *GLSLParser) processVersionDirective() error {
	// For now, just copy the version line to output
	p.processedCode.WriteString(fmt.Sprintf("#line %d\n", p.line)) // Add line directive before version
	p.processedCode.WriteString("#version ")
	p.processedCode.WriteString(p.ReadUntilEOL())
	p.processedCode.WriteString("\n")
	p.SkipLine()
	return nil
}

func (p *GLSLParser) processParamDirective() error {
	// This is a stub for now. Just skip the line.
	p.SkipLine()
	log.Println("GLSLParser: #param directive parsing is a stub.")
	return nil
}

func (p *GLSLParser) processPreprocessorDirective() error {
	p.Skip(1) // Skip '#'

	name := p.ParseName()
	if name == "" {
		return p.error("preprocessor directive name expected")
	}

	switch name {
	case "version":
		return p.processVersionDirective()
	case "include":
		return p.processIncludeDirective()
	case "param":
		return p.processParamDirective()
	default:
		// Other directives like #define, #if, #ifdef etc. are passed through for now
		// Or can be processed as needed
		// For now, simply copy the entire line to output
		p.processedCode.WriteString(fmt.Sprintf("#line %d\n", p.line))
		p.processedCode.WriteString("#" + name + " ")
		p.processedCode.WriteString(p.ReadUntilEOL())
		p.processedCode.WriteString("\n")
		p.SkipLine()
	}
	return nil
}

// Process processes the GLSL source code, handling directives.
// Corresponds to C++ GLSLExtension::process()
func (ext *GLSLExtension) Process(file io.Path, source string, isHeader bool, defines []string) (ProcessingResult, error) {
	parser := NewGLSLParser(ext, file, source, isHeader, defines)
	
	if !isHeader {
		parser.processedCode.WriteString(fmt.Sprintf("#version %s\n", GLSL_VERSION))
		// Add user-defined defines from C++ defines parameter
		for _, def := range defines {
			parser.processedCode.WriteString(fmt.Sprintf("#define %s\n", def))
		}
		// Add GLSLExtension's global defines
		for name, value := range ext.defines {
			if value != "" {
				parser.processedCode.WriteString(fmt.Sprintf("#define %s %s\n", name, value))
			} else {
				parser.processedCode.WriteString(fmt.Sprintf("#define %s\n", name))
			}
		}
		parser.processedCode.WriteString(fmt.Sprintf("#line %d\n", parser.line))
	}

	for parser.HasNext() {
		parser.SkipWhitespace(false)
		if !parser.HasNext() {
			break
		}
		
		if parser.Peek() == '# {
			startLine := parser.line
			if err := parser.processPreprocessorDirective(); err != nil {
				return ProcessingResult{}, err
			}
			// If a directive caused a newline, ensure #line is correct
			if parser.line > startLine {
				parser.processedCode.WriteString(fmt.Sprintf("#line %d\n", parser.line))
			}
		} else {
			// Copy regular GLSL code
			parser.processedCode.WriteString(parser.ReadUntilEOL())
			parser.processedCode.WriteString("\n")
			parser.Skip(1) // Skip newline
		}
	}

	result := ProcessingResult{
		Code: parser.processedCode.String(),
		Params: make(map[string]core.Param), // No param parsing yet
	}

	if ext.traceOutput {
		// trace_output(file, source, result) // Requires further io/filesystem integration
		log.Printf("GLSLExtension: Trace output for %s (stub): \n%s\n", file.String(), result.Code)
	}
	return result, nil
}