package chater3

import (
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { listener.Close() }()

	t.Logf("bound to %q", listener.Addr())

	for {
		// Accept will block until the listener detects an incoming connection and completes the TCP handshake process
		conn, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}

		go func(c net.Conn) {
			defer conn.Close()

			// write your code here
		}(conn)
	}
}
