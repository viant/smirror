package avro

import (
	"github.com/pkg/errors"
	"github.com/viant/toolbox"
	"io"
	"math"
	"strings"
)

var UnionNull = int64(0)
var UnionNotNull = int64(1)

func translateToNull(value interface{}, w io.Writer) error {
	return nil
}

func translateToString(value interface{}, w io.Writer) error {
	v := toolbox.AsString(value)
	return writeString(v, w)
}

func translateToBytes(value interface{}, w io.Writer) error {
	v, ok := value.([]byte)
	if ! ok {
		text, ok := value.(string)
		if ! ok {
			return errors.Errorf("failed to cast %T to []byte", value)
		}
		v = []byte(text)
	}

	return writeBytes(v, w)
}

func translateToLong(value interface{}, w io.Writer) error {
	if text, ok := value.(string); ok {
		value = strings.TrimSpace(text)
	}
	v, err := toolbox.ToInt(value)
	if err != nil {
		return errors.Wrapf(err, "failed to convert to float: %v", value)
	}
	return writeLong(int64(v), w)
}

func translateToDouble(value interface{}, w io.Writer) error {
	if text, ok := value.(string); ok {
		value = strings.TrimSpace(text)
	}
	v, err := toolbox.ToFloat(value)
	if err != nil {
		return errors.Wrapf(err, "failed to convert to float: %v", value)
	}
	return writeDouble(float64(v), w)
}

func translateToFloat(value interface{}, w io.Writer) error {
	if text, ok := value.(string); ok {
		value = strings.TrimSpace(text)
	}
	v, err := toolbox.ToFloat(value)
	if err != nil {
		return errors.Wrapf(err, "failed to convert to float: %v", value)
	}
	return writeFloat(float32(v), w)
}

func translateToBoolean(value interface{}, w io.Writer) error {
	if text, ok := value.(string); ok {
		value = strings.TrimSpace(text)
	}
	v, err := toolbox.ToBoolean(value)
	if err != nil {
		return errors.Wrapf(err, "failed to convert to float: %v", value)
	}
	return writeBoolean(v, w)
}

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}

func encodeInt(w io.Writer, byteCount int, encoded uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil
}

func writeBoolean(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}

func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r>>downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeNull(_ io.Writer) error {
	return nil
}

func writeBytes(r []byte, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	_, err = w.Write(r)
	return err
}

func writeString(text string, w io.Writer) error {
	err := writeLong(int64(len(text)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(text)
	} else {
		_, err = w.Write([]byte(text))
	}
	return err
}
