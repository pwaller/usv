package main

import (
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/pwaller/usv"
)

func main() {
	var (
		unitSep   = flag.Int("unit", int(usv.UnitSeparator), "ASCII char to use for unit separator")
		recordSep = flag.Int("record", int(usv.RecordSeparator), "ASCII char to use for record separator")
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

	r := usv.NewReader(os.Stdin)
	r.UnitSeparator = byte(*unitSep)
	r.RecordSeparator = byte(*recordSep)

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
