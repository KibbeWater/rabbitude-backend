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

typedef char* (*swift_greet_func)(const char*);

static char* call_swift_greet(void *handle, const char *name) {
    swift_greet_func greet = (swift_greet_func)dlsym(handle, "swift_greet");
    if (greet == NULL) {
        return 0;
    }
    return greet(name);
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
func Apple_Greet(name string) (string, error) {
	if appleHandle == nil {
		fmt.Println("Apple library not loaded")
		return "", fmt.Errorf("Apple library not loaded")
	}
	result := C.call_swift_greet(appleHandle, C.CString(name))
	fmt.Printf("Result from swift_greet: %d\n", C.GoString(result))
	return C.GoString(result), nil
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
