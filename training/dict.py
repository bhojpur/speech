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

import phonetisaurus
import pandas as pd
words = {}

for line in open("db/en.dic"):
    items = line.split()
    if items[0] not in words:
         words[items[0]] = []
    words[items[0]].append(" ".join(items[1:]))

extra_dic_words={}
for line in open("db/extra.dic"):
    items = line.split()
    if items[0] not in words:
         words[items[0]] = []
    extra_dic_words[items[0]] = []
    words[items[0]].append(" ".join(items[1:]))
    extra_dic_words[items[0]].append(" ".join(items[1:]))

new_words = set()
for line in open("db/extra.txt"):
    for w in line.split():
        if w not in words:
             new_words.add(w)


for w, phones in phonetisaurus.predict(new_words, "db/en-g2p/en.fst"):
    words[w] = []
    if w in extra_dic_words.keys():
        phones.append(extra_dic_words[w])
    words[w].append(" ".join(phones))

for w, phones in words.items():
    for p in phones:
        print (w, p)

