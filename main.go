package main

import "flag"

const (
	INFO = "PolyFloyd's LEDCube Simulator v0.1"
	UI_DETAIL     int     = 2
	UI_DRAGDIV    float32 = 240.0
	UI_FOVY       float32 = 45.0
	UI_SPACING    float32 = 8.0
	UI_WIN_H      int     = 768
	UI_WIN_W      int     = 1280
	UI_ZFAR       float32 = 640
	UI_ZNEAR      float32 = 1
	UI_ZOOMACCEL  float32 = 12.0
)

var VoxelDisplay *Display

func main() {
	l := flag.String("l", ":54746", "The TCP host and port for incoming connections")
	cx := flag.Int("cx", 16, "The width of the cube")
	cy := flag.Int("cy", 16, "The length of the cube")
	cz := flag.Int("cz", 16, "The height of the cube")
	flag.Parse()

	go StartServer(*l)
	VoxelDisplay = NewDisplay(*cx, *cy, *cz)
	VoxelDisplay.Start()
}
