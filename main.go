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
	//means    = make(map[[100]byte]float32, BILLION)
	//nums     = make(map[[100]byte]float32)
	//minimums = make(map[[100]byte]float32, BILLION)
	//maximums = make(map[[100]byte]float32, BILLION)

	currentStationIndex = int64(0)
	stations            = make(map[[100]byte]int64, 10000)
	stationLock         = sync.RWMutex{}
	means               = make([]float32, 10000)
	meanLock            = sync.RWMutex{}
	nums                = make([]float32, 10000)
	minimums            = make([]float32, 10000)
	maximums            = make([]float32, 10000)
	rangeLock           = sync.RWMutex{}

	lock = sync.RWMutex{}
	wg   = sync.WaitGroup{}
)

var FILENAME string

func main() {
	//DoIt("input.txt")
	DoIt("javashit\\measurements.txt")
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

	var keys = make([]string, len(stations))
	var keyByteStringMap = make(map[string]int64, len(stations))
	keyIndex := 0
	for k, v := range stations {
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
		keyByteStringMap[str] = v
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
	//src.Split(bufio.ScanLines)
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
				temperatureFloat += pv * float32(temperatureBytes[i]+(^byte(48)+1))
				pv *= 10
			}
		}
		if temperatureBytes[0] == 0x2D {
			temperatureFloat *= -1
		}

		//fmt.Println(offset, string(station[:]), string(temperatureBytes[:]), temperatureFloat)

		stationLock.RLock()
		stationIndex, stationExists := stations[station]
		stationLock.RUnlock()

		if !stationExists {
			stationLock.Lock()
			currentStationIndex++
			stations[station] = currentStationIndex
			stationIndex = currentStationIndex
			stationLock.Unlock()

			meanLock.Lock()
			means[stationIndex] = temperatureFloat
			nums[stationIndex] = 1
			meanLock.Unlock()

			//minLock.Lock()
			//minimums[stationIndex] = temperatureFloat
			//minLock.Unlock()
			//
			//maxLock.Lock()
			//maximums[stationIndex] = temperatureFloat
			//maxLock.Unlock()
			rangeLock.Lock()
			minimums[stationIndex] = temperatureFloat
			maximums[stationIndex] = temperatureFloat
			rangeLock.Unlock()
		} else {
			meanLock.RLock()
			prevAvg := means[stationIndex]
			prevNum := nums[stationIndex]
			meanLock.RUnlock()
			meanLock.Lock()
			means[stationIndex] = ((prevAvg * prevNum) + temperatureFloat) / (prevNum + 1)
			nums[stationIndex] = prevNum + 1
			meanLock.Unlock()

			//minLock.RLock()
			//prevMin := minimums[stationIndex]
			//minLock.RUnlock()
			//if prevMin > temperatureFloat {
			//	minLock.Lock()
			//	minimums[stationIndex] = temperatureFloat
			//	minLock.Unlock()
			//}
			//
			//maxLock.RLock()
			//prevMax := maximums[stationIndex]
			//maxLock.RUnlock()
			//if prevMax < temperatureFloat {
			//	maxLock.Lock()
			//	maximums[stationIndex] = temperatureFloat
			//	maxLock.Unlock()
			//}
			rangeLock.RLock()
			prevMin := minimums[stationIndex]
			prevMax := maximums[stationIndex]
			rangeLock.RUnlock()
			if prevMin > temperatureFloat {
				rangeLock.Lock()
				minimums[stationIndex] = temperatureFloat
				rangeLock.Unlock()
			}
			if prevMax < temperatureFloat {
				rangeLock.Lock()
				maximums[stationIndex] = temperatureFloat
				rangeLock.Unlock()
			}
		}

		//lock.RLock()
		//prevNum := nums[station]
		//prevAvg := means[station]
		//prevMax, okMax := maximums[station]
		//prevMin, okMin := minimums[station]
		//lock.RUnlock()
		//
		//lock.Lock()
		//means[station] = ((prevAvg * prevNum) + temperatureFloat) / (prevNum + 1)
		//nums[station] = prevNum + 1
		//lock.Unlock()
		//
		//if (okMax && prevMax < temperatureFloat) || !okMax {
		//	lock.Lock()
		//	maximums[station] = temperatureFloat
		//	lock.Unlock()
		//}
		//if (okMin && prevMin > temperatureFloat) || !okMin {
		//	lock.Lock()
		//	minimums[station] = temperatureFloat
		//	lock.Unlock()
		//}
	}

	wg.Done()
	_ = file.Close()
}
