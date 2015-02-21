#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import random
import time

cube          = ledcube.Cube()
particles     = []
max_particles = sum(cube.size) ** 2

class Paricle(object):

	def __init__(self):
		self.pos = ledcube.Vector(
			random.uniform(0, 1 + 1 / cube.size.x),
			random.uniform(0, 1 + 1 / cube.size.z),
			0,
		)
		self.decayrate = random.uniform(0.03, 0.05)
		self.intensity = random.uniform(.8, 1)
		self.speed     = random.uniform(0.01, 0.04)

	def update(self):
		self.pos = ledcube.Vector(
			self.pos.x,
			self.pos.y,
			self.pos.z + self.speed,
		)
		self.intensity -= self.decayrate

	def get_color(self):
		return (
			255 * (self.intensity / 2 + .5),
			255 * self.intensity,
			255 * max(self.intensity - .5, 0),
		)

while 1:
	frame = cube.make_frame()

	while len(particles) < max_particles:
		particles.append(Paricle())

	for p in particles:
		p.update()
		if p.intensity <= 0 or p.pos.z > 1:
			particles.remove(p)
		else:
			frame.setf(p.pos.x, p.pos.y, p.pos.z, p.get_color())


	cube.set_frame(frame)
	time.sleep(1 / cube.fps)
