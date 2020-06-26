package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"redis-proxy/src/resp/constants"
	"strconv"
)

type Reader struct {
	reader *bufio.Reader
}

func (r *Reader) Read() ([]byte, error) {
	buf, err := r.read()

	if r.reader.Buffered() != 0 {
		b, err := r.Read()
		buf = append(buf, b...)
		return buf, err
	}

	return buf, err
}

func (r *Reader) read() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}

	return r.do(line)
}

func (r *Reader) do(line []byte) ([]byte, error) {
	switch line[0] {
	case constants.SIMPLE_STRING, constants.INTEGER, constants.ERROR:
		return line, nil
	case constants.BULK_STRING:
		return r.readBulkString(line)
	case constants.ARRAY:
		return r.readArray(line)
	default:
		return nil, errors.New("Invalid syntax")
	}
}

func (r *Reader) readLine() (line []byte, err error) {
	line, err = r.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	if len(line) > 1 && line[len(line)-2] == '\r' {
		return line, nil
	}

	return nil, errors.New("Invalid syntax")
}

func (r *Reader) readBulkString(line []byte) ([]byte, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}
	if count == -1 {
		return line, nil
	}

	buf := make([]byte, len(line)+count+2)
	copy(buf, line)
	_, err = io.ReadFull(r.reader, buf[len(line):])
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *Reader) readArray(line []byte) ([]byte, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		buf, err := r.read()
		if err != nil {
			return nil, err
		}
		line = append(line, buf...)
	}

	return line, nil
}

func (r *Reader) getCount(line []byte) (int, error) {
	return strconv.Atoi(string(line[1:bytes.IndexByte(line, '\r')]))
}

func NewReader(r io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(r)}
}
