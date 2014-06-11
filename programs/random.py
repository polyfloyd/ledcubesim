#! /usr/bin/env python3

import ledcube
import os
import time

STEPS = 16

cube = ledcube.Cube((16, 16, 16), 3)
currentFrame = bytearray(cube.length())
targetFrame  = bytearray(cube.length())

while 1:
	sourceFrame = targetFrame
	targetFrame = os.urandom(cube.length())
	for i in range(0, STEPS):
		m = i / STEPS
		for x in range(0, 16):
			for y in range(0, 16):
				for z in range(0, 16):
					for c in range(0, cube.colors):
						pos = cube.xyz((x, y, z), c)
						currentFrame[pos] = int(sourceFrame[pos] * (1 - m) + targetFrame[pos] * m)
		cube.frame(currentFrame)
		cube.swap()
		time.sleep(1/60)
