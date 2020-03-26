package dialer

import (
	"context"
	"net"
)

type Dialer interface {
	Dial(string) (net.Conn, error)
}

type ContextDialer interface {
	Dialer
	DialContext(context.Context, string) (net.Conn, error)
}
