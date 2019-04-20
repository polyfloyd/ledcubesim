#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

from ledcube.model import Vector
import ledcube
import ledcube.util as util
import random
import time

SPEED        = 0.05
CELL_SIZE    = 2
COLOR_BODY   = (255, 255, 0)
COLOR_HEAD   = (0,   255, 0)
COLOR_TARGET = (255, 0,   0)

cube = ledcube.Cube()
grid = ledcube.Vector(
    cube.size.x // CELL_SIZE,
    cube.size.y // CELL_SIZE,
    cube.size.z // CELL_SIZE,
)
score  = 0
dir    = Vector(0, 0, 0)
target = Vector(0, 0, 0)
body   = [Vector(grid.x  // 2, grid.y // 2, grid.z // 2)]
t      = 0

def input_handler(char):
    global dir
    if char in 'da':
        dir = Vector({
            'd':  1,
            'a': -1,
        }[char], 0, 0)
    if char in 'sw':
        dir = Vector(0, {
            's':  1,
            'w': -1,
        }[char], 0)
    if char in 'eq':
        dir = Vector(0, 0, {
            'e':  1,
            'q': -1,
        }[char])

util.TTYInput(input_handler)
while util.TTYInput.worker.is_alive():
    t += SPEED
    if t >= 1:
        t = 0
        body.insert(0, Vector(
            (body[0].x + dir.x) % grid.x,
            (body[0].y + dir.y) % grid.y,
            (body[0].z + dir.z) % grid.z,
        ))
        if body.count(body[0]) == 2 and dir != Vector(0, 0, 0):
            print('\rGame Over! Your score was %s\r' % score)
            break
        if target == body[0]:
            target = Vector(random.randint(0, grid.x - 1), random.randint(0, grid.y - 1), random.randint(0, grid.z - 1))
            score += 10
        else:
            body.pop()

    frame = cube.make_frame()
    def set_cell(x, y, z, color):
        for cx in range(0, CELL_SIZE):
            for cy in range(0, CELL_SIZE):
                for cz in range(0, CELL_SIZE):
                    frame.set(x * CELL_SIZE + cx, y * CELL_SIZE + cy, z * CELL_SIZE + cz, color)

    set_cell(target.x, target.y, target.z, COLOR_TARGET)
    for part in body:
        set_cell(part.x, part.y, part.z, COLOR_BODY)
    set_cell(body[0].x, body[0].y, body[0].z, COLOR_HEAD)

    cube.set_frame(frame)
    time.sleep(1 / cube.fps)
