package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
)

var (
	//means    = make(map[[100]byte]float32, BILLION)
	//nums     = make(map[[100]byte]float32)
	//minimums = make(map[[100]byte]float32, BILLION)
	//maximums = make(map[[100]byte]float32, BILLION)

	//currentStationIndex = int64(0)
	//stations            = make(map[[100]byte]int64, 10000)
	//stationLock         = sync.RWMutex{}
	//means               = make([]float32, 10000)
	//meanLock            = sync.RWMutex{}
	//nums                = make([]float32, 10000)
	//minimums            = make([]float32, 10000)
	//maximums            = make([]float32, 10000)
	//rangeLock           = sync.RWMutex{}
	stations             []map[[100]byte]int64
	means                [][]float32
	nums                 [][]float32 // change to int and include float conversion when doing avg calc
	minimums             [][]float32
	maximums             [][]float32
	XcurrentStationIndex = int64(0)

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

	stations = make([]map[[100]byte]int64, numThreads)
	means = make([][]float32, numThreads)
	nums = make([][]float32, numThreads)
	minimums = make([][]float32, numThreads)
	maximums = make([][]float32, numThreads)

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
		go Process(i, offset, perGRlimit+adj)
	}

	_ = file.Close()
	wg.Wait()

	for Y := int64(1); Y < numThreads; Y++ {
		for Ystation, YstationIndex := range stations[Y] {
			XstationIndex, XstationExists := stations[0][Ystation]
			if !XstationExists {
				XcurrentStationIndex++
				stations[0][Ystation] = XcurrentStationIndex
				XstationIndex = XcurrentStationIndex

				means[0][XstationIndex] = means[Y][YstationIndex]
				nums[0][XstationIndex] = nums[Y][YstationIndex]

				minimums[0][XstationIndex] = maximums[Y][YstationIndex]
				maximums[0][XstationIndex] = minimums[Y][YstationIndex]
			} else {
				XAvg := means[0][XstationIndex]
				XNum := nums[0][XstationIndex]
				YAvg := means[Y][YstationIndex]
				YNum := nums[Y][YstationIndex]

				means[0][XstationIndex] = ((XAvg * XNum) + (YAvg * YNum)) / (XNum + YNum)
				nums[0][XstationIndex] = XNum + YNum

				Xmin := minimums[0][XstationIndex]
				Xmax := maximums[0][XstationIndex]
				Ymin := minimums[Y][XstationIndex]
				Ymax := maximums[Y][XstationIndex]

				if Xmin > Ymin {
					minimums[0][XstationIndex] = Ymin
				}
				if Xmax < Ymax {
					maximums[0][XstationIndex] = Ymax
				}
			}
		}
	}

	var sortedStationNames = make([]string, len(stations[0]))
	var nameIndexMap = make(map[string]int64, len(stations[0]))
	keyIndex := 0
	for k, v := range stations[0] {
		firstIndexWithZero := 0
		for i := range k {
			if k[i] == 0x00 {
				firstIndexWithZero = i
				break
			}
		}
		str := string(k[:firstIndexWithZero])

		sortedStationNames[keyIndex] = str
		keyIndex++
		nameIndexMap[str] = v
	}
	sort.Strings(sortedStationNames)

	for i := range sortedStationNames {
		bKey := nameIndexMap[sortedStationNames[i]]
		fmt.Printf("%v;%.1f;%.1f;%.1f\n", sortedStationNames[i], minimums[0][bKey], means[0][bKey], maximums[0][bKey])
	}

	//fmt.Println(len(sortedStationNames))
}

func Process(arrIndex int64, offset int64, limit int64) {
	currentStationIndex := int64(0)
	stations[arrIndex] = make(map[[100]byte]int64, 10000)
	means[arrIndex] = make([]float32, 10000)
	nums[arrIndex] = make([]float32, 10000)
	minimums[arrIndex] = make([]float32, 10000)
	maximums[arrIndex] = make([]float32, 10000)

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

		stationIndex, stationExists := stations[arrIndex][station]
		if !stationExists {
			currentStationIndex++
			stations[arrIndex][station] = currentStationIndex
			stationIndex = currentStationIndex

			means[arrIndex][stationIndex] = temperatureFloat
			nums[arrIndex][stationIndex] = 1

			minimums[arrIndex][stationIndex] = temperatureFloat
			maximums[arrIndex][stationIndex] = temperatureFloat
		} else {
			prevAvg := means[arrIndex][stationIndex]
			prevNum := nums[arrIndex][stationIndex]
			means[arrIndex][stationIndex] = ((prevAvg * prevNum) + temperatureFloat) / (prevNum + 1)
			nums[arrIndex][stationIndex] = prevNum + 1

			prevMin := minimums[arrIndex][stationIndex]
			prevMax := maximums[arrIndex][stationIndex]
			if prevMin > temperatureFloat {
				minimums[arrIndex][stationIndex] = temperatureFloat
			}
			if prevMax < temperatureFloat {
				maximums[arrIndex][stationIndex] = temperatureFloat
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

	if arrIndex == 0 {
		XcurrentStationIndex = currentStationIndex
	}

	wg.Done()
	_ = file.Close()
}
