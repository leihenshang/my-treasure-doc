package redis_lock

import (
	"context"
	"database/sql"
	"fastduck/treasure-doc/service/mall/global"
	"log"
	"testing"
	"time"
)

func setup() {
	//全局初始化
	global.GlobalInit("test")

	//同步写入日志
	defer global.Zap.Sync()
	defer global.ZapSugar.Sync()

	//关闭mysql
	db, _ := global.DbIns.DB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func TestLock(t *testing.T) {
	t.Log("start lock")
	ctx := context.Background()
	key := "hhhhlock"
	lockDuration := time.Second * 5
	lock := NewDistributeRedisLock(global.Redis, lockDuration)
	err := lock.Lock(ctx, key)
	if err != nil {
		t.Error(err)
	}

	err2 := lock.Lock(ctx, key)
	if err2 == nil {
		log.Fatal("the lock is unlocked .oh,failed!")
	}

	// unlock
	// unlockErr := unlock(key, label)
	// if unlockErr != nil {

	// }

	lock.Unlock(ctx, key)

	err3 := lock.Lock(ctx, key)
	if err3 != nil {
		log.Fatal("the lock,failed!")
	}

	t.Log("ok")
}
