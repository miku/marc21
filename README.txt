PACKAGE

package marc21
import "bitbucket.org/ww/marc21"

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

	.hg
	marc2xml
