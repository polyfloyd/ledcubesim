#! /usr/bin/env python3

import array
import ledcube
import math
import os
import time

cube = ledcube.Cube((16, 16, 16), 3)
frame = bytearray(16**3 * 3)

acc = 24
ampStep  = .1

ampIndex = 1
amp      = [0, 0, 0]

stepsPerCycle = 1 / ampStep
stepsCurrent  = 0

while (1):

	for i in range(0, 3):
		amp[i] += (1 if (ampIndex&(1<<i) != 0) else -1) * ampStep
		amp[i] = amp[i] if (amp[i] >= 0) else 0
		amp[i] = amp[i] if (amp[i] <= 1) else 1

	stepsCurrent += 1
	if stepsCurrent == stepsPerCycle:
		stepsCurrent = 0
		ampIndex = ((ampIndex + 1) % 5) + 1

	for u in range(-acc, acc):
		u /= acc
		for v in range(-acc, acc):
			v /= acc

			x = math.cos(v * math.pi) * math.sin(u * math.pi)
			y = math.sin(v * math.pi) * math.sin(u * math.pi)
			z = math.sin(u * math.pi + math.pi / 2)

			i = cube.xyz((
				int((.5 + x/2) * 15.5),
				int((.5 + y/2) * 15.5),
				int((.5 + z/2) * 15.5),
			))

			# frame[i+0] = int(amp[0] * 255)
			# frame[i+1] = int(amp[1] * 255)
			# frame[i+2] = int(amp[2] * 255)
			un = .5 + u/2
			vn = .5 + v/2
			frame[i]   = int(amp[0] * 255 * (1-un))
			frame[i+1] = int(amp[1] * 255 * (1-vn))
			frame[i+2] = int(amp[2] * 255 * vn)

	cube.frame(frame)
	cube.swap()
	time.sleep(0.05)
