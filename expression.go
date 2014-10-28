package captcha

import "math/rand"
import "strconv"
import "time"

const (
	min = 10
	max = 99
	mod = 2 // could be: 1, 2, 3, 4
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//
func randomNum() int64 {
	i := rand.Int63n(max)

	if i < min {
		return max - min + i
	}

	return i
}

//
func RandomExp() (exp []byte, result string) {
	first := randomNum()
	second := randomNum()

	ope := randomNum() % mod
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
	exp[0] = byte(first / 10)
	exp[1] = byte(first % 10)
	exp[2] = byte(10 + ope)
	exp[3] = byte(second / 10)
	exp[4] = byte(second % 10)
	exp[5] = byte(14)

	return
}
