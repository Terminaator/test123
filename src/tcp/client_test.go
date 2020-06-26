package tcp

import (
	"net"
	"testing"
	"time"
)

var (
	port string = ":8000"
)

func makeServer(handler func(*net.TCPConn)) *Server {
	server := NewServer(&port)

	go func() {
		server.Listen(handler)
	}()

	return server
}

func TestClientWrite(t *testing.T) {

	buf := []byte("+test\r\n")

	handled := false

	handler := func(c *net.TCPConn) {
		handled = true
	}

	server := makeServer(handler)
	client := NewClient(&port)

	client.Connect()

	client.Write(&buf)

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}

	if !handled {
		t.Fail()
	}

	server.Close()
}

func TestClientRead(t *testing.T) {

	buf := []byte("+test\r\n")

	response := "+OK\r\n"

	handler := func(c *net.TCPConn) {
		c.Write([]byte(response))
	}

	server := makeServer(handler)
	client := NewClient(&port)

	client.Connect()

	client.Write(&buf)

	defer server.Close()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}

	if b, e := client.Read(); e == nil && string(b) == response {
	} else {
		t.Fail()
	}
}
