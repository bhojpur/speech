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

# Bhojpur Speech: implementation of class BhojpurSpeech
#
# The class BhojpurSpeech implements the main application class and serves as model
# and controller (for generic functions)

import os, sys, json, traceback, threading
import configparser

from bhojpur import *

# --- main application class   ----------------------------------------------

class BhojpurSpeech(Base):
  """ main application class """

  VERSION = "1.0.0"

  def __init__(self,options):
    """ initialization """

    self.options    = options
    self.parser     = configparser.RawConfigParser(inline_comment_prefixes=(';',))
    self.parser.optionxform = str
    self.parser.read('/etc/webspeech.conf')

    self.read_config(options)
    self._store = os.path.join(os.path.expanduser("~"),".webspeech.json")

    self._threads    = []                   # thread-store
    self.stop_event  = threading.Event()

    # create API-object and register our own functions
    self.api = Api(self)
    self.register_apis()

    # create all objects
    if options.do_record:
      self._events  = SpeechEvents(self)
      self.backend  = None
      self.speech   = Speech(self)
      self.recorder = Recorder(self)
      self._objects = [self,self.speech,self.recorder]
    elif options.do_play:
      self._events  = SpeechEvents(self)
      self.backend  = Mpg123(self)
      self.speech   = Speech(self)
      self.player   = Player(self)
      self._objects = [self,self.speech,self.player,self.backend]
    elif options.do_list:
      self.backend  = None
      self.speech   = Speech(self)
      self._objects = [self,self.speech]
    else:
      self._events  = SpeechEvents(self)
      self._server  = WebServer(self)
      self.backend  = Mpg123(self)
      self.speech   = Speech(self)
      self.player   = Player(self)
      self.recorder = Recorder(self)
      self._objects = [self,self.speech,self.player,
                       self.recorder,self.backend]

    self._state = {'mode': 'speech'}
    self._load_state()
    if self.backend:
      self.backend.create()

  # --- read configuration   -------------------------------------------------

  def read_config(self,options):
    """ read configuration from config-file """

    # section [GLOBAL]
    if options.debug:
      self.debug = True
    else:
      self.debug  = self.get_value(self.parser,"GLOBAL", "debug","0") == "1"

  # --- register APIs   ------------------------------------------------------

  def register_apis(self):
    """ register API-functions """

    self.api.get_version      = self._get_version
    self.api.sys_restart      = self.sys_restart
    self.api.sys_stop         = self.sys_stop
    self.api.sys_reboot       = self.sys_reboot
    self.api.sys_halt         = self.sys_halt
    self.api.update_state     = self.update_state
    self.api.get_state        = self.get_state

  # --- return version   ---------------------------------------------------

  def _get_version(self):
    """ return version """

    self.msg("Bhojpur Speech: online version: %s" % BhojpurSpeech.VERSION)
    return BhojpurSpeech.VERSION

  # --- return state   -----------------------------------------------------

  def get_state(self):
    """ return state """

    return self._state

  # --- shutdown system   -----------------------------------------------------

  def sys_halt(self):
    """ shutdown system """

    self.msg("Bhojpur Speech: processing sys_halt")
    if not self.debug:
      try:
        os.system("sudo /sbin/halt &")
      except:
        pass
    else:
      self.msg("Bhojpur Speech: no shutdown in debug-mode")
    self.api._push_event({'type': 'sys', 'value': 'halt'})

  # --- reboot system   -----------------------------------------------------

  def sys_reboot(self):
    """ reboot system """

    self.msg("Bhojpur Speech: processing sys_reboot")
    if not self.debug:
      try:
        os.system("sudo /sbin/reboot &")
      except:
        pass
    else:
      self.msg("Bhojpur Speech: no reboot in debug-mode")
    self.api._push_event({'type': 'sys', 'value': 'reboot'})

  # --- restart service   -----------------------------------------------------

  def sys_restart(self):
    """ restart service """

    self.msg("Bhojpur Speech: processing sys_restart")
    if not self.debug:
      try:
        os.system("sudo /bin/systemctl restart webspeech.service &")
      except:
        pass
    else:
      self.msg("Bhojpur Speech: no application restart in debug-mode")
    self.api._push_event({'type': 'sys', 'value': 'restart'})

  # --- stop service   --------------------------------------------------------

  def sys_stop(self):
    """ stop service """

    self.msg("Bhojpur Speech: processing sys_stop")
    if not self.debug:
      try:
        os.system("sudo /bin/systemctl stop webspeech.service &")
      except:
        pass
    else:
      self.msg("Bhojpur Speech: no application stop in debug-mode")
    self.api._push_event({'type': 'sys', 'value': 'stop'})

  # --- update (and distribute) state   ---------------------------------------

  def update_state(self,state=None,section=None,key=None,value=None,publish=True):
    """ update state and publish as event """

    if state:
      # update on key-level
      for s in state.keys():
        if isinstance(s,dict):
          for k in s.keys():
            if s in self._state:
              self._state[s][k] = state[s][k]
            else:
              self._state[s] = {k:state[s][k]}
        else:
          self._state[s] = state[s]
    elif section and key:
      if section in self._state:
        self._state[section][key] = value
      else:
        self._state[section] = {key:value}
    if publish:
      self.api._push_event({'type': 'state', 'value': self._state})
    return

  # --- query state of objects and save   -------------------------------------

  def _save_state(self):
    """ query and save state of objects """

    state = {}
    for obj in self._objects:
      state[obj.__module__] = obj.get_persistent_state()

    f = open(self._store,"w")
    self.msg("Bhojpur Speech: Saving settings to %s" % self._store)
    json.dump(state,f,indent=2,sort_keys=True)
    f.close()

  # --- load state of objects   -----------------------------------------------

  def _load_state(self):
    """ load state of objects """

    try:
      if not os.path.exists(self._store):
        state = {}
      else:
        self.msg("Bhojpur Speech: Loading settings from %s" % self._store)
        f = open(self._store,"r")
        state = json.load(f)
        f.close()
      for obj in self._objects:
        if obj.__module__ in state:
          obj.set_persistent_state(state[obj.__module__])
    except:
      self.msg("Bhojpur Speech: Loading settings failed")
      if self.debug:
        traceback.print_exc()

  # --- return persistent state of this class   -------------------------------

  def get_persistent_state(self):
    """ return persistent state (overrides SRBase.get_pesistent_state()) """
    return {
      'mode': self._state['mode']
      }

  # --- restore persistent state of this class   ------------------------------

  def set_persistent_state(self,state_map):
    """ restore persistent state (overrides SRBase.set_pesistent_state()) """

    self.msg("Bhojpur Speech: restoring persistent state")
    if 'mode' in state_map:
      self._state['mode'] = state_map['mode']

  # --- setup signal handler   ------------------------------------------------

  def signal_handler(self,_signo, _stack_frame):
    """ signal-handler for clean shutdown """

    self.msg("Bhojpur Speech: received signal, stopping program ...")
    self.cleanup()

  # --- cleanup of ressources   -----------------------------------------------

  def cleanup(self):
    """ cleanup of ressources """

    if hasattr(self,'backend') and self.backend:
      self.backend.destroy()
    if hasattr(self,'_server') and self._server:
      self._server.stop()
    self.stop_event.set()
    if hasattr(self.api,'rec_stop'):
      self.api.rec_stop()
    map(threading.Thread.join,self._threads)
    self._save_state()
    self.msg("Bhojpur Speech: ... done stopping program")

  # --- run method   ----------------------------------------------------------

  def run(self):
    """ start all threads and return """

    threading.Thread(target=self._server.run).start()
    self.msg("Bhojpur Speech: started web-server")