marc21
======

A Go [MARC21](https://www.loc.gov/marc/bibliographic/) implementation.

```go
file, _ := os.Open("somedata.mrc")
record, _ := marc21.ReadRecord(file)
_ = record.WriteTo(os.Stdout)
```
