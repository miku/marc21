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

More examples
-------------

A [simple example](https://github.com/miku/marc21/blob/master/examples/simple.go).

```go
package main

import (
	"io"
	"log"
	"os"

	"github.com/miku/marc21"
)

func main() {
	for {
		r, err := marc21.ReadRecord(os.Stdin)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if _, err := r.WriteTo(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}
```

```
$ cat fixtures/sandburg.mrc | go run examples/simple.go | xmllint --format -
<?xml version="1.0"?>
<record>
  <leader>01142cam  2200301 a 4500</leader>
  <controlfield tag="001">   92005291 </controlfield>
  <controlfield tag="003">DLC</controlfield>
  <controlfield tag="005">19930521155141.9</controlfield>
  <controlfield tag="008">920219s1993    caua   j      000 0 eng  </controlfield>
  <datafield tag="010" ind1="32" ind2="32">
    <subfield code="97">   92005291 </subfield>
  </datafield>
  <datafield tag="020" ind1="32" ind2="32">
    <subfield code="97">0152038655 :</subfield>
    <subfield code="99">$15.95</subfield>
  </datafield>
  <datafield tag="040" ind1="32" ind2="32">
    <subfield code="97">DLC</subfield>
    <subfield code="99">DLC</subfield>
    <subfield code="100">DLC</subfield>
  </datafield>
  <datafield tag="042" ind1="32" ind2="32">
    <subfield code="97">lcac</subfield>
  </datafield>
  <datafield tag="050" ind1="48" ind2="48">
    <subfield code="97">PS3537.A618</subfield>
    <subfield code="98">A88 1993</subfield>
  </datafield>
  <datafield tag="082" ind1="48" ind2="48">
    <subfield code="97">811/.52</subfield>
    <subfield code="50">20</subfield>
  </datafield>
  <datafield tag="100" ind1="49" ind2="32">
    <subfield code="97">Sandburg, Carl,</subfield>
    <subfield code="100">1878-1967.</subfield>
  </datafield>
  <datafield tag="245" ind1="49" ind2="48">
    <subfield code="97">Arithmetic /</subfield>
    <subfield code="99">Carl Sandburg ; illustrated as an anamorphic adventure by Ted Rand.</subfield>
  </datafield>
  <datafield tag="250" ind1="32" ind2="32">
    <subfield code="97">1st ed.</subfield>
  </datafield>
  <datafield tag="260" ind1="32" ind2="32">
    <subfield code="97">San Diego :</subfield>
    <subfield code="98">Harcourt Brace Jovanovich,</subfield>
    <subfield code="99">c1993.</subfield>
  </datafield>
  <datafield tag="300" ind1="32" ind2="32">
    <subfield code="97">1 v. (unpaged) :</subfield>
    <subfield code="98">ill. (some col.) ;</subfield>
    <subfield code="99">26 cm.</subfield>
  </datafield>
  <datafield tag="500" ind1="32" ind2="32">
    <subfield code="97">One Mylar sheet included in pocket.</subfield>
  </datafield>
  <datafield tag="520" ind1="32" ind2="32">
    <subfield code="97">A poem about numbers and their characteristics. Features anamorphic, or distorted, drawings which can be restored to normal by viewing from a particular angle or by viewing the image's reflection in the provided Mylar cone.</subfield>
  </datafield>
  <datafield tag="650" ind1="32" ind2="48">
    <subfield code="97">Arithmetic</subfield>
    <subfield code="120">Juvenile poetry.</subfield>
  </datafield>
  <datafield tag="650" ind1="32" ind2="48">
    <subfield code="97">Children's poetry, American.</subfield>
  </datafield>
  <datafield tag="650" ind1="32" ind2="49">
    <subfield code="97">Arithmetic</subfield>
    <subfield code="120">Poetry.</subfield>
  </datafield>
  <datafield tag="650" ind1="32" ind2="49">
    <subfield code="97">American poetry.</subfield>
  </datafield>
  <datafield tag="650" ind1="32" ind2="49">
    <subfield code="97">Visual perception.</subfield>
  </datafield>
  <datafield tag="700" ind1="49" ind2="32">
    <subfield code="97">Rand, Ted,</subfield>
    <subfield code="101">ill.</subfield>
  </datafield>
</record>
```