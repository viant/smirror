package schma

// Record implements Schema and represents Avro record type.
type Record struct {
	Fields []*Field `json:"fields"`
}
