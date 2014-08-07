/*
 * Copyright (c) 2014 PolyFloyd
 */

package main

import (
	"log"
	"net"
)

func StartServer(listen string) {
	listener, err := net.Listen("tcp", listen)
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

	buf := make([]byte, VoxelDisplay.NumVoxels() * 3)
	for {
		_, err := conn.Read(buf[:3])
		if err != nil {
			log.Printf("%v disconnected", conn.RemoteAddr())
			break
		}
		switch string(buf[:3]) {
		case "ver":
			conn.Write([]byte(INFO))
		case "put":
			for completed := 0; completed < VoxelDisplay.NumVoxels() * 3; {
				read, err := conn.Read(buf[:VoxelDisplay.NumVoxels()*3 - completed])
				if err != nil {
					log.Printf("%v disconnected", conn.RemoteAddr())
					break
				}
				for i, b := range buf[:read] {
					VoxelDisplay.Buffer[completed+i] = float32(b) / 256
				}
				completed += read
			}
		case "swp":
			VoxelDisplay.SwapBuffers()
		}
	}
}
