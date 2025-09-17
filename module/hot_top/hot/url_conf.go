package hot

var UrlConfMap map[Source]*UrlConf = map[Source]*UrlConf{
	SourceITHome: {
		Url:   "https://m.ithome.com/rankm/",
		Agent: `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`,
	},
	SourceZhihu: {
		Url:   "https://api.zhihu.com/topstory/hot-lists/total?limit=50",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceWeibo: {
		Url:   "https://m.weibo.cn/api/container/getIndex?containerid=106003type%3D25%26t%3D3%26disable_hot%3D1%26filter_type%3Drealtimehot&title=%E5%BE%AE%E5%8D%9A%E7%83%AD%E6%90%9C&extparam=filter_type%3Drealtimehot%26mi_cid%3D100103%26pos%3D0_0%26c_type%3D30%26display_time%3D1540538388&luicode=10000011&lfid=231583",
		Agent: `Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1`,
	},
	SourceBilibili: {
		Url:   "https://api.bilibili.com/x/web-interface/ranking/v2?rid=0&type=all",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceBaidu: {
		Url:   "https://top.baidu.com/board?tab=realtime",
		Agent: `Mozilla/5.0 (iPhone; CPU iPhone OS 14_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/1.0 Mobile/12F69 Safari/605.1.15`,
	},
	SourceV2EX: {
		Url:   "https://www.v2ex.com/api/topics/hot.json",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceGitHub: {
		Url:   "https://github.com/trending",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceDouyin: {
		Url:   "https://www.douyin.com/aweme/v1/web/hot/search/list/?device_platform=webapp&aid=6383&channel=channel_pc_web&detail_list=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceKuaishou: {
		Url:   "https://www.kuaishou.com/?isHome=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceToutiao: {
		Url:   "https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceJuejin: {
		Url:   "https://api.juejin.cn/content_api/v1/content/article_rank?category_id=1&type=hot",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.37 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	Source36Kr: {
		Url:   "https://gateway.36kr.com/api/mis/nav/home/nav/rank/hot",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceCSDN: {
		Url:   "https://blog.csdn.net/phoenix/web/blog/hot-rank",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceTieba: {
		Url:   "https://tieba.baidu.com/hottopic/browse/topicList",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceZhihuDaily: {
		Url:   "https://daily.zhihu.com/api/4/news/latest",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceCoolapk: {
		Url:   "https://api.coolapk.com/v6/page/dataList?url=%2Ffeed%2Fdigest%3Ftype%3D12%26isIncludeTop%3D1&title=%E4%BB%8A%E6%97%A5%E7%83%AD%E9%97%A8&subTitle=&page=1",
		Agent: `Mozilla/5.0 (Linux; Android 10; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36`,
	},
	SourceHupu: {
		Url:   "https://bbs.hupu.com/all-gambia",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceHuxiu: {
		Url:   "https://www.huxiu.com/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceJianshu: {
		Url:   "https://www.jianshu.com/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceSmzdm: {
		Url:   "https://www.smzdm.com/top/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceSspai: {
		Url:   "https://sspai.com/api/v1/articles?offset=0&limit=20&sort=popular",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceNetease: {
		Url:   "https://m.163.com/fe/api/hot/news/flow",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceQQ: {
		Url:   "https://r.inews.qq.com/gw/event/hot_ranking_list?page_size=50",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceAcfun: {
		Url:   "https://www.acfun.cn/rest/pc-direct/rank/channel?channelId=&rankLimit=30&rankPeriod=DAY",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	Source51CTO: {
		Url:   "https://api-media.51cto.com/index/index/recommend",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	Source52Pojie: {
		Url:   "https://www.52pojie.cn/forum.php?mod=guide&view=digest&rss=1",
		Agent: `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36`,
	},
	SourceDoubanGroup: {
		Url:   "https://www.douban.com/group/explore",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36`,
	},
	SourceDgtle: {
		Url:   "https://opser.api.dgtle.com/v2/news/index",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceDoubanMovie: {
		Url:   "https://movie.douban.com/chart/",
		Agent: `Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15`,
	},
	SourceEarthquake: {
		Url:   "https://news.ceic.ac.cn/speedsearch.html",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceGameres: {
		Url:   "https://www.gameres.com",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceGeekpark: {
		Url:   "https://mainssl.geekpark.net/api/v2",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceGenshin: {
		Url:   "https://bbs-api-static.miyoushe.com/painter/wapi/getNewsList?client_type=4&gids=2&page_size=20",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceGuokr: {
		Url:   "https://www.guokr.com/beta/proxy/science_api/articles?limit=30",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0`,
	},
	SourceHackernews: {
		Url:   "https://news.ycombinator.com/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36`,
	},
	SourceHelloGitHub: {
		Url:   "https://abroad.hellogithub.com/v1/?sort_by=featured&tid=&page=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceHistory: {
		Url:   "https://baike.baidu.com/cms/home/eventsOnHistory/01.json",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceHonkai: {
		Url:   "https://bbs-api-static.miyoushe.com/painter/wapi/getNewsList?client_type=4&gids=1&page_size=20",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceHostloc: {
		Url:   "https://hostloc.com/forum.php?mod=guide&view=hot&rss=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceIfanr: {
		Url:   "https://sso.ifanr.com/api/v5/wp/buzz/?limit=20&offset=0",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceMiyoushe: {
		Url:   "https://bbs-api-static.miyoushe.com/painter/wapi/getNewsList?client_type=4&gids=1&page_size=30",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceNewsmth: {
		Url:   "https://wap.newsmth.net/wap/api/hot/global",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceNgabbs: {
		Url:   "https://ngabbs.com/nuke.php?__lib=load_topic&__act=load_topic_reply_ladder2&opt=1&all=1",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceNodeseek: {
		Url:   "https://www.nodeseek.com/rss",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceNytimes: {
		Url:   "https://rsshub.app/nytimes",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceProducthunt: {
		Url:   "https://www.producthunt.com/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceSinaNews: {
		Url:   "https://news.sina.com.cn/zt_d/top_news/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceSina: {
		Url:   "https://s.weibo.com/top/summary",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceStarrail: {
		Url:   "https://bbs-api-static.miyoushe.com/painter/wapi/getNewsList?client_type=4&gids=6&page_size=20",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceThepaper: {
		Url:   "https://www.thepaper.cn/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceWeatheralarm: {
		Url:   "http://www.nmc.cn/rest/findAlarm",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceWeread: {
		Url:   "https://weread.qq.com/web/bookListInCategory/rising",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceYystv: {
		Url:   "https://www.yystv.cn/",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
	SourceIthomeXijiayi: {
		Url:   "https://www.ithome.com/zt/xijiayi",
		Agent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`,
	},
}
