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

# A simple Bhojpur CLI client for the automated speech recognition software.

DEFAULT_HOST = 'localhost'
DEFAULT_PORT = 8026

import locale, os, sys, json, shlex, threading, signal, readline
from   argparse import ArgumentParser

# --- application imports   --------------------------------------------------

sys.path.append(os.path.join(
  os.path.dirname(sys.argv[0]),"../lib"))

from webspeech import SpeechClient, KeyController, have_vosk

# --- application class   ----------------------------------------------------

class SpeechCli(object):

  # --- constructor   --------------------------------------------------------

  def __init__(self):
    """ constructor """

    parser = self._get_parser()
    parser.parse_args(namespace=self)
    self._cli = SpeechClient(self.host[0],self.port[0],debug=self.debug)

  # --- cmdline-parser   -----------------------------------------------------

  def _get_parser(self):
    """ configure cmdline-parser """

    parser = ArgumentParser(prog="webspeech_cli",add_help=False,description='Bhojpur Speech is a voice processing engine with automated speech recognition')

    parser.add_argument('-H', '--host', nargs=1,
      metavar='host', default=[DEFAULT_HOST],
      dest='host',
      help='host-mask')
    parser.add_argument('-P', '--port', nargs=1,
      metavar='port', default=[DEFAULT_PORT],
      dest='port',
      help='port the server is listening on (default: %d)' % DEFAULT_PORT)

    parser.add_argument('-i', '--interactive', action='store_true',
      dest='interactive', default=False,
      help="interactive mode (read APIs from interactive shell)")
    parser.add_argument('-k', '--keyboard', action='store_true',
      dest='keyboard', default=False,
      help="key-control mode (maps keys to APIs)")
    parser.add_argument('-v', '--voice', action='store_true',
      dest='voice', default=False,
      help="voice-control mode (needs Vosk and a microphone)")

    parser.add_argument('-e', '--events', action='store_true',
      dest='events', default=False,
      help="start event-processing")
    parser.add_argument('-o', '--on', action='store_true',
      dest='on', default=False,
      help="turn speech on")

    parser.add_argument('-d', '--debug', action='store_true',
      dest='debug', default=False,
      help="force debug-mode")
    parser.add_argument('-q', '--quiet', action='store_true',
      dest='quiet', default=False,
      help="don't print messages")
    parser.add_argument('-h', '--help', action='help',
      help='print this help')

    parser.add_argument('api', nargs='?', metavar='api',
      default=0, help='api name')
    parser.add_argument('args', nargs='*', metavar='name=value',
      help='api arguments')
    return parser

  # --- return stop-event   --------------------------------------------------

  def get_stop_event(self):
    """ return stop event """

    return self._cli.get_stop_event()

  # --- setup signal handler   ------------------------------------------------

  def signal_handler(self,_signo, _stack_frame):
    """ signal-handler for clean shutdown """

    self.msg("webspeech_cli: received signal, stopping program ...")
    self.close()

  # --- close connection   ---------------------------------------------------

  def close(self):
    """ close connection """

    try:
      self._cli.close()
    except:
      pass

  # --- print message   ------------------------------------------------------

  def msg(self,text,force=False):
    """ print message """

    self._cli.msg(text,force)

  # --- dump output of API   -------------------------------------------------

  def print_response(self,response):
    """ write response to stderr and stdout """

    if self.quiet:
      return
    elif self.debug:
      sys.stderr.write("%d %s\n" % (response[0],response[1]))
      sys.stderr.flush()
    try:
      obj = json.loads(response[2])
      print(json.dumps(obj,indent=2,sort_keys=True))
    except:
      if response[2]:
        print("response: " + response[2])

  # --- print event   --------------------------------------------------------

  def handle_event(self,event):
    """ print event (depending on mode) """

    raw = self.debug or (not self.interactive and not self.keyboard)
    if raw:
      print(json.dumps(json.loads(event.data),indent=2,sort_keys=True))
    elif not self.quiet:
      ev_data = json.loads(event.data)
      if ev_data['type'] != 'keep_alive':
        print(ev_data['text'])

  # --- process single api   -------------------------------------------------

  def process_api(self,api,args=[],sync=True):
    """ process a single API-call """

    # execute api
    if api == "get_events":
      if sync:
        events = self._cli.get_events()
        for event in events:
          self.handle_event(event)
      else:
        self._cli.start_event_processing(callback=self.handle_event)
    else:
      # use synchronous calls for all other events
      params = {}
      for a in args:
        [key,value] = a.split("=",1)
        params[key] = value

      resp = self._cli.exec(api,params=params)
      self.print_response(resp)

  # --- process stdin   ------------------------------------------------------

  def process_stdin(self):
    """ check for stdin and process commands """

    # test for stdin
    try:
      _ = os.tcgetpgrp(sys.stdin.fileno())
      self.msg("webspeech_cli: no stdin ...")
      return
    except:
      pass

    # read commands from stdin
    for line in sys.stdin:
      line = line.rstrip()
      if not len(line):
        break
      cmd  = shlex.split(line)
      self.process_api(cmd[0],cmd[1:],sync=False)

  # --- completer for readline   ---------------------------------------------

  def completer(self,text,state):
    """ implement completer """

    #self.msg("SpeechCli: completer(%s,%d)" % (text,state))
    if state == 0:
      # buffer list of hits
      self._completions = [api for api in self._api_list
                           if api.startswith(text)]

    if state < len(self._completions):
      #self.msg("SpeechCli: returning %s" % self._completions[state])
      return self._completions[state]
    else:
      return True

  # --- run application   ----------------------------------------------------

  def run(self):
    """ run application """

    # setup signal-handler
    signal.signal(signal.SIGTERM, self.signal_handler)
    signal.signal(signal.SIGINT,  self.signal_handler)

    # process special options
    if self.events:
      self.process_api("get_events",sync=False)
    if self.on:
      self.process_api("bhojpur_on")

    # process cmdline
    if self.api:
      self.process_api(self.api,self.args,
                      sync=not (self.keyboard or self.interactive))

    # process stdin (if available)
    self.process_stdin()

    # process keyboard / voice / interactive input
    if self.keyboard or self.voice:
      if self.keyboard:
        ctrl = KeyController(self.get_stop_event(),self.debug)
      else:
        if have_vosk:
          from webspeech import VoskController
          ctrl = VoskController(self.get_stop_event(),self.debug)
        else:
          self.msg("webspeech_cli: voice-support not installed",force=True)
          return
      for api in ctrl.api_from_key():
        if api[0] == "_quit":
          return
        elif api[0] == "_help":
          ctrl.print_mapping()
        else:
          self.process_api(api[0],api[1:],sync=False)
        if api[0] == 'sys_stop':
          return
    elif self.interactive:
      self._api_list = self._cli.get_api_list()
      readline.set_completer(lambda text,state: self.completer(text,state))
      readline.parse_and_bind("tab: complete")
      while True:
        line = input("webspeech > ").strip()
        if not len(line):
          continue
        elif line in ['q','Q','quit','Quit']:
          return
        api  = shlex.split(line)
        self.process_api(api[0],api[1:],sync=False)
        if api[0] == 'sys_stop':
          return

# --- main program   ---------------------------------------------------------

if __name__ == '__main__':

  # set local to default from environment
  locale.setlocale(locale.LC_ALL, '')

  print("Bhojpur Speech client engine (online)")
  print("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
  print("All rights reserved.\n")

  # create client-class and parse arguments
  app = SpeechCli()
  app.run()
  app.close()