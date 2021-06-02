package smp2p_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	smp2p "github.com/idirall22/sm/p2p"
	mockc "github.com/idirall22/sm/p2p/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// TestGetConnectionParallel test get connection.
func TestGetConnectionParallel(t *testing.T) {
	t.Parallel()

	// create a pool
	p2pManager := smp2p.NewP2PManager()

	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ipAddress := gofakeit.IPv4Address()
			conn := p2pManager.GetConnection(ipAddress)
			require.NotNil(t, conn)
		}()
	}

	wg.Wait()
}

// TestOnNewRemoteIConnectionParallel test add remote connection.
func TestOnNewRemoteIConnectionParallel(t *testing.T) {
	t.Parallel()

	// create a pool
	p2pManager := smp2p.NewP2PManager()

	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()

			// add a new connection
			remoteAddress, conn := newMockConnection()
			err := p2pManager.OnNewRemoteIConnection(remoteAddress, conn)
			require.NoError(t, err)
		}()
	}
	wg.Wait()
}

// TestShutdownParallel shuttdown and clean.
func TestShutdownParallel(t *testing.T) {
	t.Parallel()

	// create a new pool
	p2pManager := smp2p.NewP2PManager()

	wg := &sync.WaitGroup{}

	// add connections concurently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			remoteAddress, conn := newMockConnection()
			p2pManager.OnNewRemoteIConnection(remoteAddress, conn)
		}()
	}

	// shutdow the peer
	p2pManager.Shutdown()
}

// generate a ransom mock connection for tests
func newMockConnection() (string, *mockc.IConnection) {
	remoteAddress := gofakeit.IPv4Address()

	mockConnection := &mockc.IConnection{}
	mockConnection.On("Close").Return(nil).Times(1)

	return remoteAddress, mockConnection
}
