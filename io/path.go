package io

import (
	"path/filepath"
	"strings"
	"time"
)

// file_time_type corresponds to C++ std::filesystem::file_time_type
type FileTimeType = time.Time

// AccessError corresponds to C++ io::access_error
type AccessError struct {
	Message string
}

func (e *AccessError) Error() string {
	return "access error: " + e.Message
}

// Path corresponds to C++ io::path class.
// It represents a project-specific path with an "entry_point:path" scheme.
type Path struct {
	entryPoint string
	pathPart   string // Path part after the entryPoint
	fullPath   string // Cached full string representation
}

// NewPath creates a new Path instance.
// Corresponds to C++ io::path(std::string str) constructor.
func NewPath(p string) Path {
	// Normalize path separators to '/' as in C++ version
	p = strings.ReplaceAll(p, "\", "/")

	colonPos := strings.Index(p, ":")
	if colonPos == -1 {
		return Path{
			entryPoint: "",
			pathPart:   p,
			fullPath:   p,
		}
	}
	return Path{
		entryPoint: p[:colonPos],
		pathPart:   p[colonPos+1:] ,
		fullPath:   p,
	}
}

// String returns the full string representation of the path.
// Corresponds to C++ io::path::string()
func (p Path) String() string {
	return p.fullPath
}

// EntryPoint returns the entry point of the path.
// Corresponds to C++ io::path::entryPoint()
func (p Path) EntryPoint() string {
	return p.entryPoint
}

// PathPart returns the path part after the entry point.
// Corresponds to C++ io::path::pathPart()
func (p Path) PathPart() string {
	return p.pathPart
}

// Name returns the last component of the path.
// Corresponds to C++ io::path::name()
func (p Path) Name() string {
	// Use filepath.Base for the pathPart to get the name
	return filepath.Base(p.pathPart)
}

// Stem returns the stem of the path (filename without extension).
// Corresponds to C++ io::path::stem()
func (p Path) Stem() string {
	name := p.Name()
	if ext := filepath.Ext(name); ext != "" {
		return name[:len(name)-len(ext)]
	}
	return name
}

// Extension returns the extension of the path.
// Corresponds to C++ io::path::extension()
func (p Path) Extension() string {
	return filepath.Ext(p.Name())
}

// Parent returns the parent path.
// Corresponds to C++ io::path::parent()
func (p Path) Parent() Path {
	if p.IsEmpty() {
		return NewPath("")
	}
	parentPathPart := filepath.Dir(p.pathPart)
	if p.entryPoint == "" {
		return NewPath(parentPathPart)
	}
	return NewPath(p.entryPoint + ":" + parentPathPart)
}

// IsEmpty checks if the path is empty.
// Corresponds to C++ io::path::empty()
func (p Path) IsEmpty() bool {
	return p.fullPath == ""
}

// IsEmptyOrInvalid checks if the path is empty or has no entry point.
// Corresponds to C++ io::path::emptyOrInvalid()
func (p Path) IsEmptyOrInvalid() bool {
	return p.fullPath == "" || p.entryPoint == ""
}

// Join concatenates a child path to the current path.
// Corresponds to C++ operator/
func (p Path) Join(child string) Path {
	child = strings.ReplaceAll(child, "\", "/")
	if p.IsEmpty() {
		return NewPath(child)
	}
	if p.entryPoint == "" {
		return NewPath(filepath.Join(p.fullPath, child))
	}
	return NewPath(p.entryPoint + ":" + filepath.Join(p.pathPart, child))
}

// PathsGenerator interface corresponds to C++ io::PathsGenerator.
type PathsGenerator interface {
	Next() (Path, bool) // Returns next path and true if successful, false if no more paths
}
