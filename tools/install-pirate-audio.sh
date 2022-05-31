#!/bin/bash

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

# This script installs packages necessary for Pirate-Audio hats. It also adds
# all necessary options to /boot/config.txt.

# --- defaults used during installation   ----------------------------------

PACKAGES="python3-rpi.gpio python3-spidev python3-pil python3-numpy"
PACKAGES_PIP="st7789"

PROJECT="webspeech-pirate-audio"

# --- basic packages   ------------------------------------------------------

if [ -n "$PACKAGES" ]; then
  apt-get update
  apt-get -y install $PACKAGES
fi

# install PIP3 packages
[ -n "$PACKAGES_PIP" ] && pip3 --disable-pip-version-check install $PACKAGES_PIP

# --- configure system   ----------------------------------------------------

# update /boot/config.txt
if ! grep -q "^dtparam=spi=on" /boot/config.txt ; then
  echo -e "[INFO] configuring SPI in /boot/config.txt" 2>&1
  echo "dtparam=spi=on" >> /boot/config.txt
fi

if grep -q "^dtparam=audio=on" /boot/config.txt ; then
  echo -e "[INFO] disabling default audio in /boot/config.txt" 2>&1
  sed -i -e  "s/dtparam=audio=on/dtparam=audio=off/" /boot/config.txt
fi

if ! grep -q "^pio=25" /boot/config.txt ; then
  echo -e "[INFO] activating DAC in /boot/config.txt" 2>&1
  echo "# activate DAC" >> /boot/config.txt
  echo "pio=25=op,dh" >> /boot/config.txt
  echo "dtoverlay=hifiberry-dac" >> /boot/config.txt
fi
    
if ! grep -q "^dtoverlay=gpio-key,gpio=5" /boot/config.txt ; then
  echo -e "[INFO] configuring keys in /boot/config.txt" 2>&1
  cat >> /boot/config.txt << EOF

#key-mapping: A->up, B->down, X->right, Y->left
#              5       6       16        24
dtoverlay=gpio-key,gpio=5,keycode=103,label="UP"
dtoverlay=gpio-key,gpio=6,keycode=108,label="DOWN"
dtoverlay=gpio-key,gpio=16,keycode=106,label="RIGHT"
dtoverlay=gpio-key,gpio=24,keycode=105,label="LEFT"
# Note: remove last line on old hardware, where Y is connected to GPIO20
EOF
fi

# --- final configuration is manual   ---------------------------------------

echo -e "\nplease reboot to activate changes"