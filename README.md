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

## TODO list
- [x] Fully configurable for each step
- [x] Auto parse: `{{ input }}`, `{{ output }}`, `{{ width }}`, `{{ height }}`, `{{ threads }}`
- [x] Keep/remove original file
- [x] Keep/remove original extension
- [x] A fallback encoderfor both image and animation
- [ ] Show time taken to convert
- [ ] Multi-threading
- [ ] Piping input/output files to hide `y4m` and `ivf` files

## Usage
- Have `ffmpeg`, `MP4Box` and an (or multiple) encoder(s) of your choice in PATH or place them in the same folder where the compiled file is.
- Modify the `config.yaml` file to fit your needs.
- Run `./main` or `./main <your-config-file.yaml>` to start.

## (Maybe) FAQ
- Why there is a fallback encoder?
  - This option was specifically made for `SVT-AV1`: [Odd image dimentions and svt-av1 encoder](https://github.com/AOMediaCodec/libavif/issues/544).
