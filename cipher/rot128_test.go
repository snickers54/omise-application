package cipher

import (
	"bytes"
	"testing"

	r "github.com/stretchr/testify/require"
)

var (
	TestBuffer        = []byte{128, 129, 130}
	ReverseTestBuffer = []byte{0, 1, 2}
)

func TestRot128Reader_Read(t *testing.T) {
	reader := NewRot128Reader(bytes.NewBuffer(TestBuffer))
	r.NotNil(t, reader)

	buf := make([]byte, 3, 3)
	n, err := reader.Read(buf)
	r.NoError(t, err)
	r.Equal(t, 3, n)
	r.Equal(t, ReverseTestBuffer, buf)
}

func TestRot128Reader_Reversible(t *testing.T) {
	reader := NewRot128Reader(bytes.NewBuffer(TestBuffer))
	r.NotNil(t, reader)

	reader = NewRot128Reader(reader)
	r.NotNil(t, reader)

	buf := make([]byte, 3, 3)
	n, err := reader.Read(buf)
	r.NoError(t, err)
	r.Equal(t, 3, n)
	r.Equal(t, TestBuffer, buf)
}
