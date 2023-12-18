package main

import (
	"io"
	"log/slog"
	"net"
	"os"
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
		buf := make([]byte, 1024)

		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			slog.Error("Error while reading: ", err)
			os.Exit(1)
			return
		}
		conn.Write([]byte("+OK\r\n"))
	}

}
