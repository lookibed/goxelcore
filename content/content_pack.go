package content

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"goxelcore/io"     // For io.Path
	"goxelcore/util"   // For Quote (if needed, or just format strings)
	// "goxelcore/data" // For dv.value (stubbed for now)
)

// ContentPackError corresponds to C++ contentpack_error class.
type ContentPackError struct {
	PackID string
	Folder io.Path
	Reason string
}

func (e *ContentPackError) Error() string {
	return fmt.Sprintf("contentpack error [%s] in folder '%s': %s", e.PackID, e.Folder.String(), e.Reason)
}

// VersionOperator corresponds to C++ VersionOperator enum.
type VersionOperator int

const (
	VersionOperatorEqual VersionOperator = iota
	VersionOperatorGreater
	VersionOperatorLess
	VersionOperatorGreaterOrEqual
	VersionOperatorLessOrEqual
)

// String returns the string representation of the VersionOperator.
func (op VersionOperator) String() string {
	switch op {
	case VersionOperatorEqual:
		return "="
	case VersionOperatorGreater:
		return ">"
	case VersionOperatorLess:
		return "<"
	case VersionOperatorGreaterOrEqual:
		return ">="
	case VersionOperatorLessOrEqual:
		return "<="
	default:
		return "UNKNOWN"
	}
}

// DependencyLevel corresponds to C++ DependencyLevel enum.
type DependencyLevel int

const (
	DependencyLevelRequired DependencyLevel = iota
	DependencyLevelOptional
	DependencyLevelWeak
)

// DependencyPack corresponds to C++ DependencyPack struct.
type DependencyPack struct {
	Level   DependencyLevel
	ID      string
	Version string
	Op      VersionOperator
}

// ContentPackStats corresponds to C++ ContentPackStats struct.
type ContentPackStats struct {
	TotalBlocks   goxelcore.size_t
	TotalItems    goxelcore.size_t
	TotalEntities goxelcore.size_t
}

// HasSavingContent checks if the pack has content that needs saving.
func (cps *ContentPackStats) HasSavingContent() bool {
	return cps.TotalBlocks+cps.TotalItems+cps.TotalEntities > 0
}

// ContentType corresponds to C++ ContentType enum from content_fwd.hpp.
type ContentType int

const (
	ContentTypeNone ContentType = iota
	ContentTypeBlock
	ContentTypeItem
	ContentTypeEntity
	ContentTypeGenerator
)

// String returns the string representation of the ContentType.
func (ct ContentType) String() string {
	switch ct {
	case ContentTypeNone:
		return "none"
	case ContentTypeBlock:
		return "block"
	case ContentTypeItem:
		return "item"
	case ContentTypeEntity:
		return "entity"
	case ContentTypeGenerator:
		return "generator"
	default:
		return "unknown"
	}
}


// ContentPack corresponds to C++ ContentPack struct.
type ContentPack struct {
	ID          string
	Title       string
	Version     string
	Creator     string
	Description string
	Folder      io.Path
	Dependencies []DependencyPack
	Source      string
}

// Static inline consts.
const (
	PackageFilename = "package.json"
	ContentFilename = "content.json"
)

// io.Path consts
var (
	BlocksFolder    = io.NewPath("blocks")
	ItemsFolder     = io.NewPath("items")
	EntitiesFolder  = io.NewPath("entities")
	GeneratorsFolder = io.NewPath("generators")
)

// RESERVED_NAMES corresponds to C++ static const std::vector<std::string> RESERVED_NAMES.
var ReservedNames = []string{"core"} // Example, add more if known

// GetContentFile returns the path to the content.json file within the pack.
// Corresponds to C++ ContentPack::getContentFile().
func (cp *ContentPack) GetContentFile() io.Path {
	return cp.Folder.Join(ContentFilename)
}

// LoadStats loads content pack statistics. (Stub for now)
// Corresponds to C++ ContentPack::loadStats().
func (cp *ContentPack) LoadStats() *ContentPackStats {
	log.Printf("ContentPack.LoadStats: Stub for pack '%s'\n", cp.ID)
	return &ContentPackStats{}
}

// IsPack checks if a folder is a valid content pack. (Stub for now)
// Corresponds to C++ static bool is_pack(const io::path& folder).
func IsPack(folder io.Path) bool {
	// A basic check would be to see if 'package.json' exists in the folder.
	packageFile := folder.Join(PackageFilename)
	return io.IsRegularFile(packageFile)
}

// Read reads a content pack from a folder. (Stub for now)
// Corresponds to C++ static ContentPack read(const io::path& folder).
func ReadContentPack(folder io.Path) (ContentPack, error) {
	log.Printf("ContentPack.Read: Stub for folder '%s'\n", folder.String())
	// Real implementation would parse package.json.
	return ContentPack{ID: folder.Name(), Folder: folder}, nil
}

// ScanFolder scans a folder for content packs. (Stub for now)
// Corresponds to C++ static void scanFolder(...).
func ScanFolder(folder io.Path) []ContentPack {
	log.Printf("ContentPack.ScanFolder: Stub for folder '%s'\n", folder.String())
	return []ContentPack{}
}

// WorldPacksList reads the list of packs from a world directory. (Stub for now)
// Corresponds to C++ static std::vector<std::string> worldPacksList(...).
func WorldPacksList(folder io.Path) []string {
	log.Printf("ContentPack.WorldPacksList: Stub for folder '%s'\n", folder.String())
	return []string{}
}

// FindPack finds a content pack in given paths. (Stub for now)
// Corresponds to C++ static io::path findPack(...).
func FindPack(paths EnginePathsInterface, worldDir io.Path, name string) io.Path {
	log.Printf("ContentPack.FindPack: Stub for name '%s'\n", name)
	return io.Path{}
}

// CreateCore creates the default "core" content pack. (Stub for now)
// Corresponds to C++ static ContentPack createCore().
func CreateCoreContentPack() ContentPack {
	log.Println("ContentPack.CreateCore: Stub")
	return ContentPack{ID: "core", Title: "Core", Version: "1.0", Description: "Core game content."}
}

// GetFolderForContentType returns the subfolder for a given ContentType.
// Corresponds to C++ static inline io::path getFolderFor(ContentType type).
func GetFolderForContentType(ct ContentType) io.Path {
	switch ct {
	case ContentTypeBlock:
		return BlocksFolder
	case ContentTypeItem:
		return ItemsFolder
	case ContentTypeEntity:
		return EntitiesFolder
	case ContentTypeGenerator:
		return GeneratorsFolder
	case ContentTypeNone:
		return io.NewPath("")
	default:
		return io.NewPath("") // Or return an error
	}
}

// ContentPackRuntime corresponds to C++ ContentPackRuntime class.
type ContentPackRuntime struct {
	Info  ContentPack
	Stats ContentPackStats
	Env   interface{} // scriptenv is std::shared_ptr<int>, use interface{}
}

// NewContentPackRuntime creates a new ContentPackRuntime instance.
// Corresponds to C++ ContentPackRuntime(ContentPack info, scriptenv env) constructor.
func NewContentPackRuntime(info ContentPack, env interface{}) *ContentPackRuntime {
	return &ContentPackRuntime{
		Info: info,
		Stats: ContentPackStats{}, // Default stats
		Env:   env,
	}
}

// GetStats returns the content pack statistics.
func (cpr *ContentPackRuntime) GetStats() *ContentPackStats {
	return &cpr.Stats
}

// GetStatsWriteable returns a writeable reference to content pack statistics.
func (cpr *ContentPackRuntime) GetStatsWriteable() *ContentPackStats {
	return &cpr.Stats
}

// GetId returns the ID of the content pack.
func (cpr *ContentPackRuntime) GetId() string {
	return cpr.Info.ID
}

// GetInfo returns the ContentPack information.
func (cpr *ContentPackRuntime) GetInfo() ContentPack {
	return cpr.Info
}

// GetEnvironment returns the script environment.
func (cpr *ContentPackRuntime) GetEnvironment() interface{} {
	return cpr.Env
}
