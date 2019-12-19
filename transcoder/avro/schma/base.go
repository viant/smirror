package schma

//Base represents base schema type
type Base struct {
	Type        string `json:"type,omitempty"`
	LogicalType string `json:"logicalType,omitempty"`
	isComplex   *bool
	isRecord    *bool
	isUnion     *bool
	isArray     *bool
	isString    *bool
	isInt       *bool
	isLong      *bool
	isFloat     *bool
	isDouble    *bool
	isBytes     *bool
	isBoolean   *bool
	isNull      *bool
}

//IsComplex return true if type is a complex type
func (b *Base) IsComplex() bool {
	if b.isComplex != nil {
		return *b.isComplex
	}
	isComplex := false
	switch b.Type {
	case typeRecord, typeEnum, typeArray, typeMap, typeFixed, typeUnion:
		isComplex = true
	}
	b.isComplex = &isComplex
	return isComplex
}

//IsRecord returns true if type is a record type
func (b *Base) IsRecord() bool {
	if b.isRecord != nil {
		return *b.isRecord
	}
	isRecord := b.Type == typeRecord
	b.isRecord = &isRecord
	return isRecord
}

func (b *Base) IsUnion() bool {
	return b.Type == typeUnion
}

func (b *Base) IsEnum() bool {
	return b.Type == typeUnion
}

func (b *Base) IsArray() bool {
	if b.isArray != nil {
		return *b.isArray
	}
	isArray := b.Type == typeArray
	b.isArray = &isArray
	return isArray
}

func (b *Base) IsMap() bool {
	return b.Type == typeMap
}

func (b *Base) IsFixed() bool {
	return b.Type == typeFixed
}

func (b *Base) IsString() bool {
	if b.isString != nil {
		return *b.isString
	}
	isString := b.Type == typeString
	b.isString = &isString
	return isString
}



func (b *Base) IsBytes() bool {
	if b.isBytes != nil {
		return *b.isBytes
	}
	isBytes := b.Type == typeBytes
	b.isBytes = &isBytes
	return isBytes
}


func (b *Base) IsInt() bool {
	if b.isInt != nil {
		return *b.isInt
	}
	isInt := b.Type == typeInt
	b.isInt = &isInt
	return isInt
}

func (b *Base) IsLong() bool {
	if b.isLong != nil {
		return *b.isLong
	}
	isLong := b.Type == typeLong
	b.isLong = &isLong
	return isLong
}

func (b *Base) IsFloat() bool {
	if b.isFloat != nil {
		return *b.isFloat
	}
	isFloat := b.Type == typeFloat
	b.isFloat = &isFloat
	return isFloat
}

func (b *Base) IsDouble() bool {
	if b.isDouble != nil {
		return *b.isDouble
	}
	isDouble := b.Type == typeDouble
	b.isDouble = &isDouble
	return isDouble
}

func (b *Base) IsBoolean() bool {
	if b.isBoolean != nil {
		return *b.isBoolean
	}
	isBoolean := b.Type == typeBoolean
	b.isBoolean = &isBoolean
	return isBoolean
}

func (b *Base) IsNull() bool {
	if b.isNull != nil {
		return *b.isNull
	}
	isNull := b.Type == typeNull
	b.isNull = &isNull
	return isNull
}
