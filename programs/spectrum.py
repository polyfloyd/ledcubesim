#! /usr/bin/env python3

#
# Copyright (c) 2014 PolyFloyd
#

import ledcube
import ledcube.audio as audio
import numpy.fft
import os
import time
import math

FIFO_FILE        = '%s/.mpd/mpd.fifo' % os.getenv('HOME')
FIFO_SAMPLE_RATE = 22000
FIFO_SAMPLE_BITS = 16

cube   = ledcube.Cube()
source = audio.PCMSource(FIFO_FILE, FIFO_SAMPLE_RATE, FIFO_SAMPLE_BITS)

# https://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
def highest_pow2(v):
	v -= 1;
	for i in range(0, 6):
		v |= v >> (1 << i)
	return (v + 1) >> 1

while True:
	signal = source.get_signal(1 / cube.fps)
	signal = signal[0:highest_pow2(len(signal))]
	(spectrum, freqs) = source.get_spectrum(signal)
	if len(spectrum) > 64:
		spectrum = spectrum[::len(spectrum) / 64]
	spectrum = [math.sqrt(v) / 2 for v in spectrum]

	frame = cube.make_frame()
	for x in range(0, 8):
		for y in range(0, 8):
			for z in range(0, int(spectrum[x * 8 + y] * cube.size.z + 1)):
				zn = z / cube.size.z
				for bx in range(0, 2):
					for by in range(0, 2):
						top    = .666 - (1 - zn)
						middle = zn * (1 - zn) * 2
						bottom = .666 - zn
						frame.set(x * 2 + bx, y * 2 + by, zn * 15, (
							(top    if top    > 0 else 0) * 255,
							(middle if middle > 0 else 0) * 255,
							(bottom if bottom > 0 else 0) * 255,
						))
	cube.set_frame(frame)
