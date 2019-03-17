// +build linux
// +build !appengine
<<<<<<< HEAD
// +build !android
=======
>>>>>>> 267135b0592e1f32e18030a2acbc5d996890bc86

package isatty

import "golang.org/x/sys/unix"

// IsTerminal return true if the file descriptor is terminal.
func IsTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	return err == nil
}

// IsCygwinTerminal return true if the file descriptor is a cygwin or msys2
// terminal. This is also always false on this environment.
func IsCygwinTerminal(fd uintptr) bool {
	return false
}
