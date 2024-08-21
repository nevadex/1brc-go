package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func Benchmark1BRC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoIt("javashit\\measurements.txt")
	}
}

func Benchmark1BRC_15K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoIt("data.txt")
	}
}

func Benchmark1BRC_1M(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoIt("million.txt")
	}
}

func Benchmark1BRCtiny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoIt("tiny.txt")
	}
}

func BenchmarkParse_Buffer1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//file, _ := os.Open("data.txt")
		file, _ := os.Open("javashit\\measurements.txt")

		file.Seek(0, 2)

		src := bufio.NewScanner(file)

		fmt.Println(src.Scan())
		fmt.Println(src.Bytes())
		continue
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

		info, _ := file.Stat()
		indexlen := info.Size() - 1
		var bytes [1]byte
		fmt.Println(file.ReadAt(bytes[:], indexlen))
		fmt.Println(bytes, string(bytes[:]))
		continue

		for {
			notEOF := src.Scan()
			if !notEOF {
				break
			}
			_ = src.Bytes()
		}
	}
}
