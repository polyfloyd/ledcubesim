package main

import "flag"

const INFO = "PolyFloyd's LEDCube Simulator v0.1\n"

var (
	SERVER_LISTEN string
	CUBE_WIDTH    int
	CUBE_HEIGHT   int
	CUBE_LENGTH   int
	UI_ZOOMACCEL  float32 = 12.0
	UI_SPACING    float32 = 8.0
	UI_FOVY       float32 = 45.0
)

func main() {
	serverListen := flag.String("p", ":54746", "The TCP host and port for incoming connections")
	cubeWidth    := flag.Int("cx", 16, "The width of the cube")
	cubeLength   := flag.Int("cy", 16, "The length of the cube")
	cubeHeight   := flag.Int("cz", 16, "The height of the cube")
	flag.Parse()
	SERVER_LISTEN = *serverListen
	CUBE_WIDTH    = *cubeWidth
	CUBE_LENGTH   = *cubeLength
	CUBE_HEIGHT   = *cubeHeight

	go StartServer()
	StartDisplay("A Simulator For LED Cubes")
}
