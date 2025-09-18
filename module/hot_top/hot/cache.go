package hot

import (
	"fastduck/treasure-doc/module/hot_top/model"
	"maps"
	"sync"
	"time"
)

type HotCache struct {
	lock  *sync.RWMutex
	cache map[model.Source]*HotCacheItem
}

var hotCache *HotCache
var hotCacheOnce *sync.Once = &sync.Once{}

type HotCacheItem struct {
	LastUpdateTime time.Time      `json:"lastUpdateTime"`
	HotData        *model.HotData `json:"hotData"`
}

func NewHotCache(len int) *HotCache {
	hotCacheOnce.Do(func() {
		hotCache = &HotCache{
			lock:  &sync.RWMutex{},
			cache: make(map[model.Source]*HotCacheItem, len),
		}
	})
	return hotCache
}

func GetHotCache() *HotCache {
	return hotCache
}

func (c *HotCache) Get(source model.Source) (*HotCacheItem, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.cache[source]
	if !ok {
		return &HotCacheItem{}, false

	}
	return item, true
}

func (c *HotCache) GetAllMap() map[model.Source]*HotCacheItem {
	c.lock.RLock()
	defer c.lock.RUnlock()
	resp := make(map[model.Source]*HotCacheItem, len(c.cache))
	maps.Copy(resp, c.cache)
	return resp
}

func (c *HotCache) GetExpired(t time.Duration) []model.Source {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var res []model.Source
	for k, v := range c.cache {
		if time.Since(v.LastUpdateTime) > t {
			res = append(res, k)
		}
	}
	return res
}

func (c *HotCache) Set(source model.Source, data *model.HotData) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if data != nil {
		c.cache[source] = &HotCacheItem{
			LastUpdateTime: time.Now(),
			HotData:        data,
		}
	}
}

func (c *HotCache) SetWithLastUpdateTime(source model.Source, data *HotCacheItem) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if data != nil {
		c.cache[source] = data
	}
}
