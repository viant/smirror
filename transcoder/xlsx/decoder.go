package xlsx

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"github.com/viant/toolbox"
	"io"
	"io/ioutil"
	"smirror/transcoder/avro/schma"
	"sort"
	"strings"
	"unicode"
)

type Decoder struct {
	sheet  *xlsx.Sheet
	fields []*schma.Field
	meta   *meta
}

func (d *Decoder) Schema() string {
	d.meta = d.getMeta()
	d.buildSchemaFields()
	fieldsDefinition := []string{}
	for _, field := range d.fields {
		fieldDefinition := ""
		switch field.Type.Name {
		case schma.TypeString, schma.TypeLong, schma.TypeFloat, schma.TypeBoolean:
			fieldDefinition = fmt.Sprintf(schma.FieldTmpl, field.Name, field.Type.Name)
		case schma.TypeTimestamp:
			fieldDefinition = fmt.Sprintf(schma.TimeFieldTmpl, field.Name)
		}
		fieldsDefinition = append(fieldsDefinition, fieldDefinition)
	}
	schema := fmt.Sprintf(schma.RootTmpl, strings.Join(fieldsDefinition, ",\n"))
	return schema
}


func (d *Decoder) HasMore() bool {
	return d.meta.hasMore()
}

//NextRecord populate supplied record if has more and returns true
func (d *Decoder) NextRecord(record map[string]interface{}) error {
	var err error
	loop:
	if 	! d.meta.hasMore() {
		return nil
	}

	row := d.sheet.Rows[d.meta.index]
	d.meta.index++
	empty := true
	for i, field := range d.fields {
		record[field.Name] = nil
		if i >= len(row.Cells) {
			continue
		}

		cell := row.Cells[i]
		if cell.Value != "" {
			empty = false
		}
		var value interface{}
		switch field.Type.Name {
		case schma.TypeBoolean:
			value = cell.Bool()
		case schma.TypeFloat:
			if cell.Value == "" {
				continue
			}
			value, err = cell.Float()
		case schma.TypeLong:
			if cell.Value == "" {
				continue
			}
			value, err = cell.Int()
		case schma.TypeString:
			value = cell.Value
		case schma.TypeTimestamp:
			if timeValue, err := cell.GetTime(false); err == nil {
				value = timeValue.UnixNano() / 1000000
			}
		}
		if err != nil {
			return errors.Wrapf(err, "failed decode %v %v", field.Name, cell.Value)
		}
		if empty {
			goto loop
		}
		record[field.Name] = value
	}
	return nil
}


func (d *Decoder) getMeta() *meta {
	var index = make(map[int]*meta)
	rowCount := len(d.sheet.Rows)
	for i := 0; i < rowCount; i++ {
		rowCells := d.sheet.Rows[i].Cells
		nonEmptyCells := 0
		for _, cell := range rowCells {
			if cell.Value != "" {
				nonEmptyCells++
			}
		}
		if nonEmptyCells == 0 {
			continue
		}
		cells := len(rowCells)
		if _, ok := index[nonEmptyCells]; ! ok {
			index[nonEmptyCells] = &meta{firstRow: i, cells: cells, nonEmptyCells:nonEmptyCells}
		}
	}
	candidates := make(metaSlice, 0)
	for _, candidate := range index {
		candidates = append(candidates, candidate)
	}
	sort.Sort(candidates)
	result := candidates[len(candidates)-1]
	result.lastRow = rowCount -1
	result.index = result.firstRow + 1
	return result
}

func (d *Decoder) buildSchemaFields() {
	meta := d.meta
	rowCount := meta.lastRow - (meta.firstRow + 1)
	if rowCount > 100 {
		rowCount = 100
	}
	row := d.sheet.Rows[meta.firstRow]
	d.fields = make([]*schma.Field, 0)
	for i, cell := range row.Cells {
		name := normalizeFieldName(cell.Value, i)
		field := &schma.Field{Name: name, Type: &schma.Schema{Name: d.getFieldType(strings.ToUpper(name), i, meta.firstRow+1, rowCount)}}
		d.fields = append(d.fields, field)
	}
}

func (d *Decoder) getFieldType(name string, cellPosition, offset, rowCount int) string {
	var typeName = schma.TypeString
	var typeCount = make(map[string]int)
	row := d.sheet.Rows[offset]
	for i := 0; i < rowCount; i++ {
		if i >= len(row.Cells) {
			continue
		}
		cell :=row.Cells[cellPosition]
		offset++
		if cell.Value == "" {
			continue
		}
		switch cell.Type() {
		case xlsx.CellTypeStringFormula, xlsx.CellTypeNumeric:
			typeName = schma.TypeFloat
			floatValue, err := toolbox.ToFloat(cell.Value)
			if err == nil && floatValue == float64(int(floatValue)) {
				typeName = schma.TypeLong
				if d.isTimestampType(name, cell, typeName) {
					typeName = schma.TypeTimestamp
				}
			}
		case xlsx.CellTypeBool:
			typeName = schma.TypeBoolean
		case xlsx.CellTypeDate:
			typeName = schma.TypeTimestamp
		}
		typeCount[typeName]++
		if typeName == schma.TypeFloat {
			return typeName
		}
		if typeCount[typeName] > 50 {
			return typeName
		}
	}
	return typeName
}

//isTimestampType returns timestamp type
func (d *Decoder) isTimestampType(name string, cell *xlsx.Cell, typeName string) bool {
	if strings.Contains(name, "DATE") || strings.Contains(name, "DAY") || strings.Contains(name, "TIME") || strings.Contains(name, "TS") {
		if timeValue, err := cell.GetTime(false); err == nil {
			if timeValue.Year() > 1900 && timeValue.Year() < 2100 {
				return true
			}
		}
	}
	return false
}

func normalizeFieldName(name string, index int) string {
	if name == "" {
		return fmt.Sprintf("Field%03d", index)
	}
	replacement := map[string]string{
		"%": "_pct",
		")": "",
		"(": "",
		"-": "_",
	}
	name = strings.TrimSpace(name)
	for k, v := range replacement {
		if count := strings.Count(name, k); count > 0 {
			name = strings.Replace(name, k, v, count)
		}
	}
	normalized := make([]byte, 0)
	for _, r := range name {
		if unicode.IsSpace(r) {
			normalized = append(normalized, '_')
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			normalized = append(normalized, byte(unicode.ToUpper(r)))
		}
	}
	name = string(normalized)
	if ! unicode.IsLetter(rune(normalized[0])) {
		name = "F_" + name
	}
	return toolbox.ToCaseFormat(name, toolbox.CaseUpperUnderscore, toolbox.CaseUpperCamel)
}

func NewDecoder(reader io.Reader) (*Decoder, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	file, err := xlsx.OpenBinary(data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open xlsx")
	}

	if len(file.Sheets) == 0 {
		return nil, errors.Errorf("sheets were empty")
	}
	return &Decoder{
		sheet: file.Sheets[0],
	}, nil
}
