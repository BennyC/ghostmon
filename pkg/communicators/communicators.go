package communicators

import (
	"fmt"
	"golang.org/x/exp/slog"
	"net"
)

type Communicator struct {
	connector Connector
	logger    *slog.Logger
}

// New will create a Communicator instance, any communications attempted
// will require an active net.Conn to be made through Connector
func New(connector Connector, logger *slog.Logger) *Communicator {
	return &Communicator{
		connector: connector,
		logger:    logger,
	}
}

// Connect will allow the Communicator to start communicating with the external
// gh-ost process. When a successful connection is achieved, the callback provided
// will receive a net.Conn to communicate with
func (a *Communicator) Connect(fn func(net.Conn) error) error {
	conn, err := a.connector.Connect()
	if err != nil {
		return err
	}

	return fn(conn)
}

// Unpostpone will communicate a cutover request with gh-ost through the io.Writer
// within the Communicator instance
func (a *Communicator) Unpostpone() error {
	return a.Connect(func(conn net.Conn) error {
		return send(conn, []byte("unpostpone"), []byte{})
	})
}

// Panic will communicate a panic request with gh-ost through the io.Writer
// within the Communicator instance
func (a *Communicator) Panic() error {
	return a.Connect(func(conn net.Conn) error {
		return send(conn, []byte("panic"), []byte{})
	})
}

type Status struct {
	Body []byte
}

// Status will return the status of the connected gh-ost process
// Any errors during communication will be returned to the caller
func (a *Communicator) Status() (*Status, error) {
	body := make([]byte, 1024)
	err := a.Connect(func(conn net.Conn) error {
		return sendReceive(conn, []byte("status"), body)
	})

	if err != nil {
		return nil, err
	}

	return &Status{
		Body: body,
	}, nil
}

// send will write the msg to a net.Conn and handle any errors associated with
// writing, reading or closing
func send(conn net.Conn, w []byte, r []byte) error {
	_, err := conn.Write(w)
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	err = conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close conn: %w", err)
	}

	return nil
}

// sendReceive will write the msg to a net.Conn and receive the response. Any errors associated with
// writing, reading or closing will be handled
func sendReceive(conn net.Conn, w []byte, r []byte) error {
	_, err := conn.Write(w)
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	_, err = conn.Read(r)
	if err != nil {
		return fmt.Errorf("failed to read after write: %w", err)
	}

	err = conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close conn: %w", err)
	}

	return nil
}
