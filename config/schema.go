package config

import (
	"encoding/csv"
	"io"
	"strings"
)

type Schema struct {
	//MaxBadRecords
	MaxBadRecords *int
	Format     string
	Delimiter  string
	LazyQuotes bool
	FieldCount int
	isJSON     *bool
	isCSV      *bool
	Fields     []*Field
}



func (r *Schema) IsJSON() bool {
	if r.isJSON != nil {
		return *r.isJSON
	}
	isJSON := strings.ToUpper(r.Format) == "JSON"
	r.isJSON = &isJSON
	return isJSON
}

func (r *Schema) IsCSV() bool {
	if r.isCSV != nil {
		return *r.isCSV
	}
	isCSV := strings.ToUpper(r.Format) == "CSV"
	r.isCSV = &isCSV
	return isCSV

}

func (r *Schema) NewCsvReader(reader io.Reader) *csv.Reader {
	result := csv.NewReader(reader)
	if r.Delimiter != "" {
		result.Comma = rune(r.Delimiter[0])
	}
	result.LazyQuotes = r.LazyQuotes
	return result
}
