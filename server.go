package main

import (
	"log"
	"net"
)

func StartServer() {
	listener, err := net.Listen("tcp", Config.String("net.listen"))
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
		_, err := conn.Read(buf[:3])
		if err != nil {
			log.Printf("%v disconnected", conn.RemoteAddr())
			break
		}
		switch string(buf[:3]) {
		case "frm":
			for completed := 0; completed < TotalVoxels * 3; {
				read, err := conn.Read(buf[:TotalVoxels*3 - completed])
				if err != nil {
					log.Printf("%v disconnected", conn.RemoteAddr())
					break
				}
				if read + completed > TotalVoxels * 3 {
					conn.Close()
					log.Printf("%v disconnected (frm overflow %v)", conn.RemoteAddr(), (completed+read) - TotalVoxels*3)
					return
				}
				for i, b := range buf[:read] {
					DisplayBackBuffer[completed+i] = float32(b) / 256
				}
				completed += read
			}
		case "swp":
			SwapDisplayBuffer()
		default:
			conn.Write([]byte("err\n"))
		}
	}
}
