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
	if _, err := NewSpider(UrlConfMap); err != nil {
		return fmt.Errorf("NewSpider failed, err: %v", err)
	}

	NewHotCache(len(UrlConfMap))
	go h.TickerGetHot()
	return nil
}

func (h *Hot) TickerGetHot() {
	var collectSources []model.Source
	for k := range UrlConfMap {
		if resp, err := h.GetHotFromFileCache(k); err != nil {
			log.Printf("get [%s] from file cache failed, err: %v\n", string(k), err)
		} else if resp != nil {
			log.Printf("get [%s] from file cache success, dataLen: %d\n", string(k), len(resp.Data))
			GetHotCache().Set(k, resp)
		} else {
			collectSources = append(collectSources, k)
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

func (h *Hot) GetExpiredHotSources() (res []model.Source) {
	cacheMap := GetHotCache().GetAllMap()
	for k := range UrlConfMap {
		if hotCache, ok := cacheMap[k]; !ok {
			res = append(res, k)
		} else if hotCache != nil && hotCache.IsUpdateTimeExpired(h.HotConf.HotPullIntervalParsed) {
			res = append(res, k)
		}
	}
	return res
}

func (h *Hot) getHotMap(sources []model.Source) map[model.Source]*model.HotData {
	res := make(map[model.Source]*model.HotData, len(sources))
	for _, k := range sources {
		if hotData, err := h.GetHotBySource(k); err != nil {
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

func (h *Hot) setHotCache(hotMap map[model.Source]*model.HotData) {
	for k, hotData := range hotMap {
		GetHotCache().Set(k, hotData)
		if err := h.SaveHotToFileCache(k, hotData); err != nil {
			log.Printf("set [%s] cache save file failed, err: %v\n", k, err)
		}
	}
}

func (h *Hot) SaveHotToFileCache(source model.Source, resp *model.HotData) error {
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

func (h *Hot) GetHotFromFileCache(source model.Source) (resp *model.HotData, err error) {
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

func (h *Hot) GetHotBySource(k model.Source) (*model.HotData, error) {
	UrlConf, ok := UrlConfMap[k]
	if !ok {
		return nil, fmt.Errorf("source: [%s], url conf not found", k)
	} else if UrlConf.Disabled {
		return nil, fmt.Errorf("source: [%s], url conf disabled, skip", k)
	}

	switch k {
	case model.SourceITHome:
		return GetSpider().GetItHome()
	case model.SourceZhihu:
		return GetSpider().GetZhihu()
	case model.SourceWeibo:
		return GetSpider().GetWeibo()
	case model.SourceBilibili:
		return GetSpider().GetBilibili()
	case model.SourceBaidu:
		return GetSpider().GetBaidu()
	case model.SourceV2EX:
		return GetSpider().GetV2EX()
	case model.SourceGitHub:
		return GetSpider().GetGitHub()
	case model.SourceDouyin:
		return GetSpider().GetDouyin()
	case model.SourceKuaishou:
		return GetSpider().GetKuaishou()
	case model.SourceToutiao:
		return GetSpider().GetToutiao()
	case model.SourceJuejin:
		return GetSpider().GetJuejin()
	case model.Source36Kr:
		return GetSpider().Get36Kr()
	case model.SourceCSDN:
		return GetSpider().GetCSDN()
	case model.SourceTieba:
		return GetSpider().GetTieba()
	case model.SourceZhihuDaily:
		return GetSpider().GetZhihuDaily()
	case model.SourceCoolapk:
		return GetSpider().GetCoolapk()
	case model.SourceHupu:
		return GetSpider().GetHupu()
	case model.SourceHuxiu:
		return GetSpider().GetHuxiu()
	case model.SourceJianshu:
		return GetSpider().GetJianshu()
	case model.SourceSmzdm:
		return GetSpider().GetSmzdm()
	case model.SourceSspai:
		return GetSpider().GetSspai()
	case model.SourceNeteaseNews:
		return GetSpider().GetNeteaseNews()
	case model.SourceQQNews:
		return GetSpider().GetQQNews()
	case model.SourceAcfun:
		return GetSpider().GetAcfun()
	case model.Source51CTO:
		return GetSpider().Get51CTO()
	case model.Source52Pojie:
		return GetSpider().Get52Pojie()
	case model.SourceDoubanGroup:
		return GetSpider().GetDoubanGroup()
	case model.SourceDgtle:
		return GetSpider().GetDgtle()
	case model.SourceDoubanMovie:
		return GetSpider().GetDoubanMovie()
	case model.SourceEarthquake:
		return GetSpider().GetEarthquake()
	case model.SourceGameres:
		return GetSpider().GetGameres()
	case model.SourceGeekpark:
		return GetSpider().GetGeekpark()
	case model.SourceGenshin:
		return GetSpider().GetGenshin()
	case model.SourceGuokr:
		return GetSpider().GetGuokr()
	case model.SourceHackernews:
		return GetSpider().GetHackernews()
	case model.SourceHelloGitHub:
		return GetSpider().GetHelloGitHub()
	case model.SourceHistory:
		return GetSpider().GetHistory()
	case model.SourceHonkai:
		return GetSpider().GetHonkai()
	case model.SourceHostloc:
		return GetSpider().GetHostloc()
	case model.SourceIfanr:
		return GetSpider().GetIfanr()
	case model.SourceIthomeXijiayi:
		return GetSpider().GetIthomeXijiayi()
	case model.SourceMiyoushe:
		return GetSpider().GetMiyoushe()
	case model.SourceNewsmth:
		return GetSpider().GetNewsmth()
	case model.SourceNgabbs:
		return GetSpider().GetNgabbs()
	case model.SourceNodeseek:
		return GetSpider().GetNodeseek()
	case model.SourceNytimes:
		return GetSpider().GetNytimes()
	case model.SourceProducthunt:
		return GetSpider().GetProducthunt()
	case model.SourceSinaNews:
		return GetSpider().GetSinaNews()
	case model.SourceSina:
		return GetSpider().GetSina()
	case model.SourceStarrail:
		return GetSpider().GetStarrail()
	case model.SourceThepaper:
		return GetSpider().GetThepaper()
	case model.SourceWeatheralarm:
		return GetSpider().GetWeatheralarm()
	case model.SourceWeread:
		return GetSpider().GetWeread()
	case model.SourceYystv:
		return GetSpider().GetYystv()
	default:
		return nil, nil
	}
}
