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
	"text/template"
)

func formatter(writer io.Writer, format string, values ...interface{}) {
    record := template.Must(template.New("record").Parse(`  <record>
    {Leader}{.repeated section Fields}
    {@}{.end}
  </record>
`))
    ldr := template.Must(template.New("ldr").Parse(`<leader>{String}</leader>`))
    cf := template.Must(template.New("cf").Parse(`<control tag="{Tag}">{Data}</control>`))
    df := template.Must(template.New("df").Parse(`<datafield tag="{Tag}" ind1="{Ind1}" ind2="{Ind2}">{.repeated section SubFields}
          {@}{.end}
        </datafield>`))
    sf := template.Must(template.New("sf").Parse(`<subfield code="{Code}">{Value}</subfield>`))

	for _, value := range values {
		switch field := value.(type) {
		case *Record:
			record.Execute(writer, field)
		case Record:
			record.Execute(writer, field)
		case *Leader:
			ldr.Execute(writer, field)
		case Leader:
			ldr.Execute(writer, field)
		case *ControlField:
			cf.Execute(writer, field)
		case ControlField:
			cf.Execute(writer, field)
		case *DataField:
			df.Execute(writer, field)
		case DataField:
			df.Execute(writer, field)
		case *SubField:
			sf.Execute(writer, field)
		case SubField:
			sf.Execute(writer, field)
		case byte:
			fmt.Fprintf(writer, "%c", field)
/*
		default:
			template.HTMLFormatter(writer, format, values...)
*/
		}
	}
}

// Write a MARC/XML representation of the record
/*
func (record Record) XML(writer io.Writer) (err error) {
	err = record.Execute(writer, record)
	return
}
*/
