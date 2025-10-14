package service

import (
	"fastduck/treasure-doc/module/reverse_index/comm/gid"
	"fastduck/treasure-doc/module/reverse_index/index"
)

func Index(content ...string) error {
	// 校验content是否为空
	if len(content) == 0 {
		return nil
	}
	for _, v := range content {
		words := index.GetSeg().Cut(v, true)
		ids := gid.BatchGenId(len(words))
		for i, w := range words {
			index.GetIndexCache().Set(w, ids[i])
			index.GetContentCache().Set(ids[i], v)
		}
	}
	return nil
}

func Search(keyword string) ([]string, error) {
	words := index.GetSeg().Cut(keyword, true)
	var results []string
	ids := index.GetIndexCache().Search(words...)
	results = index.GetContentCache().Get(ids...)
	return results, nil
}

func List() any {
	return index.GetIndexCache().List()
}
