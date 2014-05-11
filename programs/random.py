#! /usr/bin/env python3

import ledcube
import os
import socket
import time

sock = ledcube.connect()

while 1:
	sock.send(b"frm")
	sock.send(os.urandom(16**3 * 3))
	sock.send(b"swp")
	time.sleep(.5)
