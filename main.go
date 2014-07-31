package main

import "flag"

const INFO = "PolyFloyd's LEDCube Simulator v0.1"

var (
	CUBE_HEIGHT   int
	CUBE_LENGTH   int
	CUBE_WIDTH    int
	SERVER_LISTEN string
	UI_DRAGDIV    float32 = 240.0
	UI_FOVY       float32 = 45.0
	UI_SPACING    float32 = 8.0
	UI_WIN_H      int     = 768
	UI_WIN_W      int     = 1280
	UI_ZFAR       float32 = 640
	UI_ZNEAR      float32 = 1
	UI_ZOOMACCEL  float32 = 12.0
	UI_DETAIL     int     = 2
	VOXEL_TOTAL   int
)

var LEDDisplay *Display

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
	VOXEL_TOTAL = CUBE_WIDTH * CUBE_LENGTH * CUBE_HEIGHT


	LEDDisplay = NewDisplay(CUBE_WIDTH, CUBE_LENGTH, CUBE_HEIGHT)

	go StartServer()
	LEDDisplay.Start()
}
