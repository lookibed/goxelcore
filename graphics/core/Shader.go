package core

import (
	"fmt"
	"log"
	"strings"

	gl "github.com/go-gl/gl/v3.3-core/gl"
)

// Shader represents an OpenGL shader program.
// Corresponds to C++ Shader class in voxelcore/src/graphics/core/Shader.hpp
type Shader struct {
	id             uint32
	uniformLocations map[string]int32
	vertexSource   Source
	fragmentSource Source
	// C++ has `static Shader* used;`
	// In Go, we can manage the 'used' shader globally or via the Engine.
}

// Source stores shader source code and file path.
// Corresponds to C++ Shader::Source struct.
type Source struct {
	File string
	Code string
}

// TODO: Port GLSLExtension preprocessor later if needed.
// var preprocessor *GLSLExtension // static GLSLExtension* preprocessor;

// usedShader tracks the currently active shader.
// Corresponds to C++ static Shader* used;
var usedShader *Shader

// NewShader creates a new Shader instance.
// Corresponds to C++ Shader constructor.
func NewShader(id uint32, vertexSource, fragmentSource Source) *Shader {
	return &Shader{
		id:               id,
		uniformLocations: make(map[string]int32),
		vertexSource:     vertexSource,
		fragmentSource:   fragmentSource,
	}
}

// Use activates the shader program.
// Corresponds to C++ Shader::use()
func (s *Shader) Use() {
	if usedShader != s {
		gl.UseProgram(s.id)
		usedShader = s
	}
}

// getUniformLocation retrieves the location of a uniform variable.
// Corresponds to C++ Shader::getUniformLocation()
func (s *Shader) getUniformLocation(name string) int32 {
	if loc, ok := s.uniformLocations[name]; ok {
		return loc
	}
	loc := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	if loc == -1 {
		log.Printf("Warning: Uniform '%s' not found in shader program %d\n", name, s.id)
	}
	s.uniformLocations[name] = loc
	return loc
}

// UniformMatrix4fv sets a mat4 uniform.
// Corresponds to C++ Shader::uniformMatrix(const std::string&, const glm::mat4& matrix)
// For now, using a direct []float32 for matrix data (Go doesn't have glm directly).
func (s *Shader) UniformMatrix4fv(name string, matrix []float32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.UniformMatrix4fv(loc, 1, false, &matrix[0])
	}
}

// UniformMatrix3fv sets a mat3 uniform.
// Corresponds to C++ Shader::uniformMatrix(const std::string&, const glm::mat3& matrix)
func (s *Shader) UniformMatrix3fv(name string, matrix []float32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.UniformMatrix3fv(loc, 1, false, &matrix[0])
	}
}

// Uniform1i sets an int uniform.
// Corresponds to C++ Shader::uniform1i()
func (s *Shader) Uniform1i(name string, x int32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform1i(loc, x)
	}
}

// Uniform1f sets a float uniform.
// Corresponds to C++ Shader::uniform1f()
func (s *Shader) Uniform1f(name string, x float32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform1f(loc, x)
	}
}

// Uniform2f sets a vec2 uniform (float).
// Corresponds to C++ Shader::uniform2f(const std::string&, float x, float y)
func (s *Shader) Uniform2f(name string, x, y float32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform2f(loc, x, y)
	}
}

// Uniform2fv sets a vec2 uniform (float).
// Corresponds to C++ Shader::uniform2f(const std::string&, const glm::vec2& xy)
func (s *Shader) Uniform2fv(name string, xy []float32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform2fv(loc, 1, &xy[0])
	}
}

// Uniform2i sets an ivec2 uniform.
// Corresponds to C++ Shader::uniform2i(const std::string&, const glm::ivec2& xy)
func (s *Shader) Uniform2i(name string, x, y int32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform2i(loc, x, y)
	}
}

// Uniform3f sets a vec3 uniform (float).
// Corresponds to C++ Shader::uniform3f(const std::string&, float x, float y, float z)
func (s *Shader) Uniform3f(name string, x, y, z float32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform3f(loc, x, y, z)
	}
}

// Uniform3fv sets a vec3 uniform (float).
// Corresponds to C++ Shader::uniform3f(const std::string&, const glm::vec3& xyz)
func (s *Shader) Uniform3fv(name string, xyz []float32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform3fv(loc, 1, &xyz[0])
	}
}

// Uniform4fv sets a vec4 uniform (float).
// Corresponds to C++ Shader::uniform4f(const std::string&, const glm::vec4& xyzw)
func (s *Shader) Uniform4fv(name string, xyzw []float32) {
	s.Use()
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.Uniform4fv(loc, 1, &xyzw[0])
	}
}

// TODO: Implement uniform1v, uniform2v, uniform3v, uniform4v
// These require slice versions of the Uniform functions.

// Recompile (stub for now).
// Corresponds to C++ Shader::recompile()
func (s *Shader) Recompile(defines []string) {
	log.Printf("Shader.Recompile: Not fully implemented (preprocessing and recompilation needed for defines: %v).\n", defines)
	// For now, this is a stub. Real implementation would involve:
	// 1. Preprocessing vertexSource.Code and fragmentSource.Code with defines
	// 2. Compiling new shaders
	// 3. Linking new program
	// 4. Replacing old program ID if successful
}

// CreateShader compiles and links a shader program.
// Corresponds to C++ static std::unique_ptr<Shader> create(...)
func CreateShader(vertexSource, fragmentSource Source) (*Shader, error) {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, gl.Str(vertexSource.Code+"\x00"))
	gl.CompileShader(vertexShader)
	if gl.GetShaderi(vertexShader, gl.COMPILE_STATUS) == 0 {
		log.Printf("Vertex shader compilation failed for %s:\n%s\n",
			vertexSource.File, gl.GetShaderInfoLog(vertexShader))
		return nil, fmt.Errorf("vertex shader compilation failed")
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, gl.Str(fragmentSource.Code+"\x00"))
	gl.CompileShader(fragmentShader)
	if gl.GetShaderi(fragmentShader, gl.COMPILE_STATUS) == 0 {
		log.Printf("Fragment shader compilation failed for %s:\n%s\n",
			fragmentSource.File, gl.GetShaderInfoLog(fragmentShader))
		return nil, fmt.Errorf("fragment shader compilation failed")
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	if gl.GetProgrami(program, gl.LINK_STATUS) == 0 {
		log.Printf("Shader program linking failed:\n%s\n", gl.GetProgramInfoLog(program))
		return nil, fmt.Errorf("shader program linking failed")
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return NewShader(program, vertexSource, fragmentSource), nil
}

// GetUsedShader returns the currently active shader.
// Corresponds to C++ static Shader& getUsed()
func GetUsedShader() *Shader {
	return usedShader
}

// Delete cleans up the shader program.
// No direct C++ destructor equivalent call, but good practice to have.
func (s *Shader) Delete() {
	gl.DeleteProgram(s.id)
	if usedShader == s {
		usedShader = nil
	}
}
