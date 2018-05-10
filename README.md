# VCom
A very small serial terminal written in Go.

This project uses the [glide package manager](https://glide.sh/).

To build:
```bash
$ glide update
$ go build vcom.go
```

To build a smaller binary:
```
$ glide update
$ go build -ldflags="-s -w" vcom.go
```

You can then use the `upx` utility to make the binary even smaller if needed:
```
$ upx --brute vcom
```
