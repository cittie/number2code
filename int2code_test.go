package int2ucode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetPrimes(t *testing.T) {
	tests := []struct {
		prime1 int
		prime2 int
		out    error
	}{
		{0, 0, ErrInvalidPrime},
		{2, 2, ErrInvalidPrime},
		{7, 21, ErrInvalidPrime},
		{3, 7, nil},
		{23, 977, nil},
	}

	for i, test := range tests {
		assert.Equal(t, test.out, setPrimes(test.prime1, test.prime2), i)
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		in  int
		out string
	}{
		{0, "AAAAAAAAAA"},
		{1, "28N2BScrMt"},
		{56, "AAA2AAACAA"},
		{57, "28NcBSckMt"},
		{3136, "AAAAAA2CAA"},
		{2001, "g8aUwmJpir"},
		{201999, "8eUaADmLVn"},
		{2049999999, "WeccsijggV"},
	}

	for i, test := range tests {
		code, err := MarshalCode(test.in, 10)
		assert.Equal(t, test.out, string(code), i)
		assert.Nil(t, err)
	}

	for i, test := range tests {
		num, err := UnmarshalCode(test.out)
		assert.Equal(t, test.in, num, i)
		assert.Nil(t, err)
	}
}

func BenchmarkMarshal(b *testing.B) {
	num := 2049999999

	for i := 0; i < b.N; i++ {
		MarshalCode(num, 10)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	code := "WeccsijggV"

	for i := 0; i < b.N; i++ {
		UnmarshalCode(code)
	}
}
