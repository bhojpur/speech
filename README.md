# Bhojpur Speech - Processing Engine

The `Bhojpur Speech` is a high-perforance, intelligent speech engine applied within the
[Bhojpur.NET Platform](https://github.com/bhojpur/platform/) for delivery of distributed
`applicationa` or `services`.

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

## Server Engines

```bash
cd pkg/server/webrtc
python3 asr_server_webrtc.py
```
