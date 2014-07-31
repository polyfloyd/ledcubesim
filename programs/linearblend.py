#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube

cube = ledcube.Cube((16, 16, 16), 3)

cube.graph(lambda x, y, z: (x, y, z))
