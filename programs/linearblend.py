#! /usr/bin/env python3

import ledcube

cube = ledcube.Cube((16, 16, 16), 3)

cube.graph(lambda x, y, z: (x, y, z))
