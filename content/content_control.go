package content

import (
	"log"

	"goxelcore/engine" // For EnginePaths and Project
	"goxelcore/io"     // For io.Path
	"goxelcore/window" // For Input
)

// ContentControl corresponds to C++ ContentControl class.
type ContentControl struct {
	paths *engine.EnginePaths // C++ uses reference, Go uses pointer
	input *window.InputSDL    // C++ uses Input*, Go uses concrete type for now
	content *ContentStub      // C++ uses std::unique_ptr<Content>
	postContent func()        // C++ uses std::function<void()>
	basePacks []string        // C++ uses std::vector<std::string>
	manager *PacksManagerStub // C++ uses std::unique_ptr<PacksManager>
	contentPacks []ContentPack // C++ uses std::vector<ContentPack>
	allPacks []ContentPack   // C++ uses std::vector<ContentPack> (includes 'core')
}

// NewContentControl creates a new ContentControl instance.
// Corresponds to C++ ContentControl constructor.
func NewContentControl(
	project *engine.Project, // C++ uses const Project&
	paths *engine.EnginePaths,
	input *window.InputSDL, // C++ uses Input*
	postContent func(),
) *ContentControl {
	cc := &ContentControl{
		paths: paths,
		input: input,
		postContent: postContent,
		basePacks: make([]string, 0),
		contentPacks: make([]ContentPack, 0),
		allPacks: make([]ContentPack, 0),
	}
	// C++ ContentControl constructor also initializes manager and content.
	// We'll stub those for now.
	cc.manager = NewPacksManagerStub(cc) // Pass self reference
	cc.content = NewContentStub()
	return cc
}

// Get returns the Content instance.
// Corresponds to C++ ContentControl::get().
func (cc *ContentControl) Get() *ContentStub {
	return cc.content
}

// GetBasePacks returns the list of base content pack IDs.
// Corresponds to C++ ContentControl::getBasePacks().
func (cc *ContentControl) GetBasePacks() []string {
	return cc.basePacks
}

// ResetContent resets content to base packs. (Stub for now)
// Corresponds to C++ ContentControl::resetContent().
func (cc *ContentControl) ResetContent(nonReset []string) {
	log.Printf("ContentControl.ResetContent: Stub - nonReset: %v\n", nonReset)
}

// LoadContent loads specified content packs. (Stub for now)
// Corresponds to C++ ContentControl::loadContent(const std::vector<std::string>& names)
// and ContentControl::loadContent().
func (cc *ContentControl) LoadContent(names ...string) {
	if len(names) > 0 {
		log.Printf("ContentControl.LoadContent: Stub - Loading specific packs: %v\n", names)
	} else {
		log.Println("ContentControl.LoadContent: Stub - Loading all content.")
	}
}

// SetContentPacksRaw sets content packs directly. (Stub for now)
// Corresponds to C++ ContentControl::setContentPacksRaw().
func (cc *ContentControl) SetContentPacksRaw(packs []ContentPack) {
	log.Printf("ContentControl.SetContentPacksRaw: Stub - setting %d packs\n", len(packs))
	cc.contentPacks = packs
}

// GetContentPacks returns currently active content packs.
// Corresponds to C++ ContentControl::getContentPacks().
func (cc *ContentControl) GetContentPacks() []ContentPack {
	return cc.contentPacks
}

// GetAllContentPacks returns all found content packs (including core).
// Corresponds to C++ ContentControl::getAllContentPacks().
func (cc *ContentControl) GetAllContentPacks() []ContentPack {
	return cc.allPacks
}

// Scan scans for content packs. (Stub for now)
// Corresponds to C++ ContentControl::scan().
func (cc *ContentControl) Scan() *PacksManagerStub {
	log.Println("ContentControl.Scan: Stub")
	return cc.manager
}

// SetContentSources sets content source paths. (Stub for now)
// Corresponds to C++ ContentControl::setContentSources().
func (cc *ContentControl) SetContentSources(sources []io.Path) {
	log.Printf("ContentControl.SetContentSources: Stub - sources: %v\n", sources)
}

// ResetContentSources resets content source paths. (Stub for now)
// Corresponds to C++ ContentControl::resetContentSources().
func (cc *ContentControl) ResetContentSources() {
	log.Println("ContentControl.ResetContentSources: Stub")
}

// GetContentSources returns current content source paths. (Stub for now)
// Corresponds to C++ ContentControl::getContentSources().
func (cc *ContentControl) GetContentSources() []io.Path {
	log.Println("ContentControl.GetContentSources: Stub")
	return []io.Path{}
}

// PacksManagerStub is a placeholder for the C++ PacksManager class.
type PacksManagerStub struct {
	cc *ContentControl // Reference to ContentControl
}

// NewPacksManagerStub creates a new PacksManagerStub.
func NewPacksManagerStub(cc *ContentControl) *PacksManagerStub {
	return &PacksManagerStub{cc: cc}
}

// Scan for packs (stub).
func (pm *PacksManagerStub) Scan() {
	log.Println("PacksManagerStub.Scan: Stub")
}
