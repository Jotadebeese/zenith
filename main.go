package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Zenith Engine listening on :6379...")

	storage := NewStorage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn, storage)
	}
}

func handleConnection(conn net.Conn, storage *Storage) {
	defer conn.Close()
	
	resp := NewResp(conn)

	for {
		value, err := resp.Read()
		if err != nil {
			fmt.Printf("Client %s disconnected or encountered error: %v\n", conn.RemoteAddr().String(), err)
			return
		}

		if value.typ != "array" || len(value.array) == 0 {
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)

		switch command {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))

		case "SET":
			if (len(value.array) < 3) {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
				continue
			}
			key := value.array[1].bulk
			val := value.array[2].bulk
			var ttl time.Duration

			if len(value.array) >= 5 && strings.ToUpper(value.array[3].bulk) == "EX" {
				seconds, err := strconv.Atoi(value.array[4].bulk)
				if err == nil && seconds > 0 {
					ttl = time.Duration(seconds) * time.Second
				}
			}
			storage.Set(key, val, ttl)
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if (len(value.array) < 2) {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
				continue
			}
			key := value.array[1].bulk
			val, ok := storage.Get(key)
			if !ok {
				conn.Write([]byte("$-1\r\n"))
			} else {
				conn.Write(fmt.Appendf(nil, "$%d\r\n%s\r\n", len(val), val))
			}

		case "DEL":
			if (len(value.array) < 2) {
				conn.Write([]byte("-ERR wrong number of arguments for 'del' command\r\n"))
				continue
			}
			key := value.array[1].bulk
			storage.Del(key)
			conn.Write([]byte("+OK\r\n"))
			
		default:
			conn.Write(fmt.Appendf(nil, "-ERR unknown command '%s'\r\n", command))
		}
	}
}