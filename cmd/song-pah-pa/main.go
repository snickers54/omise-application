package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	// will load .env file on import, remove boilerplate code in our main
	_ "github.com/joho/godotenv/autoload"
	"github.com/snickers54/omise-application/internal"
	"github.com/snickers54/omise-application/internal/utils"
	"github.com/snickers54/omise-application/pkg/rot128"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	// want to time track, to see average performance
	defer utils.TimeTrack(time.Now(), "main")
	defer utils.PrintMemUsage()
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments, please provide the CSV file ciphered with ROT128 to consume.")
	}
	reader := rot128.NewRot128Reader(os.Args[1]) //"../../assets/fng.1000.csv.rot128"
	defer reader.Close()

	workerPool := internal.NewWorkerPool()

	workerPool.Wg.Add(workerPool.NbWorker)
	for i := 0; i < workerPool.NbWorker; i++ {
		go workerPool.Run(i)
	}
	// main loop, for each line of the file w
	reader.Scan() // we remove the first line containing the headers
	for {
		line, ok := reader.Scan()
		if ok == false {
			break
		}
		row := strings.Split(line, ",")
		// we use a pointer to avoid making copy and therefore using tons of memory for nothing
		// since the ptr here will outlive the goroutine, it should be fine
		workerPool.ChannelWork <- &row
	}
	workerPool.Close()
	workerPool.Wg.Wait()
	fmt.Println(workerPool.Stats.String())
}
