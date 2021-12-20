package converter

import "strconv"

//type  for call the convert functions from string to other types through dot
type String string

const (
	binaryTrue = "\x01"
)

func (s String) BinaryTrue() bool {
	return s == binaryTrue
}

func (s String) Uint8() uint8 {
	result := uint8(0)

	if value, err := strconv.Atoi(string(s)); err == nil {
		result = uint8(value)
	}
	return result
}

