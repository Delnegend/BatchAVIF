# Batch AVIF

I'm no longer maintaining this project since this's just for introducing myself to Golang. I'm maintaining a similar one written in Python [here](https://github.com/delnegend/scripts)

---
## ✨ Features
- [x] ⚙️ Configurable presets
- [x] 🔧 Keep/remove original file(s)
- [x] 🔧 Keep/remove original extension
- [x] ✅ Skip/overwrite converted file(s)
- [x] ♾️ Dynamic variable in `config.yaml`
  - `{{ input }}`
  - `{{ output }}`
  - `{{ width }}`
  - `{{ height }}`
  - `{{ threads }}`
- [x] 🔙 Fallback encoder
- [x] 🧵 Multi-threading (thanks to [WoofinaS](https://github.com/WoofinaS/img2avif))
- [x] ⏱️ Timer for each conversion
- [x] 🔌 Piped/non-pipe mode (thanks to [WoofinaS](https://github.com/WoofinaS/img2avif))
- [x] 📃 Export log file (single threaded, non-pipe mode only)
- [ ] 🔔 Notification when finished

## 📖 You might wanna read
- There're 3 steps on transcoding a/an video/image into AV1/AVIF:
  1. [Extractor] decode the video/image into a .y4m file
  2. [Encoder] encode y4m to AV1/AVIF (usually) into an .ivf file
  3. [Repackager] repack the .ivf to .mkv, .mp4, .avif... with appropriate headers and metadata
- Piped/non-pipe mode?
  - Piped mode: `<file> -> extractor-encoder-repackager -> <file.avif>`
  - Non-pipe mode: `<file> -> extractor -> <file.y4m> -> encoder -> <file.ivf> -> repackager -> <file.avif>`
  - Pipe mode is guarantee to be faster but it only works on linux (and macOS idk, I don't have one to test)
  - To use pipe mode, set config > mode to "pipe"
  - To use non-pipe mode, set config > mode to "file"

## 🛠️ Build
```terminal
git clone https://github.com/Delnegend/BatchAVIF.git
cd ./BatchAVIF
go build -o ./BatchAVIF main.go
```

## 📕 Usage
- Have `ffmpeg`, `MP4Box` and an (or multiple) encoder(s) of your choice in PATH.
- Modify the `config.yaml` file to fit your needs.
- Run `./main` or `./main <your-config-file.yaml>` to start.

## ❓ Questions you might ask
- Fallback encoder?<br>
  This option was exist because of `SVT-AV1`: [Odd image dimentions and svt-av1 encoder](https://github.com/AOMediaCodec/libavif/issues/544).
- Cannot scale the source with `{{ width }}` and/or `{{ height }}` in pipe mode?
  - Without ffmpeg's scale, BatchAVIF can just simply replace `{{ width }}` and/or `{{ height }}` with the resolution it read from the original file.
  - With ffmpeg's scale, BatchAVIF must read from the extracted `.y4m` file instead, piping throws it directly into the encoding stage.
  - I already had a workaround but my battery rans out, I'm not gonna touch this project for quite a while.