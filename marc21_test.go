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
	"io"
	"log"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
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
		/*		log.Print(record)
				buf := &bytes.Buffer{}
				err = record.XML(buf)
				if err != nil {
					log.Fatal(err)
				}
				log.Print(string(buf.Bytes()))
				break
		*/
	}
	if count != 85 {
		t.Errorf("Expected to read 85 records, god %d", count)
	}
	log.Printf("%d records", count)
}
