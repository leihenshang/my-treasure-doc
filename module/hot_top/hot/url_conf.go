package hot

import (
	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/model"
)

var UrlConfMap map[model.Source]*conf.UrlConf = map[model.Source]*conf.UrlConf{
	model.SourceLol: {
		Url:   "https://lol.qq.com/web/",
		Agent: `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`,
	},
	model.SourceITHome:   {},
	model.SourceZhihu:    {},
	model.SourceWeibo:    {},
	model.SourceBilibili: {},
	model.SourceBaidu:    {},
	model.SourceV2EX:     {},
	model.SourceGitHub:   {},
	model.SourceDouyin:   {},
	model.SourceKuaishou: {
		Disabled: true,
	},
	model.SourceToutiao: {},
	model.SourceJuejin:  {},
	model.Source36Kr:    {},
	model.SourceCSDN:    {},
	model.SourceTieba:   {},
	model.SourceZhihuDaily: {
		Url:   "https://daily.zhihu.com/api/4/news/latest",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	model.SourceCoolapk:     {},
	model.SourceHupu:        {},
	model.SourceHuxiu:       {},
	model.SourceJianshu:     {},
	model.SourceSmzdm:       {},
	model.SourceSspai:       {},
	model.SourceNeteaseNews: {},
	model.SourceQQNews:      {},
	model.SourceAcfun:       {},
	model.Source51CTO: {
	},
	model.Source52Pojie: {
		Url:   "https://www.52pojie.cn/forum.php?mod=guide&view=digest&rss=1",
		Agent: `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36`,
	},
	model.SourceDoubanGroup: {
		Url:   "https://www.douban.com/group/explore",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	model.SourceDgtle: {
		Url:   "https://opser.api.dgtle.com/v2/news/index",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceDoubanMovie: {
		Url:   "https://movie.douban.com/chart/",
		Agent: `Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15`,
	},
	model.SourceEarthquake: {
		Disabled: true,
		Url:      "https://news.ceic.ac.cn/speedsearch.html",
		Agent:    `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceGameres: {
		Url:   "https://www.gameres.com",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceGeekpark: {
		Url:   "https://mainssl.geekpark.net/api/v2",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceGenshin: {
		Url:   "https://bbs-api-static.miyoushe.com/painter/wapi/getNewsList?client_type=4&gids=2&page_size=20",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceGuokr: {
		Url:   "https://www.guokr.com/beta/proxy/science_api/articles?limit=30",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0`,
	},
	model.SourceHackernews: {
		Url:   "https://news.ycombinator.com/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36`,
	},
	model.SourceHelloGitHub: {
		Url:   "https://abroad.hellogithub.com/v1/?sort_by=featured&tid=&page=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceHistory: {
		Url:   "https://baike.baidu.com/cms/home/eventsOnHistory/01.json",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceHonkai: {
		Url:   "https://bbs-api-static.miyoushe.com/painter/wapi/getNewsList?client_type=4&gids=1&page_size=20",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceHostloc: {
		Url:   "https://hostloc.com/forum.php?mod=guide&view=hot&rss=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceIfanr: {
		Url:   "https://sso.ifanr.com/api/v5/wp/buzz/?limit=20&offset=0",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceMiyoushe: {
		Url:   "https://bbs-api-static.miyoushe.com/painter/wapi/getNewsList?client_type=4&gids=1&page_size=30",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceNewsmth: {
		Url:   "https://wap.newsmth.net/wap/api/hot/global",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceNgabbs: {
		Url:   "https://ngabbs.com/nuke.php?__lib=load_topic&__act=load_topic_reply_ladder2&opt=1&all=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceNodeseek: {
		Url:   "https://www.nodeseek.com",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceNytimes: {
		Url:   "https://rsshub.app/nytimes",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceProducthunt: {
		Url:   "https://www.producthunt.com/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceSinaNews: {
		Url:   "https://news.sina.com.cn/zt_d/top_news/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceSina: {
		Url:   "https://s.weibo.com/top/summary",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceStarrail: {
		Url:   "https://bbs-api-static.miyoushe.com/painter/wapi/getNewsList?client_type=4&gids=6&page_size=20",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceThepaper: {
		Url:   "https://cache.thepaper.cn/contentapi/wwwIndex/rightSidebar",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceWeatheralarm: {
		Url:   "http://www.nmc.cn/rest/findAlarm",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceWeread: {
		Url:   "https://weread.qq.com/web/bookListInCategory/rising?rank=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceYystv: {
		Url:   "https://www.yystv.cn/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67`,
	},
	model.SourceIthomeXijiayi: {},
}
