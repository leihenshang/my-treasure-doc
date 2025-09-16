package hottop

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Source string

const (
	SourceITHome     Source = "ithome"
	SourceZhihu      Source = "zhihu"
	SourceWeibo      Source = "weibo"
	SourceBilibili   Source = "bilibili"
	SourceBaidu      Source = "baidu"
	SourceV2EX       Source = "v2ex"
	SourceGitHub     Source = "github"
	SourceDouyin     Source = "douyin"
	SourceKuaishou   Source = "kuaishou"
	SourceToutiao    Source = "toutiao"
	SourceJuejin     Source = "juejin"
	Source36Kr       Source = "36kr"
	SourceCSDN       Source = "csdn"
	SourceTieba      Source = "tieba"
	SourceZhihuDaily Source = "zhihu-daily"
)

type Spider struct {
	UrlMap     map[Source]*UrlConf
	HttpClient *http.Client
}

type UrlConf struct {
	Url   string
	Agent string
}

type HotData struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Link        string     `json:"link"`
	Total       int        `json:"total"`
	Data        []*HotItem `json:"data"`
}

type HotItem struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Timestamp int64  `json:"timestamp"`
	Hot       int    `json:"hot"`
	URL       string `json:"url"`
	MobileURL string `json:"mobileUrl"`
	Author    string `json:"author,omitempty"`
	Desc      string `json:"desc,omitempty"`
}

func NewSpider() *Spider {
	return &Spider{
		UrlMap: map[Source]*UrlConf{
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
		},
		HttpClient: &http.Client{
			Timeout: time.Duration(5 * time.Second),
		},
	}
}

// ============== IT之家 ==============
func (s *Spider) GetItHome() (*HotData, error) {
	replaceItHomeLink := func(url string, getID bool) string {
		re := regexp.MustCompile(`[html|live]/(\d+)\.htm`)
		match := re.FindStringSubmatch(url)
		if len(match) > 1 {
			id := match[1]
			if getID {
				return id
			}
			if len(id) >= 6 {
				return fmt.Sprintf("https://www.ithome.com/0/%s/%s.htm", id[:3], id[3:])
			}
		}
		return url
	}

	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceITHome].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceITHome].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	document.Find(".rank-box .placeholder").Each(func(i int, selection *goquery.Selection) {
		href, exists := selection.Find("a").Attr("href")
		if !exists {
			return
		}

		title := strings.TrimSpace(selection.Find(".plc-title").Text())
		cover, _ := selection.Find("img").Attr("data-original")
		timeStr := strings.TrimSpace(selection.Find("span.post-time").Text())
		timestamp, _ := time.Parse("2006-01-02 15:04:05", timeStr)

		hotStr := strings.TrimSpace(selection.Find(".review-num").Text())
		hotStr = regexp.MustCompile(`\D`).ReplaceAllString(hotStr, "")
		hot, _ := strconv.Atoi(hotStr)

		hotItem := &HotItem{
			ID:        100000,
			Title:     title,
			Cover:     cover,
			Timestamp: timestamp.Unix(),
			Hot:       hot,
			URL:       replaceItHomeLink(href, false),
			MobileURL: replaceItHomeLink(href, false),
		}

		if idStr := replaceItHomeLink(href, true); idStr != href {
			if id, err := strconv.Atoi(idStr); err == nil {
				hotItem.ID = id
			}
		}

		listData = append(listData, hotItem)
	})

	return &HotData{
		Name:        "ithome",
		Title:       "IT之家",
		Type:        "热榜",
		Description: "爱科技，爱这里 - 前沿科技新闻网站",
		Link:        "https://m.ithome.com/rankm/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 知乎 ==============
func (s *Spider) GetZhihu() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceZhihu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceZhihu].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].([]interface{}); ok {
		for _, item := range data {
			if v, ok := item.(map[string]interface{}); ok {
				target := v["target"].(map[string]interface{})
				questionId := ""
				if url, ok := target["url"].(string); ok {
					parts := strings.Split(url, "/")
					if len(parts) > 0 {
						questionId = parts[len(parts)-1]
					}
				}

				title := ""
				if t, ok := target["title"].(string); ok {
					title = t
				}

				cover := ""
				if children, ok := v["children"].([]interface{}); ok && len(children) > 0 {
					if child, ok := children[0].(map[string]interface{}); ok {
						if thumbnail, ok := child["thumbnail"].(string); ok {
							cover = thumbnail
						}
					}
				}

				hot := 0
				if detailText, ok := v["detail_text"].(string); ok {
					if matches := regexp.MustCompile(`(\\d+(\\.\\d+)?)`).FindStringSubmatch(detailText); len(matches) > 0 {
						if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
							hot = int(val * 10000)
						}
					}
				}

				timestamp := time.Now().Unix()
				if created, ok := target["created"].(float64); ok {
					timestamp = int64(created)
				}

				listData = append(listData, &HotItem{
					ID:        int(target["id"].(float64)),
					Title:     title,
					Cover:     cover,
					Timestamp: timestamp,
					Hot:       hot,
					URL:       fmt.Sprintf("https://www.zhihu.com/question/%s", questionId),
					MobileURL: fmt.Sprintf("https://www.zhihu.com/question/%s", questionId),
				})
			}
		}
	}

	return &HotData{
		Name:  "zhihu",
		Title: "知乎",
		Type:  "热榜",
		Link:  "https://www.zhihu.com/hot",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 微博 ==============
func (s *Spider) GetWeibo() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceWeibo].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceWeibo].Agent)
	request.Header.Add("Referer", "https://s.weibo.com/top/summary?cate=realtimehot")
	request.Header.Add("MWeibo-Pwa", "1")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].(map[string]interface{}); ok {
		if cards, ok := data["cards"].([]interface{}); ok && len(cards) > 0 {
			if cardGroup, ok := cards[0].(map[string]interface{})["card_group"].([]interface{}); ok {
				for _, item := range cardGroup {
					if v, ok := item.(map[string]interface{}); ok {
						title := ""
						if desc, ok := v["desc"].(string); ok {
							title = desc
						}

						wordScheme := ""
						if ws, ok := v["word_scheme"].(string); ok {
							wordScheme = ws
						} else {
							wordScheme = fmt.Sprintf("#%s#", title)
						}

						hot := 0
						if num, ok := v["num"].(float64); ok {
							hot = int(num)
						}

						timestamp := time.Now().Unix()
						if onboardTime, ok := v["onboard_time"].(float64); ok {
							timestamp = int64(onboardTime)
						}

						listData = append(listData, &HotItem{
							ID:        int(v["itemid"].(float64)),
							Title:     title,
							Desc:      wordScheme,
							Timestamp: timestamp,
							Hot:       hot,
							URL:       fmt.Sprintf("https://s.weibo.com/weibo?q=%s&t=31&band_rank=1&Refer=top", strings.ReplaceAll(wordScheme, "#", "")),
							MobileURL: v["scheme"].(string),
						})
					}
				}
			}
		}
	}

	return &HotData{
		Name:        "weibo",
		Title:       "微博",
		Type:        "热搜榜",
		Description: "实时热点，每分钟更新一次",
		Link:        "https://s.weibo.com/top/summary/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 哔哩哔哩 ==============
func (s *Spider) GetBilibili() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceBilibili].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceBilibili].Agent)
	request.Header.Add("Referer", "https://www.bilibili.com/ranking/all")

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].(map[string]interface{}); ok {
		if list, ok := data["list"].([]interface{}); ok {
			for _, item := range list {
				if v, ok := item.(map[string]interface{}); ok {
					title := ""
					if t, ok := v["title"].(string); ok {
						title = t
					}

					desc := "该视频暂无简介"
					if d, ok := v["desc"].(string); ok && d != "" {
						desc = d
					}

					cover := ""
					if pic, ok := v["pic"].(string); ok {
						cover = strings.ReplaceAll(pic, "http:", "https:")
					}

					author := ""
					if owner, ok := v["owner"].(map[string]interface{}); ok {
						if name, ok := owner["name"].(string); ok {
							author = name
						}
					}

					timestamp := int64(0)
					if pubdate, ok := v["pubdate"].(float64); ok {
						timestamp = int64(pubdate)
					}

					hot := 0
					if stat, ok := v["stat"].(map[string]interface{}); ok {
						if view, ok := stat["view"].(float64); ok {
							hot = int(view)
						}
					}

					bvid := ""
					if b, ok := v["bvid"].(string); ok {
						bvid = b
					}

					listData = append(listData, &HotItem{
						ID:        0,
						Title:     title,
						Desc:      desc,
						Cover:     cover,
						Author:    author,
						Timestamp: timestamp,
						Hot:       hot,
						URL:       fmt.Sprintf("https://www.bilibili.com/video/%s", bvid),
						MobileURL: fmt.Sprintf("https://m.bilibili.com/video/%s", bvid),
					})
				}
			}
		}
	}

	return &HotData{
		Name:        "bilibili",
		Title:       "哔哩哔哩",
		Type:        "热榜 · 全站",
		Description: "你所热爱的，就是你的生活",
		Link:        "https://www.bilibili.com/v/popular/rank/all",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 百度 ==============
func (s *Spider) GetBaidu() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceBaidu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceBaidu].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	pattern := `<!--s-data:(.*?)-->`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return nil, fmt.Errorf("未找到数据")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(matches[1]), &data); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if cards, ok := data["cards"].([]interface{}); ok && len(cards) > 0 {
		if card, ok := cards[0].(map[string]interface{}); ok {
			if content, ok := card["content"].([]interface{}); ok {
				for _, item := range content {
					if v, ok := item.(map[string]interface{}); ok {
						title := ""
						if word, ok := v["word"].(string); ok {
							title = word
						}

						desc := ""
						if d, ok := v["desc"].(string); ok {
							desc = d
						}

						cover := ""
						if img, ok := v["img"].(string); ok {
							cover = img
						}

						author := ""
						if show, ok := v["show"].([]interface{}); ok && len(show) > 0 {
							if s, ok := show[0].(string); ok {
								author = s
							}
						}

						hot := 0
						if hotScore, ok := v["hotScore"].(float64); ok {
							hot = int(hotScore)
						}

						query := ""
						if q, ok := v["query"].(string); ok {
							query = q
						}

						rawUrl := ""
						if ru, ok := v["rawUrl"].(string); ok {
							rawUrl = ru
						}

						listData = append(listData, &HotItem{
							ID:        int(v["index"].(float64)),
							Title:     title,
							Desc:      desc,
							Cover:     cover,
							Author:    author,
							Timestamp: 0,
							Hot:       hot,
							URL:       fmt.Sprintf("https://www.baidu.com/s?wd=%s", strings.ReplaceAll(query, " ", "+")),
							MobileURL: rawUrl,
						})
					}
				}
			}
		}
	}

	return &HotData{
		Name:  "baidu",
		Title: "百度",
		Type:  "热搜",
		Link:  "https://top.baidu.com/board",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== V2EX ==============
func (s *Spider) GetV2EX() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceV2EX].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceV2EX].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var list []map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, v := range list {
		title := ""
		if t, ok := v["title"].(string); ok {
			title = t
		}

		content := ""
		if c, ok := v["content"].(string); ok {
			content = c
		}

		url := ""
		if u, ok := v["url"].(string); ok {
			url = u
		}

		replies := 0
		if r, ok := v["replies"].(float64); ok {
			replies = int(r)
		}

		author := ""
		if member, ok := v["member"].(map[string]interface{}); ok {
			if username, ok := member["username"].(string); ok {
				author = username
			}
		}

		listData = append(listData, &HotItem{
			ID:        int(v["id"].(float64)),
			Title:     title,
			Desc:      content,
			Author:    author,
			Timestamp: 0,
			Hot:       replies,
			URL:       url,
			MobileURL: url,
		})
	}

	return &HotData{
		Name:  "v2ex",
		Title: "V2EX",
		Type:  "主题榜",
		Link:  "https://www.v2ex.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== GitHub ==============
func (s *Spider) GetGitHub() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceGitHub].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceGitHub].Agent)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	document.Find("article.Box-row").Each(func(i int, selection *goquery.Selection) {
		// 仓库名称
		repoText := strings.TrimSpace(selection.Find("h2 a").Text())
		repoText = strings.ReplaceAll(repoText, "\n", "")
		repoText = regexp.MustCompile(`\s+`).ReplaceAllString(repoText, " ")
		repoParts := strings.Split(repoText, "/")

		owner := ""
		repoName := ""
		if len(repoParts) >= 2 {
			owner = strings.TrimSpace(repoParts[0])
			repoName = strings.TrimSpace(repoParts[1])
		}

		// 仓库链接
		href, exists := selection.Find("h2 a").Attr("href")
		repoURL := ""
		if exists {
			repoURL = "https://github.com" + href
		}

		// 描述
		desc := strings.TrimSpace(selection.Find("p.col-9.color-fg-muted").Text())

		// 语言
		language := strings.TrimSpace(selection.Find(`[itemprop="programmingLanguage"]`).Text())

		// Stars
		starsText := strings.TrimSpace(selection.Find(`a[href$="/stargazers"]`).Text())
		stars := 0
		if starsText != "" {
			starsStr := regexp.MustCompile(`\D`).ReplaceAllString(starsText, "")
			stars, _ = strconv.Atoi(starsStr)
		}

		listData = append(listData, &HotItem{
			ID:        i + 1,
			Title:     fmt.Sprintf("%s/%s", owner, repoName),
			Desc:      desc,
			Author:    language,
			Timestamp: time.Now().Unix(),
			Hot:       stars,
			URL:       repoURL,
			MobileURL: repoURL,
		})
	})

	return &HotData{
		Name:        "github",
		Title:       "GitHub",
		Type:        "趋势榜",
		Description: "GitHub trending repositories",
		Link:        "https://github.com/trending",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 抖音 ==============
func (s *Spider) GetDouyin() (*HotData, error) {
	// 获取抖音Cookie
	cookieURL := "https://www.douyin.com/passport/general/login_guiding_strategy/?aid=6383"
	cookieReq, err := http.NewRequest("GET", cookieURL, nil)
	if err != nil {
		return nil, err
	}
	cookieReq.Header.Add("User-Agent", s.UrlMap[SourceDouyin].Agent)

	cookieRes, err := s.HttpClient.Do(cookieReq)
	if err != nil {
		return nil, err
	}
	defer cookieRes.Body.Close()

	cookie := ""
	if setCookie := cookieRes.Header.Get("Set-Cookie"); setCookie != "" {
		re := regexp.MustCompile(`passport_csrf_token=([^;]+)`)
		matches := re.FindStringSubmatch(setCookie)
		if len(matches) > 1 {
			cookie = matches[1]
		}
	}

	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceDouyin].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceDouyin].Agent)
	if cookie != "" {
		request.Header.Add("Cookie", fmt.Sprintf("passport_csrf_token=%s", cookie))
	}

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].(map[string]interface{}); ok {
		if wordList, ok := data["word_list"].([]interface{}); ok {
			for i, item := range wordList {
				if v, ok := item.(map[string]interface{}); ok {
					title := ""
					if word, ok := v["word"].(string); ok {
						title = word
					}

					sentenceID := ""
					if sid, ok := v["sentence_id"].(string); ok {
						sentenceID = sid
					}

					hotValue := 0
					if hv, ok := v["hot_value"].(float64); ok {
						hotValue = int(hv)
					}

					eventTime := int64(0)
					if et, ok := v["event_time"].(float64); ok {
						eventTime = int64(et)
					}

					listData = append(listData, &HotItem{
						ID:        i + 1,
						Title:     title,
						Timestamp: eventTime,
						Hot:       hotValue,
						URL:       fmt.Sprintf("https://www.douyin.com/hot/%s", sentenceID),
						MobileURL: fmt.Sprintf("https://www.douyin.com/hot/%s", sentenceID),
					})
				}
			}
		}
	}

	return &HotData{
		Name:        "douyin",
		Title:       "抖音",
		Type:        "热榜",
		Description: "实时上升热点",
		Link:        "https://www.douyin.com",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 快手 ==============
func (s *Spider) GetKuaishou() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceKuaishou].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceKuaishou].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// 提取window.__APOLLO_STATE__数据
	pattern := `window\.__APOLLO_STATE__=(.*?);\(function\(\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return nil, fmt.Errorf("未找到快手数据")
	}

	var apolloData map[string]interface{}
	if err := json.Unmarshal([]byte(matches[1]), &apolloData); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if defaultClient, ok := apolloData["defaultClient"].(map[string]interface{}); ok {
		if items, ok := defaultClient["$ROOT_QUERY.visionHotRank({\"page\":\"home\"})"].([]interface{}); ok {
			for i, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if id, ok := itemMap["id"].(string); ok {
						if hotItem, ok := defaultClient[id].(map[string]interface{}); ok {
							title := ""
							if name, ok := hotItem["name"].(string); ok {
								title = name
							}

							poster := ""
							if p, ok := hotItem["poster"].(string); ok {
								poster = p
							}

							hotValue := 0
							if hv, ok := hotItem["hotValue"].(string); ok {
								// 解析中文数字
								hotValue = parseChineseNumber(hv)
							}

							photoIds := ""
							if ids, ok := hotItem["photoIds"].(map[string]interface{}); ok {
								if jsonArr, ok := ids["json"].([]interface{}); ok && len(jsonArr) > 0 {
									if pid, ok := jsonArr[0].(string); ok {
										photoIds = pid
									}
								}
							}

							listData = append(listData, &HotItem{
								ID:        i + 1,
								Title:     title,
								Cover:     poster,
								Timestamp: time.Now().Unix(),
								Hot:       hotValue,
								URL:       fmt.Sprintf("https://www.kuaishou.com/short-video/%s", photoIds),
								MobileURL: fmt.Sprintf("https://www.kuaishou.com/short-video/%s", photoIds),
							})
						}
					}
				}
			}
		}
	}

	return &HotData{
		Name:        "kuaishou",
		Title:       "快手",
		Type:        "热榜",
		Description: "快手，拥抱每一种生活",
		Link:        "https://www.kuaishou.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 今日头条 ==============
func (s *Spider) GetToutiao() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceToutiao].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceToutiao].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].([]interface{}); ok {
		for i, item := range data {
			if v, ok := item.(map[string]interface{}); ok {
				title := ""
				if t, ok := v["Title"].(string); ok {
					title = t
				}

				clusterID := ""
				if cid, ok := v["ClusterIdStr"].(string); ok {
					clusterID = cid
				}

				hotValue := 0
				if hv, ok := v["HotValue"].(float64); ok {
					hotValue = int(hv)
				}

				imageURL := ""
				if image, ok := v["Image"].(map[string]interface{}); ok {
					if url, ok := image["url"].(string); ok {
						imageURL = url
					}
				}

				listData = append(listData, &HotItem{
					ID:        i + 1,
					Title:     title,
					Cover:     imageURL,
					Timestamp: time.Now().Unix(),
					Hot:       hotValue,
					URL:       fmt.Sprintf("https://www.toutiao.com/trending/%s/", clusterID),
					MobileURL: fmt.Sprintf("https://api.toutiaoapi.com/feoffline/amos_land/new/html/main/index.html?topic_id=%s", clusterID),
				})
			}
		}
	}

	return &HotData{
		Name:  "toutiao",
		Title: "今日头条",
		Type:  "热榜",
		Link:  "https://www.toutiao.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 掘金 ==============
func (s *Spider) GetJuejin() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceJuejin].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceJuejin].Agent)
	request.Header.Add("Accept", "application/json, text/plain, */*")

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].([]interface{}); ok {
		for i, item := range data {
			if v, ok := item.(map[string]interface{}); ok {
				content := v["content"].(map[string]interface{})
				author := v["author"].(map[string]interface{})
				contentCounter := v["content_counter"].(map[string]interface{})

				title := ""
				if t, ok := content["title"].(string); ok {
					title = t
				}

				contentID := ""
				if cid, ok := content["content_id"].(string); ok {
					contentID = cid
				}

				authorName := ""
				if an, ok := author["name"].(string); ok {
					authorName = an
				}

				hotRank := 0
				if hr, ok := contentCounter["hot_rank"].(float64); ok {
					hotRank = int(hr)
				}

				listData = append(listData, &HotItem{
					ID:        i + 1,
					Title:     title,
					Author:    authorName,
					Timestamp: time.Now().Unix(),
					Hot:       hotRank,
					URL:       fmt.Sprintf("https://juejin.cn/post/%s", contentID),
					MobileURL: fmt.Sprintf("https://juejin.cn/post/%s", contentID),
				})
			}
		}
	}

	return &HotData{
		Name:  "juejin",
		Title: "稀土掘金",
		Type:  "文章榜",
		Link:  "https://juejin.cn/hot/articles",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// 辅助函数：解析中文数字
func parseChineseNumber(str string) int {
	// 简化的中文数字解析，实际需要更复杂的实现
	if str == "" {
		return 0
	}

	// 移除空格和常见字符
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, "万", "")
	str = strings.ReplaceAll(str, "千", "")
	str = strings.ReplaceAll(str, "百", "")

	// 尝试解析为普通数字
	if num, err := strconv.Atoi(str); err == nil {
		return num
	}

	return 0
}

// ... existing code ...
// ============== 36氪 ==============
func (s *Spider) Get36Kr() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[Source36Kr].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[Source36Kr].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].([]interface{}); ok {
		for i, item := range data {
			if v, ok := item.(map[string]interface{}); ok {
				title := ""
				if t, ok := v["title"].(string); ok {
					title = t
				}

				cover := ""
				if c, ok := v["cover"].(string); ok {
					cover = c
				}

				author := ""
				if a, ok := v["author"].(string); ok {
					author = a
				}

				hot := 0
				if h, ok := v["hot"].(float64); ok {
					hot = int(h)
				}

				timestamp := time.Now().Unix()
				if ts, ok := v["publishTime"].(float64); ok {
					timestamp = int64(ts)
				}

				listData = append(listData, &HotItem{
					ID:        i + 1,
					Title:     title,
					Cover:     cover,
					Author:    author,
					Timestamp: timestamp,
					Hot:       hot,
					URL:       fmt.Sprintf("https://36kr.com/p/%v", v["id"]),
					MobileURL: fmt.Sprintf("https://m.36kr.com/p/%v", v["id"]),
				})
			}
		}
	}

	return &HotData{
		Name:        "36kr",
		Title:       "36氪",
		Type:        "热榜",
		Description: "让一部分人先看到未来",
		Link:        "https://36kr.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== CSDN ==============
func (s *Spider) GetCSDN() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceCSDN].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceCSDN].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].([]interface{}); ok {
		for i, item := range data {
			if v, ok := item.(map[string]interface{}); ok {
				title := ""
				if t, ok := v["articleTitle"].(string); ok {
					title = t
				}

				cover := ""
				if c, ok := v["pic"].(string); ok {
					cover = c
				}

				author := ""
				if a, ok := v["nickName"].(string); ok {
					author = a
				}

				hot := 0
				if h, ok := v["hot"].(float64); ok {
					hot = int(h)
				}

				articleDetailURL := ""
				if url, ok := v["articleDetailUrl"].(string); ok {
					articleDetailURL = url
				}

				listData = append(listData, &HotItem{
					ID:        i + 1,
					Title:     title,
					Cover:     cover,
					Author:    author,
					Timestamp: time.Now().Unix(),
					Hot:       hot,
					URL:       articleDetailURL,
					MobileURL: articleDetailURL,
				})
			}
		}
	}

	return &HotData{
		Name:        "csdn",
		Title:       "CSDN",
		Type:        "热榜",
		Description: "专业开发者社区",
		Link:        "https://blog.csdn.net/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 百度贴吧 ==============
func (s *Spider) GetTieba() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceTieba].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceTieba].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if data, ok := result["data"].(map[string]interface{}); ok {
		if bangTopic, ok := data["bang_topic"].(map[string]interface{}); ok {
			if topicList, ok := bangTopic["topic_list"].([]interface{}); ok {
				for i, item := range topicList {
					if v, ok := item.(map[string]interface{}); ok {
						title := ""
						if t, ok := v["topic_name"].(string); ok {
							title = t
						}

						desc := ""
						if d, ok := v["topic_desc"].(string); ok {
							desc = d
						}

						cover := ""
						if c, ok := v["topic_pic"].(string); ok {
							cover = c
						}

						hot := 0
						if h, ok := v["topic_num"].(float64); ok {
							hot = int(h)
						}

						topicID := ""
						if id, ok := v["topic_id"].(string); ok {
							topicID = id
						}

						listData = append(listData, &HotItem{
							ID:        i + 1,
							Title:     title,
							Desc:      desc,
							Cover:     cover,
							Timestamp: time.Now().Unix(),
							Hot:       hot,
							URL:       fmt.Sprintf("https://tieba.baidu.com/hottopic/browse/topicList?topic_id=%s", topicID),
							MobileURL: fmt.Sprintf("https://tieba.baidu.com/hottopic/browse/topicList?topic_id=%s", topicID),
						})
					}
				}
			}
		}
	}

	return &HotData{
		Name:        "tieba",
		Title:       "百度贴吧",
		Type:        "热议榜",
		Description: "全球最大的中文社区",
		Link:        "https://tieba.baidu.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 知乎日报 ==============
func (s *Spider) GetZhihuDaily() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceZhihuDaily].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceZhihuDaily].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if stories, ok := result["stories"].([]interface{}); ok {
		for i, item := range stories {
			if v, ok := item.(map[string]interface{}); ok {
				// 只显示type为0的故事
				if storyType, ok := v["type"].(float64); ok && storyType != 0 {
					continue
				}

				title := ""
				if t, ok := v["title"].(string); ok {
					title = t
				}

				cover := ""
				if images, ok := v["images"].([]interface{}); ok && len(images) > 0 {
					if img, ok := images[0].(string); ok {
						cover = img
					}
				}

				storyID := 0
				if id, ok := v["id"].(float64); ok {
					storyID = int(id)
				}

				listData = append(listData, &HotItem{
					ID:        i + 1,
					Title:     title,
					Cover:     cover,
					Timestamp: time.Now().Unix(),
					Hot:       0,
					URL:       fmt.Sprintf("https://daily.zhihu.com/story/%d", storyID),
					MobileURL: fmt.Sprintf("https://daily.zhihu.com/story/%d", storyID),
				})
			}
		}
	}

	return &HotData{
		Name:        "zhihu-daily",
		Title:       "知乎日报",
		Type:        "日报",
		Description: "每天3次，每次7分钟",
		Link:        "https://daily.zhihu.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}
