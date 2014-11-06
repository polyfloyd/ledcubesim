#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube

cube = ledcube.Cube()

frame = cube.make_frame()
frame.graph3(lambda x, y, z: (x * 255, y * 255, z * 255))
cube.set_frame(frame)
