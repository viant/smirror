package config

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"smirror/config/transcoding"
)

//Transcoding represents transcoding
type Transcoding struct {
	Source      transcoding.Codec
	Dest        transcoding.Codec
	PathMapping transcoding.Mappings
}

//Load intialise transcoding
func (t *Transcoding) Init(ctx context.Context, fs afs.Service) error {
	if t.Dest.SchemaURL == "" {
		return nil
	}
	if t.Dest.Schema == "" {
		_, err := t.Dest.LoadSchema(ctx, fs)
		return err
	}
	return nil
}

//Load check if transcoding is valid
func (t *Transcoding) Validate() error {
	if !t.Source.IsJSON() && !t.Source.IsCSV() {
		return errors.Errorf("unsupported source format: %v", t.Source.Format)
	}
	if t.Source.IsCSV() && len(t.Source.Fields) == 0 {
		return errors.Errorf("source fields were empty: for %v format", t.Source.Format)
	}
	if !t.Dest.IsAvro() {
		return errors.Errorf("unsupported dest format: %v", t.Dest.Format)
	}
	if t.Dest.Schema == "" {
		return errors.Errorf("dest schema was empty: %v", t.Dest.SchemaURL)
	}
	return nil
}
