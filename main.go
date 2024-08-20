package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
)

const BILLION = 15000 //1000000000

var (
	means    = make(map[[100]byte]float32, BILLION)
	nums     = make(map[[100]byte]float32)
	minimums = make(map[[100]byte]float32, BILLION)
	maximums = make(map[[100]byte]float32, BILLION)
	lock     = sync.RWMutex{}
	wg       = sync.WaitGroup{}
)

var FILENAME string

func main() {
	DoIt("input.txt")
}

// DoIt
// [ ALL THIS MONEY ON THE FLOOR, TEN RACKS THROW IT UP- WATCH HOW I DO IT ]
func DoIt(fn string) {
	FILENAME = fn
	numThreads := int64(runtime.NumCPU())

	file, _ := os.Open(FILENAME)
	info, _ := file.Stat()
	perGRlimit := info.Size() / numThreads
	totalAdjustment := int64(0)

	for i := int64(0); i < numThreads; i++ {
		offset := (i * perGRlimit) + totalAdjustment
		adj := int64(0)
		for {
			b := make([]byte, 1)
			_, _ = file.ReadAt(b, offset+perGRlimit+adj)
			adj++
			if b[0] == 0x0A || b[0] == 0x00 {
				break
			}
		}
		totalAdjustment += adj

		//fmt.Println("THREAD", i, offset, perGRlimit+adj)
		wg.Add(1)
		go Process(offset, perGRlimit+adj)
	}

	_ = file.Close()
	wg.Wait()

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
		fmt.Printf("%v;%.1f;%.1f;%.1f\n", keys[i], minimums[bKey], means[bKey], maximums[bKey])
	}

	fmt.Println(len(keys))
}

func Process(offset int64, limit int64) {
	file, _ := os.Open(FILENAME)
	_, _ = file.Seek(offset, 0)
	src := bufio.NewScanner(file)
	//src.Buffer(make([]byte, 15), 105)
	//src.Buffer(make([]byte, 8192), 16384)
	src.Buffer(make([]byte, 2097152), 4194304)
	src.Split(bufio.ScanLines)
	var currentPos int64 = 0 // EOF must have CRLF, change to -2 if the generated data does not add CRLF at the end

	for {
		notEOF := src.Scan()
		if !notEOF {
			//fmt.Println(offset, "EOF")
			break
		}
		bytes := src.Bytes()
		currentPos += int64(len(bytes) + 2) // add back CRLF for counting purposes
		if currentPos > limit+2 {
			//fmt.Println(offset, "LIMIT", currentPos, limit, string(bytes))
			break
		}

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
		var temperatureFloat float32

		pv := float32(0.1)
		for i := 4; i >= 0; i-- {
			if temperatureBytes[i] != 0x00 && temperatureBytes[i] != 0x2D && temperatureBytes[i] != 0x2E {
				temperatureFloat += pv * DigitConv(temperatureBytes[i])
				pv *= 10
			}
		}
		if temperatureBytes[0] == 0x2D {
			temperatureFloat *= -1
		}

		//fmt.Println(offset, string(station[:]), string(temperatureBytes[:]), temperatureFloat)

		lock.RLock()
		prevNum := nums[station]
		prevAvg := means[station]
		prevMax, okMax := maximums[station]
		prevMin, okMin := minimums[station]
		lock.RUnlock()

		lock.Lock()
		means[station] = ((prevAvg * prevNum) + temperatureFloat) / (prevNum + 1)
		nums[station] = prevNum + 1
		lock.Unlock()

		if (okMax && prevMax < temperatureFloat) || !okMax {
			lock.Lock()
			maximums[station] = temperatureFloat
			lock.Unlock()
		}
		if (okMin && prevMin > temperatureFloat) || !okMin {
			lock.Lock()
			minimums[station] = temperatureFloat
			lock.Unlock()
		}
	}

	wg.Done()
	_ = file.Close()
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
