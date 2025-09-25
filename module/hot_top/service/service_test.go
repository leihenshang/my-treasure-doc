package service

import (
	"encoding/json"
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/hot"
	"fastduck/treasure-doc/module/hot_top/service/ai"
	"fmt"
	"log"
	"math"
	"strings"
	"testing"
)

func Test_Ai(t *testing.T) {
	if err := conf.InitConf(`D:\my-project\my-treasure-doc\module\hot_top\config.toml`); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	originConf := conf.GetConf()
	originConf.Hot.HotFileCachePath = `D:\my-project\my-treasure-doc\module\hot_top\hot_cache`
	hot := hot.NewHot(&originConf.Hot)

	var queryData [][2]string
	for k := range conf.HotConfListMap {
		res, err := hot.GetHotFromFileCache(k)
		if err != nil || res == nil {
			t.Logf("failed to get hot from file cache, source: %v, err: %v", k, err)
			continue
		}
		var list []string
		for _, v := range res.Data {
			list = append(list, fmt.Sprintf("%s(%s)", v.Title, v.URL))
		}

		queryData = append(queryData, [2]string{res.Title, strings.ReplaceAll(strings.Join(list, ", "), " ", "")})
	}
	t.Logf("hot data len: %d", len(queryData))
	jsonData, err := json.Marshal(queryData)
	if err != nil {
		t.Logf("failed to marshal hot data, err: %v", err)
		return
	}
	t.Logf("hot data json: %s,data len: %d", string(jsonData), len(jsonData))
	word := `你是一个生活方面,军事方面,世界局势方面,投资方面的专家,根据以下JSON数据帮我分析一下关键信息,并输出为JSON格式的数据返回" \n` + string(jsonData)

	if len(word) > 128*1024 {
		word = word[:128*1024]
	}
	t.Logf("word: %s,data len: %d, about: %.2f KB", word, len(word), math.Round(float64(len(word))/1024)*100/100)
	ds, err := ai.NewAiDeepSeek(originConf.Ai.DeepSeekToken)
	if err != nil {
		t.Errorf("failed to new ai deepseek, err: %v", err)
		return
	}
	answer, err := ds.Ask(word)
	if err != nil {
		t.Errorf("failed to ask ai deepseek, err: %v", err)
		return
	}
	t.Logf("ai deepseek answer: %s", answer)
}
