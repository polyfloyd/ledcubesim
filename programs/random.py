#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import os
import time

STEPS = 16

cube         = ledcube.Cube((16, 16, 16), 3)
currentFrame = bytearray(cube.length())
targetFrame  = bytearray(cube.length())

while 1:
	sourceFrame = targetFrame
	targetFrame = os.urandom(cube.length())
	for i in range(0, STEPS):
		m = i / STEPS
		for j in range(0, cube.length()):
			currentFrame[j] = int(sourceFrame[j] * (1 - m) + targetFrame[j] * m)
		cube.frame(currentFrame)
		cube.swap()
		time.sleep(1/60)
