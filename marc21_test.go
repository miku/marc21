package marc21

import (
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
		if err == os.EOF {
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
