package engine

import (
	"goxelcore/data"
	"math" // For math.MaxInt and math.MinInt
)

// AudioSettings corresponds to C++ AudioSettings struct in voxelcore/src/settings.hpp
type AudioSettings struct {
	Enabled      *data.FlagSetting
	VolumeMaster *data.NumberSetting
	VolumeRegular *data.NumberSetting
	VolumeUI     *data.NumberSetting
	VolumeAmbient *data.NumberSetting
	VolumeMusic  *data.NumberSetting
	InputDevice  *data.StringSetting
}

// NewAudioSettings creates a new AudioSettings struct with default values.
func NewAudioSettings() *AudioSettings {
	return &AudioSettings{
		Enabled:      data.NewFlagSetting(true, data.SettingFormatSimple),
		VolumeMaster: data.NewNumberSetting(1.0, 0.0, 1.0, data.SettingFormatPercent),
		VolumeRegular: data.NewNumberSetting(1.0, 0.0, 1.0, data.SettingFormatPercent),
		VolumeUI:     data.NewNumberSetting(1.0, 0.0, 1.0, data.SettingFormatPercent),
		VolumeAmbient: data.NewNumberSetting(1.0, 0.0, 1.0, data.SettingFormatPercent),
		VolumeMusic:  data.NewNumberSetting(1.0, 0.0, 1.0, data.SettingFormatPercent),
		InputDevice:  data.NewStringSetting("auto", data.SettingFormatSimple),
	}
}

// DisplaySettings corresponds to C++ DisplaySettings struct in voxelcore/src/settings.hpp
type DisplaySettings struct {
	WindowMode         *data.IntegerSetting
	Width              *data.IntegerSetting
	Height             *data.IntegerSetting
	Samples            *data.IntegerSetting
	Framerate          *data.IntegerSetting
	LimitFpsIconified  *data.FlagSetting
	AdaptiveFpsInMenu  *data.FlagSetting
}

// NewDisplaySettings creates a new DisplaySettings struct with default values.
func NewDisplaySettings() *DisplaySettings {
	return &DisplaySettings{
		WindowMode:         data.NewIntegerSetting(0, 0, 2, data.SettingFormatSimple),
		Width:              data.NewIntegerSetting(1280, 0, math.MaxInt, data.SettingFormatSimple),
		Height:             data.NewIntegerSetting(720, 0, math.MaxInt, data.SettingFormatSimple),
		Samples:            data.NewIntegerSetting(0, 0, math.MaxInt, data.SettingFormatSimple),
		Framerate:          data.NewIntegerSetting(-1, -1, 120, data.SettingFormatSimple),
		LimitFpsIconified:  data.NewFlagSetting(false, data.SettingFormatSimple),
		AdaptiveFpsInMenu:  data.NewFlagSetting(false, data.SettingFormatSimple),
	}
}

// ChunksSettings corresponds to C++ ChunksSettings struct in voxelcore/src/settings.hpp
type ChunksSettings struct {
	LoadSpeed    *data.IntegerSetting
	LoadDistance *data.IntegerSetting
	Padding      *data.IntegerSetting
}

// NewChunksSettings creates a new ChunksSettings struct with default values.
func NewChunksSettings() *ChunksSettings {
	return &ChunksSettings{
		LoadSpeed:    data.NewIntegerSetting(4, 1, 32, data.SettingFormatSimple),
		LoadDistance: data.NewIntegerSetting(22, 3, 80, data.SettingFormatSimple),
		Padding:      data.NewIntegerSetting(2, 1, 8, data.SettingFormatSimple),
	}
}

// CameraSettings corresponds to C++ CameraSettings struct in voxelcore/src/settings.hpp
type CameraSettings struct {
	FovEffects  *data.FlagSetting
	Shaking     *data.FlagSetting
	Inertia     *data.FlagSetting
	Fov         *data.NumberSetting
	Sensitivity *data.NumberSetting
}

// NewCameraSettings creates a new CameraSettings struct with default values.
func NewCameraSettings() *CameraSettings {
	return &CameraSettings{
		FovEffects:  data.NewFlagSetting(true, data.SettingFormatSimple),
		Shaking:     data.NewFlagSetting(true, data.SettingFormatSimple),
		Inertia:     data.NewFlagSetting(true, data.SettingFormatSimple),
		Fov:         data.NewNumberSetting(90.0, 10, 120, data.SettingFormatSimple),
		Sensitivity: data.NewNumberSetting(2.0, 0.1, 10.0, data.SettingFormatSimple),
	}
}

// GraphicsSettings corresponds to C++ GraphicsSettings struct in voxelcore/src/settings.hpp
type GraphicsSettings struct {
	FogCurve            *data.NumberSetting
	Gamma               *data.NumberSetting
	Backlight           *data.FlagSetting
	DenseRender         *data.FlagSetting
	FrustumCulling      *data.FlagSetting
	SkyboxResolution    *data.IntegerSetting
	ChunkMaxVertices    *data.IntegerSetting
	ChunkMaxVerticesDense *data.IntegerSetting
	ChunkMaxRenderers   *data.IntegerSetting
	AdvancedRender      *data.FlagSetting
	Ssao                *data.FlagSetting
	ShadowsQuality      *data.IntegerSetting
	DenseRenderDistance *data.IntegerSetting
	SoftLighting        *data.FlagSetting
}

// NewGraphicsSettings creates a new GraphicsSettings struct with default values.
func NewGraphicsSettings() *GraphicsSettings {
	return &GraphicsSettings{
		FogCurve:            data.NewNumberSetting(1.0, 1.0, 6.0, data.SettingFormatSimple),
		Gamma:               data.NewNumberSetting(1.0, 0.4, 1.0, data.SettingFormatSimple),
		Backlight:           data.NewFlagSetting(true, data.SettingFormatSimple),
		DenseRender:         data.NewFlagSetting(true, data.SettingFormatSimple),
		FrustumCulling:      data.NewFlagSetting(true, data.SettingFormatSimple),
		SkyboxResolution:    data.NewIntegerSetting(64+32, 64, 128, data.SettingFormatSimple),
		ChunkMaxVertices:    data.NewIntegerSetting(200000, 0, 4000000, data.SettingFormatSimple),
		ChunkMaxVerticesDense: data.NewIntegerSetting(800000, 0, 8000000, data.SettingFormatSimple),
		ChunkMaxRenderers:   data.NewIntegerSetting(6, -4, 32, data.SettingFormatSimple),
		AdvancedRender:      data.NewFlagSetting(true, data.SettingFormatSimple),
		Ssao:                data.NewFlagSetting(true, data.SettingFormatSimple),
		ShadowsQuality:      data.NewIntegerSetting(0, 0, 3, data.SettingFormatSimple),
		DenseRenderDistance: data.NewIntegerSetting(56, 0, 10000, data.SettingFormatSimple),
		SoftLighting:        data.NewFlagSetting(true, data.SettingFormatSimple),
	}
}

// PathfindingSettings corresponds to C++ PathfindingSettings struct in voxelcore/src/settings.hpp
type PathfindingSettings struct {
	StepsPerAsyncAgent *data.IntegerSetting
}

// NewPathfindingSettings creates a new PathfindingSettings struct with default values.
func NewPathfindingSettings() *PathfindingSettings {
	return &PathfindingSettings{
		StepsPerAsyncAgent: data.NewIntegerSetting(128, 1, 2048, data.SettingFormatSimple),
	}
}

// DebugSettings corresponds to C++ DebugSettings struct in voxelcore/src/settings.hpp
type DebugSettings struct {
	GeneratorTestMode *data.FlagSetting
	DoWriteLights     *data.FlagSetting
	DoTraceShaders    *data.FlagSetting
	EnableExperimental *data.FlagSetting
}

// NewDebugSettings creates a new DebugSettings struct with default values.
func NewDebugSettings() *DebugSettings {
	return &DebugSettings{
		GeneratorTestMode: data.NewFlagSetting(false, data.SettingFormatSimple),
		DoWriteLights:     data.NewFlagSetting(true, data.SettingFormatSimple),
		DoTraceShaders:    data.NewFlagSetting(false, data.SettingFormatSimple),
		EnableExperimental: data.NewFlagSetting(false, data.SettingFormatSimple),
	}
}

// UiSettings corresponds to C++ UiSettings struct in voxelcore/src/settings.hpp
type UiSettings struct {
	Language        *data.StringSetting
	WorldPreviewSize *data.IntegerSetting
}

// NewUiSettings creates a new UiSettings struct with default values.
func NewUiSettings() *UiSettings {
	return &UiSettings{
		Language:        data.NewStringSetting("auto", data.SettingFormatSimple),
		WorldPreviewSize: data.NewIntegerSetting(64, 1, 512, data.SettingFormatSimple),
	}
}

// NetworkSettings corresponds to C++ NetworkSettings struct in voxelcore/src/settings.hpp
type NetworkSettings struct {
	// Empty in C++
}

// NewNetworkSettings creates a new NetworkSettings struct.
func NewNetworkSettings() *NetworkSettings {
	return &NetworkSettings{}
}

// EngineSettings corresponds to C++ EngineSettings struct in voxelcore/src/settings.hpp
type EngineSettings struct {
	Audio      *AudioSettings
	Display    *DisplaySettings
	Chunks     *ChunksSettings
	Camera     *CameraSettings
	Graphics   *GraphicsSettings
	Debug      *DebugSettings
	Ui         *UiSettings
	Network    *NetworkSettings
	Pathfinding *PathfindingSettings
}

// NewEngineSettings creates a new EngineSettings struct with default values for all nested settings.
func NewEngineSettings() *EngineSettings {
	return &EngineSettings{
		Audio:      NewAudioSettings(),
		Display:    NewDisplaySettings(),
		Chunks:     NewChunksSettings(),
		Camera:     NewCameraSettings(),
		Graphics:   NewGraphicsSettings(),
		Debug:      NewDebugSettings(),
		Ui:         NewUiSettings(),
		Network:    NewNetworkSettings(),
		Pathfinding: NewPathfindingSettings(),
	}
}

// ResetAllToDefaults resets all settings to their default values.
func (es *EngineSettings) ResetAllToDefaults() {
	es.Audio.Enabled.ResetToDefault()
	es.Audio.VolumeMaster.ResetToDefault()
	es.Audio.VolumeRegular.ResetToDefault()
	es.Audio.VolumeUI.ResetToDefault()
	es.Audio.VolumeAmbient.ResetToDefault()
	es.Audio.VolumeMusic.ResetToDefault()
	es.Audio.InputDevice.ResetToDefault()

	es.Display.WindowMode.ResetToDefault()
	es.Display.Width.ResetToDefault()
	es.Display.Height.ResetToDefault()
	es.Display.Samples.ResetToDefault()
	es.Display.Framerate.ResetToDefault()
	es.Display.LimitFpsIconified.ResetToDefault()
	es.Display.AdaptiveFpsInMenu.ResetToDefault()

	es.Chunks.LoadSpeed.ResetToDefault()
	es.Chunks.LoadDistance.ResetToDefault()
	es.Chunks.Padding.ResetToDefault()

	es.Camera.FovEffects.ResetToDefault()
	es.Camera.Shaking.ResetToDefault()
	es.Camera.Inertia.ResetToDefault()
	es.Camera.Fov.ResetToDefault()
	es.Camera.Sensitivity.ResetToDefault()

	es.Graphics.FogCurve.ResetToDefault()
	es.Graphics.Gamma.ResetToDefault()
	es.Graphics.Backlight.ResetToDefault()
	es.Graphics.DenseRender.ResetToDefault()
	es.Graphics.FrustumCulling.ResetToDefault()
	es.Graphics.SkyboxResolution.ResetToDefault()
	es.Graphics.ChunkMaxVertices.ResetToDefault()
	es.Graphics.ChunkMaxVerticesDense.ResetToDefault()
	es.Graphics.ChunkMaxRenderers.ResetToDefault()
	es.Graphics.AdvancedRender.ResetToDefault()
	es.Graphics.Ssao.ResetToDefault()
	es.Graphics.ShadowsQuality.ResetToDefault()
	es.Graphics.DenseRenderDistance.ResetToDefault()
	es.Graphics.SoftLighting.ResetToDefault()

	es.Pathfinding.StepsPerAsyncAgent.ResetToDefault()

	es.Debug.GeneratorTestMode.ResetToDefault()
	es.Debug.DoWriteLights.ResetToDefault()
	es.Debug.DoTraceShaders.ResetToDefault()
	es.Debug.EnableExperimental.ResetToDefault()

	es.Ui.Language.ResetToDefault()
	es.Ui.WorldPreviewSize.ResetToDefault()
}
