#!/bin/sh

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

# apt-get update && \
#     apt-get install -y --no-install-recommends \
#         g++ \
#         make \
#         automake \
#         autoconf \
#         bzip2 \
#         unzip \
#         wget \
#         sox \
#         libtool \
#         git \
#         subversion \
#         python2.7 \
#         python3 \
#         zlib1g-dev \
#         gfortran \
#         ca-certificates \
#         patch \
#         ffmpeg \
#         gawk \
# 	vim && \
#     rm -rf /var/lib/apt/lists/*

# ln -s /usr/bin/python2.7 /usr/bin/python

# git clone --depth 1 https://github.com/kaldi-asr/kaldi.git /opt/kaldi && \
#     cd /opt/kaldi/tools && \
#     ./extras/install_mkl.sh && \
#     make -j $(nproc) && \
#     cd /opt/kaldi/src && \
#     ./configure --shared --use-cuda && \
#     make depend -j $(nproc) && \
#     make -j $(nproc) && \
#     find /opt/kaldi  -type f \( -name "*.o" -o -name "*.la" -o -name "*.a" \) -exec rm {} \; && \
#     find /opt/intel -type f -name "*.a" -exec rm {} \; && \
#     find /opt/intel -type f -regex '.*\(_mc.?\|_mic\|_thread\|_ilp64\)\.so' -exec rm {} \; && \
#     rm -rf /opt/kaldi/.git

# cd /opt/kaldi/tools
# ./install_srilm.sh


apt-get update && \
apt-get install -y --no-install-recommends \
    g++ \
    gawk \
    make \
    automake \
    autoconf \
    bzip2 \
    unzip \
    wget \
    sox \
    libtool \
    git \
    subversion \
    python2.7 \
    python3 \
    zlib1g-dev \
    gfortran \
    ca-certificates \
    patch \
    ffmpeg \
    vim \
    pip

git clone https://github.com/kaldi-asr/kaldi.git
cd kaldi/tools
extras/install_mkl.sh
extras/check_dependencies.sh
make
cd ../src

./configure --shared
make depend -j 8
make -j 8



