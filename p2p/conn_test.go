package smp2p_test

import (
	"fmt"
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
	p2pManager := smp2p.NewP2PManager(smp2p.NewConnection)

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// add a random connection time.
			connectionTime := time.Millisecond * time.Duration(rand.Intn(1000))
			time.Sleep(connectionTime)

			ipAddress := gofakeit.IPv4Address()
			conn := p2pManager.GetConnection(ipAddress)
			c := <-conn
			fmt.Printf("Connection ID %d Ready Address: %s- %v\n", id, c.IPAddress(), connectionTime)

		}(i)
	}
	wg.Wait()
}

// TestOnNewRemoteIConnectionParallel test add remote connection.
func TestOnNewRemoteIConnectionParallel(t *testing.T) {
	t.Parallel()

	// create a pool
	p2pManager := smp2p.NewP2PManager(smp2p.NewConnection)

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

	mockLocalConnection := func(ipAddress string) smp2p.IConnection {
		mockConnection := &mockc.IConnection{}
		mockConnection.On("Close").Return(nil).Times(1)
		return mockConnection
	}

	// create a new pool
	p2pManager := smp2p.NewP2PManager(mockLocalConnection)

	wg := &sync.WaitGroup{}

	// add connections concurently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		if i%2 == 0 {
			go func() {
				defer wg.Done()
				remoteAddress, conn := newMockConnection()
				p2pManager.OnNewRemoteIConnection(remoteAddress, conn)
			}()
		} else {
			go func() {
				defer wg.Done()
				ipAddress := gofakeit.IPv4Address()
				conn := p2pManager.GetConnection(ipAddress)
				<-conn
			}()
		}
	}

	wg.Wait()
	// shutdow the peer
	p2pManager.Shutdown()
	require.Empty(t, p2pManager.GetConnections())
}

// generate a ransom mock connection for tests
func newMockConnection() (string, *mockc.IConnection) {
	remoteAddress := gofakeit.IPv4Address()

	mockConnection := &mockc.IConnection{}
	mockConnection.On("Close").Return(nil).Times(1)

	return remoteAddress, mockConnection
}
