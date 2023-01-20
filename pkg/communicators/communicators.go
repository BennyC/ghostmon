package communicators

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/slog"
	"net"
	"regexp"
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
func (c *Communicator) Connect(fn func(net.Conn) error) error {
	conn, err := c.connector.Connect()
	if err != nil {
		return err
	}

	return fn(conn)
}

// Unpostpone will communicate a cutover request with gh-ost through the io.Writer
// within the Communicator instance
func (c *Communicator) Unpostpone() error {
	return c.Connect(func(conn net.Conn) error {
		return c.Send(conn, []byte("unpostpone"))
	})
}

// Panic will communicate a panic request with gh-ost through the io.Writer
// within the Communicator instance
func (c *Communicator) Panic() error {
	return c.Connect(func(conn net.Conn) error {
		return c.Send(conn, []byte("panic"))
	})
}

type Status struct {
	Body []byte
}

// Table lets you know which Table is being migrated within the current Status
func (s Status) Table() (string, error) {
	r, err := regexp.Compile("Migrating ([\\w`.]+)")
	if err != nil {
		return "", fmt.Errorf("Status.Table: unable to compile regex: %w", err)
	}

	matches := r.FindStringSubmatch(string(s.Body))
	if len(matches) < 2 {
		return "", fmt.Errorf("Status.Table: unable to find migration table in Body")
	}

	return matches[1], nil
}

// Status will return the status of the connected gh-ost process
// Any errors during communication will be returned to the caller
func (c *Communicator) Status() (*Status, error) {
	body := make([]byte, 1024)
	err := c.Connect(func(conn net.Conn) error {
		return c.SendAndReceive(conn, []byte("status"), body)
	})

	if err != nil {
		return nil, err
	}

	return &Status{
		Body: bytes.Trim(body, "\x00"),
	}, nil
}

// Send will write the msg to a net.Conn and handle any errors associated with
// writing, reading or closing
func (c *Communicator) Send(conn net.Conn, w []byte) error {
	c.logger.Debug("Communicator.Send: attempting to write msg to net.Conn", slog.String("msg", string(w)))
	if _, err := conn.Write(w); err != nil {
		c.logger.Error("Communicator.Send: failed to write msg to net.Conn", err, slog.String("msg", string(w)))
		return fmt.Errorf("failed to write command: %w", err)
	}

	c.logger.Debug("Communicator.Send: attempting to close net.Conn", slog.String("msg", string(w)))
	if err := conn.Close(); err != nil {
		c.logger.Error("Communicator.Send: failed to close net.Conn", err)
		return fmt.Errorf("failed to close conn: %w", err)
	}

	return nil
}

// SendAndReceive will write the msg to a net.Conn and receive the response. Any errors associated with
// writing, reading or closing will be handled
func (c *Communicator) SendAndReceive(conn net.Conn, w []byte, r []byte) error {
	c.logger.Debug("Communicator.SendAndReceive: attempting to write msg to net.Conn", slog.String("msg", string(w)))
	if _, err := conn.Write(w); err != nil {
		c.logger.Error("Communicator.SendAndReceive: failed to write msg to net.Conn", err, slog.String("msg", string(w)))
		return fmt.Errorf("failed to write command: %w", err)
	}

	c.logger.Debug("Communicator.SendAndReceive: attempting to read msg from net.Conn")
	if _, err := conn.Read(r); err != nil {
		c.logger.Error("Communicator.SendAndReceive: failed to read msg from net.Conn", err)
		return fmt.Errorf("failed to read after write: %w", err)
	}

	c.logger.Debug("Communicator.SendAndReceive: attempting to close net.Conn", slog.String("msg", string(w)))
	if err := conn.Close(); err != nil {
		c.logger.Error("Communicator.SendAndReceive: failed to close net.Conn", err)
		return fmt.Errorf("failed to close conn: %w", err)
	}

	return nil
}
