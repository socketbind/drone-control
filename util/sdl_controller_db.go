package util

//#cgo windows LDFLAGS: -lSDL2
//#cgo linux freebsd darwin pkg-config: sdl2
/*
#if defined(_WIN32)
	#include <SDL2/SDL.h>
	#include <SDL2/SDL_gamecontroller.h>
	#include <stdlib.h>
#else
	#include <SDL.h>
	#include <SDL_gamecontroller.h>
#endif
*/
import "C"
import (
	"unsafe"
)

// This is a define macro that cgo seemingly cannot quite process
func GameControllerAddMappingsFromFile(mappingFile string) int {
	_file := C.CString(mappingFile)
	_mode := C.CString("rb")
	defer C.free(unsafe.Pointer(_file))
	defer C.free(unsafe.Pointer(_mode))
	cptr := (*C.SDL_RWops)(unsafe.Pointer(C.SDL_RWFromFile(_file, _mode)))

	return int(C.SDL_GameControllerAddMappingsFromRW(cptr, 1))
}
