package main

const (
	CUBE_WIDTH  = 16
	CUBE_LENGTH = 16
	CUBE_HEIGHT = 16

	BACKGROUND   = 0.12
	LED_DISTANCE = 8

	ZOOM_ACCELERATION = 8
	FOVY              = 45.0

	RENDER_OFF = false
)

func main() {
	go StartServer()
	StartDisplay("A Simulator For LED Cubes")
}
