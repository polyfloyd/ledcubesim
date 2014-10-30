#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import ledcube.model as model
import os
import sys

model_file = sys.argv[1] if len(sys.argv) > 1 else os.path.join(os.path.dirname(__file__), 'res/teapot.obj')

cube = ledcube.Cube()
mod  = model.WavefrontModel(model_file, ('x', 'z', 'y'))

frame = ledcube.Frame(cube.size, 3)

smallest_side = min(cube.size) - 1
mod_max       = max(mod.max_x + -mod.min_x, mod.max_y + -mod.min_x, mod.max_y + -mod.min_z)
for face in mod.faces:
	for corner in face:
		frame.set(
			(corner.v.x / mod_max - mod.min_x / mod_max) * smallest_side,
			(corner.v.y / mod_max - mod.min_y / mod_max) * smallest_side,
			(corner.v.z / mod_max - mod.min_z / mod_max) * smallest_side,
			(
				int((corner.n.x + 1) / 2 * 255),
				int((corner.n.y + 1) / 2 * 255),
				int((corner.n.z + 1) / 2 * 255),
			) if corner.n is not None else (128, 128, 255),
		)

cube.frame(frame)
cube.swap()
