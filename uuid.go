package captcha

import (
//	"crypto/rand"
	"strconv"
//	"bytes"
	"crypto/sha256"
//	"encoding/base64"
	"time"
	"math/rand"
	"fmt"
)
const (
	RAND_MAX = 9999
)

//func __NewUID() (uuid string, err error) {
//	u := new([16]byte)
//	_, err = rand.Read(u[:])
//
//	if err != nil {
//		return
//	}
//
//	u[8] = (u[8] | 0x40) & 0x7F
//	u[6] = (u[6] & 0xF) | (4 << 4)
//
//	var buf bytes.Buffer
//	for _, v := range u {
//		if v < 10 {
//			buf.WriteString("00" + strconv.Itoa(int(v)))
//		} else if v < 100 {
//			buf.WriteString("0" + strconv.Itoa(int(v)))
//		} else {
//			buf.WriteString(strconv.Itoa(int(v)))
//		}
//
//	}
//	uuid = buf.String()
//
//	return
//}

func NewUID() (uuid string, err error) {
	x := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(RAND_MAX)
	key := []byte(time.Now().Format(time.StampNano) + "#" + strconv.Itoa(x))
//	fmt.Printf("- org= [%d] \n", x)
//	fmt.Printf("- key= [%s] \n", key)

	sha := sha256.Sum256(key)
	uuid = fmt.Sprintf("%x", sha)
	err = nil

	return
}
