#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import random
import time

STEP     = .005
EDGE_ANG = 2
EDGE_DOT = 4
BASE_DOT = .5
COL_EDGE = [
	(117, 103, 201),
	(194, 122, 221),
]
COL_BASE = [
	(114, 151, 22),
	(255, 0,   128),
	(0,   128, 255),
	(255, 255, 0),
]

cube   = ledcube.Cube()
colors = []
t      = 0

for i in range(0, int(cube.size.z * BASE_DOT)):
	col = lambda: COL_BASE[random.randint(0, len(COL_BASE) - 1)]
	colors.append([
	    col(),
	    col(),
	])

while 1:
	frame = bytearray(cube.length())
	t = (t + STEP) % 1
	for z in range(0, cube.size.z * EDGE_DOT):
		zn = z / cube.size.z / EDGE_DOT

		znt = (zn + t) % 1
		x = math.cos(znt * math.pi * EDGE_ANG)
		y = math.sin(znt * math.pi * EDGE_ANG)

		def index(x, y):
			return cube.index(
				x / 2 + .5,
				y / 2 + .5,
				zn,
			)

		if z % (EDGE_DOT / BASE_DOT) == 0:
			colA = colors[int(z * BASE_DOT / EDGE_DOT)][0]
			colB = colors[int(z * BASE_DOT / EDGE_DOT)][1]
			for i in range(0, 8):
				i = i / 8
				a = index(x * i, y * i)
				b = index(-x * i, -y * i)
				for j in range(0, 3):
					frame[a + j] = colA[j]
				for j in range(0, 3):
					frame[b + j] = colB[j]

		edge = index(x, y)
		for j in range(0, 3):
			frame[edge + j] = COL_EDGE[0][j]
		edge = index(-x, -y)
		for j in range(0, 3):
			frame[edge + j] = COL_EDGE[1][j]

	cube.frame(frame)
	cube.swap()
	time.sleep(1 / cube.fps)
