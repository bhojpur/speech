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

# Bhojpur Speech: implementation of class Base
#
# The class Base is the root-class of all classes and implements common methods

import sys, time

class Base:
  """ base class with common methods """

  # --- print debug messages   ------------------------------------------------

  def msg(self,text,force=False):
    """ print debug-message """

    if force:
      sys.stderr.write("%s\n" % text)
    elif self.debug:
      sys.stderr.write("[DEBUG %s] %s\n" % (time.strftime("%H:%M:%S"),text))
    sys.stderr.flush()

  # --- read configuration value   --------------------------------------------

  def get_value(self,parser,section,option,default):
    """ get value of config-variables and return given default if unset """

    if parser.has_section(section):
      try:
        value = parser.get(section,option)
      except:
        value = default
    else:
      value = default
    return value

  # --- return persistent state of this class   -------------------------------

  def get_persistent_state(self):
    """ return persistent state (implemented by subclasses) """
    return {}

  # --- set state state of this class   ---------------------------------------

  def set_persistent_state(self,state_map):
    """ set state (implemented by subclasses) """
    pass