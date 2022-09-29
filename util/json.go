package util

import (
	"bytes"
	"strconv"
	"strings"
)

type JsonNull struct{}

func (c JsonNull) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (c *JsonNull) UnmarshalJSON(in []byte) error {
	return nil
}

type JsonRaw struct {
	Value string
}

func (c JsonRaw) MarshalJSON() ([]byte, error) {
	return []byte(c.Value), nil
}

func (c *JsonRaw) UnmarshalJSON(in []byte) error {
	c.Value = string(in)
	return nil
}

type JsonIntString struct {
	Value int
}

func (c JsonIntString) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteRune('"')
	buf.WriteString(strconv.Itoa(c.Value))
	buf.WriteRune('"')
	return buf.Bytes(), nil
}

func (c *JsonIntString) UnmarshalJSON(in []byte) error {
	str := string(in)
	str = strings.Trim(str, "\"")

	val, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	c.Value = val
	
	return nil
}