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

# The Bhojpur Speech application module. The file just imports all classes into
# the webspeech namespace

from . Base           import Base           as Base
from . Api            import Api            as Api
from . EventFormatter import EventFormatter as EventFormatter
from . SpeechEvents   import SpeechEvents   as SpeechEvents
from . WebServer      import WebServer      as WebServer
from . Mpg123         import Mpg123         as Mpg123
from . Speech         import Speech         as Speech
from . Player         import Player         as Player
from . Recorder       import Recorder       as Recorder
from . BhojpurSpeech  import BhojpurSpeech  as BhojpurSpeech
from . SpeechClient   import SpeechClient   as SpeechClient
from . KeyController  import KeyController  as KeyController

# voice control with Vosk is optional
have_vosk = False
try:
  from .VoskController  import VoskController  as VoskController
  have_vosk = True
except:
  pass