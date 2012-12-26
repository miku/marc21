/*
An IO library for Go to read and write MARC21 bibliographic catalogue records.
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

Usage is straightforward. For example,

    marcfile, err := os.Open("somedata.mrc")
    record, err := marc21.ReadRecord(marcfile)
    err = record.XML(os.Stdout)

*/
package marc21

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// A leader contains structural data about the MARC record
type Leader struct {
	Length                             int
	Status, Type                       byte
	ImplementationDefined              [5]byte
	CharacterEncoding                  byte
	BaseAddress                        int
	IndicatorCount, SubfieldCodeLength int
	LengthOfLength, LengthOfStartPos   int
}

// The field interface is satisfied by Control and Data field types.
type Field interface {
	String() string
	GetTag() string
}

func (leader Leader) Bytes() (buf []byte) {
	buf = make([]byte, 24)
	copy(buf[0:5], []byte(fmt.Sprintf("%05d", leader.Length)))
	buf[5] = leader.Status
	buf[6] = leader.Type
	copy(buf[7:9], leader.ImplementationDefined[0:2])
	buf[9] = leader.CharacterEncoding
	copy(buf[10:11], fmt.Sprintf("%d", leader.IndicatorCount))
	copy(buf[11:12], fmt.Sprintf("%d", leader.SubfieldCodeLength))
	copy(buf[12:17], fmt.Sprintf("%05d", leader.BaseAddress))
	copy(buf[17:20], leader.ImplementationDefined[2:5])
	copy(buf[20:21], fmt.Sprintf("%d", leader.LengthOfLength))
	copy(buf[21:22], fmt.Sprintf("%d", leader.LengthOfStartPos))
	buf[22] = '0'
	buf[23] = '0'
	return
}

func (leader Leader) String() string {
	return string(leader.Bytes())
}

func read_leader(reader io.Reader) (leader *Leader, err error) {
	data := make([]byte, 24)
	n, err := reader.Read(data)
	if err != nil {
		return
	}
	if n != 24 {
		errs := fmt.Sprintf("MARC21: invalid leader: expected 24 bytes, read %d", n)
		err = errors.New(errs)
		return
	}
	leader = &Leader{}
	leader.Length, err = strconv.Atoi(string(data[0:5]))
	if err != nil {
		errs := fmt.Sprintf("MARC21: invalid record length: %s", err)
		err = errors.New(errs)
		return
	}
	leader.Status = data[5]
	leader.Type = data[6]
	copy(leader.ImplementationDefined[0:2], data[7:9])
	leader.CharacterEncoding = data[9]

	leader.IndicatorCount, err = strconv.Atoi(string(data[10:11]))
	if err != nil || leader.IndicatorCount != 2 {
		errs := fmt.Sprintf("MARC21: erroneous indicator count, expecte '2', got %u", data[10])
		err = errors.New(errs)
		return
	}
	leader.SubfieldCodeLength, err = strconv.Atoi(string(data[11:12]))
	if err != nil || leader.SubfieldCodeLength != 2 {
		errs := fmt.Sprintf("MARC21: erroneous subfield code length, expected '2', got %u", data[11])
		err = errors.New(errs)
		return
	}

	leader.BaseAddress, err = strconv.Atoi(string(data[12:17]))

	if err != nil {
		errs := fmt.Sprintf("MARC21: invalid base address: %s", err)
		err = errors.New(errs)
		return
	}

	copy(leader.ImplementationDefined[2:5], data[17:20])

	leader.LengthOfLength, err = strconv.Atoi(string(data[20:21]))
	if err != nil || leader.LengthOfLength != 4 {
		errs := fmt.Sprintf("MARC21: invalid length of length, expected '4', got %u", data[20])
		err = errors.New(errs)
		return
	}
	leader.LengthOfStartPos, err = strconv.Atoi(string(data[21:22]))
	if err != nil || leader.LengthOfStartPos != 5 {
		errs := fmt.Sprintf("MARC21: invalid length of starting character position, expected '5', got %u", data[21])
		err = errors.New(errs)
		return
	}
	/*
		if data[22] != '0' {
			errs := fmt.Sprintf("MARC21: invalid length of implementation defined portion, expected '0', got %u", data[22])
			err = os.ErrorString(errs)
			return
		}
		if data[23] != '0' {
			errs := fmt.Sprintf("MARC21: invalid undefined in entry map, expected '0', got %u", data[23])
			err = os.ErrorString(errs)
			return
		}
	*/
	return
}

type dirent struct {
	tag          string
	length       int
	startCharPos int
}

const RT = 0x1D
const RS = 0x1E
const DELIM = 0x1F

var ERS = errors.New("Record Separator (field terminator)")

func read_dirent(reader io.Reader) (dent *dirent, err error) {
	data := make([]byte, 12)
	_, err = reader.Read(data[0:1])
	if err != nil {
		return
	}
	if data[0] == RS {
		err = ERS
		return
	}
	n, err := reader.Read(data[1:])
	if err != nil {
		return
	}
	if n != 11 {
		errs := fmt.Sprintf("MARC21: invalid directory entry, expected 12 bytes, got %d", n)
		err = errors.New(errs)
		return
	}
	dent = &dirent{}
	dent.tag = string(data[0:3])
	dent.length, err = strconv.Atoi(string(data[3:7]))
	if err != nil {
		return
	}
	dent.startCharPos, err = strconv.Atoi(string(data[7:12]))
	if err != nil {
		return
	}

	return
}

// A control field
type ControlField struct {
	XMLName xml.Name `xml:"controlfield"`
	Tag     string   `xml:"tag,attr"`
	Data    string   `xml:",chardata"`
}

func (cf *ControlField) String() string {
	return fmt.Sprintf("%s %s", cf.Tag, cf.Data)
}

// Returns the tag for a ControlField
func (cf *ControlField) GetTag() string {
	return cf.Tag
}

func read_control(reader io.Reader, dent *dirent) (field Field, err error) {
	data := make([]byte, dent.length)
	n, err := reader.Read(data)
	if err != nil {
		return
	}
	if n != dent.length {
		errs := fmt.Sprintf("MARC21: invalid control entry, expected %d bytes, read %d", dent.length, n)
		err = errors.New(errs)
		return
	}
	if data[dent.length-1] != RS {
		errs := fmt.Sprintf("MARC21: invalid control entry, does not end with a field terminator")
		err = errors.New(errs)
		return
	}
	field = &ControlField{Tag: dent.tag, Data: string(data[:dent.length-1])}
	return
}

// A subfield within a variable data field
type SubField struct {
	XMLName xml.Name `xml:"subfield"`
	Code    byte     `xml:"code,attr"`
	Value   string   `xml:",chardata"`
}

func (sf SubField) String() string {
	return fmt.Sprintf("(%c) %s", sf.Code, sf.Value)
}

// A variable data field
type DataField struct {
	XMLName   xml.Name `xml:"datafield"`
	Tag       string   `xml:"tag,attr"`
	Ind1      byte     `xml:"ind1,attr"`
	Ind2      byte     `xml:"ind2,attr"`
	SubFields []*SubField
}

// Returns the tag for a DataField 
func (df *DataField) GetTag() string {
	return df.Tag
}

func (df *DataField) String() string {
	subfields := make([]string, 0, len(df.SubFields))
	for _, sf := range df.SubFields {
		subfields = append(subfields, "["+sf.String()+"]")
	}
	return fmt.Sprintf("%s [%c%c] %s", df.Tag, df.Ind1, df.Ind2,
		strings.Join(subfields, ", "))
}

func read_data(reader io.Reader, dent *dirent) (field Field, err error) {
	data := make([]byte, dent.length)
	n, err := reader.Read(data)
	if err != nil {
		return
	}
	if n != dent.length {
		errs := fmt.Sprintf("MARC21: invalid data entry, expected %d bytes, read %d", dent.length, n)
		err = errors.New(errs)
		return
	}
	if data[dent.length-1] != RS {
		errs := fmt.Sprintf("MARC21: invalid data entry, does not end with a field terminator")
		err = errors.New(errs)
		return
	}

	df := &DataField{Tag: dent.tag}
	df.Ind1, df.Ind2 = data[0], data[1]

	df.SubFields = make([]*SubField, 0, 1)
	for _, sfbytes := range bytes.Split(data[2:dent.length-1], []byte{DELIM}) {
		if len(sfbytes) == 0 {
			continue
		}
		sf := &SubField{Code: sfbytes[0], Value: string(sfbytes[1:])}
		df.SubFields = append(df.SubFields, sf)
	}

	field = df
	return
}

// A MARC21 record consists of a leader and a number of fields
type Record struct {
	XMLName xml.Name `xml:"record"`
	Leader  *Leader  `xml:"leader"`
	Fields  []Field
}

func (record Record) String() string {
	estrings := make([]string, len(record.Fields))
	for i, entry := range record.Fields {
		estrings[i] = entry.String()
	}
	return strings.Join(estrings, "\n")
}

// Return the fields matching the given tag
func (record Record) GetFields(tag string) (fields []Field) {
	fields = make([]Field, 0, 4)
	for _, field := range record.Fields {
		if field.GetTag() == tag {
			fields = append(fields, field)
		}
	}
	return
}

// Return the subfields matching the given tag and code
func (record Record) GetSubFields(tag string, code byte) (subfields []*SubField) {
	subfields = make([]*SubField, 0, 4)
	fields := record.GetFields(tag)
	for _, field := range fields {
		switch data := field.(type) {
		case *DataField:
			for _, subfield := range data.SubFields {
				if subfield.Code == code {
					subfields = append(subfields, subfield)
				}
			}
		}
	}
	return
}

// Read a single MARC record from a reader.
func ReadRecord(reader io.Reader) (record *Record, err error) {
	record = &Record{}
	record.Fields = make([]Field, 0, 8)

	record.Leader, err = read_leader(reader)
	if err != nil {
		return
	}
	dents := make([]*dirent, 0, 8)
	for {
		var dent *dirent
		dent, err = read_dirent(reader)
		if err == ERS {
			err = nil
			break
		}
		if err != nil {
			return
		}
		dents = append(dents, dent)
	}

	for _, dent := range dents {
		var field Field
		if strings.HasPrefix(dent.tag, "00") {
			if field, err = read_control(reader, dent); err != nil {
				return
			}
		} else {
			if field, err = read_data(reader, dent); err != nil {
				return
			}
		}
		record.Fields = append(record.Fields, field)
	}
	rtbuf := make([]byte, 1)
	_, err = reader.Read(rtbuf)
	if err != nil {
		return
	}
	if rtbuf[0] != RT {
		err = errors.New("MARC21: could not read record terminator")
	}
	return
}

// A MARCXML record
type RecordXML struct {
	XMLName xml.Name `xml:"record"`
	Leader  string   `xml:"leader"`
	Fields  []Field
}

// Write a MARCXML representation of the record
func (record *Record) XML(writer io.Writer) (err error) {
	xmlrec := &RecordXML{Leader: record.Leader.String(), Fields: record.Fields}
	output, err := xml.Marshal(xmlrec)
	writer.Write(output)
	return
}
