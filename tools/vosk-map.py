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

import locale, os, sys, json
from   argparse import ArgumentParser

# --- language-mappings   ----------------------------------------------------

from word_map_de import words_de
from word_map_en import words_en

WORDS_MAPS = {
  "de": words_de,
  "en": words_en
  }

# --- application class   ----------------------------------------------------

class App(object):

  # --- constructor   --------------------------------------------------------

  def __init__(self):
    """ constructor """

    parser = self._get_parser()
    parser.parse_args(namespace=self)

  # --- cmdline-parser   -----------------------------------------------------

  def _get_parser(self):
    """ configure cmdline-parser """

    parser = ArgumentParser(add_help=False,description='Vosk Word-Map creator')

    parser.add_argument('-d', '--debug', action='store_true',
      dest='debug', default=False,
      help="force debug-mode")
    parser.add_argument('-q', '--quiet', action='store_true',
      dest='quiet', default=False,
      help="don't print messages")
    parser.add_argument('-h', '--help', action='help',
      help='print this help')

    parser.add_argument('-L', '--language', dest="lang",
                        metavar="language", default="de",
                        choices = WORDS_MAPS.keys(),
       help='language for phrase-mappings ('+", ".join(WORDS_MAPS.keys())+')')

    parser.add_argument('file', nargs=1, metavar='channel-file',
      default=1, help='channel file name')
    return parser

  # --- convert name   -------------------------------------------------------

  def _convert_name(self,name):
    """ convert numbers within name """

    words = WORDS_MAPS[self.lang]

    name = name.lower()

    parts = name.split()
    result = []
    for part in parts:
      # check for number ...
      try:
        part_int = int(part)
        is_int = True
      except:
        is_int = False

      if is_int and part_int in words:
        # ...  and convert to text
        result.append(words[part_int])
      else:
        result.append(part)

    return " ".join(result)

  # --- read channel file   --------------------------------------------------

  def read_channels(self):
    """ read channel file """

    self._channels = []

    f = open(self.file[0],"r")
    self._channels = json.load(f)
    f.close()
    nr=1
    for channel in self._channels:
      channel['nr'] = nr
      nr += 1

  # --- print phrase-map   ---------------------------------------------------

  def print_config(self):
    """ print config for Vosk """

    words = WORDS_MAPS[self.lang]
    config = {
      "model":     "/usr/local/lib/vosk/model",
      "device_id": 1,
      "api_map": {
        words["bhojpur"]:        ["_set_cmd_mode"],
        words["on"]:           ["bhojpur_on"],
        words["off"]:          ["bhojpur_off"],
        words["pause"]:        ["bhojpur_pause"],
        words["resume"]:       ["bhojpur_resume"],
        words["mute on"]:      ["vol_mute_on"],
        words["mute off"]:     ["vol_mute_off"],
        words["volume up"]:    ["vol_up"],
        words["volume down"]:  ["vol_down"],
        words["next"]:         ["bhojpur_play_next"],
        words["previous"]:     ["bhojpur_play_prev"],
        words["record start"]: ["rec_start"],
        words["record stop"]:  ["rec_stop"],
        words["restart"]:      ["sys_restart"],
        words["stop"]:         ["sys_stop"],
        words["reboot"]:       ["sys_reboot"],
        words["shutdown"]:     ["sys_halt"],
        words["halt"]:         ["sys_halt"],
        words["quit"]:         ["_quit"]
        }
      }

    for channel in self._channels:
      # add api by channel number
      nr = channel["nr"]
      value = ["bhojpur_play_channel", "nr=%d" % nr]
      if nr in words:
        key = "%s %s" % (words["channel"], words[nr])
      else:
        key = "%s %d" % (words["channel"], nr)
      config["api_map"][key] = value

      # add api by channel name
      key = self._convert_name(channel["name"])
      config["api_map"][key] = value

    print(json.dumps(config,indent=2,ensure_ascii=False))

# --- main program   ---------------------------------------------------------

if __name__ == '__main__':

  # set local to default from environment
  locale.setlocale(locale.LC_ALL, '')

  # create client-class and parse arguments
  app = App()
  app.read_channels()
  app.print_config()