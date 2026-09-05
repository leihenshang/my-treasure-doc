package response

import (
	"encoding/json"
	"errors"
	"strings"
)

// ErrInvalidSiteModule 表示站点模块集合不满足固定模块约束。
var ErrInvalidSiteModule = errors.New("invalid site module")

// fixedSiteModules 是站点模块的固定集合。ID 与 Path 属于对外契约，
// 管理端只允许修改 icon、name、desc 和 visible。
var fixedSiteModules = []SiteModule{
	{ID: "blog", Icon: "📝", Name: "文章", Desc: "技术笔记与长文", Path: "/Blog"},
	{ID: "diary", Icon: "📔", Name: "日记", Desc: "日常碎片与随想", Path: "/Blog/Diary"},
	{ID: "portfolio", Icon: "🎨", Name: "作品", Desc: "网站、应用与开源项目", Path: "/Blog/Portfolio"},
	{ID: "tools", Icon: "🧰", Name: "工具", Desc: "自研工具与常用链接", Path: "/Blog/Tools"},
	{ID: "bookmark", Icon: "🔖", Name: "书签", Desc: "值得反复访问的资源", Path: "/Blog/Bookmark"},
	{ID: "about", Icon: "👤", Name: "关于", Desc: "个人资料与站点记录", Path: "/Blog/About"},
}

// DefaultSiteModules 返回六个固定模块的副本，默认全部可见。
func DefaultSiteModules() []SiteModule {
	result := make([]SiteModule, 0, len(fixedSiteModules))
	for _, module := range fixedSiteModules {
		module.Visible = true
		result = append(result, module)
	}
	return result
}

// NormalizeSiteModules 按固定顺序返回完整的六个模块。
//
// strict 用于管理端写入：缺少固定模块、出现未知 ID、重复 ID、路径被修改或名称为空
// 都返回 ErrInvalidSiteModule。
//
// 非 strict 用于读取历史数据：未知 ID 被丢弃，缺失的模块按默认配置补齐，
// 因此旧数据不会因为缺少模块而被解释为隐藏。
func NormalizeSiteModules(modules []SiteModule, strict bool) ([]SiteModule, error) {
	provided := make(map[string]SiteModule, len(fixedSiteModules))
	for _, module := range modules {
		id := strings.TrimSpace(module.ID)
		fixed, ok := fixedSiteModule(id)
		if !ok {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			continue
		}
		if _, exists := provided[id]; exists {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			continue
		}
		if strict && module.Path != fixed.Path {
			return nil, ErrInvalidSiteModule
		}
		if strings.TrimSpace(module.Name) == "" {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			module.Name = fixed.Name
		}
		module.ID = fixed.ID
		module.Path = fixed.Path
		provided[id] = module
	}

	result := make([]SiteModule, 0, len(fixedSiteModules))
	for _, fixed := range fixedSiteModules {
		if value, ok := provided[fixed.ID]; ok {
			result = append(result, value)
			continue
		}
		if strict {
			return nil, ErrInvalidSiteModule
		}
		defaultModule := fixed
		defaultModule.Visible = true
		result = append(result, defaultModule)
	}
	return result, nil
}

func fixedSiteModule(id string) (SiteModule, bool) {
	for _, module := range fixedSiteModules {
		if module.ID == id {
			return module, true
		}
	}
	return SiteModule{}, false
}

// UnmarshalJSON 反序列化站点模块。缺少 visible 的旧数据按 true 处理，
// 避免把缺失字段解释为隐藏模块；显式传入 false 时保持关闭。
func (m *SiteModule) UnmarshalJSON(data []byte) error {
	type module SiteModule
	value := struct {
		*module
		Visible *bool `json:"visible"`
	}{module: (*module)(m)}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	m.Visible = value.Visible == nil || *value.Visible
	return nil
}
