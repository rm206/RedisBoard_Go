package main

import "sync"

func ping(args []Value) Value {
	if len(args) != 0 {
		return Value{type_of: "string", str: args[0].bulk}
	}
	return Value{type_of: "string", str: "PONG"}
}

var SET_map = map[string]string{}
var SET_map_mutex = sync.RWMutex{} // multiple readers or one writer at a time

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{type_of: "error", str: "ERROR : wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SET_map_mutex.Lock()
	SET_map[key] = value
	SET_map_mutex.Unlock()

	return Value{type_of: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{type_of: "error", str: "ERROR : wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SET_map_mutex.RLock()
	value, ok := SET_map[key]
	SET_map_mutex.RUnlock()

	if !ok {
		return Value{type_of: "null", str: "nil"}
	}

	return Value{type_of: "bulk", bulk: value}
}

var HSETs_map = map[string]map[string]string{}
var HSETs_map_Mu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{type_of: "error", str: "ERROR : wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETs_map_Mu.Lock()
	if _, ok := HSETs_map[hash]; !ok {
		HSETs_map[hash] = map[string]string{}
	}
	HSETs_map[hash][key] = value
	HSETs_map_Mu.Unlock()

	return Value{type_of: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{type_of: "error", str: "ERROR : wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETs_map_Mu.RLock()
	value, ok := HSETs_map[hash][key]
	HSETs_map_Mu.RUnlock()

	if !ok {
		return Value{type_of: "null"}
	}

	return Value{type_of: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{type_of: "error", str: "ERROR : wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETs_map_Mu.RLock()
	value, ok := HSETs_map[hash]
	HSETs_map_Mu.RUnlock()

	if !ok {
		return Value{type_of: "null"}
	}

	values := []Value{}
	for k, v := range value {
		values = append(values, Value{type_of: "bulk", bulk: k})
		values = append(values, Value{type_of: "bulk", bulk: v})
	}

	return Value{type_of: "array", array: values}
}

func command_init(args []Value) Value {
	return Value{type_of: "null"}
}

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"COMMAND": command_init,
}

func getHSETs_mapString() string {
	HSETs_map_Mu.RLock()
	defer HSETs_map_Mu.RUnlock()

	str := ""
	for k, v := range HSETs_map {
		str += k + ":\n"
		for kk, vv := range v {
			str += "  " + kk + " : " + vv + "\n"
		}
	}

	return str
}

func clearHSETs_map() {
	HSETs_map_Mu.Lock()
	defer HSETs_map_Mu.Unlock()

	HSETs_map = map[string]map[string]string{}
}
