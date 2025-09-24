package conf

import "fastduck/treasure-doc/module/hot_top/model"

type UrlConf struct {
	Source   model.Source
	Disabled bool
}

type UrlConfs []*UrlConf

var UrlList = UrlConfs{
	{Source: model.SourceLol},
	{Source: model.SourceITHome},
	{Source: model.SourceZhihu},
	{Source: model.SourceWeibo},
	{Source: model.SourceBilibili},
	{Source: model.SourceBaidu},
	{Source: model.SourceV2EX},
	{Source: model.SourceGitHub},
	{Source: model.SourceDouyin},
	{Source: model.SourceKuaishou, Disabled: true},
	{Source: model.SourceToutiao},
	{Source: model.SourceJuejin},
	{Source: model.Source36Kr},
	{Source: model.SourceCSDN},
	{Source: model.SourceTieba},
	{Source: model.SourceZhihuDaily},
	{Source: model.SourceCoolapk},
	{Source: model.SourceHupu},
	{Source: model.SourceHuxiu},
	{Source: model.SourceJianshu},
	{Source: model.SourceSmzdm},
	{Source: model.SourceSspai},
	{Source: model.SourceNeteaseNews},
	{Source: model.SourceQQNews},
	{Source: model.SourceAcfun},
	{Source: model.Source51CTO},
	{Source: model.Source52Pojie},
	{Source: model.SourceDoubanGroup},
	{Source: model.SourceDgtle},
	{Source: model.SourceDoubanMovie},
	{Source: model.SourceEarthquake, Disabled: true},
	{Source: model.SourceGameres},
	{Source: model.SourceGeekpark},
	{Source: model.SourceGenshin},
	{Source: model.SourceGuokr},
	{Source: model.SourceHackernews},
	{Source: model.SourceHelloGitHub},
	{Source: model.SourceHistory},
	{Source: model.SourceHonkai},
	{Source: model.SourceHostloc},
	{Source: model.SourceIfanr},
	{Source: model.SourceMiyoushe},
	{Source: model.SourceNewsmth},
	{Source: model.SourceNgabbs},
	{Source: model.SourceNodeseek},
	{Source: model.SourceNytimes},
	{Source: model.SourceProducthunt},
	{Source: model.SourceSinaNews},
	{Source: model.SourceSina},
	{Source: model.SourceStarrail},
	{Source: model.SourceThepaper},
	{Source: model.SourceWeatheralarm},
	{Source: model.SourceWeread},
	{Source: model.SourceYystv},
	{Source: model.SourceIthomeXijiayi},
}
