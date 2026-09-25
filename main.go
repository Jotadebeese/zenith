package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Zenith Engine listening on :6379...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Client %s disconnected or encountered error: %v\n", conn.RemoteAddr().String(), err)
			return
		}
		fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))
		conn.Write(buf[:n])
	}

}