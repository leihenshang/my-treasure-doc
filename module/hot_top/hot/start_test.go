package hot

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	hotcache "fastduck/treasure-doc/module/hot_top/hot/hot_cache"
	"testing"
	"time"
)

func Test_GetHotFromFileCache(t *testing.T) {
	hot = NewHot(&conf.Hot{
		ExpiredCheckIntervalParsed: time.Second * 10,
		HotPullIntervalParsed:      conf.DefaultHotPullInterval,
		HotFileCachePath:           `D:\my-project\my-treasure-doc\module\hot_top\hot_cache`,
	})
	hotItem, err := hotcache.GetHotFileCache().Get(conf.Source("weibo"))
	if err != nil {
		t.Fatalf("GetHotFromFileCache failed, err: %v", err)
	} else if hotItem == nil {
		t.Logf("GetHotFromFileCache failed, items is nil")
		return
	}

	for _, v := range hotItem.Data {
		t.Logf("%#v \n", v)
	}
}
