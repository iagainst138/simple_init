#!/usr/bin/env python3

import argparse
import signal
import sys
import time


signals = {
    'int':   signal.SIGINT,
    'hup':   signal.SIGHUP,
    'term':  signal.SIGTERM,
    'winch': signal.SIGWINCH,
}


def write(s):
    sys.stdout.write(f'{s}\n')
    sys.stdout.flush()


def signal_handler(sig, frame):
    write('stopping gracefully...')
    time.sleep(3)
    write('done')
    sys.exit(0)


if __name__ == '__main__':
    a = argparse.ArgumentParser()
    a.add_argument('--signal', choices=signals.keys(), default='int', help='signal to gracefully shutdown')
    a.add_argument('--sleep', default=2, type=float, help='duration to sleep')
    args = a.parse_args()

    signal.signal(signals[args.signal], signal_handler)
    write('*** running ***')

    if args.sleep > 0:
        while True:
            write(f'now {time.time()}')
            time.sleep(args.sleep)

    signal.pause()
