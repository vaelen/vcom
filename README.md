# VCom
A very small serial terminal written in Go.

Copyright 2019, Andrew C. Young

Released under the MIT license.

To build:
```bash
$ make
```

You can then use the [upx utility](https://upx.github.io/)
to make the binary even smaller if needed:
```bash
$ upx --brute vcom
```

If you would like to build binaries for all supported platforms,
including compression:
```bash
$ make release
```

If you would like to build binaries for all supported platforms,
without compression:
```bash
$ make releases
```
