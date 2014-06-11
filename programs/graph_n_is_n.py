#! /usr/bin/env python3

import array
import ledcube
import math
import os
import time

cube = ledcube.Cube((16, 16, 16), 3)

def graphFunc(x, y, z):
	return (x, y, z)

cube.graph(graphFunc)
