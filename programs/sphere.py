#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import time

ACCURACY = 24
STEP     = .09

cube         = ledcube.Cube()
frame        = bytearray(16**3 * 3)
amp          = [1, 0, 0, 0]
ampIndex     = 1
stepsCurrent = 0

while (1):
	for i in range(0, 4):
		amp[i] += (1 if (ampIndex&(1<<i) != 0) else -1) * STEP
		amp[i] = amp[i] if (amp[i] >= 0) else 0
		amp[i] = amp[i] if (amp[i] <= 1) else 1

	stepsCurrent += 1
	if stepsCurrent >= 1 / STEP:
		stepsCurrent = 0
		ampIndex = (ampIndex) % (2**4 - 2) + 1

	for u in range(-ACCURACY, ACCURACY):
		u /= ACCURACY
		for v in range(-ACCURACY, ACCURACY):
			v /= ACCURACY
			x = math.cos(v * math.pi) * math.sin(u * math.pi)
			y = math.sin(v * math.pi) * math.sin(u * math.pi)
			z = math.sin(u * math.pi + math.pi / 2)
			i = cube.xyz((
				int((.5 + x/2) * 15.5),
				int((.5 + y/2) * 15.5),
				int((.5 + z/2) * 15.5),
			))

			sin = math.sin((.5 + v / 2) * math.pi * 2) / 2
			frame[i+0] = int(amp[1] * 255 * (.5 + sin) + amp[0] * 255 * (.5 - sin))
			frame[i+1] = int(amp[2] * 255 * (.5 - sin) + amp[1] * 255 * (.5 + sin))
			frame[i+2] = int(amp[3] * 255 * (.5 - sin) + amp[2] * 255 * (.5 + sin))

	cube.frame(frame)
	cube.swap()
	time.sleep(0.05)
