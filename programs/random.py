#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import os
import time

STEPS = 40

cube         = ledcube.Cube()
currentFrame = bytearray(cube.length())
targetFrame  = bytearray(cube.length())

while 1:
	sourceFrame = targetFrame
	targetFrame = os.urandom(cube.length())
	for i in range(0, STEPS):
		m = math.sin(i / STEPS * math.pi / 2)
		for j in range(0, cube.length()):
			currentFrame[j] = int(sourceFrame[j] * (1 - m) + targetFrame[j] * m)
		cube.frame(currentFrame)
		cube.swap()
		time.sleep(1/60)
