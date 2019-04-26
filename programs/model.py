#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import ledcube.model as model
import sys

if len(sys.argv) <= 1:
    print('Usage: %s <obj file>' % sys.argv[0])
    exit(1)

cube = ledcube.Cube()
mod  = model.WavefrontModel(sys.argv[1], ('x', 'z', 'y'))

frame = ledcube.Frame(cube.size)

smallest_side = min(cube.size) - 1
mod_max       = max(mod.max.x - mod.min.x, mod.max.y - mod.min.y, mod.max.z - mod.min.z)
for face in mod.faces:
    for vert in face:
        frame.set(
            ((vert.v.x - mod.min.x) / mod_max) * smallest_side,
            ((vert.v.y - mod.min.y) / mod_max) * smallest_side,
            ((vert.v.z - mod.min.z) / mod_max) * smallest_side,
            (
                int((vert.n.x + 1) / 2 * 255),
                int((vert.n.y + 1) / 2 * 255),
                int((vert.n.z + 1) / 2 * 255),
            ) if vert.n is not None else (128, 128, 255),
        )

cube.set_frame(frame)
