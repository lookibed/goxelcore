package assets

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"goxelcore/coders"
	"goxelcore/content"
	"goxelcore/data"
	"goxelcore/io"
	// "goxelcore/interfaces" // For Task (stubbed)
)

// TaskStub is a placeholder for interfaces/Task.hpp Task.
type TaskStub struct{}

func (t *TaskStub) WaitForEnd() {} // Placeholder method

// AssetCfg is a base interface for asset-specific configuration.
// Corresponds to C++ AssetCfg struct.
type AssetCfg interface {
	// Marker interface for now
}

// LayoutCfg corresponds to C++ LayoutCfg struct.
type LayoutCfg struct {
	GUI interface{} // Placeholder for ui.GUIStub - avoids import cycle
	Env interface{} // scriptenv is complex, use interface{} for now
}

// SoundCfg corresponds to C++ SoundCfg struct.
type SoundCfg struct {
	KeepPCM bool
}

// AtlasType corresponds to C++ AtlasType enum.
type AtlasType int

const (
	AtlasTypeAtlas AtlasType = iota
	AtlasTypeSeparate
)

// AtlasCfg corresponds to C++ AtlasCfg struct.
type AtlasCfg struct {
	Type AtlasType
}

// PostEffectCfg corresponds to C++ PostEffectCfg struct.
type PostEffectCfg struct {
	Advanced bool
}

// ALoaderFunc corresponds to C++ aloader_func.
type ALoaderFunc func(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error)

// ALoaderEntry corresponds to C++ aloader_entry struct.
type ALoaderEntry struct {
	Tag      AssetType
	Filename string
	Alias    string
	Config   AssetCfg
}

// EngineInterface abstracts engine.Engine
type EngineInterface interface {
	GetLogger() interface{} // Placeholder - would need proper logger interface
	// Add other methods from Engine that are needed
}

// ResPathsInterface abstracts engine.ResPaths
type ResPathsInterface interface {
	Find(path string) io.Path
	// Add other methods from ResPaths as needed
}

// AssetsLoader corresponds to C++ AssetsLoader class.
type AssetsLoader struct {
	engine  EngineInterface
	assets  *Assets
	loaders map[AssetType]ALoaderFunc
	entries chan ALoaderEntry // Using channel for queue
	enqueued sync.Map // std::set<std::pair<AssetType, std::string>>
	paths   ResPathsInterface

	queueMutex sync.Mutex // Protects entries and enqueued
}

// NewAssetsLoader creates a new AssetsLoader instance.
// Corresponds to C++ AssetsLoader constructor.
func NewAssetsLoader(eng EngineInterface, assets *Assets, paths ResPathsInterface) *AssetsLoader {
	al := &AssetsLoader{
		engine:  eng,
		assets:  assets,
		loaders: make(map[AssetType]ALoaderFunc),
		entries: make(chan ALoaderEntry, 1000), // Buffered channel for queue
		paths:   paths,
	}
	// Register default loaders
	al.AddLoader(AssetTypeShader, ShaderLoaderFunc)
	al.AddLoader(AssetTypeTexture, TextureLoaderFunc)
	al.AddLoader(AssetTypeFont, FontLoaderFunc)
	al.AddLoader(AssetTypeAtlas, AtlasLoaderFunc)
	al.AddLoader(AssetTypeLayout, LayoutLoaderFunc)
	al.AddLoader(AssetTypeSound, SoundLoaderFunc)
	al.AddLoader(AssetTypeModel, ModelLoaderFunc)
	al.AddLoader(AssetTypePostEffect, PostEffectLoaderFunc)
	return al
}

// AddLoader registers an asset loader function for a given asset type.
// Corresponds to C++ AssetsLoader::addLoader().
func (al *AssetsLoader) AddLoader(tag AssetType, f ALoaderFunc) {
	al.loaders[tag] = f
}

// Add enqueues an asset for loading.
// Corresponds to C++ AssetsLoader::add().
func (al *AssetsLoader) Add(tag AssetType, filename, alias string, config AssetCfg) {
	al.queueMutex.Lock()
	defer al.queueMutex.Unlock()

	key := fmt.Sprintf("%d:%s", tag, alias)
	if _, loaded := al.enqueued.Load(key); loaded {
		return // Already enqueued
	}
	al.enqueued.Store(key, true)

	al.entries <- ALoaderEntry{Tag: tag, Filename: filename, Alias: alias, Config: config}
}

// HasNext checks if there are more assets in the queue to load.
// Corresponds to C++ AssetsLoader::hasNext().
func (al *AssetsLoader) HasNext() bool {
	al.queueMutex.Lock()
	defer al.queueMutex.Unlock()
	return len(al.entries) > 0 // Check channel length
}

// LoadNext loads the next asset in the queue.
// Corresponds to C++ AssetsLoader::loadNext().
func (al *AssetsLoader) LoadNext() error {
	al.queueMutex.Lock()
	var entry ALoaderEntry
	select {
	case entry = <-al.entries:
		// Entry dequeued
	default:
		al.queueMutex.Unlock()
		return fmt.Errorf("no more assets to load")
	}
	al.queueMutex.Unlock()

	loaderFunc, ok := al.loaders[entry.Tag]
	if !ok {
		return &AssetError{
			Type:     entry.Tag,
			Filename: entry.Filename,
			Reason:   "no loader registered for this asset type",
		}
	}

	postFunc, err := loaderFunc(al, al.GetPaths(), entry.Filename, entry.Alias, entry.Config)
	if err != nil {
		return err
	}
	// Execute post-loading function if any
	if postFunc != nil {
		postFunc(al.assets)
	}
	return nil
}

// StartTask starts an asynchronous loading task. (Stub for now)
// Corresponds to C++ AssetsLoader::startTask().
func (al *AssetsLoader) StartTask(onDone func()) *TaskStub {
	log.Println("AssetsLoader.StartTask: Stub - Asynchronous loading not fully implemented.")
	// A real implementation would start a goroutine to call LoadNext repeatedly
	// and execute onDone when the queue is empty.
	go func() {
		for al.HasNext() {
			err := al.LoadNext()
			if err != nil {
				log.Printf("Error loading asset in async task: %v\n", err)
			}
		}
		if onDone != nil {
			onDone()
		}
	}()
	return &TaskStub{} // Return a stub task
}

// GetPaths returns the ResPaths instance.
// Corresponds to C++ AssetsLoader::getPaths().
func (al *AssetsLoader) GetPaths() ResPathsInterface {
	return al.paths
}

// GetLoader returns the loader function for a given asset type.
// Corresponds to C++ AssetsLoader::getLoader().
func (al *AssetsLoader) GetLoader(tag AssetType) ALoaderFunc {
	return al.loaders[tag]
}

// GetEngine returns the Engine instance.
// Corresponds to C++ AssetsLoader::getEngine().
func (al *AssetsLoader) GetEngine() EngineInterface {
	return al.engine
}

// AssetsDefFolder returns the default folder for an asset type.
// Corresponds to C++ assets_def_folder() static helper.
func AssetsDefFolder(tag AssetType) string {
	switch tag {
	case AssetTypeFont:
		return FontsFolder
	case AssetTypeShader:
		return ShadersFolder
	case AssetTypeTexture:
		return TexturesFolder
	case AssetTypeAtlas:
		return TexturesFolder
	case AssetTypeLayout:
		return LayoutsFolder
	case AssetTypeSound:
		return SoundsFolder
	case AssetTypeModel:
		return ModelsFolder
	case AssetTypePostEffect:
		return PostEffectsFolder
	}
	return "<error>"
}

// processPreload processes a single asset entry from preload.json.
// Corresponds to C++ AssetsLoader::processPreload().
func (al *AssetsLoader) processPreload(tag AssetType, name string, mapData data.Value) error {
	defFolder := AssetsDefFolder(tag)
	assetPath := defFolder + "/" + name

	var config AssetCfg = nil
	
	if mapData.IsObject() {
		// Override path if specified in JSON
		if pathVal, ok := mapData.GetAtKey("path").GetAsString(); ok {
			assetPath = pathVal
		}

		// Handle asset-specific configuration
		switch tag {
		case AssetTypeSound:
			if keepPCMVal, ok := mapData.GetAtKey("keep-pcm").GetAsBoolean(); ok {
				config = &SoundCfg{KeepPCM: keepPCMVal}
			} else {
				config = &SoundCfg{KeepPCM: false} // Default
			}
		case AssetTypeAtlas:
			if typeNameVal, ok := mapData.GetAtKey("type").GetAsString(); ok {
				atlasType := AtlasTypeAtlas
				if typeNameVal == "separate" {
					atlasType = AtlasTypeSeparate
				}
				config = &AtlasCfg{Type: atlasType}
			} else {
				config = &AtlasCfg{Type: AtlasTypeAtlas} // Default
			}
		case AssetTypePostEffect:
			if advancedVal, ok := mapData.GetAtKey("advanced").GetAsBoolean(); ok {
				config = &PostEffectCfg{Advanced: advancedVal}
			} else {
				config = &PostEffectCfg{Advanced: false} // Default
			}
		default:
			// No specific config for other types yet
		}
	}
	
	al.Add(tag, assetPath, name, config)
	return nil
}

// processPreloadList processes a list of asset entries from preload.json.
// Corresponds to C++ AssetsLoader::processPreloadList().
func (al *AssetsLoader) processPreloadList(tag AssetType, list data.Value) error {
	if list.IsNone() { // C++ handles nullptr, Go handles None type
		return nil
	}
	if !list.IsList() {
		return fmt.Errorf("expected list for asset type %s, got %s", tag.String(), list.Type.TypeName())
	}

	for i := 0; i < int(list.Size()); i++ {
		entry := list.GetAtIndex(i)
		switch entry.Type {
		case data.ValueTypeString:
			if err := al.processPreload(tag, entry.AsString(), data.Value{}); err != nil {
				return err
			}
		case data.ValueTypeObject:
			if nameVal, ok := entry.GetAtKey("name").GetAsString(); ok {
				if err := al.processPreload(tag, nameVal, entry); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("object entry for asset type %s missing 'name' field", tag.String())
			}
		default:
			return fmt.Errorf("invalid entry type for asset type %s: expected string or object, got %s", tag.String(), entry.Type.TypeName())
		}
	}
	return nil
}

// processPreloadConfig reads and processes a preload.json file.
// Corresponds to C++ AssetsLoader::processPreloadConfig().
func (al *AssetsLoader) processPreloadConfig(file io.Path) error {
	root, err := coders.ReadJSONFile(file)
	if err != nil {
		return fmt.Errorf("failed to read preload config file %s: %w", file.String(), err)
	}

	if err := al.processPreloadList(AssetTypeAtlas, root.GetAtKey("atlases")); err != nil {
		return err
	}
	if err := al.processPreloadList(AssetTypeFont, root.GetAtKey("fonts")); err != nil {
		return err
	}
	if err := al.processPreloadList(AssetTypeShader, root.GetAtKey("shaders")); err != nil {
		return err
	}
	if err := al.processPreloadList(AssetTypeTexture, root.GetAtKey("textures")); err != nil {
		return err
	}
	if err := al.processPreloadList(AssetTypeSound, root.GetAtKey("sounds")); err != nil {
		return err
	}
	if err := al.processPreloadList(AssetTypeModel, root.GetAtKey("models")); err != nil {
		return err
	}
	if err := al.processPreloadList(AssetTypePostEffect, root.GetAtKey("post-effects")); err != nil {
		return err
	}
	// layouts are loaded automatically - this seems to be handled differently in C++ (add_layouts helper)
	return nil
}

// processPreloadConfigs reads and processes preload.json files from resource roots and content packs.
// Corresponds to C++ AssetsLoader::processPreloadConfigs().
func (al *AssetsLoader) processPreloadConfigs(content *content.ContentStub) error {
	// First, check global preload.json in res:
	preloadFile := io.NewPath("res:preload.json")
	if io.Exists(preloadFile) {
		if err := al.processPreloadConfig(preloadFile); err != nil {
			return err
		}
	}

	if content == nil {
		return nil
	}
	
	// Process preload.json from each content pack
	// C++ uses content->getPacks()
	log.Println("AssetsLoader.processPreloadConfigs: Stub - Processing content packs for preload.json")
	// For now, assume a simple structure or stub this part of content
	/*
	for _, packEntry := range content.GetPacks() { // Assuming GetPacks() returns map[string]*ContentPackRuntimeStub
		if packEntry.GetId() == "core" { // Assuming GetId() exists
			continue
		}
		packInfo := packEntry.GetInfo() // Assuming GetInfo() returns ContentPack
		packPreloadFile := packInfo.Folder.Join("preload.json")
		if io.Exists(packPreloadFile) {
			if err := al.processPreloadConfig(packPreloadFile); err != nil {
				return err
			}
		}
	}
	*/
	return nil
}

// LoggerInterface provides logging functionality
type LoggerInterface interface {
	Info(format string, v ...interface{})
	// Add other logging methods as needed
}

// AddDefaults enqueues default core and content assets.
// Corresponds to C++ static void addDefaults(AssetsLoader& loader, const Content* content).
func AddDefaults(loader *AssetsLoader, content *content.ContentStub) error {
	// Need to access logger through the engine, but the engine interface already handles this
	// For now, we'll just print to console
	log.Println("AssetsLoader.AddDefaults: Adding default assets...")
	if err := loader.processPreloadConfigs(content); err != nil {
		return fmt.Errorf("failed to process preload configs: %w", err)
	}

	// The rest of the addDefaults implementation relies on various Content methods
	// (getBlockMaterials, getPacks, getSkeletons, blocks.getDefs, items.getDefs).
	// These will remain stubs for now until Content is fully ported.
	log.Println("AssetsLoader.AddDefaults: Stub - Content-specific asset enqueuing (remaining parts).")
	
	// Example of Content related loops (stubbed)
	/*
	if content != nil {
		// Example: Loop through block materials
		// for _, entry := range content.GetBlockMaterials() {
		// 	material := entry.second
		// 	loader.tryAddSound(material.stepsSound)
		// 	loader.tryAddSound(material.placeSound)
		// 	loader.tryAddSound(material.breakSound)
		// 	loader.tryAddSound(material.hitSound)
		// }

		// Example: Loop through content packs for layouts
		// for _, entry := range content.GetPacks() {
		// 	pack := entry.second
		// 	info := pack.GetInfo()
		// 	folder := info.Folder.Join("layouts")
		// 	// add_layouts(pack.GetEnvironment(), info.ID, folder, loader)
		// }
		
		// Example: Loop through skeletons for models
		// for _, entry := range content.GetSkeletons() {
		// 	skeleton := entry.second
		// 	for _, bone := range skeleton.GetBones() {
		// 		model := bone.model.name
		// 		// ... resolve model path and add
		// 	}
		// }

		// Example: Loop through blocks and items for models
		// for _, def := range content.Blocks.GetDefs() {
		// 	// ... add variants/defaults
		// }
		// for _, def := range content.Items.GetDefs() {
		// 	// ... add models
		// }
	}
	*/
	
	return nil
}

// LoadExternalTexture (Stub for now)
// Corresponds to C++ static bool loadExternalTexture(...).
func LoadExternalTexture(assets *Assets, name string, alternatives []io.Path) bool {
	log.Printf("AssetsLoader.LoadExternalTexture: Stub - name:%s alternatives:%v\n", name, alternatives)
	return false
}

// tryAddSound (Stub)
// Corresponds to C++ AssetsLoader::tryAddSound()
func (al *AssetsLoader) tryAddSound(name string) {
	log.Printf("AssetsLoader.tryAddSound: Stub - name:%s\n", name)
}