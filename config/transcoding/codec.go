package transcoding

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"io/ioutil"
	"strings"
)

//Codec represents transcoding
type Codec struct {
	Format         string
	Fields         []string
	HasHeader      bool
	Delimiter      string
	LazyQuotes     bool
	SchemaURL      string
	RecordPerBlock int64
	Schema         string
	isCSV          *bool
	isXLSX         *bool
	isJSON         *bool
	isAvro         *bool
}

func (c *Codec) LoadSchema(ctx context.Context, fs afs.Service) (string, error) {
	if c.Schema != "" {
		return c.Schema, nil
	}
	reader, err := fs.DownloadWithURL(ctx, c.SchemaURL)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load avro Schema: %v", c.SchemaURL)
	}
	defer func() {
		_ = reader.Close()
	}()
	rawSchema, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	c.Schema = string(rawSchema)
	return c.Schema, nil

}

//IsAvro returns true if AVRO format
func (c *Codec) IsAvro() bool {
	if c.isAvro != nil {
		return *c.isAvro
	}
	isAvro := strings.ToUpper(c.Format) == "AVRO"
	c.isAvro = &isAvro
	return isAvro
}

//IsCSV returns true if XLSX format
func (c *Codec) IsXLSX() bool {
	if c.isXLSX != nil {
		return *c.isXLSX
	}
	isXLSX := strings.ToUpper(c.Format) == "XLSX"
	c.isXLSX = &isXLSX
	return isXLSX
}

//IsCSV returns true if CSV format
func (c *Codec) IsCSV() bool {
	if c.isCSV != nil {
		return *c.isCSV
	}
	isCSV := strings.ToUpper(c.Format) == "CSV"
	c.isCSV = &isCSV
	return isCSV
}

//IsJSON returns true if JSON format
func (c *Codec) IsJSON() bool {
	if c.isJSON != nil {
		return *c.isJSON
	}
	isJSON := strings.ToUpper(c.Format) == "JSON"
	c.isJSON = &isJSON
	return isJSON
}
