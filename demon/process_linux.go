//go:build linux
// +build linux

package demon

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// sets the process name to "demon" on Linux systems
func setProcessName() error {
	name := "demon"
	if len(name) > 15 {
		name = name[:15]
	}

	nameBytes := make([]byte, 16)
	copy(nameBytes, name)

	err := unix.Prctl(unix.PR_SET_NAME, uintptr(unsafe.Pointer(&nameBytes[0])), 0, 0, 0)
	if err != nil {
		return err
	}

	return nil
}
