#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import random
import time

SPEED    = 0.01
INTERVAL = .005
TAIL     = 1 / 3

cube    = ledcube.Cube()
dots    = []
counter = INTERVAL

class Dot:

	def __init__(self, x, y, speed):
		self.x     = x
		self.y     = y
		self.speed = speed
		self.pos   = -TAIL

	def update(self):
		self.pos += self.speed

while 1:
	if counter > INTERVAL:
		counter = 0
		x     = random.randint(0, cube.size.x - 1)
		y     = random.randint(0, cube.size.y - 1)
		speed = SPEED + random.uniform(-SPEED, SPEED) / 2
		dots.append(Dot(x, y, speed))

	frame = bytearray(cube.length())
	for dot in dots:
		lo = 1 - dot.pos - TAIL
		hi = 1 - dot.pos
		lo = 0 if lo < 0 else lo
		hi = 1 if hi > 1 else hi
		for z in range(int(lo * cube.size.z), int(hi * cube.size.z)):
			c = 1 - z / cube.size.z - dot.pos
			frame[cube.index(dot.x, dot.y, z, 1)] = int(c * 255)
		dot.update()
		if dot.pos > 1 + TAIL:
			dots.remove(dot)

	cube.frame(frame)
	cube.swap()
	time.sleep(1 / cube.fps)
	counter += 1 / cube.fps

