package main

import (
	"log"
	"net"

	smp2p "github.com/idirall22/sm/p2p"
)

func main() {
	pool := smp2p.NewP2PPool()
	_ = pool

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()
}
