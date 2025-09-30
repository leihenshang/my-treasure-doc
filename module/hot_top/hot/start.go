package hot

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	hotcache "fastduck/treasure-doc/module/hot_top/hot/hot_cache"
	"fastduck/treasure-doc/module/hot_top/model"
	"fmt"
	"sync"
	"time"

	"log"
)

type Hot struct {
	HotConf                  *conf.Hot
	LastHotConf              *conf.Hot
	collectRefreshTicker     *time.Timer
	collectRefreshTickerLock *sync.Mutex
}

var onceHot = &sync.Once{}
var hot *Hot

func NewHot(hotConf *conf.Hot) *Hot {
	onceHot.Do(func() {
		hot = &Hot{
			HotConf:                  hotConf,
			collectRefreshTickerLock: &sync.Mutex{},
		}
		copy := *hotConf
		hot.LastHotConf = &copy
	})
	return hot
}

func GetHot() *Hot {
	return hot
}

func (h *Hot) Start() error {
	if _, err := NewSpider(); err != nil {
		return fmt.Errorf("NewSpider failed, err: %v", err)
	}
	hotcache.NewHotMemCache(len(conf.HotConfListMap))
	hotcache.NewHotFileCache(h.HotConf.HotFileCachePath)
	go h.tickerGetHot()
	return nil
}

func (h *Hot) tickerGetHot() {
	var collectSources []conf.Source
	for _, v := range conf.HotConfListMap {
		if resp, err := hotcache.GetHotFileCache().Get(v.Source); err != nil {
			log.Printf("get [%s] from file cache failed, err: %v\n", string(v.Source), err)
		} else if resp != nil {
			log.Printf("get [%s] from file cache success, dataLen: %d\n", string(v.Source), len(resp.Data))
			hotcache.GetHotMemCache().Set(v.Source, resp)
		} else {
			collectSources = append(collectSources, v.Source)
		}
	}

	h.setHotCache(h.collectHotMap(collectSources...))
	tk := time.NewTicker(h.HotConf.ExpiredCheckIntervalParsed)
	defer tk.Stop()
	for range tk.C {
		if h.HotConf.ExpiredCheckIntervalParsed != h.LastHotConf.ExpiredCheckIntervalParsed {
			tk.Reset(h.HotConf.ExpiredCheckIntervalParsed)
			h.LastHotConf.ExpiredCheckIntervalParsed = h.HotConf.ExpiredCheckIntervalParsed
			log.Printf("check TickerGetHot expire, reset interval to %s, about: %v s\n",
				h.HotConf.ExpiredCheckIntervalParsed, h.HotConf.ExpiredCheckIntervalParsed.Seconds())
		}
		log.Printf("check TickerGetHot expire,HotPullIntervalParsed is: %v, ExpiredCheckIntervalParsed is: %v \n",
			h.HotConf.HotPullIntervalParsed, h.HotConf.ExpiredCheckIntervalParsed)
		h.setHotCache(h.collectHotMap(h.getExpiredHotSources()...))
	}
}

func (h *Hot) getExpiredHotSources() (res []conf.Source) {
	cacheMap := hotcache.GetHotMemCache().GetAllMap()
	for _, v := range conf.HotConfListMap {
		if hotCache, ok := cacheMap[v.Source]; !ok {
			res = append(res, v.Source)
		} else if hotCache != nil && hotCache.IsUpdateTimeExpired(h.HotConf.HotPullIntervalParsed) {
			res = append(res, v.Source)
		}
	}
	return res
}

func (h *Hot) collectHotMap(sources ...conf.Source) map[conf.Source]*model.HotData {
	res := make(map[conf.Source]*model.HotData, len(sources))
	for _, k := range sources {
		// TODO: perfect here
		if hotConf, ok := conf.HotConfListMap[k]; ok && hotConf.Disabled {
			log.Printf("source: [%s] is disabled, skip\n", string(k))
			continue
		}
		hotData, err := GetSpider().GetHotBySource(k)
		if err != nil {
			log.Printf("get [%s] failed, err: %v,using default values to fill in\n", string(k), err)
		}
		res[k] = hotData
	}
	return res
}

func (h *Hot) setHotCache(hotMap map[conf.Source]*model.HotData) {
	for k, hotData := range hotMap {
		hotcache.GetHotMemCache().Set(k, hotData)
		if err := hotcache.GetHotFileCache().Set(k, hotData); err != nil {
			log.Printf("set [%s] cache save file failed, err: %v\n", k, err)
		}
	}
}

func (h *Hot) RefreshHotCache(source conf.Source) *model.HotData {
	if h.collectRefreshTicker != nil {
		h.collectRefreshTicker.Stop()
	}
	h.collectRefreshTicker = time.AfterFunc(time.Second*10, func() {
		log.Printf("use timer afterFunc to refresh hot cache for source: %s\n", source)
		resMap := h.collectHotMap(source)
		h.setHotCache(resMap)
	})
	hotData, _ := hotcache.GetHotMemCache().Get(source)
	return hotData
}
