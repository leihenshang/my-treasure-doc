package response

import "time"

// SitemapEntry 是 sitemap.xml 中的一条 URL。
type SitemapEntry struct {
	// Loc 站点内相对路径，例如 /Blog/hello
	Loc string
	// LastMod 最后修改时间，零值表示不输出 lastmod
	LastMod time.Time
}
