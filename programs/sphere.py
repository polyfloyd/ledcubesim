#! /usr/bin/env python3

import array
import ledcube
import math
import os
import time

cube = ledcube.Cube((16, 16, 16), 3)
frame = bytearray(16**3 * 3)

acc = 18

for u in range(-acc, acc):
	u /= acc
	for v in range(-acc, acc):
		v /= acc

		x = math.cos(v * math.pi) * math.sin(u * math.pi)
		y = math.sin(v * math.pi) * math.sin(u * math.pi)
		z = math.sin(u * math.pi + math.pi/2)

		i = cube.xyz((
			int((.5 + x/2) * 15),
			int((.5 + y/2) * 15),
			int((.5 + z/2) * 15),
		))

		frame[i+0] = 0
		frame[i+1] = 128
		frame[i+2] = 255

cube.frame(frame)
cube.swap()
time.sleep(0.1)
