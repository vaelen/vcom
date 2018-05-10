# VCom
A very small serial terminal written in Go.

This project uses the [glide package manager](https://glide.sh/).

To build:
```bash
$ glide update
$ go build vcom.go
```

To build a smaller binary:
```bash
$ glide update
$ go build -ldflags="-s -w" vcom.go
```

You can then use the [upx utility](https://upx.github.io/) to make the binary even smaller if needed:
```bash
$ upx --brute vcom
```
