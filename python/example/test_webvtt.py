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

from vosk import Model, KaldiRecognizer, SetLogLevel
from webvtt import WebVTT, Caption
import sys
import os
import subprocess
import json
import textwrap

SetLogLevel(-1)

sample_rate = 16000
model = Model(lang="en-us")
rec = KaldiRecognizer(model, sample_rate)
rec.SetWords(True)

WORDS_PER_LINE = 7


def timeString(seconds):
    minutes = seconds / 60
    seconds = seconds % 60
    hours = int(minutes / 60)
    minutes = int(minutes % 60)
    return '%i:%02i:%06.3f' % (hours, minutes, seconds)


def transcribe():
    command = ['ffmpeg', '-nostdin', '-loglevel', 'quiet', '-i', sys.argv[1],
               '-ar', str(sample_rate), '-ac', '1', '-f', 's16le', '-']
    process = subprocess.Popen(command, stdout=subprocess.PIPE)

    results = []
    while True:
        data = process.stdout.read(4000)
        if len(data) == 0:
            break
        if rec.AcceptWaveform(data):
            results.append(rec.Result())
    results.append(rec.FinalResult())

    vtt = WebVTT()
    for i, res in enumerate(results):
        words = json.loads(res).get('result')
        if not words:
            continue

        start = timeString(words[0]['start'])
        end = timeString(words[-1]['end'])
        content = ' '.join([w['word'] for w in words])

        caption = Caption(start, end, textwrap.fill(content))
        vtt.captions.append(caption)

    # save or return webvtt
    if len(sys.argv) > 2:
        vtt.save(sys.argv[2])
    else:
        print(vtt.content)


if __name__ == '__main__':
    if not (1 < len(sys.argv) < 4):
        print(f'Usage: {sys.argv[0]} audiofile [output file]')
        exit(1)
    transcribe()