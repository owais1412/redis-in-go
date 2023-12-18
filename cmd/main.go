package main

import (
	"log/slog"
	"net"
	"redis/internal/resp"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		slog.Error("Error while listening to tcp :6379: ", err)
		panic(err)
	}

	conn, err := l.Accept()
	if err != nil {
		slog.Error("Error while making connection: ", err)
		panic(err)
	}
	defer conn.Close()

	for {
		resp := resp.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			slog.Error("Error while reading: ", err)
			return
		}
		slog.Info("Value: ", value)

		conn.Write([]byte("+OK\r\n"))
	}

}
