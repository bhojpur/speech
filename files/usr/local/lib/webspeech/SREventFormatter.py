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

# Bhojpur Speech: implementation of class EventFormatter
#
# The class EventFormatter converts events to a printable form.
# Additional support for i18n is required

class EventFormatter(object):
  """ format events """

  # --- format map: type:format   ---------------------------------------------
  _FMT_MAP = {
    'version': 'Bhojpur Speech online version {value}',
    'icy_meta': '{value}',
    'icy_name': '{value}',
    'rec_start': 'recording {name} for {duration} minutes',
    'rec_stop': 'finished recording. File {file}, duration: {duration}m',
    'vol_set': 'setting current volume to {value}',
    'bhojpur_play_channel': 'start playing channel {nr} ({name})',
    'play': 'playing {value}',
    'pause': 'pausing {value}',
    'file_info': '{name}: {total_pretty}',
    'id3': '{tag}: {value}',
    'keep_alive': 'current time: {value}',
    'eof': '{name} finished',
    'dir_select': 'current directory: {value}'
    }

  # --- format event   --------------------------------------------------------

  def format(self,event):
    """ format given event """

    key = event['type']
    if key in EventFormatter._FMT_MAP:
      if isinstance(event['value'],dict):
        return EventFormatter._FMT_MAP[key].format(**event['value'])
      else:
        return EventFormatter._FMT_MAP[key].format(**event)
    else:
      return "%r" % event