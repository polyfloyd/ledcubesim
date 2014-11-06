#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import time

ACCURACY = 24
STEP     = .05

cube         = ledcube.Cube()
frame        = cube.make_frame()
amp          = [1, 0, 0, 0]
amp_index     = 1
steps_current = 0

while (1):
	for i in range(0, 4):
		amp[i] += (1 if (amp_index&(1<<i) != 0) else -1) * STEP
		amp[i] = amp[i] if (amp[i] >= 0) else 0
		amp[i] = amp[i] if (amp[i] <= 1) else 1

	steps_current += 1
	if steps_current >= 1 / STEP:
		steps_current = 0
		amp_index = (amp_index) % (2**4 - 2) + 1

	for u in range(-ACCURACY, ACCURACY):
		u /= ACCURACY
		for v in range(-ACCURACY, ACCURACY):
			v /= ACCURACY
			x = math.cos(v * math.pi) * math.sin(u * math.pi)
			y = math.sin(u * math.pi + math.pi / 2)
			z = math.sin(v * math.pi) * math.sin(u * math.pi)

			sin = math.sin((.5 + v / 2) * math.pi * 2) / 2
			frame.setf(.5 + x / 2, .5 + y / 2, .5 + z / 2, (
				amp[1] * 255 * (.5 + sin) + amp[0] * 255 * (.5 - sin),
				amp[2] * 255 * (.5 - sin) + amp[1] * 255 * (.5 + sin),
				amp[3] * 255 * (.5 - sin) + amp[2] * 255 * (.5 + sin),
			))

	cube.set_frame(frame)
	time.sleep(1 / cube.fps)
