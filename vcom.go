// Copyright 2018, Andrew C. Young
// License: MIT

// VCom is an extremely small and basic serial terminal.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jacobsa/go-serial/serial"
)

type converter func(in []byte, out []byte) []byte

func connect(device string, baudRate uint, dataBits uint, stopBits uint, skipConversion bool) {

	options := serial.OpenOptions{
		PortName:        device,
		BaudRate:        baudRate,
		DataBits:        dataBits,
		StopBits:        stopBits,
		MinimumReadSize: 1,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	log.Printf("Connected to %s at %dbps\n", device, baudRate)

	// Make sure to close it later.
	defer port.Close()

	done := make(chan bool)

	var fromConverter converter = fromUnix
	var toConverter converter = toUnix
	if skipConversion {
		fromConverter = noConversion
		toConverter = noConversion
	}

	// STDIN Reader
	go processingLoop(os.Stdin, port, done, fromConverter, "Console", "Serial port")

	// Serial Port Reader
	go processingLoop(port, os.Stdout, done, toConverter, "Serial port", "Console")

	<-done
}

func processingLoop(in io.Reader, out io.Writer, done chan bool, convert converter, readerName string, writerName string) {
	defer func() {
		e := recover()
		if e != nil {
			log.Printf("Unhandled exception caught: %v\n", e)
		}
		done <- true
	}()
	buf := make([]byte, 1024)
	outBuf := make([]byte, 2048)
	for {
		bytesRead, err := in.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("%s closed\n", readerName)
			} else {
				log.Printf("Error reading from %s: %v\n", strings.ToLower(readerName), err)
			}
			return
		}
		bytesToWrite := convert(buf[:bytesRead], outBuf)
		for len(bytesToWrite) > 0 {
			bytesWritten, err := out.Write(bytesToWrite)
			if err != nil {
				if err == io.EOF {
					log.Printf("%s closed\n", writerName)
				} else {
					log.Printf("Error reading from %s: %v\n", strings.ToLower(writerName), err)
				}
				return
			}
			bytesToWrite = bytesToWrite[bytesWritten:]
		}
	}
}

func noConversion(in []byte, out []byte) []byte {
	copy(out, in)
	return out[:len(in)]
}

func fromUnix(in []byte, out []byte) []byte {
	out = out[:0]
	for _, b := range in {
		if b == '\n' {
			out = append(out, '\r', '\n')
		} else {
			out = append(out, b)
		}
	}
	return out
}

func toUnix(in []byte, out []byte) []byte {
	out = out[:0]
	first := false
	for _, b := range in {
		if b == '\r' {
			if first {
				out = append(out, '\n')
			} else {
				first = true
			}
		} else if b == '\n' {
			out = append(out, '\n')
			first = false
		} else {
			if first {
				out = append(out, '\n')
				first = false
			}
			out = append(out, b)
		}
	}
	return out
}

func version() {
	fmt.Fprintln(flag.CommandLine.Output(), "VCom is a very small serial terminal.")
	fmt.Fprintln(flag.CommandLine.Output(), "Version 1.0")
	fmt.Fprintln(flag.CommandLine.Output(), "Copyright 2018, Andrew C. Young <andrew@vaelen.org>")
}

func usage() {
	fmt.Fprintln(flag.CommandLine.Output(), "Usage: vcom [options]")
	fmt.Fprintln(flag.CommandLine.Output(), "Options:")
	flag.PrintDefaults()
}

func main() {

	device := flag.String("d", "/dev/ttyS0", "serial device to use")
	baudRate := flag.Uint("b", 19200, "baud rate")
	dataBits := flag.Uint("data", 8, "data bits")
	stopBits := flag.Uint("stop", 1, "stop bits")
	skipConversion := flag.Bool("raw", false, "disable newline conversion")
	printVersion := flag.Bool("version", false, "print version information")

	flag.Usage = func() {
		version()
		fmt.Fprintln(flag.CommandLine.Output())
		usage()
	}

	flag.Parse()
	if flag.Parsed() {
		if *printVersion {
			version()
		} else {
			connect(*device, *baudRate, *dataBits, *stopBits, *skipConversion)
		}
	} else {
		flag.Usage()
	}

}