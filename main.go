package main

import (
	"fmt"
	"net"
	"sync"
)

// nc localhost 50005
type Server struct {
	Clients map[string]net.Conn
	sync.Mutex
}

func main() {

	l, err := net.Listen("tcp4", "localhost:50005")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer l.Close()

	clients := make(map[string]net.Conn)
	s := Server{Clients: clients}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr()

	s.Lock()
	s.Clients[remoteAddr.String()] = conn
	s.Unlock()
	fmt.Printf("Client %s connected\n", remoteAddr)

	buffer := make([]byte, 1024)

	for {
		_, err := conn.Read(buffer)
		if err != nil {
			delete(s.Clients, conn.RemoteAddr().String())
			return
		}

		for addr, client := range s.Clients {
			if addr != conn.RemoteAddr().String() {
				fmt.Printf("Sending %s", buffer)
				_, err = client.Write(buffer)
				if err != nil {
					fmt.Println("Error:", err)
					delete(s.Clients, conn.RemoteAddr().String())
					return
				}
			}
		}
	}
}
