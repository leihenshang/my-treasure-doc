package hot

import (
	"sync"
	"time"
)

type HotCache struct {
	lock  *sync.RWMutex
	Cache map[Source]*HotCacheItem
}

type HotCacheItem struct {
	LastUpdateTime time.Time
	HotData        *HotData
}

func NewHotCache() *HotCache {
	return &HotCache{
		Cache: make(map[Source]*HotCacheItem),
	}
}

func (c *HotCache) Get(source Source) (*HotCacheItem, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.Cache[source]
	if !ok {
		return &HotCacheItem{}, false

	}
	return item, true
}

func (c *HotCache) Set(source Source, data *HotData) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if data != nil {
		c.Cache[source] = &HotCacheItem{
			LastUpdateTime: time.Now(),
			HotData:        data,
		}
	}
}
