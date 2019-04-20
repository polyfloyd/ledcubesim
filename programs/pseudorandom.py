#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import math
import random
import time

cube          = ledcube.Cube()
steps         = [math.sin(i / 40 * math.pi / 2) for i in range(0, 40)]
frame_current = cube.make_frame()
frame_target  = cube.make_frame()

while 1:
    frame_source = frame_target
    frame_target = cube.make_frame()
    for i in range(0, len(frame_target)):
        frame_target[i] = random.randint(0, 255)
    for m in steps:
        for j in range(0, len(frame_target)):
            frame_current[j] = int(frame_source[j] * (1 - m) + frame_target[j] * m)
        cube.set_frame(frame_current)
        time.sleep(1 / cube.fps)
