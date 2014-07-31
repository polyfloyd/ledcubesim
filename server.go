package main

import (
	"log"
	"net"
)

func StartServer() {
	listener, err := net.Listen("tcp", SERVER_LISTEN)
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

	buf := make([]byte, VOXEL_TOTAL * 3)
	for {
		_, err := conn.Read(buf[:3])
		if err != nil {
			log.Printf("%v disconnected", conn.RemoteAddr())
			break
		}
		switch string(buf[:3]) {
		case "nfo":
			conn.Write([]byte(INFO+"\n"))
		case "frm":
			for completed := 0; completed < VOXEL_TOTAL * 3; {
				read, err := conn.Read(buf[:VOXEL_TOTAL*3 - completed])
				if err != nil {
					log.Printf("%v disconnected", conn.RemoteAddr())
					break
				}
				for i, b := range buf[:read] {
					LEDDisplay.Buffer[completed+i] = float32(b) / 256
				}
				completed += read
			}
		case "swp":
			LEDDisplay.SwapBuffers()
		}
	}
}
