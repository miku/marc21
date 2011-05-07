PACKAGE

package marc21
import "bitbucket.org/ww/marc21"

An IO library for Go to read and write MARC21 bibliographic catalogue records.
Written in 2011 by William Waites <ww@styx.org>.


CONSTANTS

const DELIM = 0x1F

const RS = 0x1E

const RT = 0x1D


VARIABLES

var ERS os.ErrorString = "Record Separator (field terminator)"


TYPES

type ControlField struct {
    Tag  string
    Data string
}
A control field

func (cf *ControlField) GetTag() string

func (cf *ControlField) String() string

type DataField struct {
    Tag        string
    Ind1, Ind2 byte
    SubFields  []*SubField
}
A variable data field

func (df *DataField) GetTag() string

func (df *DataField) String() string

type Field interface {
    String() string
    GetTag() string
}
The field interface is satisfied by Control and Data field types.

type Leader struct {
    Length                             int
    Status, Type                       byte
    ImplementationDefined              [5]byte
    CharacterEncoding                  byte
    BaseAddress                        int
    IndicatorCount, SubfieldCodeLength int
    LengthOfLength, LengthOfStartPos   int
}
A leader contains structural data about the MARC record

func (leader Leader) Bytes() (buf []byte)

func (leader Leader) String() string

type Record struct {
    Leader *Leader
    Fields []Field
}
A MARC21 record consists in a leader and a number of fields

func ReadRecord(reader io.Reader) (record *Record, err os.Error)
Read a single MARC record from a reader.

func (record Record) GetFields(tag string) (fields []Field)
Return the fields matching the given tag and code

func (record Record) GetSubFields(tag string, code byte) (subfields []*SubField)
Return the sub-fields matching the given tag and code

func (record Record) String() string

func (record Record) XML(writer io.Writer) (err os.Error)
Write a MARC/XML representation of the record

type SubField struct {
    Code  byte
    Value string
}
A subfield within a variable data field

func (sf SubField) String() string


SUBDIRECTORIES

	marc2xml
