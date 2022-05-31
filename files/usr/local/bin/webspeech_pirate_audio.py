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

# A control program for the Pimoroni's Pirate-Audio hats.

import locale, os, sys, json, traceback
from webspeech_cli import SpeechCli

try:
  from ST7789 import ST7789
  from PIL import Image, ImageDraw
  have_st7789 = True
except:
  have_st7789 = False

# --- application class   ----------------------------------------------------

class PirateAudio(SpeechCli):
  """ application class """

  # --- constructor   --------------------------------------------------------

  def __init__(self):
    """ constructor """

    super(PirateAudio,self).__init__()
    if have_st7789:
      self.msg("PirateAudio: detected ST7789")
      self._init_display()
    else:
      self.msg("PirateAudio: no ST7789")

  # --- init display   -------------------------------------------------------

  def _init_display(self):
    """ initialize the display """

    SPI_SPEED_MHZ = 80

    self._last_logo = ""

    self._screen = ST7789(
      rotation=90,  # Needed to display the right way up on Pirate Audio
      port=0,       # SPI port
      cs=1,         # SPI port Chip-select channel
      dc=9,         # BCM pin used for data/command
      backlight=13,
      spi_speed_hz=SPI_SPEED_MHZ * 1000 * 1000
      )

  # --- update display   -----------------------------------------------------

  def _update_display(self,logo):
    """ update display with logo (logo is relative to the web-root) """

    if logo == self._last_logo:
      return

    logo_file = os.path.join(self.pgm_dir,"..","lib","webspeech",
                             "web",logo)
    if not os.path.exists(logo_file):
      logo_file = os.path.join(self.pgm_dir,"..","lib","webspeech",
                             "web","images","default.png")
    self.msg("PirateAudio: logo-file: %s" % logo_file)

    try:
      img = None
      im  = Image.open(logo_file)
      img = im.resize((240,240))
      im.close()
      self._screen.display(img)
      self._last_logo = logo
    except:
      traceback.print_exc()
      self.msg("PirateAudio: failed to display %s" % logo_file)

  # --- handle event   -------------------------------------------------------

  def handle_event(self,event):
    """ override to display channel-logo """

    ev_data = json.loads(event.data)
    if ev_data['type'] != 'bhojpur_play_channel':
      self.msg("PirateAudio: ignoring event with type %s" % ev_data['type'])
      return

    logo = ev_data['value']['logo']
    self.msg("PirateAudio: logo: %s" % logo)
    if have_st7789:
      self._update_display(logo)
    
  # --- close connection   ---------------------------------------------------

  def close(self):
    """ override to blank display """

    #self._screen.set_backlight(0)        # does not work during shutdown
    self._screen.display(Image.new('RGB',(240,240),color=(0,0,0)))
    super(PirateAudio,self).close()

# --- main program   ---------------------------------------------------------

if __name__ == '__main__':

  # set local to default from environment
  locale.setlocale(locale.LC_ALL, '')

  # create client-class and parse arguments
  app = PirateAudio()
  app.pgm_dir = os.path.dirname(os.path.realpath(__file__))
  app.run()