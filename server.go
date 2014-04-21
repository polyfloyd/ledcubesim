package main

import (
	"log"
	"net"
)

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
		go communicate(conn)
	}
}

func communicate(conn net.Conn) {
	log.Printf("%v connected", conn.RemoteAddr())

	buf := make([]byte, TotalVoxels * 3)

	for {
		read, err := conn.Read(buf)
		if err != nil {
			log.Printf("%v disconnected", conn.RemoteAddr())
			break
		}
		if read >= 3 {
			payload := buf[3:read]
			switch string(buf[:3]) {
			case "frm":
				for i, b := range payload {
					if i > TotalVoxels * 3 {
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
