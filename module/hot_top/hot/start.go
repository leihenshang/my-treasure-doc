package hot

import (
	"encoding/json"
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/model"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"log"
)

type Hot struct {
	HotConf     *conf.Hot
	LastHotConf *conf.Hot
}

var onceHot = &sync.Once{}
var hot *Hot

func NewHot(hotConf *conf.Hot) *Hot {
	onceHot.Do(func() {
		hot = &Hot{
			HotConf: hotConf,
		}
		copy := *hotConf
		hot.LastHotConf = &copy
	})
	return hot
}

func (h *Hot) Start() error {
	if _, err := NewSpider(); err != nil {
		return fmt.Errorf("NewSpider failed, err: %v", err)
	}

	NewHotCache(len(conf.HotConfListMap))
	go h.TickerGetHot()
	return nil
}

func (h *Hot) TickerGetHot() {
	var collectSources []conf.Source
	for _, v := range conf.HotConfListMap {
		if resp, err := h.GetHotFromFileCache(v.Source); err != nil {
			log.Printf("get [%s] from file cache failed, err: %v\n", string(v.Source), err)
		} else if resp != nil {
			log.Printf("get [%s] from file cache success, dataLen: %d\n", string(v.Source), len(resp.Data))
			GetHotCache().Set(v.Source, resp)
		} else {
			collectSources = append(collectSources, v.Source)
		}
	}

	h.setHotCache(h.getHotMap(collectSources))
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
		h.setHotCache(h.getHotMap(h.GetExpiredHotSources()))
	}
}

func (h *Hot) GetExpiredHotSources() (res []conf.Source) {
	cacheMap := GetHotCache().GetAllMap()
	for _, v := range conf.HotConfListMap {
		if hotCache, ok := cacheMap[v.Source]; !ok {
			res = append(res, v.Source)
		} else if hotCache != nil && hotCache.IsUpdateTimeExpired(h.HotConf.HotPullIntervalParsed) {
			res = append(res, v.Source)
		}
	}
	return res
}

func (h *Hot) getHotMap(sources []conf.Source) map[conf.Source]*model.HotData {
	res := make(map[conf.Source]*model.HotData, len(sources))
	for _, k := range sources {
		// TODO: perfect here
		if hotConf, ok := conf.HotConfListMap[k]; ok && hotConf.Disabled {
			log.Printf("source: [%s] is disabled, skip\n", string(k))
			continue
		}
		if hotData, err := GetSpider().GetHotBySource(k); err != nil {
			log.Printf("get [%s] failed, err: %v,using default values to fill in\n", string(k), err)
			res[k] = &model.HotData{
				Name:       string(k),
				UpdateTime: time.Now(),
				Data:       []*model.HotItem{},
			}
		} else {
			res[k] = hotData
		}
	}
	return res
}

func (h *Hot) setHotCache(hotMap map[conf.Source]*model.HotData) {
	for k, hotData := range hotMap {
		GetHotCache().Set(k, hotData)
		if err := h.SaveHotToFileCache(k, hotData); err != nil {
			log.Printf("set [%s] cache save file failed, err: %v\n", k, err)
		}
	}
}

func (h *Hot) SaveHotToFileCache(source conf.Source, resp *model.HotData) error {
	if resp == nil {
		return fmt.Errorf("source: [%s], resp is nil", source)
	}
	savePath := filepath.Join(h.HotConf.HotFileCachePath, fmt.Sprintf("%s.json", source))
	if _, err := os.Stat(filepath.Dir(savePath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
			return fmt.Errorf("source: [%s], mkdir failed, err: %v", source, err)
		}
	}

	if f, err := os.Create(savePath); err != nil {
		return fmt.Errorf("source: [%s], create file failed, err: %v", source, err)
	} else {
		if err := json.NewEncoder(f).Encode(resp); err != nil {
			return fmt.Errorf("source: [%s], save file failed, err: %v", source, err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("source: [%s], close file failed, err: %v", source, err)
		}
	}

	return nil
}

func (h *Hot) GetHotFromFileCache(source conf.Source) (resp *model.HotData, err error) {
	resp = &model.HotData{}
	savePath := filepath.Join(h.HotConf.HotFileCachePath, fmt.Sprintf("%s.json", source))
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		return nil, nil
	}
	if f, err := os.Open(savePath); err != nil {
		return nil, fmt.Errorf("source: [%s], open file failed, err: %v", source, err)
	} else {
		defer f.Close()
		if err := json.NewDecoder(f).Decode(resp); err != nil {
			return nil, fmt.Errorf("source: [%s], decode file failed, err: %v", source, err)
		}
	}
	return resp, nil
}
