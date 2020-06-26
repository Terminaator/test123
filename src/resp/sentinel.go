package resp

import (
	"fmt"
	"log"
	"redis-proxy/src/tcp"
	"strings"
	"time"
)

type Sentinel struct {
	client *tcp.Client
	file   *string
}

func (s *Sentinel) read() ([]byte, error) {
	return s.client.Read()
}

func (s *Sentinel) write() error {
	buf := []byte(fmt.Sprintf("sentinel get-master-addr-by-name %s\n", SENTINEL_NAME))
	return s.client.Write(&buf)
}

func (s *Sentinel) connect() error {
	return s.client.Connect()
}

func (s *Sentinel) add(addr *string) {
	if *addr != REDIS_ADDRESS {
		if REDIS_INITIALIZER {
			NewInitializer(s.file).Initialize()
		}

		REDIS_ADDRESS = *addr
	}
}

func (s *Sentinel) check(buf *[]byte) {
	parts := strings.Split(string(*buf), "\r\n")

	if len(parts) > 4 {
		addr := fmt.Sprintf("%s:%s", parts[2], parts[4])
		s.add(&addr)
	}
}

func (s *Sentinel) do() {
	if err := s.connect(); err != nil {
		return
	}

	if err := s.write(); err != nil {
		return
	}

	if buf, err := s.read(); err != nil {
		return
	} else {
		s.check(&buf)
	}

}

func (s *Sentinel) Start() {
	log.Println("starting sentinel")

	for {
		s.do()
		time.Sleep(120 * time.Second)
	}
}

func NewSentinel(file *string) *Sentinel {
	return &Sentinel{client: tcp.NewClient(&SENTINEL_ADDRESS), file: file}
}
