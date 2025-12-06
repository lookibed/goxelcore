package engine

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goxelcore/io" // For io.Path, io.FileTimeType
	// "goxelcore/data" // For dv.value (stubbed for now)
)

// EnginePaths corresponds to C++ EnginePaths class.
type EnginePaths struct {
	ResPaths ResPaths // C++ has a ResPaths member

	resourcesFolder string // C++ std::filesystem::path
	userFilesFolder string // C++ std::filesystem::path
	projectFolder   string // C++ std::filesystem::path
	currentWorldFolder io.Path // C++ io::path
	scriptFolder    *string // C++ std::optional<std::filesystem::path>
	entryPoints     []PathsRoot // C++ std::vector<PathsRoot>
	writeables      map[string]string // C++ std::unordered_map<std::string, std::string> (stubbed for now)
	mounted         []string // C++ std::vector<std::string> (stubbed for now)
}

// Static inline io::path constants
var (
	ConfigDefaults = io.NewPath("config/defaults.toml")
	ControlsFile   = io.NewPath("user:controls.toml")
	SettingsFile   = io.NewPath("user:settings.toml")
)

// NewEnginePaths creates a new EnginePaths instance.
// Corresponds to C++ EnginePaths(CoreParameters& params)
func NewEnginePaths(params *CoreParameters) *EnginePaths {
	ep := &EnginePaths{}

	// Resolve resourcesFolder and userFilesFolder based on OS conventions
	// For simplicity, hardcode for now or derive from current working directory
	ep.resourcesFolder = "res" // Relative to executable
	ep.userFilesFolder = ".goxelcore" // In user home directory or current dir

	// In a real app, you'd use os.UserHomeDir() to get platform-specific user folder
	userHomeDir, err := os.UserHomeDir()
	if err == nil {
		ep.userFilesFolder = filepath.Join(userHomeDir, ".goxelcore")
	} else {
		log.Printf("Warning: Could not get user home directory: %v. Using .goxelcore in current directory.\n", err)
		ep.userFilesFolder = filepath.Join(".", ".goxelcore")
	}

	// Make sure userFilesFolder exists
	if err := os.MkdirAll(ep.userFilesFolder, 0755); err != nil {
		log.Printf("Error creating user files folder: %v\n", err)
	}

	ep.projectFolder = params.ProjectFolder // Direct use of params.ProjectFolder
	// ep.scriptFolder = nil // Not supported yet

	// Initialize ResPaths
	ep.entryPoints = []PathsRoot{
		NewPathsRoot("res", io.NewPath(ep.resourcesFolder)),
		NewPathsRoot("user", io.NewPath(ep.userFilesFolder)),
		NewPathsRoot("project", io.NewPath(ep.projectFolder)),
	}
	ep.ResPaths = *NewResPaths(ep.entryPoints)
	
	return ep
}

// GetResourcesFolder returns the resources folder path.
// Corresponds to C++ EnginePaths::getResourcesFolder()
func (ep *EnginePaths) GetResourcesFolder() string {
	return ep.resourcesFolder
}

// GetUserFilesFolder returns the user files folder path.
// Corresponds to C++ EnginePaths::getUserFilesFolder()
func (ep *EnginePaths) GetUserFilesFolder() string {
	return ep.userFilesFolder
}

// GetNewScreenshotFile generates a new screenshot file path.
// Corresponds to C++ EnginePaths::getNewScreenshotFile()
func (ep *EnginePaths) GetNewScreenshotFile(ext string) io.Path {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("screenshot_%s.%s", timestamp, ext)
	return ep.currentWorldFolder.Join(filename) // Assuming currentWorldFolder is set
}

// Mount a resource path (stub for now).
// Corresponds to C++ EnginePaths::mount()
func (ep *EnginePaths) Mount(file io.Path) string {
	log.Printf("EnginePaths.Mount: Stub for %s\n", file.String())
	ep.mounted = append(ep.mounted, file.String())
	return file.String() // Return original for stub
}

// Unmount a resource path (stub for now).
// Corresponds to C++ EnginePaths::unmount()
func (ep *EnginePaths) Unmount(name string) {
	log.Printf("EnginePaths.Unmount: Stub for %s\n", name)
	// Remove from mounted list
	for i, m := range ep.mounted {
		if m == name {
			ep.mounted = append(ep.mounted[:i], ep.mounted[i+1:]...)
			break
		}
	}
}

// CreateWriteableDevice creates a writeable device (stub for now).
// Corresponds to C++ EnginePaths::createWriteableDevice()
func (ep *EnginePaths) CreateWriteableDevice(name string) string {
	log.Printf("EnginePaths.CreateWriteableDevice: Stub for %s\n", name)
	ep.writeables[name] = "" // Placeholder
	return name
}

// CreateMemoryDevice creates a memory device (stub for now).
// Corresponds to C++ EnginePaths::createMemoryDevice()
func (ep *EnginePaths) CreateMemoryDevice() string {
	log.Println("EnginePaths.CreateMemoryDevice: Stub")
	return "memory:/" // Example memory entry point
}

// SetEntryPoints sets new entry points for resource paths.
// Corresponds to C++ EnginePaths::setEntryPoints()
func (ep *EnginePaths) SetEntryPoints(entryPoints []PathsRoot) {
	ep.entryPoints = entryPoints
	ep.ResPaths = *NewResPaths(ep.entryPoints)
}

// ScanForWorlds scans for world folders (stub for now).
// Corresponds to C++ EnginePaths::scanForWorlds()
func (ep *EnginePaths) ScanForWorlds() []io.Path {
	log.Println("EnginePaths.ScanForWorlds: Stub")
	return []io.Path{}
}

// ParsePath parses a path string into entry point and path part.
// Corresponds to C++ static std::tuple<std::string, std::string> parsePath(std::string_view view)
func ParsePath(view string) (string, string) {
	p := io.NewPath(view)
	return p.EntryPoint(), p.PathPart()
}

// GetWorldsFolder returns the path to the worlds folder. (Stub for now)
func (ep *EnginePaths) GetWorldsFolder() io.Path {
	log.Println("EnginePaths.GetWorldsFolder: Stub")
	return io.NewPath("user:worlds") // Example
}

// GetWorldFolderByName returns the path to a specific world folder. (Stub for now)
func (ep *EnginePaths) GetWorldFolderByName(name string) io.Path {
	log.Printf("EnginePaths.GetWorldFolderByName: Stub for %s\n", name)
	return ep.GetWorldsFolder().Join(name)
}

// SetCurrentWorldFolder sets the current world folder.
func (ep *EnginePaths) SetCurrentWorldFolder(folder io.Path) {
	ep.currentWorldFolder = folder
}
