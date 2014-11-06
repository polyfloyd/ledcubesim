#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import ledcube.model as model
import sys

if len(sys.argv) <= 1:
	print('Usage %s <obj file>' % sys.argv[0])
	exit(1)

cube = ledcube.Cube()
mod  = model.WavefrontModel(sys.argv[1], ('x', 'z', 'y'))

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

cube.set_frame(frame)
