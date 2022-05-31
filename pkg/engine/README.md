# Bhojpur Speech - Server-side Framework

A `speech recognition` system

## Simple Usage

To install the requirements run

```bash
pip3 install -r requirements.txt
```

To prepare the training/verification data create the following two files:

- `wav.scp` list to map uterances to wav files in filesystem
- `phones.txt` the CTM file with phonemes and timings. It could be CTM file from
the alignment or it could be a CTM file from the decoding

## Indexing

To add the data to the database run

```bash
python3 index.py wavs-train.txt phones-train.txt data.idx
```

## Verification

To verify decoding results run

```bash
python3 verify.py wavs-test.txt phones-test.txt data.idx
```
