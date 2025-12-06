package engine

// CoreParameters represents the core parameters for the engine.
// Corresponds to C++ CoreParameters struct in voxelcore/src/engine/CoreParameters.hpp
type CoreParameters struct {
	Headless        bool   // headless = false
	TestMode        bool   // testMode = false
	ResFolder       string // resFolder = "res"
	UserFolder      string // userFolder = "."
	ScriptFile      string // scriptFile
	ProjectFolder   string // projectFolder
	DebugServerString string // debugServerString
	TPS             int    // tps = 20
}

// NewCoreParameters creates a new CoreParameters struct with default values.
func NewCoreParameters() *CoreParameters {
	return &CoreParameters{
		Headless:        false,
		TestMode:        false,
		ResFolder:       "res",
		UserFolder:      ".",
		ScriptFile:      "",
		ProjectFolder:   "",
		DebugServerString: "",
		TPS:             20,
	}
}
