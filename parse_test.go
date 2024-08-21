package main

import (
	"math/rand"
	"strconv"
	"testing"
)

func IntConvMyWay() {
	bytes := []byte{0x38, 0x35}
	integer := 0
	for i := range bytes {
		num := 0
		pow := len(bytes) - 1 - i
		if bytes[i] == 0x38 {
			num = 8
		} else if bytes[i] == 0x35 {
			num = 5
		}

		if pow == 0 {
			integer += num
		} else if pow == 1 {
			integer += num * 10
		}
	}
}

func BenchmarkParse_IntConvMyWay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntConvMyWay()
	}
}

func IntConvGo() {
	bytes := []byte{0x38, 0x35}
	strconv.Atoi(string(bytes))
}

func BenchmarkParse_IntConvGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntConvGo()
	}
}

func DigitConv1(b byte) float32 {
	if b == 0x31 {
		return 1
	} else if b == 0x32 {
		return 2
	} else if b == 0x33 {
		return 3
	} else if b == 0x34 {
		return 4
	} else if b == 0x35 {
		return 5
	} else if b == 0x36 {
		return 6
	} else if b == 0x37 {
		return 7
	} else if b == 0x38 {
		return 8
	} else if b == 0x39 {
		return 9
	} else {
		return 0
	}
}
func BenchmarkParse_DigitConv1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DigitConv1(0x39)
	}
}
func DigitConv2(b byte) float32 {
	if b == 0x2D {
		return 0
	}
	return float32(b + (^byte(48) + 1))
}
func BenchmarkParse_DigitConv2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DigitConv2(0x39)
	}
}

func FloatConvMyWay(binary [5]byte) float32 {
	//binary := [5]byte{0x2D, 0x39, 0x39, 0x2E, 0x39}
	var float float32

	float += 0.1 * DigitConv1(binary[4])
	float += DigitConv1(binary[2])
	float += 10 * DigitConv1(binary[1])

	if binary[0] == 0x2D || binary[1] == 0x2D {
		float *= -1
	}

	return float
}

func BenchmarkParse_FloatConvMyWay(b *testing.B) {
	//fmt.Println(FloatConvMyWay([5]byte{0x2D, 0x39, 0x39, 0x2E, 0x39}))
	//fmt.Println(FloatConvMyWay([5]byte{0x00, 0x39, 0x39, 0x2E, 0x39}))
	//fmt.Println(FloatConvMyWay([5]byte{0x00, 0x2D, 0x39, 0x2E, 0x39}))
	//fmt.Println(FloatConvMyWay([5]byte{0x00, 0x00, 0x39, 0x2E, 0x39}))
	//fmt.Println(FloatConvMyWay([5]byte{0x00, 0x2D, 0x00, 0x2E, 0x39}))

	for i := 0; i < b.N; i++ {
		FloatConvMyWay([5]byte{0x2D, 0x39, 0x39, 0x2E, 0x39})
	}
}

func FloatConvGo(binary [5]byte) float32 {
	float, _ := strconv.ParseFloat(string(binary[:]), 32)

	return float32(float)
}

func BenchmarkParse_FloatConvGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FloatConvGo([5]byte{0x2D, 0x39, 0x39, 0x2E, 0x39})
	}
}

func BenchmarkParse_AllocArray(b *testing.B) {
	arrKeys := make([]int, b.N)
	arrValues := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		arrKeys[i] = i
		arrValues[i] = rand.Int()
		_ = arrKeys[i]
		_ = arrValues[i]
	}
}
func BenchmarkParse_ArrayAppend(b *testing.B) {
	var arrKeys []int
	var arrValues []int
	for i := 0; i < b.N; i++ {
		arrKeys = append(arrKeys, i)
		arrValues = append(arrValues, rand.Int())
		_ = arrKeys[i]
		_ = arrValues[i]
	}
}
func BenchmarkParse_Map(b *testing.B) {
	arrKeys := make(map[int]int, b.N)
	for i := 0; i < b.N; i++ {
		arrKeys[i] = rand.Int()
		_ = arrKeys[i]
	}
}
