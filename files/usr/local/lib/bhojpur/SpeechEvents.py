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

# Bhojpur Speech: implementation of class SpeechEvents
#
# The class SpeechEvents multiplexes events to multiple consumers.

import queue, threading, datetime, sys

from bhojpur import Base
from bhojpur import EventFormatter

class SpeechEvents(Base):
  """ Multiplex events to consumers """

  QUEUE_SIZE          = 20   # size of client event-queues
  KEEP_ALIVE_INTERVAL = 15   # send keep-alive every x seconds

  def __init__(self,app):
    """ initialization """

    self._api         = app.api
    self.debug        = app.debug
    self._stop_event  = app.stop_event
    self._input_queue = queue.Queue()
    self._lock        = threading.Lock()
    self._consumers   = {}
    self._formatter   = EventFormatter()
    self.register_apis()
    threading.Thread(target=self._process_events).start()

  # --- register APIs   ------------------------------------------------------

  def register_apis(self):
    """ register Bhojpur Speech API functions """

    self._api._push_event   = self.push_event
    self._api._add_consumer = self.add_consumer
    self._api._del_consumer = self.del_consumer

  # --- push an event to the input queue   -----------------------------------

  def push_event(self,event):
    """ push event to the voice input queue """

    self._input_queue.put(event)

  # --- add a consumer   -----------------------------------------------------

  def add_consumer(self,id):
    """ add a consumer to the list of consumers """

    if id in self._consumers:
      self.msg("SpeechEvents: reusing consumer-queue with id %s" % id)
      return self._consumers[id]
    else:
      self.msg("SpeechEvents: adding consumer with id %s" % id)
      with self._lock:
        self._consumers[id] = queue.Queue(SpeechEvents.QUEUE_SIZE)
      try:
        ev = {'type': 'version','value': self._api.get_version()}
        ev['text'] = self._formatter.format(ev)
        self._consumers[id].put_nowait(ev)
        ev = {'type': 'state','value': self._api.get_state()}
        ev['text'] = self._formatter.format(ev)
        self._consumers[id].put_nowait(ev)
        return self._consumers[id]
      except:
        with self._lock:
          del self._consumers[id]
        return None

  # --- remove a consumer   --------------------------------------------------

  def del_consumer(self,id):
    """ delete a consumer from the list of consumers """

    if id in self._consumers:
      self._consumers[id].put(None)
      with self._lock:
        del self._consumers[id]

  # --- multiplex events   ---------------------------------------------------

  def _process_events(self):
    """ pull events from the input-queue and distribute to the consumer queues """

    self.msg("SpeechEvents: starting event-processing")
    count = 0
    while not self._stop_event.is_set():
      try:
        event = self._input_queue.get(block=True,timeout=1)   # block 1s
        self._input_queue.task_done()
        self.msg("SpeechEvents: received event: %r" % (event,))
        count = 0
      except queue.Empty:
        count = (count+1) % SpeechEvents.KEEP_ALIVE_INTERVAL
        if count > 0:
          continue
        else:
          #self.msg("SpeechEvents: publishing keep-alive")
          event = {'type': 'keep_alive', 'value':
                   datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")}

      event['text'] = self._formatter.format(event)
      stale_consumers = []
      for id, consumer in self._consumers.items():
        try:
          consumer.put_nowait(event)
        except queue.Full:
          stale_consumers.append(id)

      # delete stale consumers
      with self._lock:
        for id in stale_consumers:
          self.msg("SpeechEvents: deleting stale queue with id %s" % id)
          del self._consumers[id]

    self.msg("SpeechEvents: stopping event-processing")
    for consumer in self._consumers.values():
      try:
        consumer.put_nowait(None)
      except:
        pass
    self.msg("SpeechEvents: event-processing finished")