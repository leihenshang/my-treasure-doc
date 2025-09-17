package hot

import (
	"sync"
	"time"
)

type HotCache struct {
	lock  *sync.RWMutex
	cache map[Source]*HotCacheItem
}

var hotCache *HotCache
var hotCacheOnce *sync.Once = &sync.Once{}

type HotCacheItem struct {
	LastUpdateTime time.Time
	HotData        *HotData
}

func NewHotCache(len int) *HotCache {
	hotCacheOnce.Do(func() {
		hotCache = &HotCache{
			lock:  &sync.RWMutex{},
			cache: make(map[Source]*HotCacheItem, len),
		}
	})
	return hotCache
}

func GetHotCache() *HotCache {
	return hotCache
}

func (c *HotCache) Get(source Source) (*HotCacheItem, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.cache[source]
	if !ok {
		return &HotCacheItem{}, false

	}
	return item, true
}

func (c *HotCache) GetExpired(t time.Duration) []Source {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var res []Source
	for k, v := range c.cache {
		if time.Since(v.LastUpdateTime) > t {
			res = append(res, k)
		}
	}
	return res
}

func (c *HotCache) Set(source Source, data *HotData) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if data != nil {
		c.cache[source] = &HotCacheItem{
			LastUpdateTime: time.Now(),
			HotData:        data,
		}
	}
}
