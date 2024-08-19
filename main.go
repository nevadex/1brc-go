package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	DoIt("input.txt")
}

func DoIt(fileName string) {
	//numThreads := runtime.NumCPU()-1 // exclude main thread

	file, _ := os.Open(fileName)
	src := bufio.NewScanner(file)
	src.Buffer(make([]byte, 15), 105)
	src.Split(bufio.ScanLines) // other option: 2 scanners, one scans for ; and \n and the split function takes the advance # and the second scanner then scans that many more and it processes it there
	// could also do split by rune and just advance and skip newline checking...

	var means = make(map[[100]byte]float32)
	var nums = make(map[[100]byte]float32)
	var minimums = make(map[[100]byte]float32)
	var maximums = make(map[[100]byte]float32)

	for {
		notEOF := src.Scan()
		if !notEOF {
			break
		}
		bytes := src.Bytes()

		var station [100]byte
		var temperatureBytes [5]byte
		iOffset := 0
		readingStation := true
		for i := range bytes {
			b := bytes[i]
			if b == 0x3b {
				readingStation = false
				iOffset = i + 1
				continue
			}

			if readingStation {
				station[i] = b
			} else {
				temperatureBytes[i-iOffset] = b
			}
		}
		//temperatureFloat := math.Float32frombits(binary.LittleEndian.Uint32(temperatureBytes))

		var temperatureFloat float32
		temperatureFloat += 0.1 * DigitConv(temperatureBytes[4])
		temperatureFloat += DigitConv(temperatureBytes[2])
		temperatureFloat += 10 * DigitConv(temperatureBytes[1])
		if temperatureBytes[0] == 0x2D || temperatureBytes[1] == 0x2D {
			temperatureFloat *= -1
		}

		//fmt.Println(string(temperatureBytes[:]), temperatureFloat)

		prevNum := nums[station]
		prevAvg := means[station]
		means[station] = ((prevAvg * prevNum) + temperatureFloat) / (prevNum + 1)
		nums[station] = prevNum + 1

		if prevMax, ok := maximums[station]; (ok && prevMax < temperatureFloat) || !ok {
			maximums[station] = temperatureFloat
		}
		if prevMin, ok := minimums[station]; (ok && prevMin > temperatureFloat) || !ok {
			minimums[station] = temperatureFloat
		}
	}

	go func() { _ = file.Close() }()

	var keys = make([]string, len(means))
	var keyByteStringMap = make(map[string][100]byte, len(means))
	keyIndex := 0
	for k := range means {
		firstIndexWithZero := 0
		for i := range k {
			if k[i] == 0 {
				firstIndexWithZero = i
				break
			}
		}
		str := string(k[:firstIndexWithZero])

		keys[keyIndex] = str
		keyIndex++
		keyByteStringMap[str] = k
	}
	sort.Strings(keys)

	for i := range keys {
		bKey := keyByteStringMap[keys[i]]
		fmt.Printf("%v;%f;%f;%f\n", keys[i], minimums[bKey], means[bKey], maximums[bKey])
	}
}

func DigitConv(b byte) float32 {
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
