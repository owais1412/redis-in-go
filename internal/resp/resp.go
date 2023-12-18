package resp

import (
	"bufio"
	"io"
	"log/slog"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),
	}
}

func (r *Resp) readLine() ([]byte, int, error) {
	var line []byte
	var n int

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		line = append(line, b)
		n++

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (int, int, error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(i), n, nil
}

func (r *Resp) Read() (Value, error) {
	typ, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typ {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		slog.Info("Unknown type: " + string(typ))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	l, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.array = make([]Value, 0)
	for i := 0; i < l; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.array = append(v.array, val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	l, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, l)

	r.reader.Read(bulk)
	v.bulk = string(bulk)

	// read the trailing \r\n
	r.readLine()

	return v, nil
}
