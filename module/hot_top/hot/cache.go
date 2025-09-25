package hot

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/model"
	"maps"
	"sync"
)

type HotCache struct {
	lock  *sync.RWMutex
	cache map[conf.Source]*model.HotData
}

var hotCache *HotCache
var hotCacheOnce *sync.Once = &sync.Once{}

func NewHotCache(len int) *HotCache {
	hotCacheOnce.Do(func() {
		hotCache = &HotCache{
			lock:  &sync.RWMutex{},
			cache: make(map[conf.Source]*model.HotData, len),
		}
	})
	return hotCache
}

func GetHotCache() *HotCache {
	return hotCache
}

func (c *HotCache) Get(source conf.Source) (*model.HotData, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.cache[source]
	if !ok {
		return &model.HotData{}, false

	}
	return item, true
}

func (c *HotCache) GetAllMap() map[conf.Source]*model.HotData {
	c.lock.RLock()
	defer c.lock.RUnlock()
	resp := make(map[conf.Source]*model.HotData, len(c.cache))
	maps.Copy(resp, c.cache)
	return resp
}

func (c *HotCache) Set(source conf.Source, data *model.HotData) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if data != nil {
		c.cache[source] = data
	}
}
