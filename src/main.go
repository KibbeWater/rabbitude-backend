package main

import (
	"fmt"
	"main/external"
	"main/providers"
	"main/server"
	"main/structures"
)

func main() {
	// Load Libraries
	// Array of LibraryInterface
	var libraries []structures.LibraryInterface = []structures.LibraryInterface{
		external.Apple_GetInterface(),
	}

	// Load all the Libraries
	for _, lib := range libraries {
		if lib.IsAvailable {
			fmt.Println("Loading", lib.Name)
			lib.Load()
		}
	}

	// Test the Apple lib
	if external.Apple_IsLoaded() {
		fmt.Println("Apple lib loaded")
		// external.Apple_Greet("Swift")
		external.Apple_Greet()
	} else {
		fmt.Println("Apple lib not loaded")
	}
	return

	// Start the server
	providers.InstallServices()
	server.StartServer()
}
