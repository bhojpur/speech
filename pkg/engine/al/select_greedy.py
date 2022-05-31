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
import math
from collections import defaultdict

counts = defaultdict(int)

def get_ent(counts, sent_counts, tot):
    ent = 0
    for w in counts:
        if w in sent_counts:
            p = float(counts[w] + sent_counts[w]) / tot
            ent += p * math.log(p)
        else:
            p = float(counts[w]) / tot
            ent += p * math.log(p)
    for w in sent_counts:
        if w not in counts:
            p = float(sent_counts[w]) / tot
            ent += p * math.log(p)

    return -ent

ent = 0
tot = 0
for line in open(sys.argv[1]):
    items = line.split()

    new_tot = tot + len(items) - 1

    sent_counts = defaultdict(int)
    for w in items[1:]:
        sent_counts[w] = sent_counts[w] + 1
    new_ent = get_ent(counts, sent_counts, new_tot)

    if new_ent > ent + 1e-8:
        print (line.strip())
        ent = new_ent
        tot = tot + len(items) - 1
        for w in items[1:]:
            counts[w] = counts[w] + 1