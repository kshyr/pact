//go:build !linux
// +build !linux

package demon

// no-op on non-Linux systems
func setProcessName() error {
	return nil
}
