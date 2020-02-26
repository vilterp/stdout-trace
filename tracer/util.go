package tracer

import (
	"io"
)

type NullWriter struct{}

func NewNullWriter() *NullWriter {
	return &NullWriter{}
}

var _ io.Writer = &NullWriter{}

func (n *NullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
