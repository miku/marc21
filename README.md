marc21
======

A Go [MARC21](https://www.loc.gov/marc/bibliographic/) implementation.

This repository started as a fork of
[gitorious.org/marc21-go/marc21](https://gitorious.org/marc21-go/marc21), but
is - as of June 2017 - the new home of the MARC21 Go library.

Usage
-----

```go
file, _ := os.Open("somedata.mrc")
record, _ := marc21.ReadRecord(file)
_ = record.WriteTo(os.Stdout)
```
