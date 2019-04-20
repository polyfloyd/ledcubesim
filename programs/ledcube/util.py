#
# Copyright (c) 2014 PolyFloyd
#

import atexit
import sys
import termios
import threading
import tty

class TTYInput(threading.Thread):

    instances = []
    worker    = None

    def __init__(self, callback):
        super(TTYInput, self).__init__(daemon=True)
        self.callback   = callback
        TTYInput.instances.append(self)
        if TTYInput.worker is None:
            TTYInput.worker = self
            self.start()

    def __del__(self):
        TTYInput.instances.remove(self)

    def run(self):
        old_term = termios.tcgetattr(sys.stdin)
        tty.setraw(sys.stdin.fileno())
        atexit.register(lambda: termios.tcsetattr(sys.stdin, termios.TCSADRAIN, old_term))
        while True:
            ch = sys.stdin.read(1)[0]
            if ord(ch) == 0x03: # SIGINT
                break
            for inst in TTYInput.instances:
                inst.callback(ch)
