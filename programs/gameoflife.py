#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import random
import time

COLOR         = (255, 255, 255)
SPEED         = 0.4
MAXREPETITION = 5

class Universe(object):

	def __init__(self, size):
		self.size = size
		self.data = []

		for x in range(cube.size.x):
			self.data.append([])
			for y in range(cube.size.y):
				self.data[x].append([])
				for z in range(cube.size.z):
					self.data[x][y].append(random.random() > 0.5)

	def set(self, x, y, z, state):
		self.data[x][y][z] = state

	def get(self, x, y, z):
		return self.data[x][y][z]

	def alive_neighbours(self, x, y, z, wrap):
		neighbours = 0
		if wrap:
			for x in range(x - 1, x + 2):
				for y in range(y - 1, y + 2):
					for z in range(z - 1, z + 2):
						if self.get(x % self.size.x, y % self.size.y, z % self.size.z):
							neighbours += 1
		else:
			for x in range(max(0, x - 1), min(self.size.x, x + 2)):
				for y in range(max(0, y - 1), min(self.size.y, y + 2)):
					for z in range(max(0, z - 1), min(self.size.z, z + 2)):
						if self.get(x, y, z):
							neighbours += 1
		return neighbours

	def map_frame(self, frame, color):
		for x in range(self.size.x):
			for y in range(self.size.y):
				for z in range(self.size.z):
					frame.set(x, y, z, color if self.get(x, y, z) else (0, 0, 0))


cube       = ledcube.Cube()
t          = 1
universe   = None
repetition = []

def reset():
	global universe, repetition
	universe   = Universe(cube.size)
	repetition = []

reset()

while 1:
	t += SPEED
	if t >= 1:
		t = 0

		frame = cube.make_frame()
		universe.map_frame(frame, COLOR)
		cube.set_frame(frame)

		changes = []

		for x in range(cube.size.x):
			for y in range(cube.size.y):
				for z in range(cube.size.z):
					neighbours = universe.alive_neighbours(x, y, z, True)
					if neighbours <= 1:
						changes.append((x, y, z, False))
					elif neighbours >= 8:
						changes.append((x, y, z, False))
					elif neighbours >= 5:
						changes.append((x, y, z, True))

		for (x, y, z, state) in changes:
			universe.set(x, y, z, state)

		if MAXREPETITION > 0:
			repetition.append(len(changes))
			if len(repetition) > MAXREPETITION:
				repetition.pop(0)

				repeating = True
				last      = repetition[0]
				for n in repetition[1:]:
					repeating = repeating and last == n
					last      = n

				if repeating:
					reset()

	time.sleep(1 / cube.fps)
