package main

import "polyfloyd/irix/util"

const INFO = "PolyFloyd's LEDCube Simulator v0.1\n"

var Config = util.Config{}

func main() {
	err := Config.LoadFile("./config.json")
	if err != nil {
		panic(err)
	}
	go StartServer()
	StartDisplay("A Simulator For LED Cubes")
}
