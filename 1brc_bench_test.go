package main

import (
	"bufio"
	"os"
	"testing"
)

func Benchmark1BRC(b *testing.B) {
	DoIt("javashit\\measurements.txt")
}

func Benchmark1BRCmini(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoIt("data.txt")
	}
}

func BenchmarkParse_Buffer1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//file, _ := os.Open("data.txt")
		file, _ := os.Open("javashit\\measurements.txt")
		src := bufio.NewScanner(file)
		//src.Buffer(make([]byte, 4096), 8192)
		//src.Buffer(make([]byte, 8192), 16384)
		//src.Buffer(make([]byte, 16384), 32768)
		//src.Buffer(make([]byte, 32768), 65536)
		//src.Buffer(make([]byte, 65536), 131072)
		//src.Buffer(make([]byte, 131072), 262144)
		//src.Buffer(make([]byte, 262144), 524288)
		//src.Buffer(make([]byte, 524288), 1048576)
		//src.Buffer(make([]byte, 1048576), 2097152)
		src.Buffer(make([]byte, 2097152), 4194304)
		src.Split(bufio.ScanLines) // other option: 2 scanners, one scans for ; and \n and the split function takes the advance # and the second scanner then scans that many more and it processes it there
		// could also do split by rune and just advance and skip newline checking...

		for {
			notEOF := src.Scan()
			if !notEOF {
				break
			}
			_ = src.Bytes()
		}
	}
}
