#
# Copyright (c) 2014 PolyFloyd
#

import numpy.fft
import os

class Source:

    def get_spectrum(self, signal):
        signal   = numpy.array([(s + 1) / 2 for s in signal], dtype=float)
        spectrum = numpy.abs(numpy.fft.rfft(signal))
        freqs    = numpy.fft.fftfreq(spectrum.size, 1 / self.get_sample_rate())
        spectrum = spectrum[1:]
        return (spectrum, freqs)

    def get_input(self):
        return self.input

    def set_input(self, input):
        if type(input) == str:
            self.input = os.fdopen(os.open(input, os.O_RDONLY), 'rb')
        else:
            self.input = input

    def get_signal(self, seconds):
        return [self.get_next_sample() for i in range(0, int(self.get_sample_rate() * seconds))]

    def get_next_sample(self):
        pass # virtual

    def get_sample_rate(self):
        pass # virtual


class PCMSource(Source):

    def __init__(self, input_file, sample_rate, sample_bits, sample_endianness='little', sample_sign='signed'):
        assert(sample_endianness == 'little' or sample_endianness == 'big')
        assert(sample_sign       == 'signed' or sample_sign       == 'unsigned')
        self.set_input(input_file)
        self.sample_rate       = sample_rate
        self.sample_bits       = sample_bits
        self.sample_endianness = sample_endianness
        self.sample_sign       = sample_sign

    def sample_from_raw_data(self, raw_data):
        intval = int.from_bytes(raw_data, self.sample_endianness, signed=self.sample_sign == 'signed')
        return intval / (2 ** (len(raw_data) * 8 - 1))

    def get_next_sample(self):
        return self.sample_from_raw_data(self.get_input().read(self.sample_bits // 8))

    def get_sample_rate(self):
        return self.sample_rate
