package smp2p

// IConnection interface.
type IConnection interface {
	Open() error
	Close() error
}

// NewConnection create a new connection.
func NewConnection(ipAddress string) *Connection {
	return &Connection{
		ipAddress: ipAddress,
	}
}

// Connection is implementation of the IConnection interface.
type Connection struct {
	ipAddress string
}

// Open new connection.
func (c *Connection) Open() error {
	return nil
}

// Close connection.
func (c *Connection) Close() error {
	return nil
}
