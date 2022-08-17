package util

import "bytes"

type JsonNull struct{}

func (c JsonNull) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`null`)
	return buf.Bytes(), nil
}

func (c *JsonNull) UnmarshalJSON(in []byte) error {
	return nil
}

type JsonRaw struct {
	Value string
}

func (c JsonRaw) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(c.Value)
	return buf.Bytes(), nil
}

func (c *JsonRaw) UnmarshalJSON(in []byte) error {
	c.Value = string(in)
	return nil
}
