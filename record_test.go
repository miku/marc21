package marc21

import "testing"

// TestAddField tests adding of a single field.
func TestAddField(t *testing.T) {
	record := &Record{}
	if len(record.Fields) != 0 {
		t.Errorf("len(record.Fields), got %v, want %v", len(record.Fields), 0)
	}
	record.AddField(&ControlField{Tag: "001", Data: "12345"})
	if len(record.Fields) != 1 {
		t.Errorf("len(record.Fields), got %v, want %v", len(record.Fields), 1)
	}
	data := record.Fields[0].(*ControlField).Data
	if data != "12345" {
		t.Errorf("field.Data, got %v, want %v", data, "12345")
	}
}
