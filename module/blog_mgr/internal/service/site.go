package service

import (
	"context"

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

// defaultSite 返回尚未保存站点配置时使用的默认对象，数组字段为空数组而不是 null。
func defaultSite() blogresponse.Site {
	return blogresponse.Site{TechStack: []string{}, Modules: defaultSiteModules(), Milestones: []blogresponse.SiteMilestone{}, Home: blogresponse.DefaultSiteHome(), Footer: blogresponse.DefaultSiteFooter(), Banner: blogresponse.DefaultSiteBanner()}
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
