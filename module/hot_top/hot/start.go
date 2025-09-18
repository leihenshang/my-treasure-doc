package hot

import (
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
	var sources []Source
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

func setHotCacheBySource(sources []Source) {
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

func GetHotBySource(k Source) (*HotData, error) {
	switch k {
	case SourceITHome:
		return GetSpider().GetItHome()
	case SourceZhihu:
		return GetSpider().GetZhihu()
	case SourceWeibo:
		return GetSpider().GetWeibo()
	case SourceBilibili:
		return GetSpider().GetBilibili()
	case SourceBaidu:
		return GetSpider().GetBaidu()
	case SourceV2EX:
		return GetSpider().GetV2EX()
	case SourceGitHub:
		return GetSpider().GetGitHub()
	case SourceDouyin:
		return GetSpider().GetDouyin()
	case SourceKuaishou:
		return GetSpider().GetKuaishou()
	case SourceToutiao:
		return GetSpider().GetToutiao()
	case SourceJuejin:
		return GetSpider().GetJuejin()
	case Source36Kr:
		return GetSpider().Get36Kr()
	case SourceCSDN:
		return GetSpider().GetCSDN()
	case SourceTieba:
		return GetSpider().GetTieba()
	case SourceZhihuDaily:
		return GetSpider().GetZhihuDaily()
	case SourceCoolapk:
		return GetSpider().GetCoolapk()
	case SourceHupu:
		return GetSpider().GetHupu()
	case SourceHuxiu:
		return GetSpider().GetHuxiu()
	case SourceJianshu:
		return GetSpider().GetJianshu()
	case SourceSmzdm:
		return GetSpider().GetSmzdm()
	case SourceSspai:
		return GetSpider().GetSspai()
	case SourceNetease:
		return GetSpider().GetNetease()
	case SourceQQ:
		return GetSpider().GetQQ()
	case SourceAcfun:
		return GetSpider().GetAcfun()
	case Source51CTO:
		return GetSpider().Get51CTO()
	case Source52Pojie:
		return GetSpider().Get52Pojie()
	case SourceDoubanGroup:
		return GetSpider().GetDoubanGroup()
	case SourceDgtle:
		return GetSpider().GetDgtle()
	case SourceDoubanMovie:
		return GetSpider().GetDoubanMovie()
	case SourceEarthquake:
		return GetSpider().GetEarthquake()
	case SourceGameres:
		return GetSpider().GetGameres()
	case SourceGeekpark:
		return GetSpider().GetGeekpark()
	case SourceGenshin:
		return GetSpider().GetGenshin()
	case SourceGuokr:
		return GetSpider().GetGuokr()
	case SourceHackernews:
		return GetSpider().GetHackernews()
	case SourceHelloGitHub:
		return GetSpider().GetHelloGitHub()
	case SourceHistory:
		return GetSpider().GetHistory()
	case SourceHonkai:
		return GetSpider().GetHonkai()
	case SourceHostloc:
		return GetSpider().GetHostloc()
	case SourceIfanr:
		return GetSpider().GetIfanr()
	case SourceIthomeXijiayi:
		return GetSpider().GetIthomeXijiayi()
	case SourceMiyoushe:
		return GetSpider().GetMiyoushe()
	case SourceNewsmth:
		return GetSpider().GetNewsmth()
	case SourceNgabbs:
		return GetSpider().GetNgabbs()
	case SourceNodeseek:
		return GetSpider().GetNodeseek()
	case SourceNytimes:
		return GetSpider().GetNytimes()
	case SourceProducthunt:
		return GetSpider().GetProducthunt()
	case SourceSinaNews:
		return GetSpider().GetSinaNews()
	case SourceSina:
		return GetSpider().GetSina()
	case SourceStarrail:
		return GetSpider().GetStarrail()
	case SourceThepaper:
		return GetSpider().GetThepaper()
	case SourceWeatheralarm:
		return GetSpider().GetWeatheralarm()
	case SourceWeread:
		return GetSpider().GetWeread()
	case SourceYystv:
		return GetSpider().GetYystv()
	default:
		return nil, nil
	}
}
