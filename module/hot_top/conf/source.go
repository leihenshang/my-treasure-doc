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

var HotConfListMap = map[Source]*HotConf{
	SourceLol:           {Source: SourceLol},
	SourceITHome:        {Source: SourceITHome},
	SourceZhihu:         {Source: SourceZhihu},
	SourceWeibo:         {Source: SourceWeibo},
	SourceBilibili:      {Source: SourceBilibili},
	SourceBaidu:         {Source: SourceBaidu},
	SourceV2EX:          {Source: SourceV2EX},
	SourceGitHub:        {Source: SourceGitHub},
	SourceDouyin:        {Source: SourceDouyin},
	SourceKuaishou:      {Source: SourceKuaishou, Disabled: true},
	SourceToutiao:       {Source: SourceToutiao},
	SourceJuejin:        {Source: SourceJuejin},
	Source36Kr:          {Source: Source36Kr},
	SourceCSDN:          {Source: SourceCSDN},
	SourceTieba:         {Source: SourceTieba},
	SourceZhihuDaily:    {Source: SourceZhihuDaily},
	SourceCoolapk:       {Source: SourceCoolapk},
	SourceHupu:          {Source: SourceHupu},
	SourceHuxiu:         {Source: SourceHuxiu},
	SourceJianshu:       {Source: SourceJianshu},
	SourceSmzdm:         {Source: SourceSmzdm},
	SourceSspai:         {Source: SourceSspai},
	SourceNeteaseNews:   {Source: SourceNeteaseNews},
	SourceQQNews:        {Source: SourceQQNews},
	SourceAcfun:         {Source: SourceAcfun},
	Source51CTO:         {Source: Source51CTO},
	Source52Pojie:       {Source: Source52Pojie},
	SourceDoubanGroup:   {Source: SourceDoubanGroup},
	SourceDgtle:         {Source: SourceDgtle},
	SourceDoubanMovie:   {Source: SourceDoubanMovie},
	SourceEarthquake:    {Source: SourceEarthquake, Disabled: true},
	SourceGameres:       {Source: SourceGameres},
	SourceGeekpark:      {Source: SourceGeekpark},
	SourceGenshin:       {Source: SourceGenshin},
	SourceGuokr:         {Source: SourceGuokr},
	SourceHackernews:    {Source: SourceHackernews},
	SourceHelloGitHub:   {Source: SourceHelloGitHub},
	SourceHistory:       {Source: SourceHistory},
	SourceHonkai:        {Source: SourceHonkai},
	SourceHostloc:       {Source: SourceHostloc},
	SourceIfanr:         {Source: SourceIfanr},
	SourceMiyoushe:      {Source: SourceMiyoushe},
	SourceNewsmth:       {Source: SourceNewsmth},
	SourceNgabbs:        {Source: SourceNgabbs},
	SourceNodeseek:      {Source: SourceNodeseek},
	SourceNytimes:       {Source: SourceNytimes},
	SourceProducthunt:   {Source: SourceProducthunt},
	SourceSinaNews:      {Source: SourceSinaNews},
	SourceSina:          {Source: SourceSina},
	SourceStarrail:      {Source: SourceStarrail},
	SourceThepaper:      {Source: SourceThepaper},
	SourceWeatheralarm:  {Source: SourceWeatheralarm},
	SourceWeread:        {Source: SourceWeread},
	SourceYystv:         {Source: SourceYystv},
	SourceIthomeXijiayi: {Source: SourceIthomeXijiayi},
}
