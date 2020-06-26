package tcp

import (
	"errors"
	"log"
	"net"
)

type Server struct {
	listener *net.TCPListener
	addr     *net.TCPAddr
	port     *string
}

func (s *Server) Resolve() error {
	if s.port == nil {
		return errors.New("tcp port null")
	}
	addr, err := net.ResolveTCPAddr("tcp", *s.port)

	if err == nil {
		s.addr = addr
	}

	return err
}

func (s *Server) Listener() error {
	if s.addr == nil {
		if err := s.Resolve(); err != nil {
			return err
		}
	}

	listener, err := net.ListenTCP("tcp", s.addr)

	if err == nil {
		s.listener = listener
	}

	return err
}

func (s *Server) Listen(handler func(*net.TCPConn)) {
	if s.listener == nil {
		if err := s.Listener(); err != nil {
			log.Fatal("TCP listen failed")
		}
	}

	for {
		conn, err := s.listener.AcceptTCP()

		if err == nil {
			go handler(conn)
		} else {
			log.Fatal("TCP listen failed")
		}
	}
}

func NewServer(port *string) *Server {
	return &Server{port: port}
}
