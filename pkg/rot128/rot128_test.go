package rot128

import (
	"testing"

	r "github.com/stretchr/testify/require"
)

var (
	TestBuffer        = []byte{128, 129, 130}
	ReverseTestBuffer = []byte{0, 1, 2}
)

func TestRot128(t *testing.T) {
	Rot128(TestBuffer)
	r.Equal(t, ReverseTestBuffer, TestBuffer)
}
