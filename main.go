package main

import (
	"log" // Keep for Fatalf initially until engine logger is fully ready
	"os"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"goxelcore/debug"  // Import our debug package
	"goxelcore/engine" // Import our new engine package
)

func main() {
	// Defer SDL library unload.
	// This is placed here as it's a global SDL operation.
	defer binsdl.Load().Unload()

	// Get engine instance
	eng := engine.GetInstance()

	// Create CoreParameters (for now, default ones)
	coreParams := engine.NewCoreParameters()
	// In future, parse command line arguments into coreParams, mirroring C++ main.cpp
	// parse_cmdline(argc, argv, coreParameters)

	// Initialize the engine
	if err := eng.Initialize(coreParams); err != nil {
		// Use standard log.Fatalf before engine's logger is fully initialized
		log.Fatalf("Engine initialization failed: %v", err)
		os.Exit(1)
	}
	// Defer engine termination, which will now also flush logs
	defer engine.Terminate() 

	// Run the engine main loop
	eng.Run()
}
