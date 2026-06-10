//go:build !unix

package port

import "syscall"

// disableReuseAddr is a no-op on non-Unix platforms, where Go does not set
// SO_REUSEADDR on listeners and plain bind probes already conflict with
// existing listeners.
func disableReuseAddr(network, address string, c syscall.RawConn) error {
	return nil
}
