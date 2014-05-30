#! /usr/bin/env python3

import ledcube
import os
import time

cube = ledcube.Cube()

while 1:
	cube.frame(os.urandom(16**3 * 3))
	cube.swap()
	time.sleep(.5)
