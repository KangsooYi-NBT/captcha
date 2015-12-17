package captcha

import "math/rand"
import "strconv"
import "time"

const (
	min = 9 // 10
	max = 500
	mod = 2 // could be: 1, 2, 3, 4
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//
func randomNum(max int64) int64 {
	i := rand.Int63n(max)
	if i < 100 {
		return max - min + i
	}

	return i
}

//
func RandomExp() (exp []byte, result string) {
	first := randomNum(max)
	second := randomNum(min)

	ope := randomNum(min) % mod
	var resultNum int64

	switch ope {
	case 0: // +
		resultNum = first + second
	case 1: // -
		resultNum = first - second
		// case 2:
		// case 3:
	}

	result = strconv.FormatInt(resultNum, 10)

	exp = make([]byte, 6)
	x := first % 100
	exp[0] = byte(first / 100)
	exp[1] = byte(x / 10)
	exp[2] = byte(x % 10)
	exp[3] = byte(10 + ope)
	exp[4] = byte(second % 10)
	exp[5] = byte(14)

	return
}
