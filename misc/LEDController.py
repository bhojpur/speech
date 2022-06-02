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

# Bhojpur Speech: sample implementation of the class LEDController
#
# The methods of this class will be called by VoskController. For your
# own version, you have to implement all methods.

import webspeech.apa102 as apa102
import time
from gpiozero import LED

class LEDController:
  """ change LED according to events """

  NUM_PIXELS = 12
  DELAY      = 0.5

  # --- constructor   --------------------------------------------------------

  def __init__(self):
    """ constructor """

    self._leds = apa102.APA102(num_led=LEDController.NUM_PIXELS,
                               global_brightness=2)
        
    self._power = LED(5)
    self._power.on()

  # --- set color (helper method)   ------------------------------------------

  def _set_color(self,r,g,b,flash=1):
    """ set color-ring and flash it """

    for f in range(flash):
      for i in range(LEDController.NUM_PIXELS):
        self._leds.set_pixel(i,r,g,b)
      self._leds.show()
      time.sleep(LEDController.DELAY)
      self._leds.clear_strip()
      time.sleep(LEDController.DELAY)

    for i in range(LEDController.NUM_PIXELS):
      self._leds.set_pixel(i,r,g,b)
    self._leds.show()

  # --- after mic is activated by wake-word   --------------------------------
  
  def active(self):
    """ active state (waiting for a command) """

    self._set_color(0,0,255)   # all blue

  # --- after mic is waiting for wake-word   ---------------------------------
  
  def inactive(self):
    """ inactive, waiting for wake word """

    self._leds.clear_strip()

  # --- after successful detection of a command   ----------------------------
  
  def success(self):
    """ detected valid command in active mode """

    self._set_color(0,255,0)   # all green

  # --- after detection of an unknown command   ------------------------------
  
  def unknown(self):
    """ detected unknown command in active mode """

    self._set_color(255,0,0,flash=2)   # all red

# --- main (test) program   --------------------------------------------------

if __name__ == '__main__':

  ctrl = LEDController()
  print("activating...")
  ctrl.active()
  time.sleep(3)

  print("deactivating...")
  ctrl.inactive()
  time.sleep(1)

  print("activating...")
  ctrl.active()
  time.sleep(2)

  print("success...")
  ctrl.success()
  time.sleep(2)

  print("unknown...")
  ctrl.unknown()
  time.sleep(2)

  print("the end")
  ctrl.inactive()