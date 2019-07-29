package int2ucode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePool(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"", false},                     // 长度1
		{"123456", false},               // 长度2
		{"1234567890abcdefghhh", false}, // 重复
		{"1234567890abcdefg=", false},   // 非法字符
		{"1234567890abcdefgh", true},
		{"A5iSqnjeZQpmHtVPBW9YaDK2sCMXUrbhuTEzNL7ywfd4vRcgG83FJ6kx", true},
	}

	for i, test := range tests {
		assert.Equal(t, test.out, validatePool(test.in), i)
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		x   int
		n   int
		out int
	}{
		{1, 0, 1},
		{11, 0, 1},
		{9, 1, 9},
		{3, 3, 27},
		{2, 16, 65536},
		{56, 7, 1727094849536},
	}

	for i, test := range tests {
		assert.Equal(t, test.out, pow(test.x, test.n), i)
	}
}

func TestShufflePool(t *testing.T) {
	pool := "A5iSqnjeZQpmHtVPBW9YaDK2sCMXUrbhuTEzNL7ywfd4vRcgG83FJ6kx"

	assert.NotEqual(t, pool, shufflePool(pool))
}

func BenchmarkPow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pow(56, 32)
	}
}
