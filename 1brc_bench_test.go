package main

import "testing"

func Benchmark1BRC(b *testing.B) {
	DoIt("javashit\\measurements.txt")
}

func Benchmark1BRCmini(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoIt("data.txt")
	}
}
