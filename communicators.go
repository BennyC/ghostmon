package ghostmon

import (
	"fmt"
	"net"
)

// Connector is how we connect with our running gh-ost instances, regardless
// of TCP usage or Unix Sockets. Will return a net.Conn which messages can sent
// across. Any errors when connecting should be returned to the caller.
type Connector interface {
	Connect() (net.Conn, error)
}

type NetConnector struct {
	addr net.Addr
}

func (n NetConnector) Connect() (net.Conn, error) {
	conn, err := net.Dial(n.addr.Network(), n.addr.String())
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}

	return conn, nil
}

func WithNetConnector(config *Config) Connector {
	return &NetConnector{addr: config}
}

type Communicator struct {
	connector Connector
	conn      net.Conn
	addr      net.Addr
}

// NewCommunicator will create a Communicator instance, any communications attempted
// will require an active net.Conn to be made through Connector
func NewCommunicator(connector Connector) *Communicator {
	return &Communicator{
		connector: connector,
		conn:      nil,
		addr:      nil,
	}
}

// Unpostpone will communicate a cutover request with gh-ost through the io.Writer
// within the Communicator instance
func (a *Communicator) Unpostpone() error {
	return a.connect(func(conn net.Conn) error {
		_, err := conn.Write([]byte("unpostpone"))
		if err != nil {
			return fmt.Errorf("failed to write command: %w", err)
		}

		err = conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close conn: %w", err)
		}

		return nil
	})
}

func (a *Communicator) connect(fn func(net.Conn) error) error {
	conn, err := a.connector.Connect()
	if err != nil {
		return err
	}

	return fn(conn)
}
