#
# Copyright (c) 2014 PolyFloyd
#

from collections import namedtuple
import re

PARSE_OBJ_V  = re.compile('^v\
	\s(?P<x>[-\d\.e]+)\
	\s(?P<y>[-\d\.e]+)\
	\s(?P<z>[-\d\.e]+)\
$' , re.M | re.X)
PARSE_OBJ_VT = re.compile('^vt\
	\s(?P<u>[-\d\.e]+)\
	\s(?P<v>[-\d\.e]+)\
$', re.M | re.X)
PARSE_OBJ_VN = re.compile('^vn\
	\s(?P<x>[-\d\.e]+)\
	\s(?P<y>[-\d\.e]+)\
	\s(?P<z>[-\d\.e]+)\
$', re.M | re.X)
PARSE_OBJ_F  = re.compile('^f\
	\s(?P<va>\d+)\/(?P<vta>\d*)\/(?P<vna>\d*)\
	\s(?P<vb>\d+)\/(?P<vtb>\d*)\/(?P<vnb>\d*)\
	\s(?P<vc>\d+)\/(?P<vtc>\d*)\/(?P<vnc>\d*)\
$', re.M | re.X)

Vector     = namedtuple('Vector',     'x y z')
TexCoord   = namedtuple('TexCoord',   'u v')
FaceVertex = namedtuple('FaceVertex', 'v t n')

class WavefrontModel:

	def __init__(self, obj_file, axismap=('x', 'y', 'z')):
		obj = open(obj_file).read()

		firstvertex_flag = True
		def parse_v(d):
			nonlocal firstvertex_flag
			x, y, z = float(d[axismap[0]]), float(d[axismap[1]]), float(d[axismap[2]])
			if firstvertex_flag:
				firstvertex_flag = False
				self.min_x = self.max_x = x
				self.min_y = self.max_y = y
				self.min_z = self.max_z = z
			self.min_x = min(self.min_x, x)
			self.min_y = min(self.min_y, y)
			self.min_z = min(self.min_z, z)
			self.max_x = max(self.max_x, x)
			self.max_y = max(self.max_y, y)
			self.max_z = max(self.max_z, z)
			return Vector(x, y, z)
		vertices = [parse_v(m.groupdict()) for m in PARSE_OBJ_V.finditer(obj)]

		def parse_vt(d):
			return TexCoord(float(d['u']), float(d['v']))
		texcoords = [parse_vt(m.groupdict()) for m in PARSE_OBJ_VT.finditer(obj)]

		def parse_vn(d):
			return Vector(float(d['x']), float(d['y']), float(d['z']))
		normals = [parse_vn(m.groupdict()) for m in PARSE_OBJ_VN.finditer(obj)]

		def parse_f(d):
			return [
				FaceVertex(
					vertices[int(d['va'])   - 1],
					texcoords[int(d['vta']) - 1] if d['vta'] != '' else None,
					normals[int(d['vna'])   - 1] if d['vna'] != '' else None,
				),
				FaceVertex(
					vertices[int(d['vb'])   - 1],
					texcoords[int(d['vtb']) - 1] if d['vtb'] != '' else None,
					normals[int(d['vnb'])   - 1] if d['vnb'] != '' else None,
				),
				FaceVertex(
					vertices[int(d['vc'])   - 1],
					texcoords[int(d['vtc']) - 1] if d['vtc'] != '' else None,
					normals[int(d['vnc'])   - 1] if d['vnc'] != '' else None,
				),
			]
		self.faces = [parse_f(m.groupdict()) for m in PARSE_OBJ_F.finditer(obj)]
