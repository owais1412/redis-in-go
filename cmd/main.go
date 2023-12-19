package main

import (
	"log/slog"
	"net"
	"os"
	"redis/internal/aof"
	"redis/internal/resp"
	"slices"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		slog.Error("Error while listening to tcp :6379: ", err)
		panic(err)
	}

	aofFile, exists := os.LookupEnv("REDIS_DB_PATH")
	if !exists {
		aofFile = "redis_db.aof"
	}

	aof, err := aof.NewAof(aofFile)
	if err != nil {
		slog.Error("Error while creating append only file: ", err)
		panic(err)
	}
	defer aof.Close()

	aof.Read(func(value resp.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := resp.Handlers[command]
		if !ok {
			slog.Error("Invalid command: " + command)
			return
		}

		handler(args)
	})

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

		// persit to the disk
		if slices.Index([]string{"SET", "HSET", "DEL"}, command) != -1 {
			aof.Write(value)
		}

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
