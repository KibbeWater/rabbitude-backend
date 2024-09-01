package external

import (
	"fmt"
	"main/structures"
)

var libraries []structures.LibraryInterface = []structures.LibraryInterface{
	Apple_GetInterface(),
}

func LoadLibraries() {
	// Load all the Libraries
	for _, lib := range libraries {
		if lib.IsAvailable {
			fmt.Println("[Library] Loading", lib.Name)
			lib.Load()
		}
	}
}

func FreeLibraries() {
	// Free all the Libraries
	for _, lib := range libraries {
		if lib.IsAvailable {
			fmt.Println("[Library] Unloading", lib.Name)
			lib.Free()
		}
	}
}
