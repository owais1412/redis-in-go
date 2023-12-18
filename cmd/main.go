package main

import (
	"log/slog"
	"net"
	"redis/internal/resp"
	"strings"
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
		r := resp.NewResp(conn)
		value, err := r.Read()
		if err != nil {
			slog.Error("Error while reading: ", err)
			return
		}

		if value.Typ != "array" {
			slog.Error("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			slog.Error("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := resp.NewWriter(conn)

		handler, ok := resp.Handlers[command]
		if !ok {
			slog.Error("Invalid command: " + command)
			writer.Write(resp.Value{Typ: "string", Str: ""})
			continue
		}

		res := handler(args)
		writer.Write(res)
	}
}
