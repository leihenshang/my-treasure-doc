package redis_lock_v1

import (
	"fastduck/treasure-doc/service/mall/global"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func lock(k string, expire time.Duration) (label string) {
	label = strconv.FormatUint(genLabel(), 10)
	_, err := global.Redis.SetNX(k, time.Now(), expire).Result()
	if err != nil {
		if err == redis.Nil {

		}
	}

	return
}

func unlock() {

}

func genLabel() uint64 {
	source := rand.NewSource(time.Now().Unix())
	ran := rand.New(source)
	return ran.Uint64()
}
