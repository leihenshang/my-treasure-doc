package hot

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type UrlConf struct {
	Url   string
	Agent string
}

type Spider struct {
	UrlMap     map[Source]*UrlConf
	HttpClient *http.Client
}

var spider *Spider
var spiderOnce *sync.Once = &sync.Once{}

func NewSpider() *Spider {
	spiderOnce.Do(func() {
		spider = &Spider{
			UrlMap: UrlConfMap,
			HttpClient: &http.Client{
				Timeout: time.Duration(5 * time.Second),
			},
		}
	})
	return spider
}

func GetSpider() *Spider {
	return spider
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
			ID:        "100000",
			Title:     title,
			Cover:     cover,
			Timestamp: timestamp.Unix(),
			Hot:       hot,
			URL:       replaceItHomeLink(href, false),
			MobileURL: replaceItHomeLink(href, false),
		}

		if idStr := replaceItHomeLink(href, true); idStr != href {
			if id, err := strconv.Atoi(idStr); err == nil {
				hotItem.ID = strconv.Itoa(id) // 使用 strconv.Itoa 将 int 转换为 string
			}
		}

		listData = append(listData, hotItem)
	})

	return &HotData{
		Code:        http.StatusOK,
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
					// 匹配纯数字、小数，以及带"万"字的数字
					if matches := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(?:万|w)?`).FindStringSubmatch(detailText); len(matches) > 0 {
						if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
							// 如果包含"万"字或"w"，乘以10000
							if strings.Contains(detailText, "万") || strings.Contains(detailText, "w") {
								hot = int(val * 10000)
							} else {
								hot = int(val)
							}
						}
					}
				}

				timestamp := time.Now().Unix()
				if created, ok := target["created"].(float64); ok {
					timestamp = int64(created)
				}

				listData = append(listData, &HotItem{
					ID:        strconv.Itoa(int(target["id"].(float64))),
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
		Code:  http.StatusOK,
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

						itemID := 0
						if id, ok := v["itemid"].(float64); ok {
							itemID = int(id)
						} else if idStr, ok := v["itemid"].(string); ok {
							if id, err := strconv.Atoi(idStr); err == nil {
								itemID = id
							}
						}

						listData = append(listData, &HotItem{
							ID:        strconv.Itoa(itemID),
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
		Code:        http.StatusOK,
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
						ID:        "0",
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
		Code:        http.StatusOK,
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
							ID:        strconv.Itoa(int(v["index"].(float64))),
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
		Code:  http.StatusOK,
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
			ID:        strconv.Itoa(int(v["id"].(float64))),
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
		Code:  http.StatusOK,
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
			ID:        strconv.Itoa(i + 1),
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
		Code:        http.StatusOK,
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
						ID:        strconv.Itoa(i + 1),
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
		Code:        http.StatusOK,
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
								ID:        strconv.Itoa(i + 1),
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
		Code:        http.StatusOK,
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
					ID:        strconv.Itoa(i + 1),
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
		Code:  http.StatusOK,
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
					ID:        strconv.Itoa(i + 1),
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
		Code:  http.StatusOK,
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
					ID:        strconv.Itoa(i + 1),
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
		Code:        http.StatusOK,
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
					ID:        strconv.Itoa(i + 1),
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
		Code:        http.StatusOK,
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
							ID:        strconv.Itoa(i + 1),
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
		Code:        http.StatusOK,
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
					ID:        strconv.Itoa(i + 1),
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
		Code:        http.StatusOK,
		Name:        "zhihu-daily",
		Title:       "知乎日报",
		Type:        "日报",
		Description: "每天3次，每次7分钟",
		Link:        "https://daily.zhihu.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 酷安 ==============
func (s *Spider) GetCoolapk() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceCoolapk].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceCoolapk].Agent)
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
	if data, ok := result["data"].([]interface{}); ok {
		for _, item := range data {
			if v, ok := item.(map[string]interface{}); ok {
				title := ""
				if t, ok := v["title"].(string); ok {
					title = t
				}

				cover := ""
				if c, ok := v["pic"].(string); ok {
					cover = c
				}

				author := ""
				if a, ok := v["username"].(string); ok {
					author = a
				}

				hot := 0
				if h, ok := v["like_num"].(float64); ok {
					hot = int(h)
				}

				id := 0
				if idFloat, ok := v["id"].(float64); ok {
					id = int(idFloat)
				}

				listData = append(listData, &HotItem{
					ID:        strconv.Itoa(id),
					Title:     title,
					Cover:     cover,
					Author:    author,
					Timestamp: time.Now().Unix(),
					Hot:       hot,
					URL:       fmt.Sprintf("https://www.coolapk.com/feed/%d", id),
					MobileURL: fmt.Sprintf("https://www.coolapk.com/feed/%d", id),
				})
			}
		}
	}

	return &HotData{
		Code:        http.StatusOK,
		Name:        "coolapk",
		Title:       "酷安",
		Type:        "热门",
		Description: "发现科技新生活",
		Link:        "https://www.coolapk.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 虎扑 ==============
func (s *Spider) GetHupu() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceHupu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceHupu].Agent)

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
	document.Find(".bbsHotPit .list-item").Each(func(i int, selection *goquery.Selection) {
		title := strings.TrimSpace(selection.Find(".textSpan").Text())
		href, exists := selection.Find("a").Attr("href")
		if !exists {
			return
		}

		author := strings.TrimSpace(selection.Find(".author").Text())
		replyStr := strings.TrimSpace(selection.Find(".reply").Text())
		replyNum := 0
		if match := regexp.MustCompile(`(\d+)`).FindStringSubmatch(replyStr); len(match) > 1 {
			replyNum, _ = strconv.Atoi(match[1])
		}

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(i + 1),
			Title:     title,
			Author:    author,
			Timestamp: time.Now().Unix(),
			Hot:       replyNum,
			URL:       fmt.Sprintf("https://bbs.hupu.com%s", href),
			MobileURL: fmt.Sprintf("https://bbs.hupu.com%s", href),
		})
	})

	return &HotData{
		Code:        http.StatusOK,
		Name:        "hupu",
		Title:       "虎扑",
		Type:        "热帖",
		Description: "虎扑篮球社区",
		Link:        "https://bbs.hupu.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 虎嗅 ==============
func (s *Spider) GetHuxiu() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceHuxiu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceHuxiu].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(bodyBytes)

	// 提取window.__INITIAL_STATE__数据
	re := regexp.MustCompile(`window\.__INITIAL_STATE__\s*=\s*({.*?});`)
	match := re.FindStringSubmatch(bodyStr)
	if len(match) < 2 {
		return nil, fmt.Errorf("无法找到INITIAL_STATE数据")
	}

	var initialState map[string]interface{}
	if err := json.Unmarshal([]byte(match[1]), &initialState); err != nil {
		return nil, err
	}

	var listData []*HotItem
	if homeData, ok := initialState["home"].(map[string]interface{}); ok {
		if hotNewsList, ok := homeData["hotNewsList"].([]interface{}); ok {
			for _, item := range hotNewsList {
				if v, ok := item.(map[string]interface{}); ok {
					title := ""
					if t, ok := v["title"].(string); ok {
						title = t
					}

					summary := ""
					if s, ok := v["summary"].(string); ok {
						summary = s
					}

					newsID := 0
					if id, ok := v["newsId"].(float64); ok {
						newsID = int(id)
					}

					listData = append(listData, &HotItem{
						ID:        strconv.Itoa(newsID),
						Title:     title,
						Desc:      summary,
						Timestamp: time.Now().Unix(),
						Hot:       0,
						URL:       fmt.Sprintf("https://www.huxiu.com/article/%d.html", newsID),
						MobileURL: fmt.Sprintf("https://m.huxiu.com/article/%d.html", newsID),
					})
				}
			}
		}
	}

	return &HotData{
		Code:        http.StatusOK,
		Name:        "huxiu",
		Title:       "虎嗅",
		Type:        "24小时",
		Description: "虎嗅网 - 商业科技新媒体",
		Link:        "https://www.huxiu.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 简书 ==============
func (s *Spider) GetJianshu() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceJianshu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceJianshu].Agent)

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
	document.Find(".note-list li").Each(func(i int, selection *goquery.Selection) {
		title := strings.TrimSpace(selection.Find(".title").Text())
		if title == "" {
			return
		}

		href, exists := selection.Find(".title").Attr("href")
		if !exists {
			return
		}

		author := strings.TrimSpace(selection.Find(".nickname").Text())
		abstract := strings.TrimSpace(selection.Find(".abstract").Text())

		// 提取文章ID
		id := 0
		if match := regexp.MustCompile(`/p/([a-f0-9]+)`).FindStringSubmatch(href); len(match) > 1 {
			// 简书的ID是字符串格式，这里用索引代替
			id = i + 1
		}

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(id),
			Title:     title,
			Desc:      abstract,
			Author:    author,
			Timestamp: time.Now().Unix(),
			Hot:       0,
			URL:       fmt.Sprintf("https://www.jianshu.com%s", href),
			MobileURL: fmt.Sprintf("https://www.jianshu.com%s", href),
		})
	})

	return &HotData{
		Code:        http.StatusOK,
		Name:        "jianshu",
		Title:       "简书",
		Type:        "热门",
		Description: "简书 - 创作你的创作",
		Link:        "https://www.jianshu.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 什么值得买 ==============
func (s *Spider) GetSmzdm() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceSmzdm].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceSmzdm].Agent)

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
	document.Find(".z-feed-content").Each(func(i int, selection *goquery.Selection) {
		title := strings.TrimSpace(selection.Find(".z-feed-title").Text())
		if title == "" {
			return
		}

		href, exists := selection.Find(".z-feed-title a").Attr("href")
		if !exists {
			return
		}

		// 提取商品ID
		id := 0
		if match := regexp.MustCompile(`/p/(\d+)`).FindStringSubmatch(href); len(match) > 1 {
			id, _ = strconv.Atoi(match[1])
		}

		price := strings.TrimSpace(selection.Find(".z-highlight").Text())
		worth := 0
		if worthStr := strings.TrimSpace(selection.Find(".z-icon-worth").Parent().Text()); worthStr != "" {
			if match := regexp.MustCompile(`(\d+)`).FindStringSubmatch(worthStr); len(match) > 1 {
				worth, _ = strconv.Atoi(match[1])
			}
		}

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(id),
			Title:     fmt.Sprintf("%s %s", title, price),
			Desc:      price,
			Timestamp: time.Now().Unix(),
			Hot:       worth,
			URL:       fmt.Sprintf("https://www.smzdm.com%s", href),
			MobileURL: fmt.Sprintf("https://m.smzdm.com%s", href),
		})
	})

	return &HotData{
		Code:        http.StatusOK,
		Name:        "smzdm",
		Title:       "什么值得买",
		Type:        "热门",
		Description: "什么值得买 - 消费决策平台",
		Link:        "https://www.smzdm.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 少数派 ==============
func (s *Spider) GetSspai() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceSspai].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceSspai].Agent)

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
	var dataList []interface{}
	if data, ok := result["data"].([]interface{}); ok {
		dataList = data
	}

	for i, item := range dataList {
		if itemMap, ok := item.(map[string]interface{}); ok {
			title := ""
			if t, ok := itemMap["title"].(string); ok {
				title = t
			}

			summary := ""
			if s, ok := itemMap["summary"].(string); ok {
				summary = s
			}

			id := 0
			if idFloat, ok := itemMap["id"].(float64); ok {
				id = int(idFloat)
			}

			author := ""
			if authorMap, ok := itemMap["author"].(map[string]interface{}); ok {
				if nickname, ok := authorMap["nickname"].(string); ok {
					author = nickname
				}
			}

			listData = append(listData, &HotItem{
				ID:        strconv.Itoa(id),
				Title:     title,
				Desc:      summary,
				Author:    author,
				Timestamp: time.Now().Unix(),
				Hot:       0,
				URL:       fmt.Sprintf("https://sspai.com/post/%d", id),
				MobileURL: fmt.Sprintf("https://sspai.com/post/%d", id),
			})
		}

		if i >= 19 { // 限制为20条
			break
		}
	}

	return &HotData{
		Code:        http.StatusOK,
		Name:        "sspai",
		Title:       "少数派",
		Type:        "热门文章",
		Description: "少数派 - 高效工作，品质生活",
		Link:        "https://sspai.com/",
		Total:       len(listData),
		Data:        listData,
	}, nil
}

// ============== 网易新闻 ==============
func (s *Spider) GetNetease() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceNetease].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceNetease].Agent)

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

					cover := ""
					if img, ok := v["imgsrc"].(string); ok {
						cover = img
					}

					author := ""
					if source, ok := v["source"].(string); ok {
						author = source
					}

					timestamp := int64(0)
					if ptime, ok := v["ptime"].(string); ok {
						if t, err := time.Parse("2006-01-02 15:04:05", ptime); err == nil {
							timestamp = t.Unix()
						}
					}

					docid := ""
					if id, ok := v["docid"].(string); ok {
						docid = id
					}

					listData = append(listData, &HotItem{
						ID:        "0",
						Title:     title,
						Cover:     cover,
						Author:    author,
						Timestamp: timestamp,
						Hot:       0,
						URL:       fmt.Sprintf("https://www.163.com/dy/article/%s.html", docid),
						MobileURL: fmt.Sprintf("https://m.163.com/dy/article/%s.html", docid),
					})
				}
			}
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "netease",
		Title: "网易新闻",
		Type:  "热点榜",
		Link:  "https://m.163.com/hot",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 腾讯新闻 ==============
func (s *Spider) GetQQ() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceQQ].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceQQ].Agent)

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
		if idlist, ok := data["idlist"].([]interface{}); ok && len(idlist) > 0 {
			if newslist, ok := idlist[0].(map[string]interface{})["newslist"].([]interface{}); ok {
				// Skip the first item as it seems to be a header
				for i := 1; i < len(newslist); i++ {
					if v, ok := newslist[i].(map[string]interface{}); ok {
						title := ""
						if t, ok := v["title"].(string); ok {
							title = t
						}

						desc := ""
						if abstract, ok := v["abstract"].(string); ok {
							desc = abstract
						}

						cover := ""
						if img, ok := v["miniProShareImage"].(string); ok {
							cover = img
						}

						author := ""
						if source, ok := v["source"].(string); ok {
							author = source
						}

						timestamp := int64(0)
						if ts, ok := v["timestamp"].(float64); ok {
							timestamp = int64(ts)
						}

						id := ""
						if itemID, ok := v["id"].(string); ok {
							id = itemID
						}

						hotScore := 0
						if hotEvent, ok := v["hotEvent"].(map[string]interface{}); ok {
							if score, ok := hotEvent["hotScore"].(float64); ok {
								hotScore = int(score)
							}
						}

						listData = append(listData, &HotItem{
							ID:        "0",
							Title:     title,
							Desc:      desc,
							Cover:     cover,
							Author:    author,
							Timestamp: timestamp,
							Hot:       hotScore,
							URL:       fmt.Sprintf("https://new.qq.com/rain/a/%s", id),
							MobileURL: fmt.Sprintf("https://view.inews.qq.com/k/%s", id),
						})
					}
				}
			}
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "qq",
		Title: "腾讯新闻",
		Type:  "热点榜",
		Link:  "https://news.qq.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 51CTO ==============
func (s *Spider) Get51CTO() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[Source51CTO].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[Source51CTO].Agent)

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
		if innerData, ok := data["data"].(map[string]interface{}); ok {
			if list, ok := innerData["list"].([]interface{}); ok {
				for _, item := range list {
					if v, ok := item.(map[string]interface{}); ok {
						title := ""
						if t, ok := v["title"].(string); ok {
							title = t
						}

						desc := ""
						if d, ok := v["abstract"].(string); ok {
							desc = d
						}

						cover := ""
						if img, ok := v["cover"].(string); ok {
							cover = img
						}

						timestamp := int64(0)
						if ts, ok := v["pubdate"].(float64); ok {
							timestamp = int64(ts)
						}

						sourceID := ""
						if id, ok := v["source_id"].(string); ok {
							sourceID = id
						}

						listData = append(listData, &HotItem{
							ID:        "0",
							Title:     title,
							Desc:      desc,
							Cover:     cover,
							Timestamp: timestamp,
							Hot:       0,
							URL:       fmt.Sprintf("https://www.51cto.com/article/%s.html", sourceID),
							MobileURL: fmt.Sprintf("https://www.51cto.com/article/%s.html", sourceID),
						})
					}
				}
			}
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "51cto",
		Title: "51CTO",
		Type:  "推荐榜",
		Link:  "https://www.51cto.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 吾爱破解 ==============
func (s *Spider) Get52Pojie() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[Source52Pojie].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[Source52Pojie].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 需要处理RSS解析，这里简化实现
	// 实际实现需要解析RSS XML内容

	return &HotData{
		Code:  http.StatusOK,
		Name:  "52pojie",
		Title: "吾爱破解",
		Type:  "最新精华",
		Link:  "https://www.52pojie.cn/",
		Total: 0,
		Data:  []*HotItem{},
	}, nil
}

// ============== 豆瓣讨论 ==============
func (s *Spider) GetDoubanGroup() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceDoubanGroup].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceDoubanGroup].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 需要HTML解析，这里简化实现
	// 实际实现需要使用goquery解析HTML

	return &HotData{
		Code:  http.StatusOK,
		Name:  "douban-group",
		Title: "豆瓣讨论",
		Type:  "讨论精选",
		Link:  "https://www.douban.com/group/explore",
		Total: 0,
		Data:  []*HotItem{},
	}, nil
}

// ============== AcFun ==============
func (s *Spider) GetAcfun() (*HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[SourceAcfun].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[SourceAcfun].Agent)
	request.Header.Add("Referer", "https://www.acfun.cn/rank/list/?cid=-1&pcid=-1&range=DAY")

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
		if rankList, ok := data["rankList"].([]interface{}); ok {
			for _, item := range rankList {
				if v, ok := item.(map[string]interface{}); ok {
					title := ""
					if t, ok := v["contentTitle"].(string); ok {
						title = t
					}

					desc := ""
					if d, ok := v["contentDesc"].(string); ok {
						desc = d
					}

					cover := ""
					if img, ok := v["coverUrl"].(string); ok {
						cover = img
					}

					author := ""
					if user, ok := v["userName"].(string); ok {
						author = user
					}

					timestamp := int64(0)
					if ts, ok := v["contributeTime"].(float64); ok {
						timestamp = int64(ts)
					}

					dougaId := ""
					if id, ok := v["dougaId"].(string); ok {
						dougaId = id
					}

					likeCount := 0
					if likes, ok := v["likeCount"].(float64); ok {
						likeCount = int(likes)
					}

					listData = append(listData, &HotItem{
						ID:        "0",
						Title:     title,
						Desc:      desc,
						Cover:     cover,
						Author:    author,
						Timestamp: timestamp,
						Hot:       likeCount,
						URL:       fmt.Sprintf("https://www.acfun.cn/v/ac%s", dougaId),
						MobileURL: fmt.Sprintf("https://m.acfun.cn/v/?ac=%s", dougaId),
					})
				}
			}
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "acfun",
		Title: "AcFun",
		Type:  "排行榜",
		Link:  "https://www.acfun.cn/rank/list/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 数字尾巴 ==============
func (s *Spider) GetDgtle() (*HotData, error) {
	urlConf := s.UrlMap[SourceDgtle]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for dgtle")
	}

	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Items []struct {
				ID        int    `json:"id"`
				Title     string `json:"title"`
				Content   string `json:"content"`
				Cover     string `json:"cover"`
				From      string `json:"from"`
				Membernum int    `json:"membernum"`
				CreatedAt string `json:"created_at"`
				Type      string `json:"type"`
			} `json:"items"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range result.Data.Items {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", item.CreatedAt)

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(item.ID),
			Title:     item.Title,
			Desc:      item.Content,
			Cover:     item.Cover,
			Author:    item.From,
			Timestamp: timestamp.Unix(),
			Hot:       item.Membernum,
			URL:       fmt.Sprintf("https://www.dgtle.com/news-%d-%s.html", item.ID, item.Type),
			MobileURL: fmt.Sprintf("https://m.dgtle.com/news-details/%d", item.ID),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "dgtle",
		Title: "数字尾巴",
		Type:  "热门文章",
		Link:  "https://www.dgtle.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 豆瓣电影 ==============
func (s *Spider) GetDoubanMovie() (*HotData, error) {
	urlConf := s.UrlMap[SourceDoubanMovie]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for douban-movie")
	}

	req, err := http.NewRequest("GET", urlConf.Url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", urlConf.Agent)

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	doc.Find(".article tr.item").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find("a").Attr("href")
		score := s.Find(".rating_nums").Text()
		if score == "" {
			score = "0.0"
		}
		title := s.Find("a").AttrOr("title", "")
		cover := s.Find("img").AttrOr("src", "")

		// 从URL中提取ID
		id := 0
		if url != "" {
			if strings.Contains(url, "/subject/") {
				parts := strings.Split(url, "/")
				if len(parts) >= 5 {
					id, _ = strconv.Atoi(parts[4])
				}
			}
		}

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(id),
			Title:     fmt.Sprintf("【%s】%s", score, title),
			Cover:     cover,
			URL:       url,
			MobileURL: url,
		})
	})

	return &HotData{
		Code:  http.StatusOK,
		Name:  "douban-movie",
		Title: "豆瓣电影",
		Type:  "新片榜",
		Link:  "https://movie.douban.com/chart",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 中国地震台 ==============
func (s *Spider) GetEarthquake() (*HotData, error) {
	urlConf := s.UrlMap[SourceEarthquake]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for earthquake")
	}

	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 使用正则表达式提取JSON数据
	re := regexp.MustCompile(`const newdata = (\[.*?\]);`)
	match := re.FindSubmatch(body)
	if match == nil {
		return nil, fmt.Errorf("no earthquake data found")
	}

	var earthquakes []struct {
		NEW_DID    string  `json:"NEW_DID"`
		LOCATION_C string  `json:"LOCATION_C"`
		M          float64 `json:"M"`
		O_TIME     string  `json:"O_TIME"`
		EPI_LAT    float64 `json:"EPI_LAT"`
		EPI_LON    float64 `json:"EPI_LON"`
		EPI_DEPTH  float64 `json:"EPI_DEPTH"`
		SAVE_TIME  string  `json:"SAVE_TIME"`
	}

	if err := json.Unmarshal(match[1], &earthquakes); err != nil {
		return nil, err
	}

	var listData []*HotItem
	mappings := map[string]string{
		"O_TIME":     "发震时刻(UTC+8)",
		"LOCATION_C": "参考位置",
		"M":          "震级(M)",
		"EPI_LAT":    "纬度(°)",
		"EPI_LON":    "经度(°)",
		"EPI_DEPTH":  "深度(千米)",
		"SAVE_TIME":  "录入时间",
	}

	for _, eq := range earthquakes {
		var contentBuilder []string
		for key, desc := range mappings {
			var value string
			switch key {
			case "O_TIME", "LOCATION_C", "SAVE_TIME":
				value = eq.O_TIME
				if key == "LOCATION_C" {
					value = eq.LOCATION_C
				} else if key == "SAVE_TIME" {
					value = eq.SAVE_TIME
				}
			case "M", "EPI_LAT", "EPI_LON", "EPI_DEPTH":
				value = fmt.Sprintf("%.1f", eq.M)
				if key == "EPI_LAT" {
					value = fmt.Sprintf("%.2f", eq.EPI_LAT)
				} else if key == "EPI_LON" {
					value = fmt.Sprintf("%.2f", eq.EPI_LON)
				} else if key == "EPI_DEPTH" {
					value = fmt.Sprintf("%.0f", eq.EPI_DEPTH)
				}
			}
			contentBuilder = append(contentBuilder, fmt.Sprintf("%s：%s", desc, value))
		}

		timestamp, _ := time.Parse("2006-01-02 15:04:05", eq.O_TIME)

		listData = append(listData, &HotItem{
			ID:        "0",
			Title:     fmt.Sprintf("%s发生%.1f级地震", eq.LOCATION_C, eq.M),
			Desc:      strings.Join(contentBuilder, "\n"),
			Timestamp: timestamp.Unix(),
			URL:       "https://news.ceic.ac.cn/",
			MobileURL: "https://news.ceic.ac.cn/",
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "earthquake",
		Title: "中国地震台",
		Type:  "地震速报",
		Link:  "https://news.ceic.ac.cn/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== GameRes游资网 ==============
func (s *Spider) GetGameres() (*HotData, error) {
	urlConf := s.UrlMap[SourceGameres]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for gameres")
	}

	req, err := http.NewRequest("GET", urlConf.Url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", urlConf.Agent)

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	doc.Find(".article-list .article-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".article-title a").Text()
		url, _ := s.Find(".article-title a").Attr("href")
		if !strings.HasPrefix(url, "http") {
			url = "https://www.gameres.com" + url
		}

		cover := s.Find(".article-cover img").AttrOr("src", "")
		if !strings.HasPrefix(cover, "http") && cover != "" {
			cover = "https://www.gameres.com" + cover
		}

		author := s.Find(".article-author").Text()
		publishTime := s.Find(".article-time").Text()

		timestamp := time.Now().Unix()
		if publishTime != "" {
			if t, err := time.Parse("2006-01-02", publishTime); err == nil {
				timestamp = t.Unix()
			}
		}

		// 从URL中提取ID
		id := 0
		if url != "" {
			if strings.Contains(url, "/article/") {
				parts := strings.Split(url, "/")
				if len(parts) >= 4 {
					idStr := parts[len(parts)-1]
					if strings.Contains(idStr, ".") {
						idStr = strings.Split(idStr, ".")[0]
					}
					id, _ = strconv.Atoi(idStr)
				}
			}
		}

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(id),
			Title:     title,
			Cover:     cover,
			Author:    author,
			Timestamp: timestamp,
			URL:       url,
			MobileURL: url,
		})
	})

	return &HotData{
		Code:  http.StatusOK,
		Name:  "gameres",
		Title: "GameRes游资网",
		Type:  "资讯",
		Link:  "https://www.gameres.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 极客公园 ==============
func (s *Spider) GetGeekpark() (*HotData, error) {
	urlConf := s.UrlMap[SourceGeekpark]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for geekpark")
	}

	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Cover       string `json:"cover"`
			Author      string `json:"author"`
			ReadCount   int    `json:"read_count"`
			CreatedAt   string `json:"created_at"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range result.Data {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", item.CreatedAt)

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(item.ID),
			Title:     item.Title,
			Desc:      item.Description,
			Cover:     item.Cover,
			Author:    item.Author,
			Timestamp: timestamp.Unix(),
			Hot:       item.ReadCount,
			URL:       fmt.Sprintf("https://www.geekpark.net/news/%d", item.ID),
			MobileURL: fmt.Sprintf("https://m.geekpark.net/news/%d", item.ID),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "geekpark",
		Title: "极客公园",
		Type:  "热门文章",
		Link:  "https://www.geekpark.net/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 原神 ==============
func (s *Spider) GetGenshin() (*HotData, error) {
	urlConf := s.UrlMap[SourceGenshin]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for genshin")
	}

	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			List []struct {
				Post struct {
					PostID    int    `json:"post_id"`
					Subject   string `json:"subject"`
					Content   string `json:"content"`
					Cover     string `json:"cover"`
					CreatedAt int64  `json:"created_at"`
					Author    struct {
						Nickname string `json:"nickname"`
					} `json:"user"`
				} `json:"post"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range result.Data.List {
		post := item.Post

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(post.PostID),
			Title:     post.Subject,
			Desc:      post.Content,
			Cover:     post.Cover,
			Author:    post.Author.Nickname,
			Timestamp: post.CreatedAt,
			URL:       fmt.Sprintf("https://bbs.mihoyo.com/ys/article/%d", post.PostID),
			MobileURL: fmt.Sprintf("https://m.bbs.mihoyo.com/ys/article/%d", post.PostID),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "genshin",
		Title: "原神",
		Type:  "动态",
		Link:  "https://bbs.mihoyo.com/ys/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// GetGuokr 获取果壳热门文章
func (s *Spider) GetGuokr() (*HotData, error) {
	urlConf := s.UrlMap[SourceGuokr]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for guokr")
	}

	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []struct {
			ID         string `json:"id"`
			Title      string `json:"title"`
			Summary    string `json:"summary"`
			SmallImage string `json:"small_image"`
			Author     struct {
				Nickname string `json:"nickname"`
			} `json:"author"`
			DateModified string `json:"date_modified"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range result.Data {
		timestamp, _ := time.Parse(time.RFC3339, item.DateModified)
		idInt, _ := strconv.Atoi(item.ID)
		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     item.Title,
			Desc:      item.Summary,
			Cover:     item.SmallImage,
			Author:    item.Author.Nickname,
			Timestamp: timestamp.Unix(),
			URL:       fmt.Sprintf("https://www.guokr.com/article/%s", item.ID),
			MobileURL: fmt.Sprintf("https://m.guokr.com/article/%s", item.ID),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "guokr",
		Title: "果壳",
		Type:  "热门文章",
		Link:  "https://www.guokr.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// GetHackernews 获取Hacker News热门文章
func (s *Spider) GetHackernews() (*HotData, error) {
	urlConf := s.UrlMap[SourceHackernews]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for hackernews")
	}

	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
		id := s.AttrOr("id", "")
		title := s.Find(".titleline a").First().Text()
		url := s.Find(".titleline a").First().AttrOr("href", "")

		if id != "" && title != "" {
			// 获取分数
			scoreText := doc.Find(fmt.Sprintf("#score_%s", id)).Text()
			var hot int
			if scoreMatch := regexp.MustCompile(`\d+`).FindString(scoreText); scoreMatch != "" {
				if score, err := strconv.Atoi(scoreMatch); err == nil {
					hot = score
				}
			}

			// 将字符串ID转换为int
			idInt, _ := strconv.Atoi(id)

			// 处理相对URL
			if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "//") {
				url = "https://news.ycombinator.com/" + url
			}

			listData = append(listData, &HotItem{
				ID:    strconv.Itoa(idInt),
				Title: title,
				Hot:   hot,
				URL:   url,
			})
		}
	})

	return &HotData{
		Code:  http.StatusOK,
		Name:  "hackernews",
		Title: "Hacker News",
		Type:  "Popular",
		Link:  "https://news.ycombinator.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// GetHelloGitHub 获取HelloGitHub热门仓库
func (s *Spider) GetHelloGitHub() (*HotData, error) {
	urlConf := s.UrlMap[SourceHelloGitHub]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for hellogithub")
	}

	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []struct {
			ItemID      string `json:"item_id"`
			Title       string `json:"title"`
			Summary     string `json:"summary"`
			Author      string `json:"author"`
			UpdatedAt   string `json:"updated_at"`
			ClicksTotal int    `json:"clicks_total"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range result.Data {
		timestamp, _ := time.Parse("2006-01-02T15:04:05", item.UpdatedAt)
		listData = append(listData, &HotItem{
			ID:        item.ItemID,
			Title:     item.Title,
			Desc:      item.Summary,
			Author:    item.Author,
			Timestamp: timestamp.Unix(),
			Hot:       item.ClicksTotal,
			URL:       fmt.Sprintf("https://hellogithub.com/repository/%s", item.ItemID),
			MobileURL: fmt.Sprintf("https://hellogithub.com/repository/%s", item.ItemID),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "hellogithub",
		Title: "HelloGitHub",
		Type:  "热门仓库",
		Link:  "https://hellogithub.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// GetHistory 获取历史上的今天
func (s *Spider) GetHistory() (*HotData, error) {
	urlConf := s.UrlMap[SourceHistory]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for history")
	}

	// 获取当前日期
	now := time.Now()
	month := now.Month()
	day := now.Day()

	// 构建URL，使用当前月份
	url := fmt.Sprintf("https://baike.baidu.com/cms/home/eventsOnHistory/%02d.json", month)

	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]map[string][]struct {
		Title    string `json:"title"`
		Desc     string `json:"desc"`
		Year     string `json:"year"`
		Link     string `json:"link"`
		PicShare string `json:"pic_share"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// 获取当前日期的数据
	monthKey := fmt.Sprintf("%02d", month)
	dayKey := fmt.Sprintf("%02d%02d", month, day)
	dayData := result[monthKey][dayKey]

	var listData []*HotItem
	for i, item := range dayData {
		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(i),
			Title:     strings.TrimSpace(item.Title),
			Desc:      strings.TrimSpace(item.Desc),
			Cover:     item.PicShare,
			Author:    "",
			Timestamp: 0,
			Hot:       0,
			URL:       item.Link,
			MobileURL: item.Link,
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "history",
		Title: "历史上的今天",
		Type:  fmt.Sprintf("%02d-%02d", month, day),
		Link:  "https://baike.baidu.com/calendar",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// GetHonkai 获取崩坏3最新动态
func (s *Spider) GetHonkai() (*HotData, error) {
	urlConf := s.UrlMap[SourceHonkai]
	if urlConf == nil {
		return nil, fmt.Errorf("url config not found for honkai")
	}

	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			List []struct {
				Post struct {
					PostID     int      `json:"post_id"`
					Subject    string   `json:"subject"`
					Content    string   `json:"content"`
					Cover      string   `json:"cover"`
					Images     []string `json:"images"`
					ViewStatus int      `json:"view_status"`
					CreatedAt  int64    `json:"created_at"`
				} `json:"post"`
				User struct {
					Nickname string `json:"nickname"`
				} `json:"user"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range result.Data.List {
		post := item.Post

		// 获取封面图，优先使用cover，如果没有则使用第一张图片
		cover := post.Cover
		if cover == "" && len(post.Images) > 0 {
			cover = post.Images[0]
		}

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(post.PostID),
			Title:     post.Subject,
			Desc:      post.Content,
			Cover:     cover,
			Author:    item.User.Nickname,
			Timestamp: post.CreatedAt,
			Hot:       post.ViewStatus,
			URL:       fmt.Sprintf("https://www.miyoushe.com/bh3/article/%d", post.PostID),
			MobileURL: fmt.Sprintf("https://m.miyoushe.com/bh3/#/article/%d", post.PostID),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "honkai",
		Title: "崩坏3",
		Type:  "最新动态",
		Link:  "https://www.miyoushe.com/bh3/home/6",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// GetHostloc 获取hostloc论坛热门帖子
func (s *Spider) GetHostloc() (*HotData, error) {
	urlConf := s.UrlMap[SourceHostloc]
	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析RSS XML
	feed, err := ParseRSS(resp.Body)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range feed.Channel.Items {
		// 将字符串Guid转换为int
		idInt, _ := strconv.Atoi(item.Guid)
		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     item.Title,
			Desc:      item.Description,
			Author:    item.Author,
			Timestamp: parseTime(item.PubDate),
			URL:       item.Link,
			MobileURL: item.Link,
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "hostloc",
		Title: "全球主机交流论坛",
		Type:  "热门帖子",
		Link:  "https://hostloc.com/forum.php?mod=guide&view=hot",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// GetIfanr 获取爱范儿快讯数据
func (s *Spider) GetIfanr() (*HotData, error) {
	urlConf := s.UrlMap[SourceIfanr]
	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
			Time    int64  `json:"time"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range result.Data {
		idInt, _ := strconv.Atoi(item.ID)
		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     item.Title,
			Desc:      item.Content,
			Timestamp: item.Time,
			URL:       fmt.Sprintf("https://www.ifanr.com/buzz/%s", item.ID),
			MobileURL: fmt.Sprintf("https://www.ifanr.com/buzz/%s", item.ID),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "ifanr",
		Title: "爱范儿",
		Type:  "快讯",
		Link:  "https://www.ifanr.com/buzz",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// GetIthomeXijiayi 获取IT之家喜加一游戏动态
func (s *Spider) GetIthomeXijiayi() (*HotData, error) {
	urlConf := s.UrlMap[SourceIthomeXijiayi]
	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	doc.Find(".blk-container .list-item").Each(func(i int, s *goquery.Selection) {
		id := s.AttrOr("data-id", "")
		title := s.Find(".title").Text()
		desc := s.Find(".desc").Text()
		timeStr := s.Find(".time").Text()

		timestamp := parseTime(timeStr)
		url := fmt.Sprintf("https://www.ithome.com/zt/xijiayi#%s", id)

		idInt, _ := strconv.Atoi(id)
		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     strings.TrimSpace(title),
			Desc:      strings.TrimSpace(desc),
			Timestamp: timestamp,
			URL:       url,
			MobileURL: url,
		})
	})

	return &HotData{
		Code:  http.StatusOK,
		Name:  "ithome-xijiayi",
		Title: "IT之家喜加一",
		Type:  "游戏动态",
		Link:  "https://www.ithome.com/zt/xijiayi",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 米游社 ==============
func (s *Spider) GetMiyoushe() (*HotData, error) {
	url := s.UrlMap[SourceMiyoushe].Url
	agent := s.UrlMap[SourceMiyoushe].Agent

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", agent)

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	var result struct {
		Data struct {
			List []struct {
				Post struct {
					PostID     string   `json:"post_id"`
					Subject    string   `json:"subject"`
					Content    string   `json:"content"`
					Cover      string   `json:"cover"`
					Images     []string `json:"images"`
					CreatedAt  int64    `json:"created_at"`
					ViewStatus int      `json:"view_status"`
				} `json:"post"`
				User struct {
					Nickname string `json:"nickname"`
				} `json:"user"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	var listData []*HotItem
	for _, item := range result.Data.List {
		post := item.Post
		cover := post.Cover
		if cover == "" && len(post.Images) > 0 {
			cover = post.Images[0]
		}

		idInt, _ := strconv.Atoi(post.PostID)
		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     strings.TrimSpace(post.Subject),
			Desc:      strings.TrimSpace(post.Content),
			Cover:     cover,
			Author:    item.User.Nickname,
			Timestamp: post.CreatedAt,
			Hot:       post.ViewStatus,
			URL:       fmt.Sprintf("https://www.miyoushe.com/ys/article/%s", post.PostID),
			MobileURL: fmt.Sprintf("https://m.miyoushe.com/ys/#/article/%s", post.PostID),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "miyoushe",
		Title: "米游社",
		Type:  "最新公告",
		Link:  "https://www.miyoushe.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 水木社区 ==============
func (s *Spider) GetNewsmth() (*HotData, error) {
	url := s.UrlMap[SourceNewsmth].Url
	agent := s.UrlMap[SourceNewsmth].Agent

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("User-Agent", agent)

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	var result struct {
		Data struct {
			Topics []struct {
				FirstArticleID string `json:"firstArticleId"`
				Article        struct {
					TopicID  string `json:"topicId"`
					Subject  string `json:"subject"`
					Body     string `json:"body"`
					PostTime int64  `json:"postTime"`
					Account  struct {
						Name string `json:"name"`
					} `json:"account"`
				} `json:"article"`
				Board struct {
					Title string `json:"title"`
				} `json:"board"`
			} `json:"topics"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	var listData []*HotItem
	for _, topic := range result.Data.Topics {
		post := topic.Article
		url := fmt.Sprintf("https://wap.newsmth.net/article/%s?title=%s&from=home", post.TopicID, topic.Board.Title)

		idInt, _ := strconv.Atoi(topic.FirstArticleID)
		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     strings.TrimSpace(post.Subject),
			Desc:      strings.TrimSpace(post.Body),
			Author:    post.Account.Name,
			Timestamp: post.PostTime,
			URL:       url,
			MobileURL: url,
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "newsmth",
		Title: "水木社区",
		Type:  "热门话题",
		Link:  "https://www.newsmth.net/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== NGA ==============
func (s *Spider) GetNgabbs() (*HotData, error) {
	url := s.UrlMap[SourceNgabbs].Url
	// agent := s.UrlMap[SourceNgabbs].Agent

	reqBody := strings.NewReader(`{ __output: "14"}`)
	// reqBody := strings.NewReader(`__output=14`)
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	// req.Header.Set("User-Agent", agent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "ngabbs.com")
	req.Header.Set("Referer", "https://ngabbs.com/")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "11")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-User-Agent", "NGA_skull/7.3.1(iPhone13,2;iOS 17.2.1)")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	// 处理gzip压缩的响应体
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("创建gzip读取器失败: %v", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	var result struct {
		Result [][]struct {
			Tid      string `json:"tid"`
			Subject  string `json:"subject"`
			Author   string `json:"author"`
			Replies  int    `json:"replies"`
			Postdate int64  `json:"postdate"`
			Tpcurl   string `json:"tpcurl"`
		} `json:"result"`
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	var listData []*HotItem
	for _, item := range result.Result[0] {
		tidInt, _ := strconv.Atoi(item.Tid)
		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(tidInt),
			Title:     strings.TrimSpace(item.Subject),
			Author:    item.Author,
			Hot:       item.Replies,
			Timestamp: item.Postdate,
			URL:       fmt.Sprintf("https://bbs.nga.cn%s", item.Tpcurl),
			MobileURL: fmt.Sprintf("https://bbs.nga.cn%s", item.Tpcurl),
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "ngabbs",
		Title: "NGA",
		Type:  "论坛热帖",
		Link:  "https://ngabbs.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== NodeSeek ==============
func (s *Spider) GetNodeseek() (*HotData, error) {
	url := s.UrlMap[SourceNodeseek].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssData struct {
		Channel struct {
			Item []struct {
				Title       string `xml:"title"`
				Link        string `xml:"link"`
				Description string `xml:"description"`
				PubDate     string `xml:"pubDate"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	err = xml.Unmarshal(body, &rssData)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range rssData.Channel.Item {
		timestamp, _ := time.Parse(time.RFC1123, item.PubDate)

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(len(listData) + 1),
			Title:     item.Title,
			Desc:      item.Description,
			Timestamp: timestamp.Unix(),
			URL:       item.Link,
			MobileURL: item.Link,
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "nodeseek",
		Title: "NodeSeek",
		Type:  "技术社区",
		Link:  "https://www.nodeseek.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 纽约时报 ==============
func (s *Spider) GetNytimes() (*HotData, error) {
	url := s.UrlMap[SourceNytimes].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssData struct {
		Channel struct {
			Item []struct {
				Title       string `xml:"title"`
				Link        string `xml:"link"`
				Description string `xml:"description"`
				PubDate     string `xml:"pubDate"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	err = xml.Unmarshal(body, &rssData)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range rssData.Channel.Item {
		timestamp, _ := time.Parse(time.RFC1123, item.PubDate)

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(len(listData) + 1),
			Title:     item.Title,
			Desc:      item.Description,
			Timestamp: timestamp.Unix(),
			URL:       item.Link,
			MobileURL: item.Link,
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "nytimes",
		Title: "纽约时报",
		Type:  "国际新闻",
		Link:  "https://www.nytimes.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== Product Hunt ==============
func (s *Spider) GetProducthunt() (*HotData, error) {
	url := s.UrlMap[SourceProducthunt].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 使用正则表达式提取产品信息
	re := regexp.MustCompile(`<div data-test="post-item"[^>]*>.*?<a[^>]*href="([^"]*)"[^>]*>.*?<h3[^>]*>([^<]*)</h3>.*?<div[^>]*>([^<]*)</div>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var listData []*HotItem
	for i, match := range matches {
		if len(match) >= 4 {
			url := "https://www.producthunt.com" + match[1]

			listData = append(listData, &HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Desc:      strings.TrimSpace(match[3]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "producthunt",
		Title: "Product Hunt",
		Type:  "产品发现",
		Link:  "https://www.producthunt.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 新浪新闻 ==============
func (s *Spider) GetSinaNews() (*HotData, error) {
	url := s.UrlMap[SourceSinaNews].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 使用正则表达式提取新闻信息
	re := regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var listData []*HotItem
	for i, match := range matches {
		if len(match) >= 3 && strings.Contains(match[1], "news.sina.com.cn") {
			url := match[1]
			if !strings.HasPrefix(url, "http") {
				url = "https:" + url
			}

			listData = append(listData, &HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "sina-news",
		Title: "新浪新闻",
		Type:  "新闻资讯",
		Link:  "https://news.sina.com.cn/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 新浪微博 ==============
func (s *Spider) GetSina() (*HotData, error) {
	url := s.UrlMap[SourceSina].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 使用正则表达式提取微博热搜信息
	re := regexp.MustCompile(`<a[^>]*href="([^"]*weibo.com[^"]*)"[^>]*>([^<]*)</a>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var listData []*HotItem
	for i, match := range matches {
		if len(match) >= 3 {
			url := match[1]
			if !strings.HasPrefix(url, "http") {
				url = "https:" + url
			}

			listData = append(listData, &HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "sina",
		Title: "微博热搜",
		Type:  "社交媒体",
		Link:  "https://s.weibo.com/top/summary",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 星穹铁道 ==============
func (s *Spider) GetStarrail() (*HotData, error) {
	url := s.UrlMap[SourceStarrail].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData struct {
		Data struct {
			List []struct {
				Post struct {
					PostId    int64  `json:"post_id"`
					Subject   string `json:"subject"`
					Content   string `json:"content"`
					CreatedAt int64  `json:"created_at"`
					ViewCnt   int    `json:"view_cnt"`
				} `json:"post"`
			} `json:"list"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range jsonData.Data.List {
		post := item.Post
		url := fmt.Sprintf("https://bbs.miyoushe.com/detail/%d", post.PostId)

		listData = append(listData, &HotItem{
			ID:        strconv.Itoa(int(post.PostId)),
			Title:     post.Subject,
			Desc:      post.Content,
			Timestamp: post.CreatedAt,
			Hot:       post.ViewCnt,
			URL:       url,
			MobileURL: url,
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "starrail",
		Title: "星穹铁道",
		Type:  "游戏资讯",
		Link:  "https://bbs.miyoushe.com/ys/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 澎湃新闻 ==============
func (s *Spider) GetThepaper() (*HotData, error) {
	url := s.UrlMap[SourceThepaper].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 使用正则表达式提取新闻信息
	re := regexp.MustCompile(`<a[^>]*href="([^"]*thepaper.cn[^"]*)"[^>]*>([^<]*)</a>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var listData []*HotItem
	for i, match := range matches {
		if len(match) >= 3 {
			url := match[1]
			if !strings.HasPrefix(url, "http") {
				url = "https:" + url
			}

			listData = append(listData, &HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "thepaper",
		Title: "澎湃新闻",
		Type:  "新闻资讯",
		Link:  "https://www.thepaper.cn/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 气象预警 ==============
func (s *Spider) GetWeatheralarm() (*HotData, error) {
	url := s.UrlMap[SourceWeatheralarm].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData struct {
		Data struct {
			Page struct {
				List []struct {
					Alertid   string `json:"alertid"`
					Title     string `json:"title"`
					Issuetime string `json:"issuetime"`
					Pic       string `json:"pic"`
				} `json:"list"`
			} `json:"page"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, item := range jsonData.Data.Page.List {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", item.Issuetime)
		url := fmt.Sprintf("http://www.nmc.cn/publish/alarm.html?alertid=%s", item.Alertid)

		listData = append(listData, &HotItem{
			ID:        item.Alertid,
			Title:     item.Title,
			Desc:      fmt.Sprintf("发布时间: %s", item.Issuetime),
			Cover:     item.Pic,
			Timestamp: timestamp.Unix(),
			URL:       url,
			MobileURL: url,
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "weatheralarm",
		Title: "气象预警",
		Type:  "气象信息",
		Link:  "http://www.nmc.cn/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 微信读书 ==============
func (s *Spider) GetWeread() (*HotData, error) {
	url := s.UrlMap[SourceWeread].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData struct {
		Books []struct {
			BookId       string `json:"bookId"`
			Title        string `json:"title"`
			Author       string `json:"author"`
			Intro        string `json:"intro"`
			Cover        string `json:"cover"`
			PublishTime  int64  `json:"publishTime"`
			ReadingCount int    `json:"readingCount"`
		} `json:"books"`
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	var listData []*HotItem
	for _, book := range jsonData.Books {
		cover := strings.Replace(book.Cover, "_s.jpg", "_l.jpg", 1)
		url := fmt.Sprintf("https://weread.qq.com/web/bookDetail/%s", book.BookId)

		listData = append(listData, &HotItem{
			ID:        book.BookId,
			Title:     book.Title,
			Author:    book.Author,
			Desc:      book.Intro,
			Cover:     cover,
			Timestamp: book.PublishTime,
			Hot:       book.ReadingCount,
			URL:       url,
			MobileURL: url,
		})
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "weread",
		Title: "微信读书",
		Type:  "图书排行",
		Link:  "https://weread.qq.com/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// ============== 游研社 ==============
func (s *Spider) GetYystv() (*HotData, error) {
	url := s.UrlMap[SourceYystv].Url
	resp, err := s.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 使用正则表达式提取文章信息
	re := regexp.MustCompile(`<a[^>]*href="([^"]*yystv.cn[^"]*)"[^>]*>([^<]*)</a>`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var listData []*HotItem
	for i, match := range matches {
		if len(match) >= 3 {
			url := match[1]
			if !strings.HasPrefix(url, "http") {
				url = "https:" + url
			}

			listData = append(listData, &HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &HotData{
		Code:  http.StatusOK,
		Name:  "yystv",
		Title: "游研社",
		Type:  "游戏资讯",
		Link:  "https://www.yystv.cn/",
		Total: len(listData),
		Data:  listData,
	}, nil
}

// RSSFeed RSS订阅数据结构
type RSSFeed struct {
	Channel struct {
		Items []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem RSS项目结构
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
	Author      string `xml:"author"`
}

// ParseRSS 解析RSS XML内容
func ParseRSS(body io.Reader) (*RSSFeed, error) {
	var feed RSSFeed
	decoder := xml.NewDecoder(body)
	err := decoder.Decode(&feed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS: %v", err)
	}
	return &feed, nil
}

// parseTime 解析各种时间格式为Unix时间戳
func parseTime(timeStr string) int64 {
	if timeStr == "" {
		return time.Now().Unix()
	}

	// 尝试多种时间格式
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC3339,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t.Unix()
		}
	}

	// 如果所有格式都失败，返回当前时间
	return time.Now().Unix()
}
