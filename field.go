package marc21

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Field defines an interface that is satisfied by the Control and Data field
// types.
type Field interface {
	String() string
	GetTag() string
}

// ControlField represents a control field, which contains only a tag and data.
type ControlField struct {
	Tag  string `xml:"tag,attr"`
	Data string `xml:",chardata"`
}

// String returns the ControlField as a string.
func (cf *ControlField) String() string {
	return fmt.Sprintf("%s %s", cf.Tag, cf.Data)
}

// GetTag returns the tag for a ControlField.
func (cf *ControlField) GetTag() string {
	return cf.Tag
}

func readControl(reader io.Reader, dent *dirent) (field Field, err error) {
	data := make([]byte, dent.length)
	n, err := reader.Read(data)
	if err != nil {
		return
	}
	if n != dent.length {
		err = fmt.Errorf("MARC21: invalid control entry, expected %d bytes, read %d", dent.length, n)
		return
	}
	if data[dent.length-1] != RS {
		err = fmt.Errorf("MARC21: invalid control entry, does not end with a field terminator")
		return
	}
	field = &ControlField{Tag: dent.tag, Data: string(data[:dent.length-1])}
	return
}

// SubField represents a subfield, containing a single-byte code and
// associated data.
type SubField struct {
	Code  byte   `xml:"code,attr"`
	Value string `xml:",chardata"`
}

// String returns the subfield as a string.
func (sf SubField) String() string {
	return fmt.Sprintf("(%c) %s", sf.Code, sf.Value)
}

// MarshalXML customized XML serialization.
func (sf *SubField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "subfield"}
	start.Attr = []xml.Attr{
		xml.Attr{Name: xml.Name{Local: "code"}, Value: string(sf.Code)},
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	data := xml.CharData([]byte(sf.Value))
	if err := e.EncodeToken(data); err != nil {
		return err
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// DataField represents a variable data field, containing a tag, two
// single-byte indicators, and one or more subfields.
type DataField struct {
	XMLName   xml.Name `xml:"datafield"`
	Tag       string   `xml:"tag,attr"`
	Ind1      byte     `xml:"ind1,attr"`
	Ind2      byte     `xml:"ind2,attr"`
	SubFields []*SubField
}

// MarshalXML customized XML serialization.
func (df *DataField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "datafield"}
	start.Attr = []xml.Attr{
		xml.Attr{Name: xml.Name{Local: "tag"}, Value: df.Tag},
		xml.Attr{Name: xml.Name{Local: "ind1"}, Value: string(df.Ind1)},
		xml.Attr{Name: xml.Name{Local: "ind2"}, Value: string(df.Ind2)},
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := e.Encode(df.SubFields); err != nil {
		return err
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// GetTag returns the tag for a DataField.
func (df *DataField) GetTag() string {
	return df.Tag
}

// String returns the DataField as a string.
func (df *DataField) String() string {
	subfields := make([]string, 0, len(df.SubFields))
	for _, sf := range df.SubFields {
		subfields = append(subfields, "["+sf.String()+"]")
	}
	return fmt.Sprintf("%s [%c%c] %s", df.Tag, df.Ind1, df.Ind2,
		strings.Join(subfields, ", "))
}

func readData(reader io.Reader, dent *dirent) (field Field, err error) {
	data := make([]byte, dent.length)
	n, err := reader.Read(data)
	if err != nil {
		return
	}
	if n != dent.length {
		err = fmt.Errorf("MARC21: invalid data entry, expected %d bytes, read %d", dent.length, n)
		return
	}
	if data[dent.length-1] != RS {
		err = fmt.Errorf("MARC21: invalid data entry, does not end with a field terminator")
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
