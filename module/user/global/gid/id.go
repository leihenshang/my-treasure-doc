package gid

import (
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

func genSid() (string, error) {
	id, err := sf.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}

func GenId() string {
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

func BatchGenId(count int) []string {
	res := make([]string, 0, count)
	for i := 0; i < count; i++ {
		res = append(res, GenId())
	}
	return res
}
