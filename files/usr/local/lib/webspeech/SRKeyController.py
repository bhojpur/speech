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

# Bhojpur Speech: implementation of class KeyController
#
# The class KeyController maps key-events to api-calls

import os, sys, evdev, select, tty, termios
from webspeech import Base

class KeyController(Base):
  """ map key-events to api-calls """

  TIMEOUT = 5             # check stop-event every x seconds

  # --- builtin-keymaps   ----------------------------------------------------

  KEYMAP_SPEECH_EVENT = {
    "KEY_1":     ["bhojpur_play_channel", "nr=1"],
    "KEY_2":     ["bhojpur_play_channel", "nr=2"],
    "KEY_3":     ["bhojpur_play_channel", "nr=3"],
    "KEY_4":     ["bhojpur_play_channel", "nr=4"],
    "KEY_5":     ["bhojpur_play_channel", "nr=5"],
    "KEY_6":     ["bhojpur_play_channel", "nr=6"],
    "KEY_7":     ["bhojpur_play_channel", "nr=7"],
    "KEY_8":     ["bhojpur_play_channel", "nr=8"],
    "KEY_9":     ["bhojpur_play_channel", "nr=9"],
    "KEY_0":     ["bhojpur_play_channel", "nr=10"],
    "KEY_D":     ["player_play_dir"],                   # TBD
    "KEY_F":     ["player_play_file"],                  # TBD
    "KEY_H":     ["_help"],                             # TBD
    "KEY_I":     ["bhojpur_state"],                     # TBD
    "KEY_O":     ["bhojpur_toggle"],
    "KEY_P":     ["player_mode_toggle"],                # TBD
    "KEY_Q":     ["_quit"],
    "KEY_R":     ["rec_toggle"],
    "KEY_S":     ["sys_stop"],
    "KEY_L":     ["bhojpur_get_channels"],
    "KEY_M":     ["vol_mute_toggle"],
    "KEY_LEFT":  ["bhojpur_play_prev"],
    "KEY_RIGHT": ["bhojpur_play_next"],
    "KEY_UP":    ["vol_up"],
    "KEY_DOWN":  ["vol_down"],
    "KEY_ENTER": ["player_select"]                      # TBD
    }

  KEYMAP_SPEECH_TERM = {
    "31":     ["bhojpur_play_channel", "nr=1"],
    "32":     ["bhojpur_play_channel", "nr=2"],
    "33":     ["bhojpur_play_channel", "nr=3"],
    "34":     ["bhojpur_play_channel", "nr=4"],
    "35":     ["bhojpur_play_channel", "nr=5"],
    "36":     ["bhojpur_play_channel", "nr=6"],
    "37":     ["bhojpur_play_channel", "nr=7"],
    "38":     ["bhojpur_play_channel", "nr=8"],
    "39":     ["bhojpur_play_channel", "nr=9"],
    "30":     ["bhojpur_play_channel", "nr=10"],
    "64":     ["player_play_dir"],                   # TBD
    "66":     ["player_play_file"],                  # TBD
    "69":     ["bhojpur_state"],                     # TBD
    "6f":     ["bhojpur_toggle"],
    "70":     ["player_mode_toggle"],                # TBD
    "68":     ["_help"],
    "71":     ["_quit"],
    "72":     ["rec_toggle"],
    "73":     ["sys_stop"],
    "6c":     ["bhojpur_get_channels"],
    "6d":     ["vol_mute_toggle"],
    "1b5b44": ["bhojpur_play_prev"],
    "1b5b43": ["bhojpur_play_next"],
    "1b5b41": ["vol_up"],
    "1b5b42": ["vol_down"],
    "0a":     ["player_select"]                      # TBD
    }

  KEY_SPECIAL = ['KEY_LEFTCTRL','KEY_LEFTALT','KEY_LEFTSHIFT',
                 'KEY_RIGHTCTRL','KEY_RIGHTALT','KEY_RIGHTSHIFT']

  # --- constructor   --------------------------------------------------------

  def __init__(self,stop,debug=False):
    """ constructor """

    self._stop    = stop
    self.debug    = debug

    # test for terminal
    self._have_term = False
    try:
      _ = os.tcgetpgrp(sys.stdin.fileno())
      self.msg("KeyController: have terminal")
      self._have_term = True
      self._kmap    = KeyController.KEYMAP_SPEECH_TERM
    except:
      self._kmap    = KeyController.KEYMAP_SPEECH_EVENT

  # --- yield api from key-event   -------------------------------------------

  def _api_from_key_event(self):
    """ monitor key-events and yield mapped API-name """

    devices = [evdev.InputDevice(fn) for fn in evdev.list_devices()]
    if not len(devices):
      self.msg("no input devices available")
      return
    devices = {dev.fd: dev for dev in devices}

    special = 0                    # to ignore combinations with special keys
    while True:
      fds, _1, _2 = select.select(devices,[],[],KeyController.TIMEOUT)
      if self._stop.is_set():
        break
      elif not len(fds):
        # timeout condition, try again
        continue
      for fd in fds:
        for event in devices[fd].read():
          event = evdev.util.categorize(event)
          if not isinstance(event, evdev.events.KeyEvent):
            continue
          self.msg("KeyController: processing %s (%d)" %
                   (event.keycode,event.keystate))
          if event.keystate == event.key_down:
            if event.keycode in KeyController.KEY_SPECIAL:
              special += 1
              continue
            elif special > 0:
              self.msg("KeyController: ignoring %s" % event.keycode)
              continue
            if event.keycode in self._kmap:
              # key is mapped, yield api-name
              self.msg("KeyController: mapping %s to %s" %
                       (event.keycode,self._kmap[event.keycode]))
              yield self._kmap[event.keycode]
            else:
              # key is not mapped, ignore
              self.msg("KeyController: ignoring %s" % event.keycode)
          elif event.keystate == event.key_up:
            if event.keycode in KeyController.KEY_SPECIAL:
              special = max(0,special-1)

  # --- yield api from key-event   -------------------------------------------

  def _api_from_term(self):
    """ monitor chars from terminal and yield mapped API-name """

    old_settings = termios.tcgetattr(sys.stdin)
    tty.setcbreak(sys.stdin.fileno())
    devices = [sys.stdin.fileno()]
    try:
      while True:
        fds, _1, _2 = select.select(devices,[],[],KeyController.TIMEOUT)
        if self._stop.is_set():
          break
        elif not len(fds):
          # timeout condition, try again
          continue

        keycode = os.read(sys.stdin.fileno(), 3).hex()
        self.msg("KeyController: processing %s" % keycode)
        if keycode in self._kmap:
          # key is mapped, yield api-name
          self.msg("KeyController: mapping %s to %s" %
                   (keycode,self._kmap[keycode]))
          yield self._kmap[keycode]
        else:
          # key is not mapped, ignore
          self.msg("KeyController: ignoring %s" % keycode)
    finally:
      termios.tcsetattr(sys.stdin, termios.TCSADRAIN, old_settings)

  # --- yield api from key-event   -------------------------------------------

  def api_from_key(self):
    """ monitor key-events and yield mapped API-name """

    if self._have_term:
      return self._api_from_term()
    else:
      return self._api_from_key_event()

  # --- print key-mapping   --------------------------------------------------

  def print_mapping(self):
    """ print key-mapping """

    print("key-mapping:")
    for key,value in KeyController.KEYMAP_SPEECH_EVENT.items():
      if len(value) > 1:
        print("%9s -> %s %s" % (key,*value))
      else:
        print("%9s -> %s" % (key,value[0]))