package marc21

/*
   Go Language MARC21 Library
   Copyright (C) 2011 William Waites

   This program is free software: you can redistribute it and/or
   modify it under the terms of the GNU Lesser General Public License
   as published by the Free Software Foundation, either version 3 of
   the License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU Lesser General Public
   License and the GNU General Public License along with this program
   (the files COPYING and GPL3 respectively).  If not, see
   <http://www.gnu.org/licenses/>.
*/

import (
    "bytes"
	"io"
	"log"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
    exp := 85
	data, err := os.Open("test.mrc")
	if err != nil {
		t.Fatal(err)
	}
	defer data.Close()

	count := 0
	for {
		_, err := ReadRecord(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		count++
		/*
			r.XML(os.Stdout)
			log.Print(record)
				buf := &bytes.Buffer{}
				err = record.XML(buf)
				if err != nil {
					log.Fatal(err)
				}
				log.Print(string(buf.Bytes()))
				break
		*/
	}
	if count != exp {
		t.Errorf("Expected to read %d records, got %d", exp, count)
	}
	log.Printf("%d records", count)
}

func TestToXML(t *testing.T) {
    const exp = `<record><leader>00819cam a2200289   4500</leader><controlfield tag="001">50001</controlfield><controlfield tag="005">20010903131819.0</controlfield><controlfield tag="008">701012s1970    moua     b    001 0 eng  </controlfield><datafield tag="010" ind1="32" ind2="32"><subfield code="97">   73117956 </subfield></datafield><datafield tag="035" ind1="32" ind2="32"><subfield code="97">ocm00094426 </subfield></datafield><datafield tag="035" ind1="32" ind2="32"><subfield code="57">7003024381</subfield></datafield><datafield tag="040" ind1="32" ind2="32"><subfield code="97">DLC</subfield><subfield code="99">DLC</subfield><subfield code="100">OKO</subfield></datafield><datafield tag="020" ind1="32" ind2="32"><subfield code="97">0801657024</subfield></datafield><datafield tag="050" ind1="48" ind2="48"><subfield code="97">RC78.7.C9</subfield><subfield code="98">Z83</subfield></datafield><datafield tag="060" ind1="32" ind2="32"><subfield code="97">QS 504 Z94d 1970</subfield></datafield><datafield tag="082" ind1="48" ind2="48"><subfield code="97">616.07/583</subfield></datafield><datafield tag="049" ind1="32" ind2="32"><subfield code="97">CUDA</subfield></datafield><datafield tag="100" ind1="49" ind2="32"><subfield code="97">Zugibe, Frederick T.</subfield><subfield code="113">(Frederick Thomas),</subfield><subfield code="100">1928-</subfield></datafield><datafield tag="245" ind1="49" ind2="48"><subfield code="97">Diagnostic histochemistry</subfield><subfield code="99">[by] Frederick T. Zugibe.</subfield></datafield><datafield tag="260" ind1="32" ind2="32"><subfield code="97">Saint Louis,</subfield><subfield code="98">Mosby,</subfield><subfield code="99">1970.</subfield></datafield><datafield tag="300" ind1="32" ind2="32"><subfield code="97">xiv, 366 p.</subfield><subfield code="98">illus.</subfield><subfield code="99">25 cm.</subfield></datafield><datafield tag="504" ind1="32" ind2="32"><subfield code="97">Bibliography: p. 332-349.</subfield></datafield><datafield tag="650" ind1="32" ind2="48"><subfield code="97">Cytodiagnosis.</subfield></datafield><datafield tag="650" ind1="32" ind2="48"><subfield code="97">Histochemistry</subfield><subfield code="120">Technique.</subfield></datafield><datafield tag="650" ind1="32" ind2="50"><subfield code="97">Histocytochemistry.</subfield></datafield><datafield tag="650" ind1="32" ind2="50"><subfield code="97">Histological Techniques.</subfield></datafield><datafield tag="994" ind1="32" ind2="32"><subfield code="97">92</subfield><subfield code="98">CUD</subfield></datafield></record>`

    data, err := os.Open("test.mrc")
    if err != nil {
        t.Fatal(err)
    }
    defer data.Close()
    r, err := ReadRecord(data)
    buf := &bytes.Buffer{}
    err = r.XML(buf)
    if err != nil {
        t.Fatal(err)
    }
    if exp != string(buf.Bytes()) {
        t.Errorf("Output XML did not match expected XML")
    }
    log.Printf("XML output")
}
