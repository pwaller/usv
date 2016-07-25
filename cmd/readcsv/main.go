package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	var (
		unitSep = flag.Int("comma", int(','), "ASCII char to use for unit separator")
	)
	flag.Parse()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	defer func() {
		allocStart := m.Mallocs
		runtime.ReadMemStats(&m)
		allocEnd := m.Mallocs
		log.Printf("Mallocs: %v", allocEnd-allocStart)
	}()

	r := csv.NewReader(os.Stdin)
	r.FieldsPerRecord = -1
	r.Comma = rune(*unitSep)

	var n int
	start := time.Now()
	defer func() {
		dur := time.Since(start)
		log.Printf("Read %d rows in %v: %.0f/sec", n, dur, float64(n)/dur.Seconds())
	}()

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		_ = row
		n++
	}
}
