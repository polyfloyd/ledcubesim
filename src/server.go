/*
 * Copyright (c) 2014 PolyFloyd
 */

package main

import (
	"log"
	"net"
	"unsafe"
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
	buf := make([]byte, VoxelDisplay.NumVoxels() * 3)
	main: for {
		_, err := conn.Read(buf[:3])
		if err != nil {
			break
		}
		switch string(buf[:3]) {
		case "ver":
			conn.Write([]byte(INFO))
		case "inf":
			x := *(*[4]byte)(unsafe.Pointer(&VoxelDisplay.CubeWidth))
			conn.Write(x[:])
			y := *(*[4]byte)(unsafe.Pointer(&VoxelDisplay.CubeLength))
			conn.Write(y[:])
			z := *(*[4]byte)(unsafe.Pointer(&VoxelDisplay.CubeHeight))
			conn.Write(z[:])
			conn.Write([]byte{ 3 })
			conn.Write([]byte{ 60 })
		case "put":
			for completed := 0; completed < VoxelDisplay.NumVoxels() * 3; {
				read, err := conn.Read(buf[:VoxelDisplay.NumVoxels()*3 - completed])
				if err != nil {
					break main
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