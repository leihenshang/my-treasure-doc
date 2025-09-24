package hot

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/model"
	"testing"
	"time"
)

func Test_GetHotFromFileCache(t *testing.T) {
	hot = NewHot(&conf.Hot{
		ExpiredCheckIntervalParsed: time.Second * 10,
		HotPullIntervalParsed:      conf.DefaultHotPullInterval,
		HotFileCachePath:           `D:\my-project\my-treasure-doc\module\hot_top\hot_cache`,
	})
	hotItem, err := hot.GetHotFromFileCache(model.Source("weibo"))
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
