// Copyright 2016 Peter Waller <p@pwaller.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package usv

import (
	"bufio"
	"io"
)

// RecordSeparator as defined in ASCII
const (
	RecordSeparator byte = 0x1e
	UnitSeparator   byte = 0x1f
)

// Reader implements a unit-separator separated values reader.
// It works for any data where the field separator and the record separator
// cannot appear in the data.
// It avoids allocations for performance. It's around 6x faster than the
// builtin CSV reader, at the expense of not being able to handle quoted fields.
type Reader struct {
	r *bufio.Reader

	RecordSeparator, UnitSeparator byte

	records [][]byte
}

// Read reads one line of a CSV into bss.
func (r *Reader) Read() ([][]byte, error) {
	return r.ReadInto(&r.records)
}

// ReadInto reads the next line into buf without allocating.
func (r *Reader) ReadInto(buf *[][]byte) ([][]byte, error) {
	line, err := r.r.ReadSlice(r.RecordSeparator)
	if err != nil {
		return nil, err
	}
	if len(r.records) > 0 {
		reset(buf)
	}

	line = line[:len(line)-1] // Remove record separator

	// Allocate storage for column 'i' if it doesn't already have it.
	ensureCol := func(i int) {
		if i < len(*buf) {
			// *buf is wide enough.
			return
		}
		if i >= cap(*buf) {
			// Allocate.
			for i >= len(*buf) {
				*buf = append(*buf, []byte{})
			}
			return
		}
		// Reuse previously allocated.
		*buf = (*buf)[:i+1]
	}

	col := 0
	for _, b := range line {
		if b == r.UnitSeparator {
			col++
			continue
		}
		ensureCol(col)
		(*buf)[col] = append((*buf)[col], b)
	}
	return *buf, nil

}

func reset(bss *[][]byte) {
	for i := range *bss {
		(*bss)[i] = (*bss)[i][:0]
	}
	(*bss) = (*bss)[:0]
}

// NewReader constructs a Reader
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:               ensureBuffered(r),
		UnitSeparator:   UnitSeparator,
		RecordSeparator: RecordSeparator,
		records:         nil,
	}
}

// TSV makes Reader read TSV files.
func (r *Reader) TSV() *Reader {
	r.RecordSeparator = '\n'
	r.UnitSeparator = '\t'
	return r
}

// CSV makes Reader read simple (no spacing, no quoting) CSV files.
func (r *Reader) CSV() *Reader {
	r.RecordSeparator = '\n'
	r.UnitSeparator = ','
	return r
}

func (r *Reader) Skip(n int) *Reader {
	for i := 0; i < n; i++ {
		// The error can be raised later.
		_, _ = r.Read()
	}
	return r
}

func ensureBuffered(r io.Reader) *bufio.Reader {
	if bufR, ok := r.(*bufio.Reader); ok {
		return bufR
	}
	const MegaByte = 1 << 20
	return bufio.NewReaderSize(r, MegaByte)
}
