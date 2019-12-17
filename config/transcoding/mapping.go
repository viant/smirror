package transcoding

//Mappings represents mappings
type Mappings []*Mapping

//Mapping represent a mapping
type Mapping struct {
	From string
	To   string
}
