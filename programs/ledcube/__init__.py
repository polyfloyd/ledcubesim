#
# Copyright (c) 2014 PolyFloyd
#

from collections import namedtuple
import socket
import sys

Vector = namedtuple('Vector', 'x y z')

class Cube(socket.socket):
    def __init__(self):
        self.size   = Vector(16, 16, 16)
        self.fps    = 60

    def set_frame(self, data, swap=True):
        self._data = data
        if swap:
            self.swap()

    def make_frame(self):
        return Frame(self.size)

    def swap(self):
        sys.stdout.buffer.write(self._data)


class Frame(bytearray):
    def __init__(self, size):
        super(Frame, self).__init__(size.x * size.y * size.z * 3)
        self.size = size

    def index(self, x, y, z):
        x, y, z = int(x), int(y), int(z)
        if 0 <= x < self.size.x and 0 <= y < self.size.y and 0 <= z < self.size.z:
            return (x * self.size.y * self.size.z + y * self.size.z + z) * 3
        return -1

    def get(self, x, y, z):
        i = self.index(x, y, z)
        if i == -1:
            raise IndexError("(%s, %s, %s) is outside screenspace" % (x, y, z))
        return (self[i], self[i + 1], self[i + 2])

    def set(self, x, y, z, voxel, clip=True):
        i = self.index(x, y, z)
        if i != -1:
            for j in range(0, 3):
                self[i + j] = int(voxel[j])
        elif not clip:
            raise IndexError("(%s, %s, %s) is outside screenspace" % (x, y, z))

    def indexf(self, x, y, z):
        return self.indexf(x * (self.size.x - 1), y * (self.size.y - 1), z * (self.size.z - 1))

    def getf(self, x, y, z):
        return self.get(x * (self.size.x - 1), y * (self.size.y - 1), z * (self.size.z - 1))

    def setf(self, x, y, z, voxel):
        self.set(x * (self.size.x - 1), y * (self.size.y - 1), z * (self.size.z - 1), voxel)

    def graph2(self, func):
        for x in range(0, self.size.x + 1):
            xn = x / self.size.x
            for y in range(0, self.size.y + 1):
                yn = y / self.size.y
                (zn, vox) = func(xn, yn)
                self.setf(xn, yn, zn, vox)

    def graph3(self, func):
        for x in range(0, self.size.x + 1):
            xn = x / self.size.x
            for y in range(0, self.size.y + 1):
                yn = y / self.size.y
                for z in range(0, self.size.z + 1):
                    zn = z / self.size.z
                    self.setf(xn, yn, zn, func(xn, yn, zn))
