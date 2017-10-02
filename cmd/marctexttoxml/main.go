// marctexttoxml parses a textual marc format and converts it to XML.
//
// =LDR  00586ngm a2200000i 4500
// =001  JoVEBiology50
// =008  161201s2016\\\\xx\\\\\g\\\\\\s\\\\\eng\d
// =035  \\$a(OCoLC)733445860
// =040  \\$aVaAlASP$cVaAlASP
// =100  1\$aTrish,Erin
// =245  12$aFreezing Human ES Cells$h[electronic resource]
// =260  \\$aCambridge, MA$bMyJoVE Corp$c2016
// =300  \\$aonline resource (480 seconds)
// =490  0\$aBiology
// =500  \\$aTitle from resource description page
// =520  \\$aHere we demonstrate how our lab freezes HuES human embryonic stem cell lines.
// =521  \\$aFor undergraduate, graduate, and professional students
// =546  \\$aEnglish
// =650  \0$aBiology
// =856  40$uhttps://www.jove.com/video/50
// =945  \\$aThe Whole World
//
// Notes (from http://www.loc.gov/marc/makrbrkr.html):
//
// * An "=" (equal sign) in front of each field signals the start of a new field.
// * The blank space between the MARC tag and the rest of the field is to enhance
// readability.
// * The character "\" (the reverse solidus or "backslash") is used to represent
// the spaces that sometimes occupy indicator positions at the beginning of a MARC
// field or in a fixed field.
// * The first two examples of MARC fields have more than one subfield. The third
// example has only one subfield. The last example, a field beginning with "00"
// has no indicators or subfields.
// * The fill character "|" is for unused spaces in field 008 (since it must
// always be 40 characters long). Fill characters cannot be used in indicator
// positions.
//
package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/miku/marc21"
)

var (
	declaration     = `<?xml version="1.0" encoding="utf-8" ?>`
	subfieldPattern = regexp.MustCompile(`([$][a-z0-9])(.*?)`)
)

// stickyErrWriter keeps an error around, so you can *occasionally* check if an
// error occured.
type stickyErrWriter struct {
	w   io.Writer
	err *error
}

// Write writes the given bytes to the underlying writer. If the writer
// encountered an error already, this method does nothing.
func (sew stickyErrWriter) Write(p []byte) (n int, err error) {
	if *sew.err != nil {
		return 0, *sew.err
	}
	n, err = sew.w.Write(p)
	*sew.err = err
	return
}

// parseControlField parses bytes containing a control field.
func parseControlField(b []byte) (*marc21.ControlField, error) {
	return &marc21.ControlField{
		Tag:  string(b[1:4]),
		Data: string(b[6:]),
	}, nil
}

// decodeWhitespace decodes marc maker/breaker whitespace encoding.
func decodeWhitespace(b byte) byte {
	if b == 92 {
		return ' '
	}
	return b
}

// parseSubfields parses bytes only containing subfield information.
func parseSubfields(b []byte) (subfields []*marc21.SubField) {
	if len(b) == 0 {
		return
	}
	indices := subfieldPattern.FindAllIndex(b, -1)
	for i, index := range indices {
		if i < len(indices)-1 {
			subfields = append(subfields, &marc21.SubField{
				Code:  b[index[1]-1],
				Value: string(b[index[1]:indices[i+1][0]]),
			})
		} else {
			subfields = append(subfields, &marc21.SubField{
				Code:  b[index[1]-1],
				Value: string(b[index[1]:]),
			})
		}
	}
	return
}

// parseDataField parses a line containing a data field.
func parseDataField(b []byte) (*marc21.DataField, error) {
	field := &marc21.DataField{
		Tag:       string(b[1:4]),
		Ind1:      decodeWhitespace(b[6]),
		Ind2:      decodeWhitespace(b[7]),
		SubFields: parseSubfields(b[8:]),
	}
	return field, nil
}

// parseField dispatches parsing for control and datafields.
func parseField(b []byte) (marc21.Field, error) {
	switch string(b[1:4]) {
	case "001", "002", "003", "004", "005", "006", "007", "008", "009":
		return parseControlField(b)
	default:
		return parseDataField(b)
	}
}

// parseRecord parses a textual marc record.
func parseRecord(b []byte) *marc21.Record {
	record := &marc21.Record{}
	br := bufio.NewReader(bytes.NewReader(b))
	for {
		b, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		b = bytes.TrimSpace(b)
		if len(b) == 0 {
			continue
		}
		switch string(b[1:4]) {
		case "LDR":
			leader, err := marc21.ParseLeader(bytes.NewReader(b[6:]))
			if err != nil {
				log.Fatal(err)
			}
			record.Leader = leader
		default:
			f, err := parseField(b)
			if err != nil {
				continue
			}
			record.AddField(f)
		}
	}
	return record
}

func main() {
	br := bufio.NewReader(os.Stdin)
	fieldPattern := regexp.MustCompile(`(^=[0-9][0-9][0-9]|^=LDR)`)

	var buf bytes.Buffer
	var once sync.Once
	var err error

	w := &stickyErrWriter{os.Stdout, &err}

	for {
		b, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if fieldPattern.Match(b) {
			// This is true, when we have a complete record buffered.
			if len(buf.Bytes()) > 0 && bytes.HasPrefix(b, []byte("=LDR")) {
				record := parseRecord(buf.Bytes())
				once.Do(func() {
					io.WriteString(w, declaration)
					io.WriteString(w, `<collection xmlns="http://www.loc.gov/MARC21/slim">`)
				})
				record.WriteTo(w)
				buf.Reset()
			}
		}
		if _, err := buf.Write(b); err != nil {
			log.Fatal(err)
		}
	}
	io.WriteString(w, "</collection>\n")
	if *w.err != nil {
		log.Fatal(err)
	}
}
