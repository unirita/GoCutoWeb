package pathutil

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"
)

const maxPathLength = 1024

var modulePath = getModulePath()

// GetRootPath creates application root path from module path.
func GetRootPath() string {
	return filepath.Dir(filepath.Dir(modulePath))
}

func getModulePath() string {
	procGetModuleFileNameW := loadProc()

	// To use pointer, buf must be array, not slice.
	var buf [maxPathLength]byte
	procGetModuleFileNameW.Call(0, uintptr(unsafe.Pointer(&buf)), (uintptr)(maxPathLength))

	// Remove zero value in second byte from Unicode string.
	path := make([]byte, maxPathLength/2)
	var j int
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 {
			path[j] = buf[i]
			j++
		}
	}
	return fmt.Sprintf("%s", path)
}

func loadProc() *syscall.Proc {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		panic(err)
	}
	proc, err := dll.FindProc("GetModuleFileNameW")
	if err != nil {
		panic(err)
	}
	return proc
}
