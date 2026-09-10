package api

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"fastduck/treasure-doc/module/blog/data/response"

	"github.com/gin-gonic/gin"
)

const feedLimit = 50

// Robots 输出 robots.txt 并指向 sitemap。
func (h *Handler) Robots(c *gin.Context) {
	c.String(http.StatusOK, "User-agent: *\nAllow: /\nSitemap: %s/sitemap.xml\n", baseURL(c))
}

type sitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type sitemapSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// Sitemap 输出 sitemap.xml。
func (h *Handler) Sitemap(c *gin.Context) {
	entries, err := h.service.Sitemap(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "服务内部错误")
		return
	}
	base := baseURL(c)
	set := sitemapSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", URLs: make([]sitemapURL, 0, len(entries))}
	for _, entry := range entries {
		item := sitemapURL{Loc: base + entry.Loc}
		if !entry.LastMod.IsZero() {
			item.LastMod = entry.LastMod.Format("2006-01-02")
		}
		set.URLs = append(set.URLs, item)
	}
	writeXML(c, set)
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate,omitempty"`
	Description string `xml:"description"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

// RSS 输出文章订阅源。
func (h *Handler) RSS(c *gin.Context) {
	ctx := c.Request.Context()
	site, err := h.service.Site(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "服务内部错误")
		return
	}
	posts, err := h.service.Feed(ctx, feedLimit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "服务内部错误")
		return
	}

	base := baseURL(c)
	channel := rssChannel{Title: site.Name, Link: base, Description: site.Intro, Items: make([]rssItem, 0, len(posts))}
	for _, post := range posts {
		link := base + "/Blog/" + post.ID
		channel.Items = append(channel.Items, rssItem{
			Title:       post.Title,
			Link:        link,
			GUID:        link,
			PubDate:     rssPubDate(post.Date),
			Description: post.Summary,
		})
	}
	writeXML(c, rssFeed{Version: "2.0", Channel: channel})
}

func rssPubDate(date string) string {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}
	return parsed.Format(time.RFC1123Z)
}

// baseURL 依据当前请求推断站点根地址，便于同一份代码部署到不同域名。
func baseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

func writeXML(c *gin.Context, value interface{}) {
	data, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "服务内部错误")
		return
	}
	c.Data(http.StatusOK, "application/xml; charset=utf-8", append([]byte(xml.Header), data...))
}
