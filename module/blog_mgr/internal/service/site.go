package service

import (
	"context"
	"log"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	blogresponse "fastduck/treasure-doc/module/blog/data/response"
	"fastduck/treasure-doc/module/blog_mgr/data/response"
)

// defaultSiteModules 返回固定站点模块的默认配置。
func defaultSiteModules() []blogresponse.SiteModule {
	return blogresponse.DefaultSiteModules()
}

// normalizeSiteModules 复用共享的站点模块约束，并把错误映射为本服务的 ErrInvalid。
// strict 为 true 时校验管理端写入载荷，为 false 时按默认配置补齐历史数据。
func normalizeSiteModules(modules []blogresponse.SiteModule, strict bool) ([]blogresponse.SiteModule, error) {
	result, err := blogresponse.NormalizeSiteModules(modules, strict)
	if err != nil {
		return nil, ErrInvalid
	}
	return result, nil
}

// 站点名与站长名是必填字段（validateSite / validateProfile 会拒绝空值）。
// 默认对象带上与 seed 数据一致的名称，保证全新库上「GET → 改一项 → PUT 回去」这条链路可用。
const (
	defaultSiteName    = "Treasure Blog"
	defaultProfileName = "Treasure"
)

// defaultSite 返回尚未保存站点配置时使用的默认对象，数组字段为空数组而不是 null。
func defaultSite() blogresponse.Site {
	return blogresponse.Site{Name: defaultSiteName, TechStack: []string{}, Modules: defaultSiteModules(), Milestones: []blogresponse.SiteMilestone{}, Home: blogresponse.DefaultSiteHome(), Footer: blogresponse.DefaultSiteFooter(), Banner: blogresponse.DefaultSiteBanner()}
}

// Stats 汇总后台仪表盘需要的资源数量与总浏览量。
func (s *Service) Stats(ctx context.Context) (response.Stats, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.Stats{}, err
	}

	var result response.Stats
	targets := []struct {
		model interface{}
		out   *response.StatusCounts
	}{
		{model: &blogmodel.Post{}, out: &result.Posts},
		{model: &blogmodel.Diary{}, out: &result.Diaries},
		{model: &blogmodel.PortfolioItem{}, out: &result.Portfolio},
		{model: &blogmodel.Tool{}, out: &result.Tools},
		{model: &blogmodel.Bookmark{}, out: &result.Bookmarks},
	}
	statuses := []struct {
		status string
		out    func(*response.StatusCounts) *int64
	}{
		{status: blogmodel.StatusDraft, out: func(c *response.StatusCounts) *int64 { return &c.Draft }},
		{status: blogmodel.StatusPublished, out: func(c *response.StatusCounts) *int64 { return &c.Published }},
		{status: blogmodel.StatusArchived, out: func(c *response.StatusCounts) *int64 { return &c.Archived }},
	}

	for _, target := range targets {
		if err := db.Model(target.model).Count(&target.out.Total).Error; err != nil {
			return response.Stats{}, err
		}
		for _, status := range statuses {
			if err := db.Model(target.model).Where("publish_status = ?", status.status).Count(status.out(target.out)).Error; err != nil {
				return response.Stats{}, err
			}
		}
	}

	if err := db.Model(&blogmodel.Category{}).Count(&result.Categories).Error; err != nil {
		return response.Stats{}, err
	}
	if err := db.Model(&blogmodel.Tag{}).Count(&result.Tags).Error; err != nil {
		return response.Stats{}, err
	}

	for _, model := range []interface{}{&blogmodel.Post{}, &blogmodel.Diary{}, &blogmodel.PortfolioItem{}} {
		var views int64
		if err := db.Model(model).Select("COALESCE(SUM(view_count), 0)").Scan(&views).Error; err != nil {
			return response.Stats{}, err
		}
		result.TotalViews += views
	}
	return result, nil
}

// VisitorStats 统计最近 days 天内访问公开博客的 IP：总量、去重 IP 数与明细。
// days 由调用方限制在 [1, 365]。
func (s *Service) VisitorStats(ctx context.Context, days int) (response.VisitorStats, error) {
	db, err := s.database(ctx)
	if err != nil {
		return response.VisitorStats{}, err
	}

	since := time.Now().AddDate(0, 0, -days)
	var result response.VisitorStats
	result.Days = int64(days)

	if err := db.Model(&blogmodel.VisitorLog{}).Where("created_at >= ?", since).Count(&result.TotalCount).Error; err != nil {
		return response.VisitorStats{}, err
	}
	if err := db.Model(&blogmodel.VisitorLog{}).Where("created_at >= ?", since).Distinct("ip").Count(&result.DistinctIP).Error; err != nil {
		return response.VisitorStats{}, err
	}

	var list []response.VisitorStat
	type visitorRow struct {
		IP        string
		Count     int64
		FirstSeen string
		LastSeen  string
	}
	// SQLite 下 MIN/MAX(created_at) 聚合返回字符串，不能用 time.Time 直接 Scan，
	// 先收进字符串再按常见时间格式解析。
	var rows []visitorRow
	if err := db.Model(&blogmodel.VisitorLog{}).
		Where("created_at >= ?", since).
		Select("ip, COUNT(*) AS count, MIN(created_at) AS first_seen, MAX(created_at) AS last_seen").
		Group("ip").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		log.Printf("[visitor-stats] 聚合查询失败: %v", err)
		return response.VisitorStats{}, err
	}
	list = make([]response.VisitorStat, 0, len(rows))
	// 附加每个 IP 访问过的去重分类：SQLite/PG 的 group_concat 方言不同，改用一次 DISTINCT 查询在 Go 里聚合。
	type ipCategory struct {
		IP       string
		Category string
	}
	var ipCats []ipCategory
	if err := db.Model(&blogmodel.VisitorLog{}).
		Where("created_at >= ?", since).
		Distinct("ip", "category").
		Scan(&ipCats).Error; err != nil {
		log.Printf("[visitor-stats] 分类查询失败: %v", err)
		return response.VisitorStats{}, err
	}
	catsByIP := make(map[string][]string, len(ipCats))
	for _, row := range ipCats {
		if row.Category == "" {
			continue
		}
		catsByIP[row.IP] = append(catsByIP[row.IP], row.Category)
	}
	for _, row := range rows {
		entry := response.VisitorStat{
			IP:         row.IP,
			Count:      row.Count,
			FirstSeen:  parseVisitorTime(row.FirstSeen),
			LastSeen:   parseVisitorTime(row.LastSeen),
			Categories: catsByIP[row.IP],
		}
		if entry.Categories == nil {
			entry.Categories = []string{}
		}
		list = append(list, entry)
	}
	result.List = list
	return result, nil
}

// parseVisitorTime 解析访客日志里的时间字符串，兼容 SQLite 存储与 RFC3339 两种形态。
func parseVisitorTime(raw string) time.Time {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02 15:04:05.999999999-07:00", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t
		}
	}
	return time.Time{}
}
