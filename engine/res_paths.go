package engine

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"goxelcore/io" // For io.Path
	// "goxelcore/data" // For dv.value (stubbed for now)
)

// PathsRoot corresponds to C++ PathsRoot struct.
type PathsRoot struct {
	Name string
	Path io.Path
}

// NewPathsRoot creates a new PathsRoot instance.
func NewPathsRoot(name string, path io.Path) PathsRoot {
	return PathsRoot{
		Name: name,
		Path: path,
	}
}

// ResPaths corresponds to C++ ResPaths class.
// Manages a collection of resource roots.
type ResPaths struct {
	roots []PathsRoot
}

// NewResPaths creates a new ResPaths instance with given roots.
// Corresponds to C++ ResPaths(std::vector<PathsRoot> roots)
func NewResPaths(roots []PathsRoot) *ResPaths {
	return &ResPaths{
		roots: roots,
	}
}

// Find attempts to resolve a filename (possibly with an entry point) to an actual io.Path.
// Corresponds to C++ ResPaths::find()
func (rp *ResPaths) Find(filename string) io.Path {
	p := io.NewPath(filename)
	if p.IsEmpty() {
		return io.Path{}
	}

	if p.EntryPoint() != "" { // Has an entry point
		for _, root := range rp.roots {
			if root.Name == p.EntryPoint() {
				resolvedPath := root.Path.Join(p.PathPart())
				if io.Exists(resolvedPath) {
					return resolvedPath
				}
			}
		}
	} else { // No entry point, search all roots
		for _, root := range rp.roots {
			resolvedPath := root.Path.Join(p.PathPart())
			if io.Exists(resolvedPath) {
				return resolvedPath
			}
		}
		// Also check relative to current working directory if not found in roots
		if io.Exists(io.NewPath(p.PathPart())) {
			return io.NewPath(p.PathPart())
		}
	}
	log.Printf("ResPaths: Could not find '%s'\n", filename)
	return io.Path{} // Return empty path if not found
}

// FindRaw attempts to resolve a filename to its raw string path.
// Corresponds to C++ ResPaths::findRaw()
func (rp *ResPaths) FindRaw(filename string) string {
	return rp.Find(filename).String()
}

// Listdir lists contents of a directory. (Stub for now, needs io.Device integration)
// Corresponds to C++ ResPaths::listdir()
func (rp *ResPaths) Listdir(folder string) []io.Path {
	log.Printf("ResPaths: Listdir for '%s' is a stub.\n", folder)
	return []io.Path{}
}

// ListdirRaw lists raw string contents of a directory. (Stub for now)
// Corresponds to C++ ResPaths::listdirRaw()
func (rp *ResPaths) ListdirRaw(folder string) []string {
	log.Printf("ResPaths: ListdirRaw for '%s' is a stub.\n", folder)
	return []string{}
}

// CollectRoots returns all known root paths.
// Corresponds to C++ ResPaths::collectRoots()
func (rp *ResPaths) CollectRoots() []io.Path {
	var collected []io.Path
	for _, root := range rp.roots {
		collected = append(collected, root.Path)
	}
	return collected
}

// ReadCombinedList reads and combines lists from all roots. (Stub)
// Corresponds to C++ ResPaths::readCombinedList()
func (rp *ResPaths) ReadCombinedList(file string) interface{} { // dv::value in C++
	log.Printf("ResPaths: ReadCombinedList for '%s' is a stub.\n", file)
	return nil // Placeholder for dv.value
}

// ReadCombinedObject reads and combines objects from all roots. (Stub)
// Corresponds to C++ ResPaths::readCombinedObject()
func (rp *ResPaths) ReadCombinedObject(file string, deep bool) interface{} { // dv::value in C++
	log.Printf("ResPaths: ReadCombinedObject for '%s' (deep: %t) is a stub.\n", file, deep)
	return nil // Placeholder for dv.value
}
