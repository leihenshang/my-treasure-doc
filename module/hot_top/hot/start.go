package hot

import (
	"fastduck/treasure-doc/module/hot_top/model"
	"fmt"
	"sync"
	"time"
)

type Hot struct {
	HotExpiredTime time.Duration
}

var onceHot = &sync.Once{}
var hot *Hot

func NewHot(hotExpiredTime time.Duration) *Hot {
	onceHot.Do(func() {
		hot = &Hot{
			HotExpiredTime: hotExpiredTime,
		}
	})
	return hot
}

func (h *Hot) Start() {
	NewSpider()
	NewHotCache(len(UrlConfMap))
	go TickerGetHot(time.Hour)
}

func TickerGetHot(expireTime time.Duration) {
	fmt.Println("TickerGetHot start!")
	var sources []model.Source
	for k := range UrlConfMap {
		sources = append(sources, k)
	}
	setHotCacheBySource(sources)

	tk := time.NewTicker(time.Second * 10)
	defer tk.Stop()
	for t := range tk.C {
		current := t.Format(time.DateTime)
		fmt.Printf("check TickerGetHot expire time: %s\n", current)
		setHotCacheBySource(GetHotCache().GetExpired(expireTime))
	}
}

func setHotCacheBySource(sources []model.Source) {
	for _, k := range sources {
		resp, err := GetHotBySource(k)
		if err != nil {
			fmt.Printf("TickerGet [%s] failed, err: %v\n", k, err)
			continue
		}
		fmt.Printf("TickerGet [%s] success, dataLen: %d\n", k, len(resp.Data))
		GetHotCache().Set(k, resp)
	}
}

func GetHotBySource(k model.Source) (*model.HotData, error) {
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
	case model.SourceNetease:
		return GetSpider().GetNetease()
	case model.SourceQQ:
		return GetSpider().GetQQ()
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
