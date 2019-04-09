package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/snickers54/omise-application/pkg/rot128"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
func main() {
	defer timeTrack(time.Now(), "main")
	PrintMemUsage()
	reader := rot128.NewRot128Reader("../../assets/fng.1000.csv.rot128")
	PrintMemUsage()
	defer reader.Close()
	for {
		line, ok := reader.Scan()
		if ok == false {
			break
		}
		fmt.Println(line)
	}
	PrintMemUsage()
}
