package assets

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug" // For stack trace on asset not found
	"sync"

	"goxelcore/graphics/core" // For TextureAnimation
	"goxelcore/util"          // For stringutil.quote (stub for now)
)

// AssetType corresponds to C++ AssetType enum in voxelcore/src/assets/Assets.hpp
type AssetType int

const (
	AssetTypeTexture AssetType = iota
	AssetTypeShader
	AssetTypeFont
	AssetTypeAtlas
	AssetTypeLayout
	AssetTypeSound
	AssetTypeModel
	AssetTypePostEffect
)

// String returns the string representation of the AssetType.
func (at AssetType) String() string {
	switch at {
	case AssetTypeTexture:
		return "TEXTURE"
	case AssetTypeShader:
		return "SHADER"
	case AssetTypeFont:
		return "FONT"
	case AssetTypeAtlas:
		return "ATLAS"
	case AssetTypeLayout:
		return "LAYOUT"
	case AssetTypeSound:
		return "SOUND"
	case AssetTypeModel:
		return "MODEL"
	case AssetTypePostEffect:
		return "POST_EFFECT"
	default:
		return "UNKNOWN"
	}
}

// AssetError corresponds to C++ assetload::error.
type AssetError struct {
	Type     AssetType
	Filename string
	Reason   string
}

func (e *AssetError) Error() string {
	return fmt.Sprintf("assetload error [%s] %s: %s", e.Type.String(), e.Filename, e.Reason)
}

// PostFunc corresponds to C++ assetload::postfunc.
type PostFunc func(assets *Assets)

// SetupFunc corresponds to C++ assetload::setupfunc.
type SetupFunc func(assets *Assets)

// Assets is a collection of various game assets.
// Corresponds to C++ Assets class in voxelcore/src/assets/Assets.hpp
type Assets struct {
	animations []core.TextureAnimation
	
	// C++ uses std::unordered_map<std::type_index, std::unordered_map<std::string, std::shared_ptr<void>>> assets;
	// In Go, we'll use a map of asset type to a map of name to interface{}
	assets map[AssetType]map[string]interface{}

	setupFuncs []SetupFunc
	mu         sync.RWMutex // For concurrent access to assets map
}

// NewAssets creates a new Assets instance.
// Corresponds to C++ Assets() constructor.
func NewAssets() *Assets {
	return &Assets{
		assets: make(map[AssetType]map[string]interface{}),
	}
}

// Delete is a placeholder for asset cleanup.
// Corresponds to C++ ~Assets() destructor.
func (a *Assets) Delete() {
	a.mu.Lock()
	defer a.mu.Unlock()
	log.Println("Assets.Delete: Stub for cleaning up assets (e.g., OpenGL textures, meshes)")
	// In a full implementation, iterate through stored assets and call their Delete/Destroy methods.
}

// GetAnimations returns the list of texture animations.
// Corresponds to C++ Assets::getAnimations().
func (a *Assets) GetAnimations() []core.TextureAnimation {
	return a.animations
}

// StoreAnimation stores a texture animation.
// Corresponds to C++ Assets::store(const TextureAnimation& animation).
func (a *Assets) StoreAnimation(animation core.TextureAnimation) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.animations = append(a.animations, animation)
}

// storeInternal stores an asset of a given type.
// This is a helper for the generic C++ template store functions.
func (a *Assets) storeInternal(assetType AssetType, asset interface{}, name string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.assets[assetType]; !ok {
		a.assets[assetType] = make(map[string]interface{})
	}
	a.assets[assetType][name] = asset
}

// StoreShader stores a shader asset.
func (a *Assets) StoreShader(asset *core.Shader, name string) { a.storeInternal(AssetTypeShader, asset, name) }

// StoreTexture stores a texture asset.
func (a *Assets) StoreTexture(asset *core.Texture, name string) { a.storeInternal(AssetTypeTexture, asset, name) }

// Get retrieves an asset by its type and name.
// Corresponds to C++ Assets::get<T>(const std::string& name) const.
func (a *Assets) Get(assetType AssetType, name string) (interface{}, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if typeMap, ok := a.assets[assetType]; ok {
		if asset, found := typeMap[name]; found {
			return asset, nil
		}
	}
	return nil, &AssetError{Type: assetType, Filename: name, Reason: "asset not found"}
}

// Require retrieves an asset by its type and name, panicking if not found.
// Corresponds to C++ Assets::require<T>(const std::string& name) const.
func (a *Assets) Require(assetType AssetType, name string) interface{} {
	asset, err := a.Get(assetType, name)
	if err != nil {
		log.Printf("Asset require failed: %v\n%s", err, debug.Stack())
		panic(err) // Panic as per C++ behavior
	}
	return asset
}

// GetShader retrieves a Shader asset.
func (a *Assets) GetShader(name string) (*core.Shader, error) {
	asset, err := a.Get(AssetTypeShader, name)
	if err != nil {
		return nil, err
	}
	if shader, ok := asset.(*core.Shader); ok {
		return shader, nil
	}
	return nil, &AssetError{Type: AssetTypeShader, Filename: name, Reason: "asset found but wrong type"}
}

// GetTexture retrieves a Texture asset.
func (a *Assets) GetTexture(name string) (*core.Texture, error) {
	asset, err := a.Get(AssetTypeTexture, name)
	if err != nil {
		return nil, err
	}
	if texture, ok := asset.(*core.Texture); ok {
		return texture, nil
	}
	return nil, &AssetError{Type: AssetTypeTexture, Filename: name, Reason: "asset found but wrong type"}
}

// Setup runs all registered setup functions.
// Corresponds to C++ Assets::setup().
func (a *Assets) Setup() {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, setupFunc := range a.setupFuncs {
		setupFunc(a)
	}
}

// AddSetupFunc adds a setup function.
// Corresponds to C++ Assets::addSetupFunc().
func (a *Assets) AddSetupFunc(setupFunc SetupFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.setupFuncs = append(a.setupFuncs, setupFunc)
}

// AssetsSetup (template-like function)
// Corresponds to C++ template <class T> void assetload::assets_setup(const Assets* assets);
// In Go, you'd typically have specific setup functions for each asset type.
// For now, it's just a general placeholder if specific setup logic is needed.
func AssetsSetup(assets *Assets, assetType AssetType, setupLogic func(asset interface{})) {
	assets.mu.RLock()
	defer assets.mu.RUnlock()
	if typeMap, ok := assets.assets[assetType]; ok {
		for _, asset := range typeMap {
			setupLogic(asset)
		}
	}
}

// Asset types will need specific getters and setters (StoreX, GetX, RequireX)
// to maintain type safety and avoid frequent type assertions.
// Example: GetFont, StoreFont, etc.
