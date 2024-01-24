# goflut
Pixelflut client written in Go

# Features

- Image: Multithreaded image drawing with scaling and position control.
- Video/Gif: Multithreaded real time video streaming with position control (no scaling you scale the video womp womp).
- Text: Text rendering with size, position and font control (ttf font files only for now).
- Wipe: Just wipes the canvas with a boring grey color.

# Installation

```
git clone https://github.com/vymvn/goflut

# compile to a temp binary and run
go run main.go

# Or compile to a persistant binary and run
go build
./goflut

# You can also just install the binary in your GOPATH (you would need to set a path for a font file when using text command)
go install https://github.com/vymvn/goflut@latest
goflut
```

# Usage

```
A humble pixelflut client

Usage:
  goflut [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  image       Image drawing mode.
  text        Text rendering mode
  video       Video streaming mode.
  wipe        Wipes the canvas.

Flags:
  -h, --help           help for goflut
  -H, --host string    The pixelflut server hostname or ip.
      --loop           Loops duh.
  -p, --port int       Server port.
      --threads int    Number of threads to use. (default 1)
  -x, --x-offset int   X axis offset.
  -y, --y-offset int   Y axis offset.

Use "goflut [command] --help" for more information about a command.
```

# TODO

- [ ] Pre-processed video frames.

# Known bugs

- Multithreaded video not as smooth as single thread.
- Text background needs better calculations.

