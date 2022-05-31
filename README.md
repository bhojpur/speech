# Bhojpur Speech - Processing Engine

The `Bhojpur Speech` is a high-perforance, intelligent speech engine applied within the
[Bhojpur.NET Platform](https://github.com/bhojpur/platform/) for delivery of distributed
`applicationa` or `services`.

## Prerequisites

You need `OpenFST`, `Kaldi` and `Vosk` libraries to be able to build the `server engine`.

```bash
brew install openfst automake sox subversion
fstinfo --help
```

or

```bash
sudo pip3 install openfst
fstinfo --help
```

On macOS, you need to copy the `libvosk.dynlib` into `/usr/local/lib` folder so that the
`Go` programs can detect the library.

## Simple Usage

Firstly, issue the following command in a new Terminal window to install the `webspeech`
server engine. Also, download the `Vosk` [models](https://alphacephei.com/vosk/models)

```bash
python3 -m pip install aiortc aiohttp aiorpc vosk
sudo tools/install.sh [username]
sudo tools/install-vosk.sh
```

To see the command line `Help` options, type the following commands

```bash
webspeech.py -h
```

## Server-side Processor

Firstly, the `automated speech recognition` engine must be built using `Python`.

```bash
cd pkg/server
pip3 install -r requirements.txt
```

then, it should be started in a new Terminal window.

```bash
./pkg/server/websocket/asr_server.py /usr/local/lib/vosk/vosk-model-small-en-us-0.15
```

Typically, it listens at the `ws:localhost:2700` IP address/port.

## Client-side Application

You could try to connect the personal computing device's microphone directly by
running the following command in a new Terminal window.

```bash
./pkg/server/websocket/test_microphone.py -u ws://localhost:2700
```

Also, you could run the following `Go` program to test automated transcription

```bash
go run internal/transcribe/main.go -f ./python/example/test.wav
```
