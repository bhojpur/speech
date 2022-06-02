#!/usr/bin/env python3

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

import json
# import librosa
import subprocess
import sys
import tempfile


def main():
    if len(sys.argv) != 2:
        sys.exit('Usage:  {} <sound-file>'.format(sys.argv[0]))
    filename = sys.argv[1]

    print('Loading file ...')
    soundfile = load_soundfile(filename)

    print('Extrating data ...', )
    data = extract_data(soundfile)

    # Pretty print
    pretty_print(data)


def pretty_print(data):
    print(json.dumps(data, indent=4))
    # f = tempfile.NamedTemporaryFile(delete=False)
    # f.write(json.dumps(data).encode('ascii'))
    # print("jq '.' {}".format(f.name))
    # subprocess.call("jq '.' {}".format(f.name), shell=True)


def extract_data(sound):
    # TODO
    databytes = b'{"tests":{"net.wifi":{"result":1,"resultText":"SSID"},"net.gateway":{"result":1,"resultText":"ping 1.2.3.4 ok"},"net.inet":{"result":1,"resultText":"ping 8.8.8.8 ok"},"net.dns":{"result":1,"resultText":"vpn.eneco.toon.eu: 63.35.124.51"},"net.time":{"result":1,"resultText":"2019-05-16T09:16:08+0000"}}}'
    return json.loads(databytes)


def load_soundfile(filepath):
    # return librosa.load(filepath)
    pass

if __name__ == '__main__':
    main()