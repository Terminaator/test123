package tcp

import (
	"errors"
	"net"
	"redis-proxy/src/resp/util"
)

type Client struct {
	addr     *string
	dialed   *net.TCPConn
	resolved *net.TCPAddr
	reader   *util.Reader
}

func (c *Client) Write(buf *[]byte) error {
	if c.dialed != nil {
		_, err := c.dialed.Write(*buf)
		return err
	} else {
		return errors.New("No connection")
	}
}

func (c *Client) Read() ([]byte, error) {
	if c.dialed != nil {
		return c.reader.Read()
	} else {
		return nil, errors.New("No connection")
	}
}

func (c *Client) Close(buf *[]byte) {
	if c.dialed != nil {
		c.dialed.Close()
	}
}

func (c *Client) Connect() error {
	if c.addr == nil {
		return errors.New("addr null")
	}

	if c.resolved == nil {
		if err := c.Resolve(); err != nil {
			return errors.New("client failed to resolve")
		}
	}

	if c.dialed == nil {
		if err := c.Dial(); err != nil {
			return errors.New("client failed to dial")
		}
	}

	return nil
}

func (c *Client) Resolve() error {
	if c.addr == nil {
		return errors.New("client addr is null")
	}

	addr, err := net.ResolveTCPAddr("tcp", *c.addr)

	if err == nil {
		c.resolved = addr
	}

	return err
}

func (c *Client) Dial() error {
	if c.resolved == nil {
		if err := c.Resolve(); err != nil {
			return err
		}
	}

	conn, err := net.DialTCP("tcp", nil, c.resolved)

	if err == nil {
		c.dialed = conn
		c.reader = util.NewReader(conn)
	}

	return err
}

func (c *Client) Clear(addr *string, buf *[]byte) error {
	c.Close(buf)

	c.dialed = nil
	c.reader = nil
	c.resolved = nil
	c.addr = addr

	return c.Dial()
}

func NewClient(addr *string) *Client {
	return &Client{addr: addr}
}
