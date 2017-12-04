# NOTICE: do not use this

[Recent versions of Go (>=1.9)](https://github.com/golang/go/commit/2181653be637cdcc7a6efee8ec0a719df1d83c00) now have [(encoding/csv.Reader).ReuseRecord](https://tip.golang.org/pkg/encoding/csv/#Reader.ReuseRecord) which causes far fewer allocations. This makes USV less of a win than it used to be (maybe 1-2x, rather than 5x), and I would advise using the vanilla reader if possible.

# USV: A parser for Unit Separator separated values

An efficient (zero-alloc) parser for CSV-like data.

Why use USV? Well, ASCII has a unit separator, which is unlikely to appear in
your data. If you can avoid having fields which need quoting, you can parse it
much more efficiently (typically 5x or more).

This parser can also be used with TSV, so long as the tab does not appear in
the data, for example:

```go
r := usv.NewReader(os.Stdin)
r.RecordSeparator = '\n'
r.UnitSeparator = '\t'

for {
	row, err := r.Read()
	if err == io.EOF {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
    // This prints bytes.
    log.Println(row)
}
```

# License

BSD 3-clause. See LICENSE.
