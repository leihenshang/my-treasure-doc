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

	"fastduck/treasure-doc/module/hot_top/conf"
	"fastduck/treasure-doc/module/hot_top/hot/token"
	"fastduck/treasure-doc/module/hot_top/model"

	"log"

	"github.com/PuerkitoBio/goquery"
)

type Spider struct {
	UrlMap     map[model.Source]*conf.UrlConf
	HttpClient *http.Client
}

var spider *Spider
var spiderOnce *sync.Once = &sync.Once{}

func NewSpider(urlMap map[model.Source]*conf.UrlConf) (*Spider, error) {
	if len(urlMap) == 0 {
		return nil, fmt.Errorf("urlMap is empty")
	}
	spiderOnce.Do(func() {
		spider = &Spider{
			UrlMap: urlMap,
			HttpClient: &http.Client{
				Timeout: time.Duration(5 * time.Second),
			},
		}
	})
	return spider, nil
}

func GetSpider() *Spider {
	return spider
}

// ============== IT之家 ==============
func (s *Spider) GetItHome() (*model.HotData, error) {
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
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceITHome].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceITHome].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

		hotItem := &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "ithome",
		Title:       "IT之家",
		Type:        "热榜",
		Description: "爱科技，爱这里 - 前沿科技新闻网站",
		Link:        "https://m.ithome.com/rankm/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 知乎 ==============
func (s *Spider) GetZhihu() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceZhihu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceZhihu].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

				listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "zhihu",
		Title:      "知乎",
		Type:       "热榜",
		Link:       "https://www.zhihu.com/hot",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 微博 ==============
func (s *Spider) GetWeibo() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceWeibo].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceWeibo].Agent)
	request.Header.Add("Referer", "https://s.weibo.com/top/summary?cate=realtimehot")
	request.Header.Add("MWeibo-Pwa", "1")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Data struct {
			Cards []struct {
				CardGroup []struct {
					Desc   string  `json:"desc"`
					Num    float64 `json:"num"`
					ItemID string  `json:"itemid"`
					Scheme string  `json:"scheme"`
					Pic    string  `json:"pic"`
				} `json:"card_group"`
			} `json:"cards"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	for _, v := range result.Data.Cards {
		for _, vv := range v.CardGroup {
			listData = append(listData, &model.HotItem{
				ID:        vv.ItemID,
				Title:     vv.Desc,
				Desc:      vv.Desc,
				Timestamp: time.Now().Unix(),
				Hot:       int(vv.Num),
				Cover:     vv.Pic,
				URL:       fmt.Sprintf("https://s.weibo.com/weibo?q=%s&t=31&band_rank=1&Refer=top", strings.ReplaceAll(vv.Desc, "#", "")),
				MobileURL: vv.Scheme,
			})
		}
	}

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "weibo",
		Title:       "微博",
		Type:        "热搜榜",
		Description: "实时热点，每分钟更新一次",
		Link:        "https://s.weibo.com/top/summary/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 哔哩哔哩 ==============
func (s *Spider) GetBilibili() (*model.HotData, error) {
	wbi, err := token.GetBilibiliWbi()
	if err != nil {
		return nil, err
	}

	url := s.UrlMap[model.SourceBilibili].Url + "&" + wbi
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceBilibili].Agent)
	request.Header.Add("Referer", "https://www.bilibili.com/ranking/all")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	request.Header.Add("Accept-Encoding", "gzip, deflate, br")
	request.Header.Add("Sec-Ch-Ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	request.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	request.Header.Add("Sec-Ch-Ua-Platform", `"Windows"`)
	request.Header.Add("Sec-Fetch-Dest", "document")
	request.Header.Add("Sec-Fetch-Mode", "navigate")
	request.Header.Add("Sec-Fetch-Site", "same-origin")
	request.Header.Add("Sec-Fetch-User", "?1")
	request.Header.Add("Upgrade-Insecure-Requests", "1")

	resp, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("unmarshal Bilibili data failed, err: %v \n", err)
		bakRequest, err := http.NewRequest("GET", `https://api.bilibili.com/x/web-interface/ranking?jsonp=jsonp?rid=${type}&type=all&callback=__jp0`, nil)
		if err != nil {
			return nil, err
		}
		bakRequest.Header.Set("User-Agent", s.UrlMap[model.SourceBilibili].Agent)
		bakRequest.Header.Set("Referer", "https://www.bilibili.com/ranking/all")
		bakResp, err := s.HttpClient.Do(bakRequest)
		if err != nil {
			return nil, err
		}
		defer bakResp.Body.Close()
		bakRespBody, err := io.ReadAll(bakResp.Body)
		if err != nil {
			return nil, err
		}

		var bakResult struct {
			Data struct {
				List []struct {
					Bvid      string `json:"bvid"`
					Title     string `json:"title"`
					Desc      string `json:"desc"`
					Pic       string `json:"pic"`
					Author    string `json:"author"`
					VideoView int64  `json:"video_view"`
				} `json:"list"`
			} `json:"data"`
		}

		if err := json.Unmarshal(bakRespBody, &bakResult); err != nil {
			log.Printf("GetBilibili failed, err: %v", err)
			return nil, err
		}
		for _, v := range bakResult.Data.List {
			listData = append(listData, &model.HotItem{
				ID:        v.Bvid,
				Title:     v.Title,
				Desc:      v.Desc,
				Cover:     v.Pic,
				Author:    v.Author,
				Timestamp: time.Now().Unix(),
				Hot:       int(v.VideoView),
				URL:       fmt.Sprintf("https://www.bilibili.com/video/%s", v.Bvid),
				MobileURL: fmt.Sprintf("https://m.bilibili.com/video/%s", v.Bvid),
			})
		}

	} else {
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

						listData = append(listData, &model.HotItem{
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
	}

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "bilibili",
		Title:       "哔哩哔哩",
		Type:        "热榜 · 全站",
		Description: "你所热爱的，就是你的生活",
		Link:        "https://www.bilibili.com/v/popular/rank/all",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 百度 ==============
func (s *Spider) GetBaidu() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceBaidu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceBaidu].Agent)

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

	listData := []*model.HotItem{}
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

						listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "baidu",
		Title:      "百度",
		Type:       "热搜",
		Link:       "https://top.baidu.com/board",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== V2EX ==============
func (s *Spider) GetV2EX() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceV2EX].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceV2EX].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var list []map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "v2ex",
		Title:      "V2EX",
		Type:       "主题榜",
		Link:       "https://www.v2ex.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== GitHub ==============
func (s *Spider) GetGitHub() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceGitHub].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceGitHub].Agent)
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

	listData := []*model.HotItem{}
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

		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "github",
		Title:       "GitHub",
		Type:        "趋势榜",
		Description: "GitHub trending repositories",
		Link:        "https://github.com/trending",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 抖音 ==============
func (s *Spider) GetDouyin() (*model.HotData, error) {
	// 获取抖音Cookie
	cookieURL := "https://www.douyin.com/passport/general/login_guiding_strategy/?aid=6383"
	cookieReq, err := http.NewRequest("GET", cookieURL, nil)
	if err != nil {
		return nil, err
	}
	cookieReq.Header.Add("User-Agent", s.UrlMap[model.SourceDouyin].Agent)

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
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceDouyin].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceDouyin].Agent)
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

	listData := []*model.HotItem{}
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

					listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "douyin",
		Title:       "抖音",
		Type:        "热榜",
		Description: "实时上升热点",
		Link:        "https://www.douyin.com",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 快手 ==============
func (s *Spider) GetKuaishou() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceKuaishou].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceKuaishou].Agent)

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

	listData := []*model.HotItem{}
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

							listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "kuaishou",
		Title:       "快手",
		Type:        "热榜",
		Description: "快手，拥抱每一种生活",
		Link:        "https://www.kuaishou.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 今日头条 ==============
func (s *Spider) GetToutiao() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceToutiao].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceToutiao].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

				listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "toutiao",
		Title:      "今日头条",
		Type:       "热榜",
		Link:       "https://www.toutiao.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 掘金 ==============
func (s *Spider) GetJuejin() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceJuejin].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceJuejin].Agent)
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

	listData := []*model.HotItem{}
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

				listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "juejin",
		Title:      "稀土掘金",
		Type:       "文章榜",
		Link:       "https://juejin.cn/hot/articles",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
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
func (s *Spider) Get36Kr() (*model.HotData, error) {
	//todo: 支持多个分类
	// typeMap := map[string]string{
	// 	"hot":     "人气榜",
	// 	"video":   "视频榜",
	// 	"comment": "热议榜",
	// 	"collect": "收藏榜",
	// }

	var Body io.Reader

	type RequestBody struct {
		PartnerID string `json:"partner_id"`
		Param     struct {
			SiteID     int `json:"siteId"`
			PlatformID int `json:"platformId"`
		} `json:"param"`
		Timestamp int64 `json:"timestamp"`
	}

	reqBody, err := json.Marshal(RequestBody{
		PartnerID: "wap",
		Param: struct {
			SiteID     int `json:"siteId"`
			PlatformID int `json:"platformId"`
		}{
			SiteID:     1,
			PlatformID: 2,
		},
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	Body = strings.NewReader(string(reqBody))
	request, err := http.NewRequest("POST", s.UrlMap[model.Source36Kr].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.Source36Kr].Agent)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			HotRankList []struct {
				TemplateMaterial struct {
					ItemId      int64  `json:"itemId"`
					WidgetTitle string `json:"widgetTitle"`
					WidgetImage string `json:"widgetImage"`
					AuthorName  string `json:"authorName"`
					StatRead    int    `json:"statRead"`
					PublishTime int64  `json:"publishTime"`
				} `json:"templateMaterial"`
				Route string `json:"route"`
			} `json:"hotRankList"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	for _, item := range resp.Data.HotRankList {
		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(int(item.TemplateMaterial.ItemId)),
			Title:     item.TemplateMaterial.WidgetTitle,
			Cover:     item.TemplateMaterial.WidgetImage,
			Author:    item.TemplateMaterial.AuthorName,
			Timestamp: item.TemplateMaterial.PublishTime,
			Hot:       item.TemplateMaterial.StatRead,
			URL:       fmt.Sprintf("https://36kr.com/p/%d", item.TemplateMaterial.ItemId),
			MobileURL: fmt.Sprintf("https://m.36kr.com/p/%d", item.TemplateMaterial.ItemId),
		})
	}

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "36kr",
		Title:       "36氪",
		Type:        "热榜",
		Description: "让一部分人先看到未来",
		Link:        "https://36kr.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== CSDN ==============
func (s *Spider) GetCSDN() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceCSDN].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceCSDN].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

				listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "csdn",
		Title:       "CSDN",
		Type:        "热榜",
		Description: "专业开发者社区",
		Link:        "https://blog.csdn.net/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 百度贴吧 ==============
func (s *Spider) GetTieba() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceTieba].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceTieba].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

						listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "tieba",
		Title:       "百度贴吧",
		Type:        "热议榜",
		Description: "全球最大的中文社区",
		Link:        "https://tieba.baidu.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 知乎日报 ==============
func (s *Spider) GetZhihuDaily() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceZhihuDaily].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceZhihuDaily].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

				listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "zhihu-daily",
		Title:       "知乎日报",
		Type:        "日报",
		Description: "每天3次，每次7分钟",
		Link:        "https://daily.zhihu.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 酷安 ==============
func (s *Spider) GetCoolapk() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceCoolapk].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceCoolapk].Agent)
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

	listData := []*model.HotItem{}
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

				listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "coolapk",
		Title:       "酷安",
		Type:        "热门",
		Description: "发现科技新生活",
		Link:        "https://www.coolapk.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 虎扑 ==============
func (s *Spider) GetHupu() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceHupu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceHupu].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(i + 1),
			Title:     title,
			Author:    author,
			Timestamp: time.Now().Unix(),
			Hot:       replyNum,
			URL:       fmt.Sprintf("https://bbs.hupu.com%s", href),
			MobileURL: fmt.Sprintf("https://bbs.hupu.com%s", href),
		})
	})

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "hupu",
		Title:       "虎扑",
		Type:        "热帖",
		Description: "虎扑篮球社区",
		Link:        "https://bbs.hupu.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 虎嗅 ==============
func (s *Spider) GetHuxiu() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceHuxiu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceHuxiu].Agent)

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

	listData := []*model.HotItem{}
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

					listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "huxiu",
		Title:       "虎嗅",
		Type:        "24小时",
		Description: "虎嗅网 - 商业科技新媒体",
		Link:        "https://www.huxiu.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 简书 ==============
func (s *Spider) GetJianshu() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceJianshu].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceJianshu].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "jianshu",
		Title:       "简书",
		Type:        "热门",
		Description: "简书 - 创作你的创作",
		Link:        "https://www.jianshu.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 什么值得买 ==============
func (s *Spider) GetSmzdm() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceSmzdm].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceSmzdm].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(id),
			Title:     fmt.Sprintf("%s %s", title, price),
			Desc:      price,
			Timestamp: time.Now().Unix(),
			Hot:       worth,
			URL:       fmt.Sprintf("https://www.smzdm.com%s", href),
			MobileURL: fmt.Sprintf("https://m.smzdm.com%s", href),
		})
	})

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "smzdm",
		Title:       "什么值得买",
		Type:        "热门",
		Description: "什么值得买 - 消费决策平台",
		Link:        "https://www.smzdm.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 少数派 ==============
func (s *Spider) GetSspai() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceSspai].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceSspai].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Error int    `json:"error"`
		Msg   string `json:"msg"`
		Data  []struct {
			ID                int    `json:"id"`
			Title             string `json:"title"`
			Banner            string `json:"banner"`
			Summary           string `json:"summary"`
			CommentCount      int    `json:"comment_count"`
			LikeCount         int    `json:"like_count"`
			ViewCount         int    `json:"view_count"`
			Free              bool   `json:"free"`
			PostType          int    `json:"post_type"`
			Important         int    `json:"important"`
			ReleasedTime      int64  `json:"released_time"`
			MorningPaperTitle []any  `json:"morning_paper_title"`
			AdvertisementURL  string `json:"advertisement_url"`
			Series            []any  `json:"series"`
			Author            struct {
				ID       int    `json:"id"`
				Slug     string `json:"slug"`
				Avatar   string `json:"avatar"`
				Nickname string `json:"nickname"`
			} `json:"author"`
			Corner struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL   string `json:"url"`
				Icon  string `json:"icon"`
				Memo  string `json:"memo"`
				Color string `json:"color"`
			} `json:"corner"`
			SpecialColumns        []any  `json:"special_columns"`
			Status                int    `json:"status"`
			CreatedTime           int64  `json:"created_time"`
			ModifyTime            int64  `json:"modify_time"`
			IsMatrix              bool   `json:"is_matrix"`
			IsRecommendToHome     bool   `json:"is_recommend_to_home"`
			Slug                  string `json:"slug"`
			BelongToMember        bool   `json:"belong_to_member"`
			Issue                 string `json:"issue"`
			Tags                  []any  `json:"tags"`
			PodcastDuration       int    `json:"podcast_duration"`
			ArticleLimitFree      bool   `json:"article_limit_free"`
			ArticleLimitFreeStime int64  `json:"article_limit_free_stime"`
			ArticleLimitFreeEtime int64  `json:"article_limit_free_etime"`
			UserMemberCardShowOn  bool   `json:"user_member_card_show_on"`
			ObjectType            int    `json:"object_type"`
			IDHash                string `json:"id_hash"`
			RecommendToHomeAt     int64  `json:"recommend_to_home_at"`
			IsPreRecommendToHome  bool   `json:"is_pre_recommend_to_home"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	for _, item := range result.Data {
		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(item.ID),
			Title:     item.Title,
			Desc:      item.Summary,
			Author:    item.Author.Nickname,
			Timestamp: time.Unix(item.ReleasedTime, 0).Unix(),
			Hot:       0,
			URL:       fmt.Sprintf("https://sspai.com/post/%d", item.ID),
			MobileURL: fmt.Sprintf("https://m.sspai.com/post/%d", item.ID),
		})
	}

	return &model.HotData{
		Code:        http.StatusOK,
		Name:        "sspai",
		Title:       "少数派",
		Type:        "热门文章",
		Description: "少数派 - 高效工作，品质生活",
		Link:        "https://sspai.com/",
		Total:       len(listData),
		Data:        listData,
		UpdateTime:  time.Now(),
	}, nil
}

// ============== 网易新闻 ==============
func (s *Spider) GetNeteaseNews() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceNeteaseNews].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceNeteaseNews].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Data struct {
			List []struct {
				Title  string `json:"title"`
				Imgsrc string `json:"imgsrc"`
				Source string `json:"source"`
				Ptime  string `json:"ptime"`
				Docid  string `json:"docid"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	for _, item := range result.Data.List {
		listData = append(listData, &model.HotItem{
			ID:        item.Docid,
			Title:     item.Title,
			Cover:     item.Imgsrc,
			Author:    item.Source,
			Timestamp: parseTime(item.Ptime),
			Hot:       0,
			URL:       fmt.Sprintf("https://www.163.com/dy/article/%s.html", item.Docid),
			MobileURL: fmt.Sprintf("https://m.163.com/dy/article/%s.html", item.Docid),
		})

	}
	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "netease",
		Title:      "网易新闻",
		Type:       "热点榜",
		Link:       "https://m.163.com/hot",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 腾讯新闻 ==============
func (s *Spider) GetQQNews() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceQQNews].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceQQNews].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Idlist []struct {
			Newslist []struct {
				Title             string  `json:"title"`
				Abstract          string  `json:"abstract"`
				MiniProShareImage string  `json:"miniProShareImage"`
				Source            string  `json:"source"`
				Timestamp         float64 `json:"timestamp"`
				ID                string  `json:"id"`
				HotEvent          struct {
					HotScore int64 `json:"hotScore"`
				} `json:"hotEvent"`
			} `json:"newslist"`
		} `json:"idlist"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}

	for _, v := range result.Idlist {
		for _, vv := range v.Newslist {
			if strings.Contains(vv.Title, "腾讯新闻用户最关注的热点") {
				continue
			}
			listData = append(listData, &model.HotItem{
				ID:        vv.ID,
				Title:     vv.Title,
				Desc:      vv.Abstract,
				Cover:     vv.MiniProShareImage,
				Author:    vv.Source,
				Timestamp: int64(vv.Timestamp),
				Hot:       int(vv.HotEvent.HotScore),
				URL:       fmt.Sprintf("https://new.qq.com/rain/a/%s", vv.ID),
				MobileURL: fmt.Sprintf("https://view.inews.qq.com/k/%s", vv.ID),
			})
		}

	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "qq",
		Title:      "腾讯新闻",
		Type:       "热点榜",
		Link:       "https://news.qq.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 51CTO ==============
func (s *Spider) Get51CTO() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.Source51CTO].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.Source51CTO].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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

						listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "51cto",
		Title:      "51CTO",
		Type:       "推荐榜",
		Link:       "https://www.51cto.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 吾爱破解 ==============
func (s *Spider) Get52Pojie() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.Source52Pojie].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.Source52Pojie].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 需要处理RSS解析，这里简化实现
	// 实际实现需要解析RSS XML内容

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "52pojie",
		Title:      "吾爱破解",
		Type:       "最新精华",
		Link:       "https://www.52pojie.cn/",
		Total:      0,
		Data:       []*model.HotItem{},
		UpdateTime: time.Now(),
	}, nil
}

// ============== 豆瓣讨论 ==============
func (s *Spider) GetDoubanGroup() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceDoubanGroup].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceDoubanGroup].Agent)

	res, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 使用goquery解析HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	doc.Find(".article .channel-item").Each(func(i int, s *goquery.Selection) {
		url := s.Find("h3 a").AttrOr("href", "")
		title := s.Find("h3 a").Text()
		cover := s.Find(".pic-wrap img").AttrOr("src", "")
		desc := s.Find(".block p").Text()
		timeStr := s.Find("span.pubtime").Text()

		// 提取ID
		id := 0
		if url != "" {
			// 使用正则表达式提取topic后面的数字ID
			re := regexp.MustCompile(`/topic/(\d+)`)
			matches := re.FindStringSubmatch(url)
			if len(matches) > 1 {
				id, _ = strconv.Atoi(matches[1])
			}
		}

		// 解析时间
		timestamp := parseTime(timeStr)

		// 构建URL
		finalURL := url
		if finalURL == "" && id > 0 {
			finalURL = fmt.Sprintf("https://www.douban.com/group/topic/%d", id)
		}

		mobileURL := ""
		if id > 0 {
			mobileURL = fmt.Sprintf("https://m.douban.com/group/topic/%d/", id)
		}

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(id),
			Title:     strings.TrimSpace(title),
			Desc:      strings.TrimSpace(desc),
			Cover:     cover,
			Timestamp: timestamp,
			Hot:       0,
			URL:       finalURL,
			MobileURL: mobileURL,
		})
	})

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "douban-group",
		Title:      "豆瓣讨论",
		Type:       "讨论精选",
		Link:       "https://www.douban.com/group/explore",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== AcFun ==============
func (s *Spider) GetAcfun() (*model.HotData, error) {
	var Body io.Reader
	request, err := http.NewRequest("GET", s.UrlMap[model.SourceAcfun].Url, Body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", s.UrlMap[model.SourceAcfun].Agent)
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

	listData := []*model.HotItem{}
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

					listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "acfun",
		Title:      "AcFun",
		Type:       "排行榜",
		Link:       "https://www.acfun.cn/rank/list/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 数字尾巴 ==============
func (s *Spider) GetDgtle() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceDgtle]
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

	listData := []*model.HotItem{}
	for _, item := range result.Data.Items {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", item.CreatedAt)

		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "dgtle",
		Title:      "数字尾巴",
		Type:       "热门文章",
		Link:       "https://www.dgtle.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 豆瓣电影 ==============
func (s *Spider) GetDoubanMovie() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceDoubanMovie]
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

	listData := []*model.HotItem{}
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

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(id),
			Title:     fmt.Sprintf("【%s】%s", score, title),
			Cover:     cover,
			URL:       url,
			MobileURL: url,
		})
	})

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "douban-movie",
		Title:      "豆瓣电影",
		Type:       "新片榜",
		Link:       "https://movie.douban.com/chart",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 中国地震台 ==============
func (s *Spider) GetEarthquake() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceEarthquake]
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
		NEW_DID    string `json:"NEW_DID"`
		LOCATION_C string `json:"LOCATION_C"`
		M          string `json:"M"`
		O_TIME     string `json:"O_TIME"`
		EPI_LAT    string `json:"EPI_LAT"`
		EPI_LON    string `json:"EPI_LON"`
		EPI_DEPTH  int64  `json:"EPI_DEPTH"`
		SAVE_TIME  string `json:"SAVE_TIME"`
	}

	if err := json.Unmarshal(match[1], &earthquakes); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
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
				value = eq.M
				if key == "EPI_LAT" {
					value = eq.EPI_LAT
				} else if key == "EPI_LON" {
					value = eq.EPI_LON
				} else if key == "EPI_DEPTH" {
					value = strconv.FormatInt(eq.EPI_DEPTH, 10)
				}
			}
			contentBuilder = append(contentBuilder, fmt.Sprintf("%s：%s", desc, value))
		}

		timestamp, _ := time.Parse("2006-01-02 15:04:05", eq.O_TIME)

		listData = append(listData, &model.HotItem{
			ID:        "0",
			Title:     fmt.Sprintf("%s发生%s级地震", eq.LOCATION_C, eq.M),
			Desc:      strings.Join(contentBuilder, "\n"),
			Timestamp: timestamp.Unix(),
			URL:       "https://news.ceic.ac.cn/",
			MobileURL: "https://news.ceic.ac.cn/",
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "earthquake",
		Title:      "中国地震台",
		Type:       "地震速报",
		Link:       "https://news.ceic.ac.cn/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== GameRes游资网 ==============
func (s *Spider) GetGameres() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceGameres]
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

	listData := []*model.HotItem{}
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

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(id),
			Title:     title,
			Cover:     cover,
			Author:    author,
			Timestamp: timestamp,
			URL:       url,
			MobileURL: url,
		})
	})

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "gameres",
		Title:      "GameRes游资网",
		Type:       "资讯",
		Link:       "https://www.gameres.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 极客公园 ==============
func (s *Spider) GetGeekpark() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceGeekpark]
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

	listData := []*model.HotItem{}
	for _, item := range result.Data {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", item.CreatedAt)

		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "geekpark",
		Title:      "极客公园",
		Type:       "热门文章",
		Link:       "https://www.geekpark.net/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 原神 ==============
func (s *Spider) GetGenshin() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceGenshin]
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

	listData := make([]*model.HotItem, 0, len(result.Data.List))
	for _, item := range result.Data.List {
		post := item.Post

		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "genshin",
		Title:      "原神",
		Type:       "动态",
		Link:       "https://bbs.mihoyo.com/ys/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// GetGuokr 获取果壳热门文章
func (s *Spider) GetGuokr() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceGuokr]
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

	var result []struct {
		ID         int    `json:"id"`
		Title      string `json:"title"`
		Summary    string `json:"summary"`
		SmallImage string `json:"small_image"`
		Authors    []struct {
			Nickname string `json:"nickname"`
		} `json:"authors"`
		DateCreated string `json:"date_created"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	for _, item := range result {
		timestamp, _ := time.Parse(time.RFC3339, item.DateCreated)
		listData = append(listData, &model.HotItem{
			ID:    strconv.Itoa(item.ID),
			Title: item.Title,
			Desc:  item.Summary,
			Cover: item.SmallImage,
			Author: func() string {
				if len(item.Authors) > 0 {
					return item.Authors[0].Nickname
				}
				return ""
			}(),
			Timestamp: timestamp.Unix(),
			URL:       fmt.Sprintf("https://www.guokr.com/article/%d", item.ID),
			MobileURL: fmt.Sprintf("https://m.guokr.com/article/%d", item.ID),
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "guokr",
		Title:      "果壳",
		Type:       "热门文章",
		Link:       "https://www.guokr.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// GetHackernews 获取Hacker News热门文章
func (s *Spider) GetHackernews() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceHackernews]
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

	listData := []*model.HotItem{}
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

			listData = append(listData, &model.HotItem{
				ID:    strconv.Itoa(idInt),
				Title: title,
				Hot:   hot,
				URL:   url,
			})
		}
	})

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "hackernews",
		Title:      "Hacker News",
		Type:       "Popular",
		Link:       "https://news.ycombinator.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// GetHelloGitHub 获取HelloGitHub热门仓库
func (s *Spider) GetHelloGitHub() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceHelloGitHub]
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

	listData := []*model.HotItem{}
	for _, item := range result.Data {
		timestamp, _ := time.Parse("2006-01-02T15:04:05", item.UpdatedAt)
		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "hellogithub",
		Title:      "HelloGitHub",
		Type:       "热门仓库",
		Link:       "https://hellogithub.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// GetHistory 获取历史上的今天
func (s *Spider) GetHistory() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceHistory]
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

	listData := []*model.HotItem{}
	for i, item := range dayData {
		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "history",
		Title:      "历史上的今天",
		Type:       fmt.Sprintf("%02d-%02d", month, day),
		Link:       "https://baike.baidu.com/calendar",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// GetHonkai 获取崩坏3最新动态
func (s *Spider) GetHonkai() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceHonkai]
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

	listData := []*model.HotItem{}
	for _, item := range result.Data.List {
		post := item.Post

		// 获取封面图，优先使用cover，如果没有则使用第一张图片
		cover := post.Cover
		if cover == "" && len(post.Images) > 0 {
			cover = post.Images[0]
		}

		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "honkai",
		Title:      "崩坏3",
		Type:       "最新动态",
		Link:       "https://www.miyoushe.com/bh3/home/6",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// GetHostloc 获取hostloc论坛热门帖子
func (s *Spider) GetHostloc() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceHostloc]
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

	listData := []*model.HotItem{}
	for _, item := range feed.Channel.Items {
		// 将字符串Guid转换为int
		idInt, _ := strconv.Atoi(item.Guid)
		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     item.Title,
			Desc:      item.Description,
			Author:    item.Author,
			Timestamp: parseTime(item.PubDate),
			URL:       item.Link,
			MobileURL: item.Link,
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "hostloc",
		Title:      "全球主机交流论坛",
		Type:       "热门帖子",
		Link:       "https://hostloc.com/forum.php?mod=guide&view=hot",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// GetIfanr 获取爱范儿快讯数据
func (s *Spider) GetIfanr() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceIfanr]
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

	listData := []*model.HotItem{}
	for _, item := range result.Data {
		idInt, _ := strconv.Atoi(item.ID)
		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     item.Title,
			Desc:      item.Content,
			Timestamp: item.Time,
			URL:       fmt.Sprintf("https://www.ifanr.com/buzz/%s", item.ID),
			MobileURL: fmt.Sprintf("https://www.ifanr.com/buzz/%s", item.ID),
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "ifanr",
		Title:      "爱范儿",
		Type:       "快讯",
		Link:       "https://www.ifanr.com/buzz",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// GetIthomeXijiayi 获取IT之家喜加一游戏动态
func (s *Spider) GetIthomeXijiayi() (*model.HotData, error) {
	urlConf := s.UrlMap[model.SourceIthomeXijiayi]
	resp, err := s.HttpClient.Get(urlConf.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	counter := 1
	doc.Find(".newslist li").Each(func(i int, s *goquery.Selection) {
		href := s.Find("a").AttrOr("href", "")
		timeStr := s.Find("span.time").Text()
		title := s.Find(".newsbody h2").Text()
		desc := s.Find(".newsbody p").Text()
		cover := s.Find("img").AttrOr("data-original", "")
		hotStr := s.Find(".comment").Text()

		// 提取时间戳
		timestamp := parseTime(timeStr)

		// 提取热度值
		hot := 0
		if hotStr != "" {
			hot, _ = strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(hotStr, `\D`, "")))
		}

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(counter),
			Title:     strings.TrimSpace(title),
			Desc:      strings.TrimSpace(desc),
			Cover:     cover,
			Timestamp: timestamp,
			Hot:       hot,
			URL:       href,
			MobileURL: href,
		})
		counter++
	})

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "ithome-xijiayi",
		Title:      "IT之家喜加一",
		Type:       "游戏动态",
		Link:       "https://www.ithome.com/zt/xijiayi",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 米游社 ==============
func (s *Spider) GetMiyoushe() (*model.HotData, error) {
	url := s.UrlMap[model.SourceMiyoushe].Url
	agent := s.UrlMap[model.SourceMiyoushe].Agent

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

	listData := []*model.HotItem{}
	for _, item := range result.Data.List {
		post := item.Post
		cover := post.Cover
		if cover == "" && len(post.Images) > 0 {
			cover = post.Images[0]
		}

		idInt, _ := strconv.Atoi(post.PostID)
		listData = append(listData, &model.HotItem{
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

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "miyoushe",
		Title:      "米游社",
		Type:       "最新公告",
		Link:       "https://www.miyoushe.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 水木社区 ==============
func (s *Spider) GetNewsmth() (*model.HotData, error) {
	url := s.UrlMap[model.SourceNewsmth].Url
	agent := s.UrlMap[model.SourceNewsmth].Agent

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

	listData := []*model.HotItem{}
	for _, topic := range result.Data.Topics {
		post := topic.Article
		url := fmt.Sprintf("https://wap.newsmth.net/article/%s?title=%s&from=home", post.TopicID, topic.Board.Title)

		idInt, _ := strconv.Atoi(topic.FirstArticleID)
		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(idInt),
			Title:     strings.TrimSpace(post.Subject),
			Desc:      strings.TrimSpace(post.Body),
			Author:    post.Account.Name,
			Timestamp: post.PostTime,
			URL:       url,
			MobileURL: url,
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "newsmth",
		Title:      "水木社区",
		Type:       "热门话题",
		Link:       "https://www.newsmth.net/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== NGA ==============
func (s *Spider) GetNgabbs() (*model.HotData, error) {
	url := s.UrlMap[model.SourceNgabbs].Url
	// agent := s.UrlMap[model.SourceNgabbs].Agent

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

	listData := []*model.HotItem{}
	for _, item := range result.Result[0] {
		tidInt, _ := strconv.Atoi(item.Tid)
		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(tidInt),
			Title:     strings.TrimSpace(item.Subject),
			Author:    item.Author,
			Hot:       item.Replies,
			Timestamp: item.Postdate,
			URL:       fmt.Sprintf("https://bbs.nga.cn%s", item.Tpcurl),
			MobileURL: fmt.Sprintf("https://bbs.nga.cn%s", item.Tpcurl),
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "ngabbs",
		Title:      "NGA",
		Type:       "论坛热帖",
		Link:       "https://ngabbs.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== NodeSeek ==============
func (s *Spider) GetNodeseek() (*model.HotData, error) {
	url := s.UrlMap[model.SourceNodeseek].Url
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

	listData := []*model.HotItem{}
	for _, item := range rssData.Channel.Item {
		timestamp, _ := time.Parse(time.RFC1123, item.PubDate)

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(len(listData) + 1),
			Title:     item.Title,
			Desc:      item.Description,
			Timestamp: timestamp.Unix(),
			URL:       item.Link,
			MobileURL: item.Link,
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "nodeseek",
		Title:      "NodeSeek",
		Type:       "技术社区",
		Link:       "https://www.nodeseek.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 纽约时报 ==============
func (s *Spider) GetNytimes() (*model.HotData, error) {
	url := s.UrlMap[model.SourceNytimes].Url
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

	listData := []*model.HotItem{}
	for _, item := range rssData.Channel.Item {
		timestamp, _ := time.Parse(time.RFC1123, item.PubDate)

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(len(listData) + 1),
			Title:     item.Title,
			Desc:      item.Description,
			Timestamp: timestamp.Unix(),
			URL:       item.Link,
			MobileURL: item.Link,
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "nytimes",
		Title:      "纽约时报",
		Type:       "国际新闻",
		Link:       "https://www.nytimes.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== Product Hunt ==============
func (s *Spider) GetProducthunt() (*model.HotData, error) {
	url := s.UrlMap[model.SourceProducthunt].Url
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

	listData := []*model.HotItem{}
	for i, match := range matches {
		if len(match) >= 4 {
			url := "https://www.producthunt.com" + match[1]

			listData = append(listData, &model.HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Desc:      strings.TrimSpace(match[3]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "producthunt",
		Title:      "Product Hunt",
		Type:       "产品发现",
		Link:       "https://www.producthunt.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 新浪新闻 ==============
func (s *Spider) GetSinaNews() (*model.HotData, error) {
	url := s.UrlMap[model.SourceSinaNews].Url
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

	listData := []*model.HotItem{}
	for i, match := range matches {
		if len(match) >= 3 && strings.Contains(match[1], "news.sina.com.cn") {
			url := match[1]
			if !strings.HasPrefix(url, "http") {
				url = "https:" + url
			}

			listData = append(listData, &model.HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "sina-news",
		Title:      "新浪新闻",
		Type:       "新闻资讯",
		Link:       "https://news.sina.com.cn/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 新浪微博 ==============
func (s *Spider) GetSina() (*model.HotData, error) {
	url := s.UrlMap[model.SourceSina].Url
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

	listData := []*model.HotItem{}
	for i, match := range matches {
		if len(match) >= 3 {
			url := match[1]
			if !strings.HasPrefix(url, "http") {
				url = "https:" + url
			}

			listData = append(listData, &model.HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "sina",
		Title:      "微博热搜",
		Type:       "社交媒体",
		Link:       "https://s.weibo.com/top/summary",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 星穹铁道 ==============
func (s *Spider) GetStarrail() (*model.HotData, error) {
	url := s.UrlMap[model.SourceStarrail].Url
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

	listData := make([]*model.HotItem, 0, len(jsonData.Data.List))
	for _, item := range jsonData.Data.List {
		post := item.Post
		url := fmt.Sprintf("https://bbs.miyoushe.com/detail/%d", post.PostId)

		listData = append(listData, &model.HotItem{
			ID:        strconv.Itoa(int(post.PostId)),
			Title:     post.Subject,
			Desc:      post.Content,
			Timestamp: post.CreatedAt,
			Hot:       post.ViewCnt,
			URL:       url,
			MobileURL: url,
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "starrail",
		Title:      "星穹铁道",
		Type:       "游戏资讯",
		Link:       "https://bbs.miyoushe.com/ys/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 澎湃新闻 ==============
func (s *Spider) GetThepaper() (*model.HotData, error) {
	url := s.UrlMap[model.SourceThepaper].Url
	resp, err := s.HttpClient.Get(url)
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
			HotNews []struct {
				Name        string `json:"name"`
				Pic         string `json:"pic"`
				ContId      string `json:"contId"`
				PariseTimes string `json:"praiseTimes"`
				PubTimeLong int64  `json:"pubTimeLong"`
				NodeInfo    struct {
					Desc string `json:"desc"`
				} `json:"nodeInfo"`
			} `json:"hotNews"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	for _, item := range result.Data.HotNews {
		hot, _ := strconv.Atoi(item.PariseTimes)
		listData = append(listData, &model.HotItem{
			ID:        item.ContId,
			Title:     item.Name,
			Hot:       hot,
			Cover:     item.Pic,
			Desc:      item.NodeInfo.Desc,
			Timestamp: time.Unix(item.PubTimeLong, 0).Unix(),
			URL:       fmt.Sprintf("https://www.thepaper.cn/newsDetail_forward_%s", item.ContId),
			MobileURL: fmt.Sprintf("https://m.thepaper.cn/newsDetail_forward_%s", item.ContId),
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "thepaper",
		Title:      "澎湃新闻",
		Type:       "新闻资讯",
		Link:       "https://www.thepaper.cn/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 气象预警 ==============
func (s *Spider) GetWeatheralarm() (*model.HotData, error) {
	url := s.UrlMap[model.SourceWeatheralarm].Url
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

	listData := []*model.HotItem{}
	for _, item := range jsonData.Data.Page.List {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", item.Issuetime)
		url := fmt.Sprintf("http://www.nmc.cn/publish/alarm.html?alertid=%s", item.Alertid)

		listData = append(listData, &model.HotItem{
			ID:        item.Alertid,
			Title:     item.Title,
			Desc:      fmt.Sprintf("发布时间: %s", item.Issuetime),
			Cover:     item.Pic,
			Timestamp: timestamp.Unix(),
			URL:       url,
			MobileURL: url,
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "weatheralarm",
		Title:      "气象预警",
		Type:       "气象信息",
		Link:       "http://www.nmc.cn/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 微信读书 ==============
func (s *Spider) GetWeread() (*model.HotData, error) {
	url := s.UrlMap[model.SourceWeread].Url
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
			ReadingCount int `json:"readingCount"`
			BookInfo     struct {
				BookId      string `json:"bookId"`
				Title       string `json:"title"`
				Author      string `json:"author"`
				Intro       string `json:"intro"`
				Cover       string `json:"cover"`
				PublishTime string `json:"publishTime"`
			} `json:"bookInfo"`
		} `json:"books"`
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	listData := []*model.HotItem{}
	for _, book := range jsonData.Books {
		cover := strings.Replace(book.BookInfo.Cover, "_s.jpg", "_l.jpg", 1)
		realBookId, err := token.GetWereadID(book.BookInfo.BookId)
		if err != nil {
			log.Printf("GetWereadID failed, err: %v", err)
			continue
		}

		listData = append(listData, &model.HotItem{
			ID:        book.BookInfo.BookId,
			Title:     book.BookInfo.Title,
			Author:    book.BookInfo.Author,
			Desc:      book.BookInfo.Intro,
			Cover:     cover,
			Timestamp: parseTime(book.BookInfo.PublishTime),
			Hot:       book.ReadingCount,
			URL:       fmt.Sprintf("https://weread.qq.com/web/bookDetail/%s", realBookId),
			MobileURL: fmt.Sprintf("https://m.weread.qq.com/web/bookDetail/%s", realBookId),
		})
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "weread",
		Title:      "微信读书",
		Type:       "图书排行",
		Link:       "https://weread.qq.com/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
	}, nil
}

// ============== 游研社 ==============
func (s *Spider) GetYystv() (*model.HotData, error) {
	url := s.UrlMap[model.SourceYystv].Url
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

	listData := []*model.HotItem{}
	for i, match := range matches {
		if len(match) >= 3 {
			url := match[1]
			if !strings.HasPrefix(url, "http") {
				url = "https:" + url
			}

			listData = append(listData, &model.HotItem{
				ID:        strconv.Itoa(i + 1),
				Title:     strings.TrimSpace(match[2]),
				Timestamp: time.Now().Unix(),
				URL:       url,
				MobileURL: url,
			})
		}
	}

	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "yystv",
		Title:      "游研社",
		Type:       "游戏资讯",
		Link:       "https://www.yystv.cn/",
		Total:      len(listData),
		Data:       listData,
		UpdateTime: time.Now(),
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

func (s *Spider) GetLol() (*model.HotData, error) {
	return &model.HotData{
		Code:       http.StatusOK,
		Name:       "lol",
		Title:      "lol",
		Type:       "lol",
		Link:       "https://lol.qq.com/",
		Total:      0,
		Data:       model.HotItems{},
		UpdateTime: time.Now(),
	}, nil
}
