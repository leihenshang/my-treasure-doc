package service

import (
	blogresponse "fastduck/treasure-doc/module/blog/data/response"
)

// defaultSiteModules 返回六个固定站点模块的默认配置。
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
	return blogresponse.Site{TechStack: []string{}, Modules: defaultSiteModules(), Milestones: []blogresponse.SiteMilestone{}}
}
