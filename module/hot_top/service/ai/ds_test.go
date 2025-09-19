package ai

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"testing"
)

func Test_DeepSeek_Call(t *testing.T) {
	if err := conf.InitConf("conf.toml"); err != nil {
		t.Fatal(err)
	}
	ds, err := NewAiDeepSeek(conf.GetConf().DeepSeekToken)
	if err != nil {
		t.Fatal(err)
	}
	answer, err := ds.Ask("你好")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(answer)
}
