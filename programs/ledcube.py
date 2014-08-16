#
# Copyright (c) 2014 PolyFloyd
#

from collections import namedtuple
import os
import socket
import sys

Dimension = namedtuple("Dimension", "x y z")

def determineConnection():
	addr = os.getenv("CUBE_ADDR")
	port = os.getenv("CUBE_PORT")

	for (i, arg) in enumerate(sys.argv[1:]):
		if arg == "-a":
			addr = sys.argv[i + 2]
		elif arg == "-p":
			port = int(sys.argv[i + 2])

	if not addr:
		addr = "127.0.0.1"
	if not port:
		port = 54746

	return (addr, port)

class Cube(socket.socket):

	size   = Dimension(0, 0, 0)
	colors = 3
	fps    = 0

	def __init__(self, server=determineConnection()):
		super(Cube, self).__init__(socket.AF_INET, socket.SOCK_STREAM)
		self.connect(server)

		self.send(b"inf")
		data = self.recv(4 * 3 + 1 + 1, socket.MSG_WAITALL)
		getInt = lambda offset: int.from_bytes(data[offset:offset + 4], byteorder="little")
		self.size = Dimension(
		    getInt(4 * 0),
		    getInt(4 * 1),
		    getInt(4 * 2),
		)
		self.colors = int(data[4 * 3])
		self.fps    = int(data[4 * 3 + 1])

	def swap(self):
		self.send(b"swp")

	def frame(self, data):
		self.send(b"put")
		self.send(data)

	def index(self, x, y, z, rgb=0):
		return (x*self.size.y*self.size.z + y*self.size.z + z) * self.colors + rgb

	def graph2(self, func, send=True, swap=True):
		frame = bytearray(self.length())
		for x in range(0, self.size.x):
			for y in range(0, self.size.y):
				dot = func(x / self.size.x, y / self.size.y)
				i   = self.index(x, y, int(dot[0] * self.size.z - .5))
				for c in range(0, 3):
					frame[i + c] = int(dot[1][c] * 255)
		if send:
			self.frame(frame)
			if swap:
				self.swap()
		return frame

	def graph3(self, func, send=True, swap=True):
		frame = bytearray(self.length())
		for x in range(0, self.size.x):
			for y in range(0, self.size.y):
				for z in range(0, self.size.z):
					vox = func(x / self.size.x, y / self.size.y, z / self.size.z)
					for c in range(0, self.colors):
						frame[self.index(x, y, z, c)] = int(vox[c] * 255)
		if send:
			self.frame(frame)
			if swap:
				self.swap()
		return frame

	def length(self):
		return self.size.x * self.size.y * self.size.z * self.colors
