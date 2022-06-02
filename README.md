# Bhojpur Speech - Processing Engine

The `Bhojpur Speech` is an advanced, high-performance, audio data processing engine using
artificial intelligence techniques for `speech recognition` and `speech synthesis`. It is
applied within [Bhojpur.NET Platform](https://github.com/bhojpur/platform/) for delivery
of distributed `applications` or `services` in various fields (e.g. voice pathology). It
leverages [Vosk](http://alphacephei.com/vosk/) framework that works in offline mode too.

## Key Features

- Offline mode automated `speech recognition` and `speech synthesis`
- A web-based application (using Python) for online speech recognition
- `Python`-based and `Go`-based software framework using C/C++ libraries
- Utilities to build speech training models

## Prerequisites

Please note that this software framework is based on `Python` >= 3.0 and `Go` >= 1.16. So,
please install these runtimes, if you plan to build any custom applications.

You need [OpenFST](https://www.openfst.org), [Kaldi](https://github.com/kaldi-asr/kaldi),
and [Vosk](http://alphacephei.com/vosk/) software libraries to build the `server engine`.

```bash
brew install openfst automake sox subversion
sudo pip3 install numpy flask openfst pyttsx3
fstinfo --help
```

On macOS, you need to copy `libvosk.dynlib` into the `/usr/local/lib` folder so that the
`Go` programs could detect the library.

## Speech Recognition Framework

### Installation

Firstly, issue the following command in a new Terminal window to install the `webspeech`
server engine. Also, download the `Vosk` [models](https://alphacephei.com/vosk/models).

```bash
python3 -m pip install aiortc aiohttp aiorpc vosk
```

#### WebSpeech Application

```bash
sudo tools/install.sh [username]
sudo tools/install-vosk.sh
```

To check all the `webspeech` command line `Help` options, type the following command

```bash
webspeech.py -h
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

### Client-side Speech Synthesis

#### Python-based Speech Synthesis

A sample `Python` program (i.e. [speaker](internal/speaker/main.py)) is included in this
repository. It is based on `pyttsx3` library.

```bash
python3 ./internal/speaker/main.py
```
