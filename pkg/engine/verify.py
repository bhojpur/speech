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

from index import SegmentGenerator, Segment

def verify_data():

    try:
        inf = open(sys.argv[3], "rb")
        database = pickle.load(inf)
    except:
        database = {}

    for mhash, start, end, segments in SegmentGenerator(sys.argv[1], sys.argv[2]):
        if mhash in database:
            target = " ".join([x.name for x in segments])
            source = [" ".join([x.name for x in chunk[0]]) for chunk in database[mhash]]
            if target in source:
                 print ("+", target, source, segments[0].utt, start, end)
            else:
                 print ("-", target, source, segments[0].utt, start, end)

if __name__ == '__main__':
    verify_data()