package proxy

import (
	"log"
	"net"
	"redis-proxy/src/proxy/connection"
	"redis-proxy/src/tcp"
)

type Proxy struct {
	port *string
}

func (p *Proxy) handler(conn *net.TCPConn) {
	log.Println("proxy new connection", conn.RemoteAddr().String())

	go connection.NewConnection(conn).Listen()
}

func (p *Proxy) Start() {
	log.Println("starting proxy", *p.port)

	go tcp.NewServer(p.port).Listen(p.handler)

}

func NewProxy(port *string) *Proxy {
	return &Proxy{port: port}
}
