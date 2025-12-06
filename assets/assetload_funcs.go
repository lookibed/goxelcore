package assets

import (
	"fmt"
	"log"
	// "goxelcore/graphics/core" // For Shader, Texture
)

// ShaderLoaderFunc is a stub loader for shaders.
// Corresponds to C++ assetload::shader.
func ShaderLoaderFunc(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error) {
	log.Printf("ShaderLoaderFunc: Stub - Loading shader %s (alias: %s)\n", filename, alias)
	// Actual implementation would:
	// - Read shader source (vertex and fragment)
	// - Use GLSLExtension to process it
	// - Compile and link with core.CreateShader
	// - Store with assets.StoreShader
	return nil, nil
}

// TextureLoaderFunc is a stub loader for textures.
// Corresponds to C++ assetload::texture.
func TextureLoaderFunc(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error) {
	log.Printf("TextureLoaderFunc: Stub - Loading texture %s (alias: %s)\n", filename, alias)
	// Actual implementation would:
	// - Read image data using imageio (once ported)
	// - Create texture with core.TextureFromImageData
	// - Store with assets.StoreTexture
	return nil, nil
}

// FontLoaderFunc is a stub loader for fonts.
// Corresponds to C++ assetload::font.
func FontLoaderFunc(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error) {
	log.Printf("FontLoaderFunc: Stub - Loading font %s (alias: %s)\n", filename, alias)
	return nil, nil
}

// AtlasLoaderFunc is a stub loader for atlases.
// Corresponds to C++ assetload::atlas.
func AtlasLoaderFunc(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error) {
	log.Printf("AtlasLoaderFunc: Stub - Loading atlas %s (alias: %s)\n", filename, alias)
	return nil, nil
}

// LayoutLoaderFunc is a stub loader for UI layouts.
// Corresponds to C++ assetload::layout.
func LayoutLoaderFunc(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error) {
	log.Printf("LayoutLoaderFunc: Stub - Loading layout %s (alias: %s)\n", filename, alias)
	return nil, nil
}

// SoundLoaderFunc is a stub loader for sounds.
// Corresponds to C++ assetload::sound.
func SoundLoaderFunc(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error) {
	log.Printf("SoundLoaderFunc: Stub - Loading sound %s (alias: %s)\n", filename, alias)
	return nil, nil
}

// ModelLoaderFunc is a stub loader for 3D models.
// Corresponds to C++ assetload::model.
func ModelLoaderFunc(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error) {
	log.Printf("ModelLoaderFunc: Stub - Loading model %s (alias: %s)\n", filename, alias)
	return nil, nil
}

// PostEffectLoaderFunc is a stub loader for post-effects.
// Corresponds to C++ assetload::posteffect.
func PostEffectLoaderFunc(
	loader *AssetsLoader,
	paths *engine.ResPaths,
	filename string,
	alias string,
	config AssetCfg,
) (PostFunc, error) {
	log.Printf("PostEffectLoaderFunc: Stub - Loading post-effect %s (alias: %s)\n", filename, alias)
	return nil, nil
}
