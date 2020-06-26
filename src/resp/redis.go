package resp

import (
	"redis-proxy/src/resp/constants"
	"redis-proxy/src/tcp"
	"time"
)

type Redis struct {
	client *tcp.Client
}

func (r *Redis) read() (*[]byte, error) {
	buf, err := r.client.Read()
	return &buf, err

}

func (r *Redis) write(buf *[]byte) error {
	err := r.client.Write(buf)

	if err != nil {
		if err := r.connect(); err == nil {
			for x := 0; x < 5; x++ {
				if err := r.client.Write(buf); err == nil {
					return nil
				} else {
					time.Sleep(time.Second)
				}
			}
		}
	}

	return err
}

func (r *Redis) connect() error {
	err := r.client.Connect()

	if err != nil {
		for x := 0; x < 10; x++ {
			if err := r.client.Clear(&REDIS_ADDRESS, &constants.REDIS_CLOSE); err == nil {
				return nil
			} else {
				time.Sleep(time.Second)
			}
		}
	}

	return err
}

func (r *Redis) Do(buf *[]byte) *[]byte {
	if err := r.connect(); err != nil {
		return &constants.CONNECTION_ERROR
	}

	if err := r.write(buf); err != nil {
		return &constants.WRITING_ERROR
	}

	if buf, err := r.read(); err != nil {
		return &constants.READING_ERROR
	} else {
		return buf
	}
}

func (r *Redis) Close() {
	r.client.Close(&constants.REDIS_CLOSE)
}

func NewRedis() *Redis {
	return &Redis{client: tcp.NewClient(&REDIS_ADDRESS)}
}
