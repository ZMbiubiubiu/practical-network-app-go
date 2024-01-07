package chater3

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	var done = make(chan struct{})

	// create a new goroutine for work together with client's side
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("accepting a new conn...")

		go func(c net.Conn) {
			defer func() { done <- struct{}{} }()

			var buf = make([]byte, 1024)
			for {
				n, err := c.Read(buf)
				if err != nil {
					if err != io.EOF {
						t.Fatal(err)
					}
					return
				}
				t.Logf("receive: %q", buf[:n])
			}
		}(conn)
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	conn.Close()
	<-done
}
