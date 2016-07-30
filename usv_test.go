package usv

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"strings"
	"testing"
)

const TestTSVRow = "a\tb\tc\td\te\tf\tg\th\ti\tj\tk\tl\tm\tn\to\tp\tq\n"

const dataSize = 1 * (1 << 20) // 1 MiB

var TestInput = bytes.Repeat([]byte(TestTSVRow), dataSize/len(TestTSVRow))

func TestUSV(t *testing.T) {
	rs := strings.NewReader(TestTSVRow)
	ru := NewReader(rs)
	ru.RecordSeparator = '\n'
	ru.UnitSeparator = '\t'

	for {
		row, err := ru.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		// t.Log(row)
		if !bytes.Equal(bytes.Join(row, []byte("\t")), []byte(TestTSVRow)[:len(TestTSVRow)-1]) {
			log.Fatal("a != b")
		}
	}
}

func readUSV(r io.Reader) error {
	ru := NewReader(r)
	ru.RecordSeparator = '\n'
	ru.UnitSeparator = '\t'

	for {
		row, err := ru.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_ = row
	}
	return nil
}

func readCSV(r io.Reader) error {
	ru := csv.NewReader(r)
	ru.Comma = '\t'

	for {
		row, err := ru.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_ = row
	}
	return nil
}

func BenchmarkUSV(b *testing.B) {
	b.SetBytes(int64(len(TestInput)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		br := bytes.NewReader(TestInput)
		err := readUSV(br)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkCSV(b *testing.B) {
	b.SetBytes(int64(len(TestInput)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		br := bytes.NewReader(TestInput)
		err := readCSV(br)
		if err != nil {
			log.Fatal(err)
		}
	}
}
