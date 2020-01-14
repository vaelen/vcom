// Copyright 2019, Andrew C. Young
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

func connect(device string, baudRate uint, dataBits uint, stopBits uint, eol string) {

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

	var fromConverter = eolConverter(eol)
	var toConverter = toUnix

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

func eolConverter(eol string) converter {
	eolBytes := []byte(eol)
	return func(in []byte, out []byte) []byte {
		out = out[:0]
		for _, b := range in {
			if b == '\n' {
				out = append(out, eolBytes...)
			} else {
				out = append(out, b)
			}
		}
		return out
	}
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
	fmt.Fprintln(flag.CommandLine.Output(), "VCom v1.1, a very small serial terminal.")
	fmt.Fprintln(flag.CommandLine.Output(), "Copyright 2019, Andrew C. Young <andrew@vaelen.org>")
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
	cr := flag.Bool("r", true, "send CR(\\r) for end of line")
	lf := flag.Bool("n", false, "send LF(\\n) for end of line")
	crlf := flag.Bool("rn", false, "send CRLF (\\r\\n) for end of line")
	printVersion := flag.Bool("version", false, "print version information")
	flag.Usage = func() {
		version()
		fmt.Fprintln(flag.CommandLine.Output())
		usage()
	}

	eol := ""
	if *cr {
		eol = "\r"
	}
	if *lf {
		eol = "\n"
	}
	if *crlf {
		eol = "\r\n"
	}

	flag.Parse()
	if flag.Parsed() {
		if *printVersion {
			version()
		} else {
			connect(*device, *baudRate, *dataBits, *stopBits, eol)
		}
	} else {
		flag.Usage()
	}

}
