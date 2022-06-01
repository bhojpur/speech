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

# set base image (host OS)
FROM ubuntu:20.04

ENV TZ=Asia/Kolkata \
    DEBIAN_FRONTEND=noninteractive

# set the working directory in the container
WORKDIR /vosk_training

# copy the dependencies file to the working directory
COPY requirements.txt  install_kaldi.sh ./
RUN pwd
RUN ls
# install dependencies
RUN chmod +x install_kaldi.sh
RUN ./install_kaldi.sh
RUN pip install -r requirements.txt
RUN ls
WORKDIR /vosk_training/kaldi/tools
#RUN cd /kaldi/tools
RUN extras/install_opengrm.sh
RUN ./install_srilm.sh lolz lorg abc@gmail.com

# copy the content of the local src directory to the working directory
#WORKDIR /vosk_training/vosk-model-en-us-0.22-compile/
#COPY  vosk-model-en-us-0.22-compile .
RUN wget https://alphacephei.com/vosk/models/vosk-model-en-us-0.22-compile.zip
RUN unzip vosk-model-en-us-0.22-compile.zip

# command to run on container start
CMD [ "sh", "-c", "/bin/sh"]
