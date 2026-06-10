package port

import (
	"context"
	"net"
	"strconv"
)

// IsFree reports whether a TCP port is available for a service to bind.
//
// The probe binds the wildcard address on both stacks with SO_REUSEADDR
// disabled. Both details matter on BSD/macOS:
//   - A 127.0.0.1 probe succeeds even while another process holds the
//     wildcard (e.g. a Node server on "::"), making busy ports look free.
//   - Go enables SO_REUSEADDR on listeners by default, which lets wildcard
//     and specific-address binds coexist, so even a wildcard probe misses
//     specific-address listeners unless SO_REUSEADDR is turned off.
//
// There is still an inherent TOCTOU race between this check and the moment
// the child process binds the port. This is mitigated by (1) the file-level
// lock in state.FileStore serializing port allocation across concurrent
// portree invocations, and (2) early-exit detection when a service dies
// shortly after starting.
func IsFree(port int) bool {
	lc := net.ListenConfig{Control: disableReuseAddr}
	addr := ":" + strconv.Itoa(port)
	for _, network := range []string{"tcp4", "tcp6"} {
		ln, err := lc.Listen(context.Background(), network, addr)
		if err != nil {
			return false
		}
		_ = ln.Close()
	}
	return true
}
