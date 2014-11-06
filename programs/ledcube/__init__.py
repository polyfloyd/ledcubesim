#
# Copyright (c) 2014 PolyFloyd
#

from collections import namedtuple
import os
import socket
import sys

def determine_connection():
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

Vector = namedtuple('Vector', 'x y z')


class Cube(socket.socket):

	size   = Vector(0, 0, 0)
	colors = 3
	fps    = 0

	def __init__(self, server=determine_connection()):
		super(Cube, self).__init__(socket.AF_INET, socket.SOCK_STREAM)
		self.connect(server)

		self.send(b"inf")
		data = self.recv(4 * 3 + 1 + 1, socket.MSG_WAITALL)
		get_int = lambda offset: int.from_bytes(data[offset:offset + 4], byteorder="little")
		self.size = Vector(
			get_int(4 * 0),
			get_int(4 * 1),
			get_int(4 * 2),
		)
		self.colors = int(data[4 * 3])
		self.fps    = int(data[4 * 3 + 1])

	def set_frame(self, data, swap=True):
		self.send(b"put")
		self.send(data)
		if swap:
			self.swap()

	def make_frame(self):
		return Frame(self.size, self.colors)

	def swap(self):
		self.send(b"swp")


class Frame(bytearray):

	def __init__(self, size, bytes_per_voxel):
		super(Frame, self).__init__(size.x * size.y * size.z * bytes_per_voxel)
		self.size            = size
		self.bytes_per_voxel = bytes_per_voxel

	def set(self, x, y, z, voxel, clip=True):
		x, y, z = int(x), int(y), int(z)
		visible = 0 <= x < self.size.x and 0 <= y < self.size.y and 0 <= z < self.size.z
		if not visible and not clip:
			raise IndexError("(%s, %s, %s) is outside screenspace" % (x, y, z))
		if visible:
			i = (x * self.size.y * self.size.z + y * self.size.z + z) * self.bytes_per_voxel
			for j in range(0, self.bytes_per_voxel):
				self[i + j] = int(voxel[j])

	def setf(self, x, y, z, voxel):
		self.set(x * (self.size.x - 1), y * (self.size.y - 1), z * (self.size.z - 1), voxel)

	def graph2(self, func):
		for x in range(0, self.size.x + 1):
			xn = x / self.size.x
			for y in range(0, self.size.y + 1):
				yn = y / self.size.y
				(zn, vox) = func(xn, yn)
				self.setf(xn, yn, zn, vox)

	def graph3(self, func):
		for x in range(0, self.size.x + 1):
			xn = x / self.size.x
			for y in range(0, self.size.y + 1):
				yn = y / self.size.y
				for z in range(0, self.size.z + 1):
					zn = z / self.size.z
					self.setf(xn, yn, zn, func(xn, yn, zn))
