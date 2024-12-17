package gid

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sony flake not created")
	}
}

func genSid() (int64, error) {
	id, err := sf.NextID()
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func GenId() int64 {
	id, err := genSid()
	if err != nil {
		log.Printf("failed to gen sonyflake: %v", err)
		time.Sleep(time.Millisecond)
		if id, err = genSid(); err != nil {
			log.Printf("failed to retry gen sonyflake: %v", err)
		}
	}
	return id
}

func BatchGenId(count int) []int64 {
	res := make([]int64, 0, count)
	for i := 0; i < count; i++ {
		res = append(res, GenId())
	}
	return res
}

type Gid string

func (s Gid) Int64() int64 {
	num, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		fmt.Printf("failed to convert gid %s to int", s)
		return 0
	}
	return num
}
