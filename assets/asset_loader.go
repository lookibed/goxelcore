package assets

import (
	"fmt"
	"log"
	"sync"

	"goxelcore/engine"   // For Engine and ResPaths
	"goxelcore/io"      // For io.Path
	"goxelcore/graphics/core" // For Shader, Texture, TextureAnimation
	"goxelcore/content" // For ContentStub
	// "goxelcore/logic"   // For scriptenv (stubbed)
	"goxelcore/graphics/ui" // For GUIStub
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
	GUI *ui.GUIStub
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

// AssetsLoader corresponds to C++ AssetsLoader class.
type AssetsLoader struct {
	engine  *engine.Engine
	assets  *Assets
	loaders map[AssetType]ALoaderFunc
	entries chan ALoaderEntry // Using channel for queue
	enqueued sync.Map // std::set<std::pair<AssetType, std::string>>
	paths   *engine.ResPaths
	
	queueMutex sync.Mutex // Protects entries and enqueued
}

// NewAssetsLoader creates a new AssetsLoader instance.
// Corresponds to C++ AssetsLoader constructor.
func NewAssetsLoader(eng *engine.Engine, assets *Assets, paths *engine.ResPaths) *AssetsLoader {
	return &AssetsLoader{
		engine:  eng,
		assets:  assets,
		loaders: make(map[AssetType]ALoaderFunc),
		entries: make(chan ALoaderEntry, 1000), // Buffered channel for queue
		paths:   paths,
	}
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

	loader, ok := al.loaders[entry.Tag]
	if !ok {
		return &AssetError{
			Type:     entry.Tag,
			Filename: entry.Filename,
			Reason:   "no loader registered for this asset type",
		}
	}

	postFunc, err := loader(al, al.getPaths(), entry.Filename, entry.Alias, entry.Config)
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
func (al *AssetsLoader) GetPaths() *engine.ResPaths {
	return al.paths
}

// GetLoader returns the loader function for a given asset type.
// Corresponds to C++ AssetsLoader::getLoader().
func (al *AssetsLoader) GetLoader(tag AssetType) ALoaderFunc {
	return al.loaders[tag]
}

// GetEngine returns the Engine instance.
// Corresponds to C++ AssetsLoader::getEngine().
func (al *AssetsLoader) GetEngine() *engine.Engine {
	return al.engine
}

// AddDefaults enqueues default core and content assets. (Stub for now)
// Corresponds to C++ static void addDefaults(AssetsLoader& loader, const Content* content).
func AddDefaults(loader *AssetsLoader, content *content.ContentStub) {
	log.Println("AssetsLoader.AddDefaults: Stub")
	// This would typically involve reading a configuration file (like preload.json)
	// and adding assets based on it.
	// Example: loader.Add(AssetTypeTexture, "texture.png", "texture_alias", nil)
}

// LoadExternalTexture (Stub for now)
// Corresponds to C++ static bool loadExternalTexture(...).
func LoadExternalTexture(assets *Assets, name string, alternatives []io.Path) bool {
	log.Printf("AssetsLoader.LoadExternalTexture: Stub - name:%s alternatives:%v\n", name, alternatives)
	return false
}

// processPreload (Stub)
// Corresponds to C++ AssetsLoader::processPreload()
func (al *AssetsLoader) processPreload(tag AssetType, name string, mapData interface{}) { // dv::value in C++
	log.Printf("AssetsLoader.processPreload: Stub - tag:%s name:%s mapData:%v\n", tag.String(), name, mapData)
}

// processPreloadList (Stub)
// Corresponds to C++ AssetsLoader::processPreloadList()
func (al *AssetsLoader) processPreloadList(tag AssetType, listData interface{}) { // dv::value in C++
	log.Printf("AssetsLoader.processPreloadList: Stub - tag:%s listData:%v\n", tag.String(), listData)
}

// processPreloadConfig (Stub)
// Corresponds to C++ AssetsLoader::processPreloadConfig()
func (al *AssetsLoader) processPreloadConfig(file io.Path) {
	log.Printf("AssetsLoader.processPreloadConfig: Stub - file:%s\n", file.String())
}

// processPreloadConfigs (Stub)
// Corresponds to C++ AssetsLoader::processPreloadConfigs()
func (al *AssetsLoader) processPreloadConfigs(content *content.ContentStub) {
	log.Printf("AssetsLoader.processPreloadConfigs: Stub - content:%v\n", content)
}

// tryAddSound (Stub)
// Corresponds to C++ AssetsLoader::tryAddSound()
func (al *AssetsLoader) tryAddSound(name string) {
	log.Printf("AssetsLoader.tryAddSound: Stub - name:%s\n", name)
}
