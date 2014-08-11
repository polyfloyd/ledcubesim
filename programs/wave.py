#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import time

STEPS = 160

cube = ledcube.Cube()

def wave(x, y, t):
	x = x / 2
	y = y / 2
	n = lambda a: (a + t) % 1 * 2 - 1
	zx  = math.cos(n(x) * math.pi)
	zy  = math.sin(n(y) * math.pi)
	z = (zx * zy + 1) / 2
	col = (0, z, (1 - z))
	return (z, col)

while (1):
	for i in range(0, STEPS):
		cube.graph2(lambda x, y: wave(x, y, i / STEPS))
		time.sleep(1 / cube.fps)
