package main

import (
	"log"
	"os"

	"github.com/Zyko0/go-sdl3/bin/binsdl" // Still deferring this
	"goxelcore/engine" // Import our new engine package
)

func main() {
	defer binsdl.Load().Unload() // Still deferring this here, as it's from the original example

	// Get engine instance
	eng := engine.GetInstance()

	// Create CoreParameters (for now, default ones)
	coreParams := engine.NewCoreParameters()
	// In future, parse command line arguments into coreParams, mirroring C++ main.cpp
	// parse_cmdline(argc, argv, coreParameters)

	// Initialize the engine
	if err := eng.Initialize(coreParams); err != nil {
		log.Fatalf("Engine initialization failed: %v", err)
		os.Exit(1)
	}
	defer engine.Terminate() // Defer engine termination

	// Run the engine main loop
	eng.Run()
}