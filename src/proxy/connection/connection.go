package connection

import (
	"log"
	"net"
	"redis-proxy/src/resp"
	"redis-proxy/src/resp/util"
)

type Connection struct {
	conn   *net.TCPConn
	reader *util.Reader
	redis  *resp.Redis
}

func (c *Connection) close(err *error) {
	log.Println([]byte("*1\r\n$4\r\nquit\r\n"))
	log.Println("proxy listening closed", c.conn.RemoteAddr().String(), *err)
	c.redis.Close()
	c.conn.Close()
}

func (c *Connection) Out(buf *[]byte) {
	c.conn.Write(*buf)
}

func (c *Connection) In(buf *[]byte) {
	c.Out(c.redis.Do(buf))
}

func (c *Connection) Listen() {
	log.Println("proxy listening new connection", c.conn.RemoteAddr().String())

	for {
		buf, err := c.reader.Read()

		if err == nil {
			c.In(&buf)
		} else {
			c.close(&err)
			break
		}
	}
}

func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{conn: conn, reader: util.NewReader(conn), redis: resp.NewRedis()}
}
