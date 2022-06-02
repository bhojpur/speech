#!/usr/bin/python3
# -*- coding: utf-8 -*-

# Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

# Main application program for the automated speech recognition software.
#
# This program starts the application either in synchronous mode or in
# server mode. The latter is usually done from a systemd-service.
# Synchronous mode is for listing channels, direct recording and direct
# playing. Note that direct playing does not allow any interaction, so
# this feature is mainly useful for development and debugging.

import locale, os, sys, signal, queue, threading
from   argparse import ArgumentParser

# --- application imports   --------------------------------------------------

sys.path.append(os.path.join(
  os.path.dirname(sys.argv[0]),"../lib"))

from bhojpur import *

# --- helper class for options   --------------------------------------------

class Options(object):
  pass

# --- cmdline-parser   ------------------------------------------------------

def get_parser():
  """ configure cmdline-parser """

  parser = ArgumentParser(prog="webspeech",add_help=False,description='Bhojpur Speech is a voice processing engine with automated speech recognition')

  parser.add_argument('-p', '--play', action='store_true',
    dest='do_play', default=False,
    help="play speech/file (direct, no web-interface, needs channel/file as argument)")

  parser.add_argument('-l', '--list', action='store_true',
    dest='do_list', default=False,
    help="display speech-channels")

  parser.add_argument('-r', '--record', action='store_true',
    dest='do_record', default=False,
    help="record speech (direct, no web interface, needs channel as argument)")
  parser.add_argument('-t', '--tdir', nargs=1,
    metavar='target directory', default=None,
    dest='target_dir',
    help='target directory for recordings')

  parser.add_argument('-d', '--debug', action='store_true',
    dest='debug', default=False,
    help="force debug-mode (overrides config-file)")
  parser.add_argument('-q', '--quiet', action='store_true',
    dest='quiet', default=False,
    help="don't print messages")
  parser.add_argument('-h', '--help', action='help',
    help='print this help')

  parser.add_argument('channel', nargs='?', metavar='channel',
    default=0, help='channel number/filename')
  parser.add_argument('duration', nargs='?', metavar='duration',
    default=0, help='duration of recording')
  return parser

# --- validate and fix options   ---------------------------------------------

def check_options(options):
  """ validate and fix options """

  # record needs a channel number
  if options.do_record and not options.channel:
    print("[ERROR] record-option (-r) needs channel nummber as argument")
    sys.exit(3)

# --- process events   -------------------------------------------------------

def process_events(app,options,queue):
  while True:
    ev = queue.get()
    if ev:
      if not options.quiet and not ev['type'] == 'keep_alive':
        print(ev['text'])
      queue.task_done()
      if ev['type'] == 'eof' and options.do_play:
        break
      if ev['type'] == 'sys':
        break
    else:
      break
  app.msg("Bhojpur Speech: server engine finished processing events")
  try:
    os.kill(os.getpid(), signal.SIGTERM)
  except:
    pass

# --- main program   ----------------------------------------------------------

if __name__ == '__main__':

  # set local to default from environment
  locale.setlocale(locale.LC_ALL, '')

  # parse commandline-arguments
  opt_parser     = get_parser()
  options        = opt_parser.parse_args(namespace=Options)
  options.pgm_dir = os.path.dirname(os.path.abspath(__file__))
  check_options(options)

  print("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
  print("All rights reserved.\n")

  app = BhojpurSpeech(options)

  # setup signal-handler
  signal.signal(signal.SIGTERM, app.signal_handler)
  signal.signal(signal.SIGINT,  app.signal_handler)

  if options.do_list:
    if not options.quiet:
      app.msg("Bhojpur Speech server engine (online v%s)" % app.api.get_version(),force=True)
    channels = app.api.bhojpur_get_channels()
    PRINT_CHANNEL_FMT="{0:2d}: {1}"
    for channel in channels:
      print(PRINT_CHANNEL_FMT.format(channel['nr'],channel['name']))
  else:
    ev_queue = app.api._add_consumer("main")
    threading.Thread(target=process_events,args=(app,options,ev_queue)).start()
    if options.do_record:
      app.api.rec_start(nr=int(options.channel),sync=True)
      app.cleanup()
    elif options.do_play:
      try:
        nr = int(options.channel)
        app.api.bhojpur_play_channel(nr)
      except ValueError:
        app.api.player_play_file(options.channel) # assume argument is a filename
      signal.pause()
    else:
      app.run()
      signal.pause()
    app.api._del_consumer("main")
  sys.exit(0)