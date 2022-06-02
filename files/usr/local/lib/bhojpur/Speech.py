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

# Bhojpur Speech: implementation of class Speech
#
# The class Speech implements the core functionality of the Bhojpur Speech.

import os, time, datetime, shlex, json
import queue, collections
import threading, signal, subprocess, traceback

from bhojpur import *

class Speech(Base):
  """ Speech-controller """

  def __init__(self,app):
    """ initialization """

    self._app          = app
    self._api          = app.api
    self.debug         = app.debug
    self._backend      = app.backend

    self._channel_nr   = 0                  # current channel number
    self._last_channel = 0                  # last active channel number
    self.stop_event    = app.stop_event
    self.read_config()
    self.register_apis()
    self.read_channels()

  # --- read configuration   --------------------------------------------------

  def read_config(self):
    """ read configuration from config-file """

    # section [GLOBAL]
    default_path        = "/etc/webspeech.channels"
    self._channel_file  = self.get_value(self._app.parser,"GLOBAL","channel_file",
                                         default_path)

    # section [WEB]
    default_web_root = os.path.realpath(
      os.path.join(self._app.options.pgm_dir,"..","lib","bhojpur","web"))
    self._web_root  = self.get_value(self._app.parser,"WEB","web_root",
                                         default_web_root)

  # --- register APIs   ------------------------------------------------------

  def register_apis(self):
    """ register Bhojpur Speech API functions """

    self._api.bhojpur_on             = self.bhojpur_on
    self._api.bhojpur_off            = self.bhojpur_off
    self._api.bhojpur_pause          = self.bhojpur_pause
    self._api.bhojpur_resume         = self.bhojpur_resume
    self._api.bhojpur_toggle         = self.bhojpur_toggle
    self._api.bhojpur_get_channels   = self.bhojpur_get_channels
    self._api.bhojpur_get_channel    = self.bhojpur_get_channel
    self._api.bhojpur_play_channel   = self.bhojpur_play_channel
    self._api.bhojpur_play_next      = self.bhojpur_play_next
    self._api.bhojpur_play_prev      = self.bhojpur_play_prev

  # --- return persistent state of this class   -------------------------------

  def get_persistent_state(self):
    """ return persistent state (overrides SRBase.get_pesistent_state()) """
    return {
      'channel_nr': self._last_channel
      }

  # --- restore persistent state of this class   ------------------------------

  def set_persistent_state(self,state_map):
    """ restore persistent state (overrides SRBase.set_pesistent_state()) """

    self.msg("Speech: restoring persistent state")
    if 'channel_nr' in state_map:
      self._last_channel = state_map['channel_nr']

    self._api.update_state(section="radio",key="channel_nr",
                           value=self._last_channel,publish=False)

  # --- read channels   -------------------------------------------------------

  def read_channels(self):
    """ read channels into a list """

    self._channels = []
    try:
      self.msg("Speech: Loading channels from %s" % self._channel_file)
      f = open(self._channel_file,"r")
      self._channels = json.load(f)
      f.close()
      nr=1
      for channel in self._channels:
        channel['nr'] = nr
        logo_path = os.path.join(self._web_root,"images",channel['logo'])
        if os.path.exists(logo_path):
          channel['logo'] = os.path.join("images",channel['logo'])
        else:
          channel['logo'] = None
        nr += 1
    except:
      self.msg("Speech: Loading channels failed")
      if self.debug:
        traceback.print_exc()

  # --- get channel info   ----------------------------------------------------

  def bhojpur_get_channel(self,nr=0):
    """ return info-dict {name,url,logo} for channel nr """

    try:
      nr = int(nr)
    except:
      nr = 0
    if nr == 0:
      if self._last_channel == 0:
        nr = 1
      else:
        nr = self._last_channel

    return dict(self._channels[nr-1])

  # --- return channel-list   ------------------------------------------------

  def bhojpur_get_channels(self):
    """ return complete channel-list """

    return [dict(c) for c in self._channels]

  # --- play given channel   --------------------------------------------------

  def bhojpur_play_channel(self,nr=0):
    """ switch to given channel """

    channel = self.bhojpur_get_channel(int(nr))
    nr      = channel['nr']
    self.msg("Speech: start playing channel %d (%s)" % (nr,channel['name']))

    # check if we have to do anything
    if self._backend.play(channel['url']):
      self._api.update_state(section="radio",key="channel_nr",
                             value=channel,publish=False)
      self._api._push_event({'type': 'bhojpur_play_channel', 'value': channel})
      self._channel_nr   = nr
      self._last_channel = self._channel_nr
    else:
      self.msg("Speech: already on channel %d" % nr)
      # theoretically we could also have lost our backend
    return channel

  # --- switch to next channel   ----------------------------------------------

  def bhojpur_play_next(self):
    """ switch to next channel """

    self.msg("Speech: switch to next channel")
    if self._channel_nr == 0:
      return self.bhojpur_play_channel(0)
    elif self._channel_nr == len(self._channels):
      return self.bhojpur_play_channel(1)
    else:
      return self.bhojpur_play_channel(1+self._channel_nr)

  # --- switch to previous channel   ------------------------------------------

  def bhojpur_play_prev(self):
    """ switch to previous channel """

    self.msg("Speech: switch to previous channel")
    if self._channel_nr == 0:
      return self.bhojpur_play_channel(0)
    if self._channel_nr == 1:
      return self.bhojpur_play_channel(len(self._channels))
    else:
      return self.bhojpur_play_channel(self._channel_nr-1)

  # --- turn speech off   -----------------------------------------------------

  def bhojpur_off(self):
    """ turn speech off """

    self.msg("Speech: turning speech off")
    self._channel_nr = 0
    self._backend.stop()

  # --- turn speech on   ------------------------------------------------------

  def bhojpur_on(self):
    """ turn speech on """

    if self._channel_nr == 0:
      self.msg("Speech: turning speech on")
      self.bhojpur_play_channel()
    else:
      self.msg("Speech: ignoring command, speech already on")

  # --- pause speech   --------------------------------------------------------

  def bhojpur_pause(self):
    """ pause playing """

    self.msg("Speech: pause playing")
    self._backend.pause()

  # --- pause speech   --------------------------------------------------------

  def bhojpur_resume(self):
    """ resume playing """

    self.msg("Speech: resume playing")
    self._backend.resume()

  # --- toggle speech state   -------------------------------------------------

  def bhojpur_toggle(self):
    """ toggle playing """

    self.msg("Speech: toggle playing")
    self._backend.toggle()