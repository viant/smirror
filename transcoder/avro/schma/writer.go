package schma

import "io"

type Translator func(value interface{}, w io.Writer) error
