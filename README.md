# Batch AVIF

## Build
```terminal
git clone https://github.com/Delnegend/BatchAVIF.git

cd ./BatchAVIF

go build -o ./BatchAVIF main.go
```

## Basic understanding
- There're 3 steps on transcoding a/an video/image into AV1/AVIF:
  1. [Extractor] Decode the video/image into Y4M container
  2. [Encodoer] Encode to AV1/AVIF
  3. [Repackager] Repackaging to .mkv, .mp4, .avif... with the appropriate headers and metadata

## Usage
- Have `ffmpeg`, `MP4Box` and an (or multiple) encoder(s) of your choice in PATH or place them in the same folder where the compiled file is.
- Modify the `config.yaml` file to fit your needs.
- Run `./main` or `./main <your-config-file.yaml>` to start.

## (Maybe) FAQ
- Why there is a fallback encoder for animation?
  - `SVT-AV1`, the one I have been using to encode animation, is having this issue: [Odd image dimentions and svt-av1 encoder](https://github.com/AOMediaCodec/libavif/issues/544), so I decied to add a fallback encoder for it.
- Will there be a fallback encoder for image?
  - Yes (when I have the mood to do it).