package main

import "polyfloyd/irix/util"

var Config = util.Config{}

func main() {
	err := Config.LoadString(`{
		"cube": {
			"width":  16,
			"length": 16,
			"height": 16
		},
		"ui": {
			"background": 0.12,
			"spacing":    8.0,
			"zoomAccel":  8.0,
			"showOff":    true,
			"fovy":       45.0
		}
	}`)
	if err != nil {
		panic(err)
	}
	go StartServer()
	StartDisplay("A Simulator For LED Cubes")
}
