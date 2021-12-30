package converter

import (
	"strconv"
)

//type  for call the convert functions from string to other types through dot
type String string

const (
	binaryTrue = "\x01"
	bimaryFalse = "\x00"
)

func (s String) IsBinaryBool() bool {
	return s == binaryTrue || s == bimaryFalse
}

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

func (s String) Int() int {
	result := int(0)

	if value, err := strconv.Atoi(string(s)); err == nil {
		result = int(value)
	}
	return result
}

func (s String) Float32() float32 {
	if tm, err := strconv.ParseFloat(string(s), 32); err != nil {
		return 0
	} else {
		return float32(tm)
	}
}

func (s String) BoolAsString() string{
	if s.IsBinaryBool() {
		if s.BinaryTrue() {
			return "true"
		} else {
			return "false"
		}
	}
	return "false"
}
 