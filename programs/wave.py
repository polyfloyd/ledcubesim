#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import time

SPEED_X = 0.005
SPEED_Y = 0.01

cube = ledcube.Cube()

def wave(x, y, xt, yt):
	n = lambda a, t: (a + t) % 1 * 2 - 1
	zx  = math.cos(n(x / 2, xt) * math.pi)
	zy  = math.sin(n(y / 2, yt) * math.pi)
	z = (zx * zy + 1) / 2
	return (z, (0, z, (1 - z)))

xt = 0
yt = 0
while (1):
	xt = (xt + SPEED_X) % 1
	yt = (yt + SPEED_Y) % 1
	cube.graph2(lambda x, y: wave(x, y, xt, yt))
	time.sleep(1 / cube.fps)
