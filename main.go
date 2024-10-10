package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Stats struct {
	files uint64
}

func createFiles(filenumber <-chan int, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	for filenumber := range filenumber {
		filename := fmt.Sprintf("%d.bin", filenumber)
		fp, err := os.Create(filename)
		if err != nil {
			log.Print(err)
		} else {
			fp.Close()
			atomic.AddUint64(&stats.files, 1)
		}
	}
}

func main() {
	workers := flag.Int("w", runtime.NumCPU(), "number of workers")
	startNum := flag.Int("s", 0, "starting number")
	count := flag.Int("f", 0, "number files to create")
	flag.Parse()

	var stats Stats

	filenumbers := make(chan int)

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go createFiles(filenumbers, &stats, &wg)
	}

	go func() {
		for i := *startNum; i < (*startNum + *count); i++ {
			filenumbers <- i
		}
	}()

	go func() {
		var last uint64
		for {
			time.Sleep(2 * time.Second)
			curr := atomic.LoadUint64(&stats.files)
			created_per_s := (curr - last) / 2
			fmt.Printf("files created: %12d | files/s: %12d\n", curr, created_per_s)
			last = curr
		}
	}()

	wg.Wait()
}
