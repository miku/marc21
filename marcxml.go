package marc21

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
