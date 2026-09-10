package response

import (
	"encoding/json"
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	MaxSiteModuleNameLength     = 50
	MaxSiteModuleTitleLength    = 50
	MaxSiteModuleSubtitleLength = 50
	MaxSiteModuleMarkerLength   = 50
	MaxSiteHomeSubtitleLength   = 50
)

// ErrInvalidSiteModule 表示站点模块集合不满足固定模块约束。
var ErrInvalidSiteModule = errors.New("invalid site module")

// fixedSiteModules 是站点模块的固定集合。ID 与 Path 属于对外契约，
// 管理端可修改 icon、name、title、desc、marker 和 visible。
var fixedSiteModules = []SiteModule{
	{ID: "home", Icon: "⌘", Name: "主页", Title: "创造力改变世界", Desc: "AI 编程软件的产品介绍、核心能力与开始入口", Path: "/Blog/Home", Marker: "HOME"},
	{ID: "blog", Icon: "📝", Name: "文章", Title: "我的文章", Desc: "技术笔记与长文", Path: "/Blog", Marker: "BLOG"},
	{ID: "diary", Icon: "📔", Name: "日记", Title: "日记", Desc: "日常碎片与随想", Path: "/Blog/Diary", Marker: "DIARY"},
	{ID: "portfolio", Icon: "🎨", Name: "作品", Title: "作品集", Desc: "网站、应用与开源项目", Path: "/Blog/Portfolio", Marker: "PORTFOLIO"},
	{ID: "tools", Icon: "🧰", Name: "工具", Title: "利器", Desc: "自研工具与常用链接", Path: "/Blog/Tools", Marker: "TOOLS"},
	{ID: "bookmark", Icon: "🔖", Name: "书签", Title: "收藏集", Desc: "值得反复访问的资源", Path: "/Blog/Bookmark", Marker: "BOOKMARK"},
	{ID: "about", Icon: "👤", Name: "关于", Title: "关于我", Desc: "个人资料与站点记录", Path: "/Blog/About", Marker: "ABOUT"},
}

func DefaultSiteHome() SiteHome {
	return SiteHome{
		Title:    "创造力改变世界",
		Subtitle: "用技术、设计与 AI，把想法变成真实的作品。",
		AI: SiteHomeAI{
			Eyebrow:     "AI CREATIVE LAB",
			Title:       "让 AI 成为创造力的放大器",
			Description: "探索 AI 如何帮助思考、表达、设计与构建，让每一个灵感更快抵达现实。",
			LinkText:    "探索 AI 创作",
			LinkURL:     "/Blog/Tools",
		},
	}
}

func DefaultSiteFooter() SiteFooter {
	return SiteFooter{Text: "© 2026 Treasure Doc · 创造力改变世界"}
}

func DefaultSiteBanner() SiteBanner {
	return SiteBanner{BackgroundColor: "#1769ff", TextColor: "#ffffff"}
}

// DefaultSiteModules 返回固定模块的副本，默认全部可见。
func DefaultSiteModules() []SiteModule {
	result := make([]SiteModule, 0, len(fixedSiteModules))
	for _, module := range fixedSiteModules {
		module.Visible = true
		result = append(result, module)
	}
	return result
}

// NormalizeSiteModules 按固定顺序返回完整的模块集合。
//
// strict 用于管理端写入：缺少固定模块、出现未知 ID、重复 ID、路径被修改、名称为空
// 或标题超长都返回 ErrInvalidSiteModule；标题可以为空。
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
		module.Name = strings.TrimSpace(module.Name)
		if module.Name == "" {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			module.Name = fixed.Name
		}
		if utf8.RuneCountInString(module.Name) > MaxSiteModuleNameLength {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			module.Name = fixed.Name
		}
		module.Title = strings.TrimSpace(module.Title)
		if utf8.RuneCountInString(module.Title) > MaxSiteModuleTitleLength {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			module.Title = fixed.Title
		}
		module.Desc = strings.TrimSpace(module.Desc)
		if utf8.RuneCountInString(module.Desc) > MaxSiteModuleSubtitleLength {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			module.Desc = fixed.Desc
		}
		module.Marker = strings.TrimSpace(module.Marker)
		if module.Marker == "" {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			module.Marker = fixed.Marker
		}
		if utf8.RuneCountInString(module.Marker) > MaxSiteModuleMarkerLength {
			if strict {
				return nil, ErrInvalidSiteModule
			}
			module.Marker = fixed.Marker
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
