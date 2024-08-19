package main

import "testing"

func Benchmark1BRC(b *testing.B) {
	DoIt("javashit\\measurements.txt")
}

func Benchmark1BRCmini(b *testing.B) {
	DoIt("data.txt")
}
