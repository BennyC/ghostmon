package communicators

import (
	"fmt"
	"net"
	"time"
)

const DeadlineIn = time.Duration(5) * time.Second

// Connector is how we Connect with our running gh-ost instances, regardless
// of TCP usage or Unix Sockets. Will return a net.Conn which messages can sent
// across. Any errors when connecting should be returned to the caller.
type Connector interface {
	Connect() (net.Conn, error)
}

type DialConnector struct {
	addr net.Addr
}

func (n DialConnector) Connect() (net.Conn, error) {
	conn, err := net.Dial(n.addr.Network(), n.addr.String())

	if err != nil {
		return nil, fmt.Errorf("unable to Connect: %w", err)
	}

	if err = conn.SetDeadline(time.Now().Add(DeadlineIn)); err != nil {
		return nil, fmt.Errorf("unable to SetDeadline: %w", err)
	}

	return conn, nil
}

func WithDialConnector(addr net.Addr) Connector {
	return &DialConnector{addr: addr}
}
