package index

import (
	"maps"
	"slices"
	"sync"

	"github.com/fumiama/jieba"
)

type IndexCache struct {
	cache map[string][]string
	lock  *sync.RWMutex
}

type ContentCache struct {
	cache map[string]string
	lock  *sync.RWMutex
}

var indexCache *IndexCache = &IndexCache{
	cache: make(map[string][]string),
	lock:  &sync.RWMutex{},
}
var contentCache *ContentCache = &ContentCache{
	cache: make(map[string]string),
	lock:  &sync.RWMutex{},
}

var seg *jieba.Segmenter

func GetSeg() *jieba.Segmenter {
	return seg
}

func GetContentCache() *ContentCache {
	return contentCache
}

func GetIndexCache() *IndexCache {
	return indexCache
}

func InitDict() {
	var err error
	seg, err = jieba.LoadDictionaryAt("conf/dict.txt")
	if err != nil {
		panic(err)
	}
}

func (i *IndexCache) Set(key, value string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.cache[key] = append(i.cache[key], value)
}

func (i *IndexCache) Get(key string) (values []string, ok bool) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	values, ok = i.cache[key]
	return
}

func (i *IndexCache) List() map[string]any {
	results := make(map[string]any)
	i.lock.RLock()
	defer i.lock.RUnlock()
	indexCacheMap := make(map[string][]string)
	maps.Copy(indexCacheMap, i.cache)
	contentCache.lock.RLock()
	defer contentCache.lock.RUnlock()
	contentCacheMap := make(map[string]string)
	maps.Copy(contentCacheMap, contentCache.cache)
	results["index"] = indexCacheMap
	results["content"] = contentCacheMap
	return results
}

func (i *IndexCache) Search(keyword ...string) (ids []string) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	for _, key := range keyword {
		ids = append(ids, i.cache[key]...)
	}
	return
}

func (i *ContentCache) Set(key, value string) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.cache[key] = value
}

func (i *ContentCache) Get(key ...string) (values []string) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	for _, k := range slices.Compact(key) {
		values = append(values, i.cache[k])
	}
	return
}
