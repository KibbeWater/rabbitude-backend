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

typedef void (*swift_requestSpeechPermission_func)(int*);
typedef char* (*swift_speechRecognition_func)(void*, int);

static void call_swift_requestSpeechPermissions(void *handle, int *status) {
    swift_requestSpeechPermission_func requestSpeechPermission = (swift_requestSpeechPermission_func)dlsym(handle, "swift_requestSpeechPermissions");
    if (requestSpeechPermission == NULL) {
        return;
    }
    requestSpeechPermission(status);
}

static char* call_swift_speechRecognition(void *handle, void *data, int length) {
    swift_speechRecognition_func speechRecognition = (swift_speechRecognition_func)dlsym(handle, "swift_speechRecognition");
    if (speechRecognition == NULL) {
        return NULL;
    }
    return speechRecognition(data, length);
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
func Apple_RequestSpeechPermissions(status *int) {
	if appleHandle == nil {
		fmt.Println("Apple library not loaded")
		return
	}

	C.call_swift_requestSpeechPermissions(appleHandle, (*C.int)(unsafe.Pointer(status)))
}

func Apple_SpeechRecognition(data []byte) (string, error) {
	if appleHandle == nil {
		fmt.Println("Apple library not loaded")
		return "", fmt.Errorf("Apple library not loaded")
	}

	ret := C.call_swift_speechRecognition(appleHandle, unsafe.Pointer(&data[0]), C.int(len(data)))
	return C.GoString(ret), nil
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
