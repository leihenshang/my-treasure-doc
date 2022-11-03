package order

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"
)

// TODO 还有点问题，会生成长度不一样的订单号 31,32,33 都有
func GenerateOrderNo() (number string, err error) {
	timeNow := time.Now()
	dateTimeStr := timeNow.Format("20060102150405")
	number = dateTimeStr
	b := new(big.Int).SetInt64(timeNow.UnixNano())
	if i, err := rand.Int(rand.Reader, b); err != nil {
		return "", err
	} else {
		number = number + strconv.FormatInt(i.Int64(), 10)
	}

	return
}
