package io

import (
	"errors"
	"fmt"
	"io/ioutil" // For ReadFile, WriteFile
	"log"
	"os"
	"path/filepath" // For filesystem operations
	"time"
)

var (
	// ErrFileNotFound indicates that a file was not found.
	ErrFileNotFound = errors.New("file not found")
	// ErrAccessDenied indicates an access violation.
	ErrAccessDenied = errors.New("access denied")
	// ErrUnsupportedOperation indicates an unsupported I/O operation.
	ErrUnsupportedOperation = errors.New("unsupported I/O operation")
)

// Device interface corresponds to C++ io::Device.
// For now, this is a placeholder.
type Device interface {
	// Read(p Path) (bytes []byte, err error)
	// Write(p Path, data []byte) (err error)
	// Exists(p Path) bool
	// IsDir(p Path) bool
	// Stat(p Path) (os.FileInfo, error)
}

// Global device manager (stub for now)
var devices = make(map[string]Device)

// SetDevice corresponds to C++ io::set_device.
func SetDevice(name string, device Device) {
	devices[name] = device
	log.Printf("IO: Device '%s' set (stub).\n", name)
}

// RemoveDevice corresponds to C++ io::remove_device.
func RemoveDevice(name string) {
	delete(devices, name)
	log.Printf("IO: Device '%s' removed (stub).\n", name)
}

// GetDevice corresponds to C++ io::get_device.
func GetDevice(name string) Device {
	return devices[name]
}

// RequireDevice corresponds to C++ io::require_device.
func RequireDevice(name string) (Device, error) {
	if dev, ok := devices[name]; ok {
		return dev, nil
	}
	return nil, fmt.Errorf("required device '%s' not found", name)
}

// CreateSubdevice corresponds to C++ io::create_subdevice.
func CreateSubdevice(name string, parent string, root Path) {
	log.Printf("IO: Create subdevice '%s' under '%s' with root '%s' (stub).\n", name, parent, root.String())
}

// ReadString reads the entire content of a file as a string.
// Corresponds to C++ io::read_string.
func ReadString(p Path) (string, error) {
	// For now, assume a direct filesystem path if no entry point
	// A real implementation would resolve 'entryPoint:pathPart' to a real filesystem path
	// or interact with the corresponding io::Device.
	filePath := resolvePath(p)
	
data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrFileNotFound
		}
		if os.IsPermission(err) {
			return "", ErrAccessDenied
		}
		return "", err
	}
	return string(data), nil
}

// WriteString writes content to a file.
// Corresponds to C++ io::write_string.
func WriteString(p Path, content string) error {
	filePath := resolvePath(p)
	
	// Ensure parent directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", filePath, err)
	}

	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		if os.IsPermission(err) {
			return ErrAccessDenied
		}
	}
	return err
}

// Exists checks if a file or directory exists.
// Corresponds to C++ io::exists.
func Exists(p Path) bool {
	filePath := resolvePath(p)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// IsRegularFile checks if a path points to a regular file.
// Corresponds to C++ io::is_regular_file.
func IsRegularFile(p Path) bool {
	filePath := resolvePath(p)
	info, err := os.Stat(filePath)
	return err == nil && info.Mode().IsRegular()
}

// IsDirectory checks if a path points to a directory.
// Corresponds to C++ io::is_directory.
func IsDirectory(p Path) bool {
	filePath := resolvePath(p)
	info, err := os.Stat(filePath)
	return err == nil && info.Mode().IsDir()
}

// CreateDirectory creates a single directory.
// Corresponds to C++ io::create_directory.
func CreateDirectory(p Path) error {
	filePath := resolvePath(p)
	return os.Mkdir(filePath, 0755)
}

// CreateDirectories creates a directory and any necessary parent directories.
// Corresponds to C++ io::create_directories.
func CreateDirectories(p Path) error {
	filePath := resolvePath(p)
	return os.MkdirAll(filePath, 0755)
}

// Remove removes a file or empty directory.
// Corresponds to C++ io::remove.
func Remove(p Path) error {
	filePath := resolvePath(p)
	return os.Remove(filePath)
}

// Copy copies a file from src to dst.
// Corresponds to C++ io::copy.
func Copy(src, dst Path) error {
	srcPath := resolvePath(src)
	dstPath := resolvePath(dst)
	
	input, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}
	
dir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", dstPath, err)
	}

	err = ioutil.WriteFile(dstPath, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

// resolvePath is a helper to convert custom Path to OS-specific absolute path for basic file ops.
// In a full implementation, this would involve a proper virtual filesystem or device resolution.
func resolvePath(p Path) string {
	// For simplicity, assume "res:" maps to a "res" directory in the current working directory.
	// And other entry points map directly or use a base path.
	if p.EntryPoint() == "res" {
		return filepath.Join("res", p.PathPart())
	}
	// For empty entry point or other entry points, treat pathPart as relative to current dir for now
	return p.PathPart()
}

// Other io functions (rafile, directory_iterator, JSON/TOML/BJSON parsing, Device management)
// are complex and will be stubbed or omitted for now.
// For example:
/*
func WriteJSON(p Path, obj interface{}, nice bool) error {
	log.Printf("IO: WriteJSON stub for %s\n", p.String())
	return ErrUnsupportedOperation
}
func ReadJSON(p Path) (interface{}, error) {
	log.Printf("IO: ReadJSON stub for %s\n", p.String())
	return nil, ErrUnsupportedOperation
}
// ... etc.
*/
