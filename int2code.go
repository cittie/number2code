package int2ucode

import (
	"errors"
	"log"
	"math"
	"math/big"
)

const (
	minPoolLength = 12                                                         // 随机字符池最小长度
	minCodeLength = 6                                                          // 码最小长度
	defaultPool   = "A5iSqnjeZQpmHtVPBW9YaDK2sCMXUrbhuTEzNL7ywfd4vRcgG83FJ6kx" // 去掉了"I1lOo0"的26个字母大小写以及数字2到9的随机顺序字符串
	defaultPrime1 = 23                                                         // 对number混淆，不要太大，否则容易超过上限
	defaultPrime2 = 977                                                        // 对code顺序混淆
)

var (
	stringPool []byte
	prime1     int
	prime2     int
	pBox       map[int]int
)

var (
	ErrInvalidPool     = errors.New("invalid pool")
	ErrInvalidPrime    = errors.New("invalid prime")
	ErrNumberExceedMax = errors.New("number exceed max")
	ErrCodeInvalid     = errors.New("code invalid")
	ERRCodeLessThanMin = errors.New("code length not enough")
)

func init() {
	if err := InitWithParams(defaultPool, defaultPrime1, defaultPrime2); err != nil {
		panic(err)
	}
}

// InitWithParams 初始化
// pool为随机字符池, 如有重复字符会报错
// prime1用于混淆number，prime2用于混淆code位置，都必须为大于2的质数
func InitWithParams(pool string, prime1, prime2 int) error {
	err := newPool(pool)
	if err != nil {
		return err
	}

	err = setPrimes(prime1, prime2)
	if err != nil {
		return err
	}

	size := GetPoolSize()
	pBox = make(map[int]int, size)
	for i := 0; i < size; i++ {
		pBox[i] = (i * prime2) % GetPoolSize()
	}

	log.Printf("init success with %v, %v, %v\n", pool, prime1, prime2)

	return nil
}

// GetPoolSize 获得字符池长度
func GetPoolSize() int {
	return len(stringPool)
}

// newPool 设置新字符池
func newPool(pool string) error {
	if !validatePool(pool) {
		return ErrInvalidPool
	}

	stringPool = []byte(pool)

	return nil
}

// setPrimes prime值修改
func setPrimes(p1, p2 int) error {
	if p1 < 3 || p1 > math.MaxInt64 || !big.NewInt(int64(p1)).ProbablyPrime(0) {
		return ErrInvalidPrime
	}

	if p2 < 3 || p2 > math.MaxInt64 || !big.NewInt(int64(p2)).ProbablyPrime(0) {
		return ErrInvalidPrime
	}

	prime1 = p1
	prime2 = p2

	return nil
}

// MarshalCode 将数字转为n位字符池长度进制逆序形式，最右位为校验码，n最小值为6
// 首先将数字加prime1混淆后，转换成n进制
// 第二次混淆：最低位 * i 加到高位，取余
// 第三次混淆：prime2 p-box
func MarshalCode(num, n int) ([]byte, error) {
	if n < minCodeLength {
		return nil, ERRCodeLessThanMin
	}
	size := GetPoolSize()

	// 判定num是否超过了n位可以表示的上限
	if num*prime1 > pow(size, n-1) {
		return nil, ErrNumberExceedMax
	}

	// 进制转换
	num *= prime1
	nums := make([]int, n)

	for i := 0; i < n-1; i++ {
		nums[i] = (num + nums[0]*i) % size
		num /= size
		nums[n-1] += nums[i]
	}

	// 校验码
	nums[n-1] *= prime1
	nums[n-1] %= size

	// P-Box混淆
	ret := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		ret = append(ret, stringPool[nums[(i*prime2)%n]])
	}

	return ret, nil
}

// UnmarshalCode 将n位码转回数字
// 是上面Marshal的逆运算
func UnmarshalCode(code string) (int, error) {
	size := GetPoolSize()
	n := len(code)

	// 混淆复原
	pos := make([]int, n)
	for i := 0; i < n; i++ {
		pos[(i*prime2)%n] = i
	}

	originCode := make([]byte, n)
	for i := 0; i < n; i++ {
		originCode[i] = byte(code[pos[i]])
	}

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		num := getIdx(byte(originCode[i]))
		if num == -1 {
			// 出现了不存在于池子中的字符
			return 0, ErrCodeInvalid
		}
		nums[i] = num
	}

	// 校验
	vNum := 0
	for i := 0; i < n-1; i++ {
		vNum += nums[i]
	}
	if vNum*prime1%size != nums[n-1] {
		return 0, ErrCodeInvalid
	}

	// 逆运算
	num, base := 0, 1
	for i := 0; i < n-1; i++ {
		nums[i] += (size - nums[0]) * i
		nums[i] %= size
		num += nums[i] * base
		base *= size
	}

	return num / prime1, nil
}
