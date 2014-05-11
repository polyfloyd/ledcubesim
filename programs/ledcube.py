import os
import socket
import sys

def connect():
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
	s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	s.connect((addr, port))
	return s
