package external

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"main/structures"
	"main/utils"
)

// #cgo LDFLAGS: -ldl
/*
#include <stdlib.h>
#include <dlfcn.h>

typedef int (*swift_greet_func)();

static int call_swift_greet(void *handle) {
    swift_greet_func greet = (swift_greet_func)dlsym(handle, "swift_greet");
    if (greet == NULL) {
        return -1;
    }
    return greet();
}
*/
import "C"

var appleHandle unsafe.Pointer

func loadAppleLib() bool {
	handle := C.dlopen(C.CString("apple.dylib"), C.RTLD_LAZY)
	if handle == nil {
		errorMsg := C.GoString(C.dlerror())
		fmt.Printf("Failed to load apple.dylib: %s\n", errorMsg)
		return false
	}

	appleHandle = handle

	return true
}

func freeAppleLib() bool {
	C.dlclose(appleHandle)
	appleHandle = nil
	return true
}

func Apple_IsLoaded() bool {
	return appleHandle != nil
}

// Function to call swift_greet from apple.dylib
func Apple_Greet() {
	if appleHandle == nil {
		fmt.Println("Apple library not loaded")
		return
	}
	result := C.call_swift_greet(appleHandle)
	if result == -1 {
		fmt.Println("Failed to call swift_greet")
	} else {
		fmt.Printf("Result from swift_greet: %d\n", result)
	}
}

func Apple_GetInterface() structures.LibraryInterface {
	libName := "apple.dylib"

	exeDir, err := utils.GetExecutableDir()
	if err != nil {
		return structures.LibraryInterface{
			Name:        libName,
			IsAvailable: false,
		}
	}

	// Check if apple.dylib exists
	if _, err := os.Stat(filepath.Join(exeDir, "apple.dylib")); os.IsNotExist(err) {
		fmt.Println("apple.dylib not found at path", filepath.Join(exeDir, "apple.dylib"))
		return structures.LibraryInterface{
			Name:        libName,
			IsAvailable: false,
		}
	}

	return structures.LibraryInterface{
		Name:        libName,
		IsAvailable: true,
		Load:        loadAppleLib,
		Free:        freeAppleLib,
	}
}
