# USV: A parser for Unit Separator separated values

An efficient (zero-alloc) parser for CSV-like data.

Why use USV? Well, ASCII has a unit separator, which is unlikely to appear in
your data. If you can avoid having fields which need quoting, you can parse it
much more efficiently (typically 5x or more).

This parser can also be used with TSV, for example:

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
