//go:build unix

package port

import "syscall"

// disableReuseAddr clears SO_REUSEADDR (set by Go's listener defaults) so the
// probe bind conflicts with every existing listener on the port, regardless
// of which address that listener bound.
func disableReuseAddr(network, address string, c syscall.RawConn) error {
	var sockErr error
	if err := c.Control(func(fd uintptr) {
		sockErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 0)
	}); err != nil {
		return err
	}
	return sockErr
}
