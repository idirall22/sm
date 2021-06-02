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
func NewP2PManager() *P2PManager {
	return &P2PManager{
		connections: make(map[string]IConnection),
		m:           &sync.RWMutex{},
	}
}

// P2PManager implementation of the p2p manager interface.
type P2PManager struct {
	connections map[string]IConnection
	m           *sync.RWMutex
}

// GetConnection return existing
// connection or create one and store it, this caller should be
// blocked until the connection is returned (thread safe).
func (p *P2PManager) GetConnection(ipAddress string) IConnection {
	p.m.Lock()
	defer p.m.Unlock()

	// check if a connection already exists.
	conn, ok := p.connections[ipAddress]
	if ok {
		return conn
	}

	// create a new connection
	newConn := NewConnection(ipAddress)
	p.connections[ipAddress] = newConn

	return newConn
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
