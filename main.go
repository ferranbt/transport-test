package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/yamux"
	"github.com/umbracle/minimal/helper/enode"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/minimal/network/transport/rlpx"
)

func main() {

}

func yamuxIt(conn0, conn1 net.Conn) (net.Conn, net.Conn) {
	conns := make(chan net.Conn, 2)

	// server
	go func() {
		session, err := yamux.Server(conn0, nil)
		if err != nil {
			panic(err)
		}

		stream, err := session.Accept()
		if err != nil {
			panic(err)
		}
		conns <- stream
	}()

	// client
	go func() {
		session, err := yamux.Client(conn1, nil)
		if err != nil {
			panic(err)
		}

		stream, err := session.Open()
		if err != nil {
			panic(err)
		}
		conns <- stream
	}()

	c0, c1 := <-conns, <-conns
	return c0, c1
}

func pipeTLS() (net.Conn, net.Conn) {
	c0, c1 := net.Pipe()
	conns := make(chan net.Conn, 2)

	// Instructions to generate the X509 certificates
	// https://gist.github.com/jim3ma/00523f865b8801390475c4e2049fe8c3

	// server
	go func() {
		cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		config.Rand = rand.Reader
		conns <- tls.Server(c0, config)
	}()

	// client
	go func() {
		cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
		if err != nil {
			panic(fmt.Errorf("server: loadkeys: %s", err))
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		conns <- tls.Client(c1, config)
	}()

	conn0, conn1 := <-conns, <-conns
	return yamuxIt(conn0, conn1)
}

func pipeRLPX() (*rlpx.Stream, *rlpx.Stream) {
	conn0, conn1 := net.Pipe()

	prv0, _ := crypto.GenerateKey()
	prv1, _ := crypto.GenerateKey()

	errs := make(chan error, 2)
	var c0, c1 *rlpx.Session

	go func() {
		c0 = rlpx.Server(conn0, prv0, mockInfo(prv0))
		errs <- c0.Handshake()
	}()
	go func() {
		c1 = rlpx.Client(conn1, prv1, &prv0.PublicKey, mockInfo(prv1))
		errs <- c1.Handshake()
	}()

	for i := 0; i < 2; i++ {
		if err := <-errs; err != nil {
			panic(err)
		}
	}

	s0 := c0.OpenStream(5, 10)
	s1 := c1.OpenStream(5, 10)

	return s0, s1
}

func mockInfo(prv *ecdsa.PrivateKey) *rlpx.Info {
	return &rlpx.Info{
		Version:    1,
		Name:       "test-mock",
		ListenPort: 30303,
		Caps:       rlpx.Capabilities{&rlpx.Cap{Name: "eth", Version: 1}},
		ID:         enode.PubkeyToEnode(&prv.PublicKey),
	}
}
