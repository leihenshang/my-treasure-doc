package hot

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/model"
	"testing"
)

func Test_GetHotFromFileCache(t *testing.T) {
	hotItem, err := GetHotFromFileCache(`D:\my-project\my-treasure-doc\module\hot_top\hot_cache`, model.Source("weibo"), conf.DefaultHotPullInterval)
	if err != nil {
		t.Fatalf("GetHotFromFileCache failed, err: %v", err)
	} else if hotItem == nil {
		t.Logf("GetHotFromFileCache failed, items is nil")
		return
	}

	for _, v := range hotItem.HotData.Data {
		t.Logf("%#v \n", v)
	}
}
