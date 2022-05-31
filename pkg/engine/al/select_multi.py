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

import sys
import math

from collections import defaultdict

def main():

    counts = defaultdict(float)
    tot = 0.0
    step = 8
    lines = open(sys.argv[1]).readlines()

    for i in range(0, len(lines), step):

        sublines = lines[i : i + step]

        allw = set()
        new_tot = [0.0] * step
        new_counts = [None] * step
        for j in range(0, step):
            items = sublines[j].split()
            new_tot[j] = tot + len(items) - 1
            new_counts[j] = defaultdict(float)
            for w in items[1:]:
                new_counts[j][w] = new_counts[j][w] + 1
                allw.add(w)
        allwf = frozenset(allw)

        ent = [0.0] * step
        for w, wcnt in counts.items():
            if w in allwf:
                for j in range(0, step):
                    prob = (wcnt + new_counts[j].get(w, 0.0)) / new_tot[j]
                    ent[j] -= prob * math.log2(prob)
            else:
                for j in range(0, step):
                    prob = wcnt / new_tot[j]
                    ent[j] -= prob * math.log2(prob)
        for j in range(0, step):
            for w in new_counts[j]:
                if w not in counts:
                    prob = new_counts[j][w] / new_tot[j]
                    ent[j] -= prob * math.log2(prob)

        new_index = max(range(len(ent)), key=ent.__getitem__)
        print (sublines[new_index].strip())
        items = sublines[new_index].split()
        tot = tot + len(items) - 1
        for w in items[1:]:
            counts[w] = counts[w] + 1

main()