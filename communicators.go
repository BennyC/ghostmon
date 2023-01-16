package ghostmon

import (
	"fmt"
	"net"
)

type Communicator struct {
	conn net.Conn
}

// NewNetCommunicator will create a new *Communicator from the desired network type
// and address location. Any errors when connecting the *Communicator to gh-ost will be
// returned to the caller
func NewNetCommunicator(addr net.Addr) (*Communicator, error) {
	conn, err := net.Dial(addr.Network(), addr.Network())
	if err != nil {
		return nil, fmt.Errorf("unable to dial addr: %w", err)
	}

	return NewCommunicator(conn), nil
}

// NewCommunicator will create a Communicator instance, based upon the io.Writer
// provided. This will normally be a net.Conn
func NewCommunicator(conn net.Conn) *Communicator {
	return &Communicator{conn}
}

// Unpostpone will communicate a cutover request with gh-ost through the io.Writer
// within the Communicator instance
func (a *Communicator) Unpostpone() error {
	err := a.send([]byte("unpostpone"))
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	return nil
}

// send will attempt to write to the connected net.Conn, any errors received
// during Write or Close will be returned to the caller
func (a *Communicator) send(b []byte) error {
	_, err := a.conn.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write to conn: %w", err)
	}

	err = a.conn.Close()
	if err != nil {
		return fmt.Errorf("unable to close conn: %w", err)
	}

	return nil
}
