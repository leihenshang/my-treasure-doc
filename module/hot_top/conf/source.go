package conf

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
	SourceLol           Source = "lol"
)

func (s Source) String() string {
	return string(s)
}

type HotConf struct {
	Source   Source
	Disabled bool
}

var UrlList = []*HotConf{
	{Source: SourceLol},
	{Source: SourceITHome},
	{Source: SourceZhihu},
	{Source: SourceWeibo},
	{Source: SourceBilibili},
	{Source: SourceBaidu},
	{Source: SourceV2EX},
	{Source: SourceGitHub},
	{Source: SourceDouyin},
	{Source: SourceKuaishou, Disabled: true},
	{Source: SourceToutiao},
	{Source: SourceJuejin},
	{Source: Source36Kr},
	{Source: SourceCSDN},
	{Source: SourceTieba},
	{Source: SourceZhihuDaily},
	{Source: SourceCoolapk},
	{Source: SourceHupu},
	{Source: SourceHuxiu},
	{Source: SourceJianshu},
	{Source: SourceSmzdm},
	{Source: SourceSspai},
	{Source: SourceNeteaseNews},
	{Source: SourceQQNews},
	{Source: SourceAcfun},
	{Source: Source51CTO},
	{Source: Source52Pojie},
	{Source: SourceDoubanGroup},
	{Source: SourceDgtle},
	{Source: SourceDoubanMovie},
	{Source: SourceEarthquake, Disabled: true},
	{Source: SourceGameres},
	{Source: SourceGeekpark},
	{Source: SourceGenshin},
	{Source: SourceGuokr},
	{Source: SourceHackernews},
	{Source: SourceHelloGitHub},
	{Source: SourceHistory},
	{Source: SourceHonkai},
	{Source: SourceHostloc},
	{Source: SourceIfanr},
	{Source: SourceMiyoushe},
	{Source: SourceNewsmth},
	{Source: SourceNgabbs},
	{Source: SourceNodeseek},
	{Source: SourceNytimes},
	{Source: SourceProducthunt},
	{Source: SourceSinaNews},
	{Source: SourceSina},
	{Source: SourceStarrail},
	{Source: SourceThepaper},
	{Source: SourceWeatheralarm},
	{Source: SourceWeread},
	{Source: SourceYystv},
	{Source: SourceIthomeXijiayi},
}
