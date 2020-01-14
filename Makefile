LDFLAGS = "-w"

releases = \
	bin/vcom-linux-386 bin/vcom-linux-amd64 \
	bin/vcom-linux-arm bin/vcom-linux-arm64 \
	bin/vcom-windows-386.exe bin/vcom-windows-amd64.exe \
	bin/vcom-darwin-amd64

word-split = $(word $2,$(subst -, , $(subst ., ,$1)))

all: vcom

release: releases compress

releases: $(releases)

compress:
	upx --brute bin/vcom-*-*

vcom: vcom.go
	go build

clean:
	rm -rf vcom bin

bin:
	@mkdir -p bin

bin/vcom-%: vcom.go | bin
	GOOS=$(call word-split,$@,2) GOARCH=$(call word-split,$@,3) go build -ldflags="$(LDFLAGS)" -o $@ 
