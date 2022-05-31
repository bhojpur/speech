#!/usr/bin/python3

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

import sys
import librosa
import numpy
import phash
import pickle


def get_hash(mel_basis, wavfn, start, end):
     y, sr = librosa.load(wavfn, sr=16000)
     y = y[int(start * sr):int(end * sr)]
     # Get the hop size so we get about 50 frames

     hop = int(y.shape[0] / 64) + 1

     S, n_fft = librosa.core.spectrum._spectrogram(y=y, n_fft=512, hop_length=hop, power=2.0)



     S = librosa.power_to_db(numpy.dot(mel_basis, S), ref=numpy.max)
     h = phash.hash(S)
     return h

class Segment():
    def __init__(self, utt, start, dur, name):
        self.utt = utt
        self.start = float(start)
        self.dur = float(dur)
        self.name = name

    def __repr__(self):
        return "[%s %s %.3f %.3f]" % (self.utt, self.name, self.start, self.dur)

def SegmentGenerator(wav_list, phone_list):
    wavs = {}
    for line in open(wav_list):
        items = line.split()
        wavs[items[0]] = items[1]

    segments={}
    for line in open(phone_list):
        utt, channel, start, dur, pn = line.split()
        if utt not in segments:
            segments[utt] = []
        segments[utt].append(Segment(utt, start, dur, pn))

    # Build a Mel filter
    mel_basis = librosa.filters.mel(16000, n_fft=512, n_mels=32)

    for utt in segments:
        utt_segments = segments[utt]
        for i, phone in enumerate(utt_segments):
            # End should be approximately + 0.5 seconds from start
            j = i
            start = phone.start
            end = phone.start
            while end < start + 0.5 and j < len(utt_segments):
                end = end + utt_segments[j].dur
                j = j + 1
            if j - i < 3 or end - start < 0.4: # Ignore this
                continue

            mhash = get_hash(mel_basis, wavs[utt], start, end)
            yield (mhash, start, end, utt_segments[i:j + 1])


def index_data():
    try:
        inf = open(sys.argv[3], "rb")
        database = pickle.load(inf)
    except:
        database = {}

    for mhash, start, end, segments in SegmentGenerator(sys.argv[1], sys.argv[2]):
        if mhash not in database:
            database[mhash] = []
#        print (mhash, start, end, segments)
        database[mhash].append((segments, start, end))
    pickle.dump(database, open(sys.argv[3], "wb"))

if __name__ == '__main__':
    index_data()