package resp

import "sync"

var Handlers = map[string]func([]Value) Value{
	"COMMAND": command,
	"PING":    ping,
	"GET":     get,
	"SET":     set,
	"DEL":     del,
	"EXISTS":  exists,
	"HGET":    hget,
	"HSET":    hset,
	"HGETALL": hgetall,
}

func command(args []Value) Value {
	return Value{Typ: "string", Str: "OK"}
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}
	}

	return Value{Typ: "string", Str: args[0].Bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{Typ: "string", Str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: value}
}

func del(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"}
	}

	keysDeleted := 0

	SETsMu.Lock()
	for _, k := range args {
		if _, ok := SETs[k.Bulk]; ok {
			keysDeleted++
		}
		delete(SETs, k.Bulk)
	}
	SETsMu.Unlock()

	return Value{Typ: "integer", Num: keysDeleted}
}

func exists(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'exists' command"}
	}

	keysFound := 0

	SETsMu.RLock()
	for _, k := range args {
		if _, ok := SETs[k.Bulk]; ok {
			keysFound++
		}
	}
	SETsMu.RUnlock()

	return Value{Typ: "integer", Num: keysFound}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{Typ: "string", Str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{Typ: "null"}
	}

	values := []Value{}
	for k, v := range value {
		values = append(values, Value{Typ: "bulk", Bulk: k})
		values = append(values, Value{Typ: "bulk", Bulk: v})
	}

	return Value{Typ: "array", Array: values}
}
