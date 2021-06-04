package smp2p

import (
	"log"
	"sync"
)

// P2P Peer to peer manager interface.
type P2P interface {
	GetConnection(ipAddress string) IConnection
	OnNewRemoteIConnection(remotePeer string, newConn IConnection)
	Shutdown()
}

// NewP2PManager create a new p2p manager.
func NewP2PManager(localConnection LocalConnection) *P2PManager {
	return &P2PManager{
		localConnection: localConnection,
		connections:     map[string]IConnection{},
		m:               &sync.RWMutex{},
	}
}

// P2PManager implementation of the p2p manager interface.
type P2PManager struct {
	localConnection LocalConnection
	connections     map[string]IConnection
	m               *sync.RWMutex
}

// GetConnection return existing
// connection or create one and store it, this caller should be
// blocked until the connection is returned (thread safe).
func (p *P2PManager) GetConnection(ipAddress string) chan IConnection {

	// add a connection channel
	connChan := make(chan IConnection)

	{
		p.m.RLock()
		defer p.m.RUnlock()

		// check if a connection already exists.
		conn, ok := p.connections[ipAddress]
		if ok {
			connChan <- conn
		}
	}

	// run a concurrent function that create a new connection then store it.
	go func(ipAddress string) {
		conn := p.localConnection(ipAddress)

		p.m.Lock()
		defer p.m.Unlock()
		p.connections[ipAddress] = conn

		connChan <- conn
	}(ipAddress)

	return connChan
}

// OnNewRemoteIConnection A callback function that is called whenever a remote peer
// establishes a connection with the local node (thread safe).
func (p *P2PManager) OnNewRemoteIConnection(remotePeer string, newConn IConnection) error {
	p.m.Lock()
	defer p.m.Unlock()

	// check if there is already a connection with the same ipAddress, if yes
	// the the new connection is closed, else the connection is cached.
	_, ok := p.connections[remotePeer]
	if ok {
		err := newConn.Close()
		if err != nil {
			return err
		}
	}

	p.connections[remotePeer] = newConn

	return nil
}

// Shutdown Graceful shutdown, all background workers
// should be stopped before this method returns.
func (p *P2PManager) Shutdown() {
	p.m.Lock()
	defer p.m.Unlock()

	for remoteAdd, conn := range p.connections {
		// close connection
		err := conn.Close()
		if err != nil {
			log.Fatalf("Error to close a Connection with remote Peer %s: %v", remoteAdd, err)
		}
		delete(p.connections, remoteAdd)
	}
}

func (p *P2PManager) GetConnections() map[string]IConnection {
	return p.connections
}
