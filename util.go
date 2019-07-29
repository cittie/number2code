package int2ucode

import (
	"math/rand"
	"unicode"
)

// validatePool 检查一个字符串是否为合法的池子
// 规则：1. 只使用数字和字母 2. 没有重复 3. 至少包括x位
func validatePool(pool string) bool {
	if len(pool) < minPoolLength {
		return false
	}

	unique := make(map[rune]struct{}, len(pool))

	for _, r := range pool {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}

		if _, ok := unique[r]; ok {
			return false
		}

		unique[r] = struct{}{}
	}

	return true
}

// getIdx 获取字符池中对应字符的index，如果不存在返回-1
func getIdx(b byte) int {
	for i, c := range stringPool {
		if b == c {
			return i
		}
	}

	return -1
}

// pow 就是简单算乘方的
func pow(x, n int) int {
	r := 1
	for n != 0 {
		if n&1 == 1 {
			r *= x
		}
		n /= 2
		x *= x
	}
	return r
}

// shufflePool 将随机字符池洗牌，可以用来生成你自己的默认随机字符池
func shufflePool(pool string) string {
	pb := []byte(pool)

	for i := len(pb) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		pb[i], pb[j] = pb[j], pb[i]
	}

	return string(pb)
}
