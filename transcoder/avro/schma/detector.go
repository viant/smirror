package schma

const (
	//RootTmpl schema root template
	RootTmpl = `{"type":"record","name": "Root","fields": [%v]}`
	//Field template
	FieldTmpl = `{ "name":"%v","type":["null", "%v"],"default": null}`

	//TimeFieldTmpl template
	TimeFieldTmpl = `{ "name":"%v","type":["null",{"type":"long","logicalType":"timestamp-millis"}], "default": null}`



	TypeString = "string"
	TypeBoolean = "boolean"
	TypeLong = "long"
	TypeFloat = "float"
	TypeTimestamp = "timestamp"
)


