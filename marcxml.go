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
	"fmt"
	"io"
	"os"
	"template"
)


var record_xml = `  <record>
    {Leader}{.repeated section Fields}
    {@}{.end}
  </record>
`
var record_template *template.Template

var leader_xml = `<leader>{String}</leader>`
var leader_template *template.Template
var control_xml = `<control tag="{Tag}">{Data}</control>`
var control_template *template.Template
var data_xml = `<datafield tag="{Tag}" ind1="{Ind1}" ind2="{Ind2}">{.repeated section SubFields}
      {@}{.end}
    </datafield>`
var data_template *template.Template
var subfield_xml = `<subfield code="{Code}">{Value}</subfield>`
var subfield_template *template.Template

func init() {
	formatters := make(template.FormatterMap)
	formatters[""] = formatter
	leader_template = template.MustParse(leader_xml, formatters)
	control_template = template.MustParse(control_xml, formatters)
	data_template = template.MustParse(data_xml, formatters)
	subfield_template = template.MustParse(subfield_xml, formatters)
	record_template = template.MustParse(record_xml, formatters)
}

func formatter(writer io.Writer, format string, values ...interface{}) {
	for _, value := range values {
		switch field := value.(type) {
		case *Record:
			record_template.Execute(writer, field)
		case Record:
			record_template.Execute(writer, field)
		case *Leader:
			leader_template.Execute(writer, field)
		case Leader:
			leader_template.Execute(writer, field)
		case *ControlField:
			control_template.Execute(writer, field)
		case ControlField:
			control_template.Execute(writer, field)
		case *DataField:
			data_template.Execute(writer, field)
		case DataField:
			data_template.Execute(writer, field)
		case *SubField:
			subfield_template.Execute(writer, field)
		case SubField:
			subfield_template.Execute(writer, field)
		case byte:
			fmt.Fprintf(writer, "%c", field)
		default:
			template.HTMLFormatter(writer, format, values...)
		}
	}
}

// Write a MARC/XML representation of the record
func (record Record) XML(writer io.Writer) (err os.Error) {
	err = record_template.Execute(writer, record)
	return
}
