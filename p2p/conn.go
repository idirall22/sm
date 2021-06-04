package smp2p

// IConnection interface.
type IConnection interface {
	Open() error
	Close() error
	IPAddress() string
}

type LocalConnection func(ipAddress string) IConnection

// NewConnection create a new connection.
func NewConnection(ipAddress string) IConnection {
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

// Open new connection.
func (c *Connection) IPAddress() string {
	return c.ipAddress
}

// Close connection.
func (c *Connection) Close() error {
	return nil
}
