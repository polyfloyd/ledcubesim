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
EDGES = [
	(117, 103, 201),
	(194, 122, 221),
]
BASES = [
	(114, 151, 22),
	(255, 0,   128),
	(0,   128, 255),
	(255, 255, 0),
]

cube   = ledcube.Cube()
colors = []
t      = 0

for i in range(0, int(cube.size.z * BASE_DOT)):
	col = lambda: BASES[random.randint(0, len(BASES) - 1)]
	colors.append([col() for _ in EDGES])

while 1:
	frame = bytearray(cube.length())
	t = (t + STEP) % 1
	for z in range(0, cube.size.z * EDGE_DOT):
		zn = z / cube.size.z / EDGE_DOT

		for i in range(0, len(EDGES)):
			ang = i / len(EDGES)

			znt = (zn + t) % 1
			x = math.cos((znt + ang) * math.pi * EDGE_ANG)
			y = math.sin((znt + ang) * math.pi * EDGE_ANG)
			index = lambda x, y: cube.index(x / 2 + .5, y / 2 + .5, zn)

			if z % (EDGE_DOT / BASE_DOT) == 0:
				col = colors[int(z * BASE_DOT / EDGE_DOT)][i]
				for j in range(0, 8):
					j = j / 8
					a = index(x * j, y * j)
					for k in range(0, 3):
						frame[a + k] = col[k]

			edge = index(x, y)
			for j in range(0, 3):
				frame[edge + j] = EDGES[i][j]

	cube.frame(frame)
	cube.swap()
	time.sleep(1 / cube.fps)
