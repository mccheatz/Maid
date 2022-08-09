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
