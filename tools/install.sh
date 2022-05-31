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

# This script installs files and services specific to Bhojpur Speech project.

# --- defaults used during installation   ----------------------------------

PACKAGES="python3-pip python3-flask mpg123 mp3info python3-evdev"
PACKAGES_PIP="sseclient-py"

PROJECT="webspeech"
USERNAME="${1:-pi}"

# --- basic packages   ------------------------------------------------------

if [ -n "$PACKAGES" ]; then
  apt-get update
  apt-get -y install $PACKAGES
fi

# install PIP3 packages
[ -n "$PACKAGES_PIP" ] && pip3 --disable-pip-version-check install $PACKAGES_PIP

# --- install specific files   ----------------------------------------------

rand="$RANDOM"
if [ -f /etc/${PROJECT}.conf ]; then
  # save current configuration
  mv /etc/${PROJECT}.conf /etc/${PROJECT}.conf.$rand
fi
if [ -f /etc/${PROJECT}.channels ]; then
  # save channel-list
  mv /etc/${PROJECT}.channels /etc/${PROJECT}.channels.$rand
fi

for f in `find $(dirname "$0")/../files/ -type f -not -name "*.pyc"`; do
  target="${f#*files}"
  target_dir="${target%/*}"
  [ ! -d "$target_dir" ] && mkdir -p "$target_dir"
  cp "$f" "$target"
  chown root:root "$target"
  chmod 644       "$target"
done

chmod 755 /usr/local/bin/${PROJECT}.py
chmod 755 /usr/local/bin/${PROJECT}_cli.py
chmod 755 /usr/local/bin/${PROJECT}_pirate_audio.py
chmod 755 /usr/local/bin/${PROJECT}_chrome.sh
chmod 644 /etc/${PROJECT}.conf
chmod 644 /etc/${PROJECT}.channels
chmod 644 /etc/systemd/system/${PROJECT}.service

# restore old configuration
if [ -f /etc/${PROJECT}.conf.$rand ]; then
  mv -f /etc/${PROJECT}.conf /etc/${PROJECT}.conf.new
  mv /etc/${PROJECT}.conf.$rand /etc/${PROJECT}.conf
  echo -e "\nnew version of configuration file: /etc/${PROJECT}.conf.new"
fi
if [ -f /etc/${PROJECT}.channels.$rand ]; then
  mv -f /etc/${PROJECT}.channels /etc/${PROJECT}.channels.new
  mv -f /etc/${PROJECT}.channels.$rand /etc/${PROJECT}.channels
fi

# --- fix user of service   -------------------------------------------------

sed -i -e "/User=/s/=.*/=$USERNAME/" \
  /etc/systemd/system/${PROJECT}.service \
  /etc/systemd/system/${PROJECT}-cli.service \
  /etc/systemd/system/${PROJECT}-pirate-audio.service
usermod -a -G audio "$USERNAME"

# --- activate service   ----------------------------------------------------

#systemctl enable ${PROJECT}.service

# --- final configuration is manual   ---------------------------------------

echo -e "\nPlease edit /etc/${PROJECT}.conf and start ${PROJECT}.service"