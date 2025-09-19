package token

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fastduck/treasure-doc/module/hot_top/model"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type EncodedKeys struct {
	ImgKey string `json:"img_key"`
	SubKey string `json:"sub_key"`
}

type NavResponse struct {
	Data struct {
		WbiImg struct {
			ImgUrl string `json:"img_url"`
			SubUrl string `json:"sub_url"`
		} `json:"wbi_img"`
	} `json:"data"`
}

func GetWbiKeys() (EncodedKeys, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/web-interface/nav", nil)
	if err != nil {
		return EncodedKeys{}, err
	}

	req.Header.Set("Cookie", "SESSDATA=xxxxxx")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := client.Do(req)
	if err != nil {
		return EncodedKeys{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return EncodedKeys{}, err
	}

	var result NavResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return EncodedKeys{}, err
	}

	imgUrl := result.Data.WbiImg.ImgUrl
	subUrl := result.Data.WbiImg.SubUrl

	return EncodedKeys{
		ImgKey: extractKeyFromUrl(imgUrl),
		SubKey: extractKeyFromUrl(subUrl),
	}, nil
}

func extractKeyFromUrl(url string) string {
	if url == "" {
		return ""
	}
	lastSlash := strings.LastIndex(url, "/")
	lastDot := strings.LastIndex(url, ".")

	if lastSlash == -1 || lastDot == -1 || lastSlash >= lastDot {
		return ""
	}

	return url[lastSlash+1 : lastDot]
}

var mixinKeyEncTab = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49, 33, 9, 42, 19, 29, 28,
	14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40, 61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54,
	21, 56, 59, 6, 63, 57, 62, 11, 36, 20, 34, 44, 52,
}

// WbiParams 定义请求参数类型
type WbiParams map[string]interface{}

// GetMixinKey 对 imgKey 和 subKey 进行字符顺序打乱编码
func GetMixinKey(orig string) string {
	var result []byte
	for _, n := range mixinKeyEncTab {
		if n < len(orig) {
			result = append(result, orig[n])
		}
	}
	return string(result)[:32]
}

// EncWbi 为请求参数进行 wbi 签名
func EncWbi(params WbiParams, imgKey string, subKey string) string {
	mixinKey := GetMixinKey(imgKey + subKey)
	currTime := time.Now().Unix()

	// 添加 wts 字段
	params["wts"] = currTime

	// 按照 key 重排参数
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var queryParts []string
	chrFilter := regexp.MustCompile(`[!'()*]`)

	for _, key := range keys {
		filteredValue := ""
		if value, ok := params[key].(string); ok {
			// 过滤 value 中的 "!'\(\)*" 字符
			filteredValue = chrFilter.ReplaceAllString(value, "")
		} else if value, ok := params[key].(int64); ok {
			filteredValue = strconv.FormatInt(value, 10)
		} else if value, ok := params[key].(int); ok {
			filteredValue = strconv.FormatInt(int64(value), 10)
		}
		queryParts = append(queryParts, key+"="+filteredValue)
	}

	query := strings.Join(queryParts, "&")

	// 计算 w_rid
	hash := md5.Sum([]byte(query + mixinKey))
	wbiSign := hex.EncodeToString(hash[:])

	return query + "&w_rid=" + wbiSign
}

func GetBilibiliWbi() (string, error) {
	keys, err := GetWbiKeys()
	if err != nil {
		return "", err
	}
	if token, ok := GetTokenCache().GetToken(model.SourceBilibili.String()); ok {
		return token, nil
	}
	token := EncWbi(WbiParams{
		"foo": "114", "bar": "514", "baz": 1919810,
	}, keys.ImgKey, keys.SubKey)
	GetTokenCache().SetToken(model.SourceBilibili.String(), token)
	return token, nil
}

type tokenCache struct {
	lock  *sync.Mutex
	cache map[string]any
}

var tokenCacheInstance = &tokenCache{
	lock:  &sync.Mutex{},
	cache: make(map[string]any),
}

func GetTokenCache() *tokenCache {
	return tokenCacheInstance
}

func (c *tokenCache) GetToken(key string) (string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if token, ok := c.cache[key]; ok {
		return token.(string), true
	}
	return "", false
}

func (c *tokenCache) SetToken(key string, token string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache[key] = token
}
