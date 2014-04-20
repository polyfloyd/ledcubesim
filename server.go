package main

import (
	"log"
	"net"
)

var isConnected = false

func StartServer() {
	listener, err := net.Listen("tcp", ":54746")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		if isConnected {
			conn.Write([]byte("A pink unicorn is already connected\n"))
			conn.Close()
		} else{
			go communicate(conn)
		}
	}
}

func communicate(conn net.Conn) {
	isConnected = true
	log.Printf("Connected to %v", conn.(*net.TCPConn).RemoteAddr())

	buf := make([]byte, CUBE_TOTAL_VOXELS * 3)

	for {
		read, err := conn.Read(buf)
		if err != nil {
			isConnected = false
			log.Printf("Client disconnected")
			break
		}
		if read >= 3 {
			payload := buf[3:read]
			switch string(buf[:3]) {
			case "frm":
				for i, b := range payload {
					if i > CUBE_TOTAL_VOXELS * 3 {
						break
					}
					DisplayBackBuffer[i] = float32(b) / 256
				}
			case "swp":
				SwapDisplayBuffer()
			default:
				conn.Write([]byte("Invalid command\n"))
			}
		}
	}
}
