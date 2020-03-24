package base



//SchemaError represent schema validation error
type SchemaError struct {
	err error
}

func (v SchemaError) Error() string {
	return v.err.Error()
}

//NewSchemaError creates a schema error
func NewSchemaError(err error) *SchemaError {
	return &SchemaError{err:err}
}


//IsSchemaError returns true if schema error
func IsSchemaError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*SchemaError)
	return ok
}