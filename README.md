# Bhojpur Speech - Processing Engine

The `Bhojpur Speech` is an advanced, high-performance, audio data processing engine using
artificial intelligence techniques for *speech recognition* and *speech synthesis*. It is
applied within [Bhojpur.NET Platform](https://github.com/bhojpur/platform/) for delivery
of distributed `applications` or `services` in various fields (e.g. voice pathology). It
leverages [Vosk](http://alphacephei.com/vosk/) framework that works in offline mode too.

## Key Features

- Offline mode automated *speech recognition* and *speech synthesis*
- A web-based application (using Python) for online speech recognition
- `Python`-based and `Go`-based software framework using C/C++ libraries
- Advanced tools (e.g. Oscilloscope, Recorder, Player) for data processing
- Utilities to build speech training models

## Prerequisites

Please note that this software framework is based on `Python` >= 3.8 and `Go` >= 1.17. So,
please install these runtimes, if you plan to build any custom applications.

It is assumed that `portaudio` will be used to capture audio inputs from your local machine.
However, we have Go libraries to support `serial port`, `portmidi`, and `miniport` as well.

You need [OpenFST](https://www.openfst.org), [Kaldi](https://github.com/kaldi-asr/kaldi),
and [Vosk](http://alphacephei.com/vosk/) software libraries to build the `server engine`.
These libraries could be utilised during custom development of speech training models.

On `macOS`, you could run the following commands to install these key dependencies. Perhaps,
you can use `apt-get` or `yum` command on a Linux server to install the same.

```bash
brew install openfst automake sox subversion ffmpeg portaudio portmidi mpg123
sudo pip3 install numpy flask openfst pyttsx3 flask sseclient
fstinfo --help
```

On `macOS`, a software developer nneds to copy `libvosk.dynlib` into the `/usr/local/lib`
folder so that the `Go` programs could detect the library.

## Speech Recognition Framework

### Installation

Firstly, issue the following command in a new Terminal window to install the `webspeech`
server engine. Also, download the `Vosk` [models](https://alphacephei.com/vosk/models).

```bash
python3 -m pip install aiortc aiohttp aiorpc vosk
```

#### WebSpeech Application

Firstly, please note that `evdev` dependency is available on a Linux operating system only.
It is required for Keyboard device events.

```bash
sudo pip3 install evdev
```

Also, you must have [mpg123](https://www.mpg123.de) >= v1.29.3 MPEG audio player installed.

```bash
sudo apt-get install -y mpg123
sudo tools/install.sh [username]
sudo tools/install-vosk.sh
```

Type the following command in a new Termianl window to run the `webspeech` server engine.

```bash
webspeech.py
```

You can open `http://localhost:8026` URL in a web browser to access the application.

Type the following command in a new Terminal window to run the `webspeech` command line.

```bash
webspeech_cli.py -H localhost -P 8026 -o
```

### Server-side Speech Recognition

Firstly, the `automated speech recognition` engine must be built using `Python`.

```bash
cd pkg/server
pip3 install -r requirements.txt
```

then, it should be started in a new Terminal window.

```bash
python3 ./pkg/server/websocket/asr_server.py /usr/local/lib/vosk/vosk-model-small-en-us-0.15
```

Please note that `vosk-model-small-en-us-0.15` [model](https://alphacephei.com/vosk/models)
is downloaded and installed on your system. Otherwise, please specify your own PATH.

Typically, the `automated speech recognition` engine listens at the `ws:localhost:2700` IP
address/port.

### Speech Recognition Training

Firstly, download the [Kaldi](https://kaldi-asr.org/doc/tutorial.html) source code and run
the following commands in a new Terminal window.

```bash
git clone https://github.com/kaldi-asr/kaldi.git
cd kaldi/tools/; make; cd ../src; ./configure; make
./configure --use-cuda=no
```

Now, edit the `cmd.sh` file under `kaldi/egs/mini_librispeech/s5` (for example). In fact, you could
choose any other folder under `kaldi/egs` or choose to make something of your own.

For data preparation, please refer [here](https://kaldi-asr.org/doc/data_prep.html)

For additional datasets, you can
find some more [models](https://sourceforge.net/projects/cmusphinx/files/Acoustic%20and%20Language%20Models/) to practice.

### Client-side Speech Recognition

You could try to connect the personal computing device's microphone directly by running
the following command in a new Terminal window.

```bash
python3 ./pkg/server/websocket/test_microphone.py -u ws://localhost:2700
```

#### Go-based Speech Transcription

Perhaps, you could run following `Go` program (i.e. [transcribe](internal/transcribe/main.go))
to test `automated speech transcription` methods using [Vosk](http://alphacephei.com/vosk/).

```bash
go run internal/transcribe/main.go -f ./python/example/test.wav
```

## Speech Synthesis Framework

Our `Speech-to-Text` framework is designed to work using `Python` and `Go` bindings. During
training model development, you could also look into [eSpeak](http://espeak.sourceforge.net/)
or [eSpeak NG](https://github.com/espeak-ng/espeak-ng/) frameworks too. Firstly, you need to
install required software libraries on your local machine. For example

```bash
brew install espeak
```

### Client-side Speech Synthesis

You can try our a *web-based* user interface of remote speaker `Go` application built using
[eSpeak](http://espeak.sourceforge.net/) framework.

```bash
go run ./internal/espeak/web/main.go
```

#### Go-based Speech Synthesis

It is performed in *offline* mode using [eSpeak](http://espeak.sourceforge.net/) framework
and our `Go` language bindings.

```bash
speechtext "मेरा नाम भोजपुर कंसल्टिंग है"
speechplay audios/test_hi.wav
```

#### Python-based Speech Synthesis

A sample `Python` program (i.e. [speaker](internal/speaker/main.py)) is included in this
repository. It is based on `pyttsx3` or [coqui STT](https://github.com/coqui-ai/STT) library.

```bash
sudo pip3 install pyttsx3 stt
python3 ./internal/speaker/main.py
```

## Speech Translation Framework

We have Google enabled language translation capabilities integrated. Perhaps, you could try
the following [program](/internal/translate/main.go). Also, it detects the language used.

```bash
./internal/translate/main.go "मैं भोजपुर कंसल्टिंग के लिए काम कर रहा हूँ"
```
