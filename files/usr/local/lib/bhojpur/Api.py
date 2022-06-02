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

# Bhojpur Speech: implementation of class API
# Collect all API-functions

from bhojpur import Base

class Api(Base):
  """ The class holds references to all Bhojpur Speech API functions """

  def __init__(self,app):
    """ initialization """

    self._app          = app
    self.debug         = app.debug

  # --- execute Bhojpur Speech API by name   ----------------------------------

  def _exec(self,name,**args):
    """ execute a Bhojpur Speech API by name """

    if hasattr(self,name):
      self.msg("executing: %s(%r)" % (name,dict(**args)))
      return getattr(self,name)(**args)
    else:
      self.msg("unknown Bhojpur Speech API method %s" % name)
      raise NotImplementedError("Bhojpur Speech API %s not implemented" % name)

  # --- return list of APIs   ------------------------------------------------

  def get_api_list(self):
    """ return list of Bhojpur Speech APIs """

    return [func for func in dir(self)
            if callable(getattr(self, func)) and not func.startswith("_")
            and func not in Base.__dict__ ]