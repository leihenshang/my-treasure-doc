package service

import (
	"encoding/json"
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/hot"
	"fastduck/treasure-doc/module/hot_top/model"
	"fastduck/treasure-doc/module/hot_top/service/ai"
	"log"
)

func ThinkWithDeepSeek(question string) (answer string, err error) {
	ds, err := ai.NewAiDeepSeek(conf.GetConf().Ai.DeepSeekToken)
	if err != nil {
		return "", err
	}
	resMap := hot.GetHotCache().GetAllMap()
	totalLen := len(resMap) / 2
	filterMap := make(map[model.Source]*hot.HotCacheItem, totalLen)
	var counter = 0
	for k, v := range resMap {
		if counter < totalLen {
			filterMap[k] = v
			break
		}
		counter++
	}
	data, err := json.Marshal(filterMap)
	if err != nil {
		return "", err
	}

	word := `你是一个生活方面,军事方面,世界局势方面,投资方面的专家,根据以下JSON数据帮我分析一下关键信息,并输出为JSON格式的数据返回" \n` + string(data)
	log.Println(word)
	answer, err = ds.Ask(word)
	if err != nil {
		return "", err
	}
	return answer, nil
}
