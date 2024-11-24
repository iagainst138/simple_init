#!/usr/bin/env python3

import signal
import sys
import time


def write(s):
    sys.stdout.write(f'{s}\n')
    sys.stdout.flush()


def signal_handler(sig, frame):
    write('stopping gracefully...')
    time.sleep(3)
    write('done')
    sys.exit(0)


if __name__ == '__main__':
    signal.signal(signal.SIGINT, signal_handler)
    write('*** running ***')
    signal.pause()
