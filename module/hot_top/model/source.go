package model

import (
	"encoding/json"
	"time"
)

type Source string

const (
	SourceITHome        Source = "ithome"
	SourceZhihu         Source = "zhihu"
	SourceWeibo         Source = "weibo"
	SourceBilibili      Source = "bilibili"
	SourceBaidu         Source = "baidu"
	SourceV2EX          Source = "v2ex"
	SourceGitHub        Source = "github"
	SourceDouyin        Source = "douyin"
	SourceKuaishou      Source = "kuaishou"
	SourceToutiao       Source = "toutiao"
	SourceJuejin        Source = "juejin"
	Source36Kr          Source = "36kr"
	SourceCSDN          Source = "csdn"
	SourceTieba         Source = "tieba"
	SourceZhihuDaily    Source = "zhihu-daily"
	SourceCoolapk       Source = "coolapk"
	SourceHupu          Source = "hupu"
	SourceHuxiu         Source = "huxiu"
	SourceJianshu       Source = "jianshu"
	SourceSmzdm         Source = "smzdm"
	SourceSspai         Source = "sspai"
	SourceNeteaseNews   Source = "netease-news"
	SourceQQNews        Source = "qq-news"
	SourceAcfun         Source = "acfun"
	Source51CTO         Source = "51cto"
	Source52Pojie       Source = "52pojie"
	SourceDoubanGroup   Source = "douban-group"
	SourceDgtle         Source = "dgtle"
	SourceDoubanMovie   Source = "douban-movie"
	SourceEarthquake    Source = "earthquake"
	SourceGameres       Source = "gameres"
	SourceGeekpark      Source = "geekpark"
	SourceGenshin       Source = "genshin"
	SourceGuokr         Source = "guokr"
	SourceHackernews    Source = "hackernews"
	SourceHelloGitHub   Source = "hellogithub"
	SourceHistory       Source = "history"
	SourceHonkai        Source = "honkai"
	SourceHostloc       Source = "hostloc"
	SourceIfanr         Source = "ifanr"
	SourceIthomeXijiayi Source = "ithome-xijiayi"
	SourceMiyoushe      Source = "miyoushe"
	SourceNewsmth       Source = "newsmth"
	SourceNgabbs        Source = "ngabbs"
	SourceNodeseek      Source = "nodeseek"
	SourceNytimes       Source = "nytimes"
	SourceProducthunt   Source = "producthunt"
	SourceSinaNews      Source = "sina-news"
	SourceSina          Source = "sina"
	SourceStarrail      Source = "starrail"
	SourceThepaper      Source = "thepaper"
	SourceWeatheralarm  Source = "weatheralarm"
	SourceWeread        Source = "weread"
	SourceYystv         Source = "yystv"
)

type HotData struct {
	Code        int       `json:"code"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Total       int       `json:"total"`
	Data        HotItems  `json:"data"`
	Params      any       `json:"params,omitempty"`
	UpdateTime  time.Time `json:"updateTime,omitempty"`
}

type HotItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Timestamp int64  `json:"timestamp"`
	Hot       int    `json:"hot"`
	URL       string `json:"url"`
	MobileURL string `json:"mobileUrl"`
	Author    string `json:"author,omitempty"`
	Desc      string `json:"desc,omitempty"`
}

type HotItems []*HotItem

func (h HotItems) ToJsonWithErr() string {
	jsonBytes, err := json.Marshal(h)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}

func (s Source) String() string {
	return string(s)
}

func (h *HotData) IsUpdateTimeExpired(t time.Duration) bool {
	return time.Since(h.UpdateTime) > t
}
