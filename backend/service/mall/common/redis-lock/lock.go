package redis_lock

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type DistributeRedisLock struct {
	lockVal  string
	duration time.Duration
	client   *redis.Client
}

type RedisLock interface {
	Lock(ctx context.Context, key string) error
	Unlock(ctx context.Context, key string) error
}

func NewDistributeRedisLock(client *redis.Client, duration time.Duration) RedisLock {
	label := GenLockVal()
	return &DistributeRedisLock{client: client, duration: duration, lockVal: label}
}

func (d *DistributeRedisLock) Lock(ctx context.Context, key string) (err error) {

	res, setErr := d.client.SetNX(ctx, key, d.lockVal, d.duration).Result()
	if setErr != nil {
		if setErr != redis.Nil {
			err = setErr
			return
		}
	}
	if res == false {
		err = errors.New("locking failure! ")
		return
	}

	return
}

func (d *DistributeRedisLock) Unlock(ctx context.Context, key string) (err error) {
	// TODO 获取和删除任然不是原子化的，最好还是用lua脚本
	// get := d.client.Get(ctx, key)
	// if get.Err() != nil {
	// 	return get.Err()
	// }
	// if get.Val() == d.lockVal {
	// 	delRes := d.client.Del(ctx, key)
	// 	if delRes.Err() != nil {
	// 		return delRes.Err()
	// 	}
	// }

	script := redis.NewScript(`
	if redis.call('get',KEYS[1]) == ARGV[1] then 
   		return redis.call('del',KEYS[1]) 
	else
		return 0
	end;
	`)

	delRes := script.Run(ctx, d.client, []string{key}, []interface{}{d.lockVal})
	if delRes.Err() != nil {
		return delRes.Err()
	}

	return
}

func GenLockVal() string {
	source := rand.NewSource(time.Now().Unix())
	ran := rand.New(source)
	label := ran.Uint64()
	return strconv.FormatUint(label, 10)
}
