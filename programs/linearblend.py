#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube

cube = ledcube.Cube()

cube.graph(lambda x, y, z: (x, y, z))
