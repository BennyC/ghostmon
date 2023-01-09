package ghostmon

import (
	"fmt"
	"io"
)

type Communicator struct {
	connection io.Writer
}

// NewCommunicator will create a Communicator instance, based upon the io.Writer
// provided. This will normally be a net.Conn
func NewCommunicator(connection io.Writer) *Communicator {
	return &Communicator{
		connection: connection,
	}
}

// CutOver will communicate a cutover request with gh-ost through the io.Writer
// within the Communicator instance
func (a *Communicator) CutOver() error {
	_, err := a.connection.Write([]byte("cutover"))
	if err != nil {
		return fmt.Errorf("unable to send command: %w", err)
	}

	return nil
}
