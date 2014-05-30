import os
import socket
import sys

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

	size   = (0, 0, 0)
	colors = 1

	def __init__(self, size, colors, connectTo=determineConnection()):
		super(Cube, self).__init__(socket.AF_INET, socket.SOCK_STREAM)
		self.size   = size
		self.colors = colors
		self.connect(connectTo)

	def swap(self):
		self.send(b"swp")

	def frame(self, data):
		self.send(b"frm")
		self.send(data)

	def xyz(self, xyz, rgb=0):
		return (xyz[0]*self.size[0]*self.size[1] + xyz[1]*self.size[1] + xyz[2]) * self.colors + rgb
