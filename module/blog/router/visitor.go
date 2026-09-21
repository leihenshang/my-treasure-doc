package router

import (
	"log"
	"strings"

	"fastduck/treasure-doc/module/blog/data/model"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

// VisitorLogger 记录公开博客 API 的一次访客访问（IP + 路径）。
//
// 只挂载在 /api/blog/* 路由组上，后台与鉴权请求不经过这里；
// 部署在前置反代后面时，通过 X-Real-IP / X-Forwarded-For 取真实访客 IP。
// 个人博客流量下逐请求同步插入可接受，插入失败只记日志、不影响接口响应。
func VisitorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if global.Db == nil {
			return
		}
		entry := &model.VisitorLog{
			IP:        clientIP(c),
			Path:      c.Request.URL.Path,
			Category:  visitorCategory(c.Request.URL.Path),
			UserAgent: c.GetHeader("User-Agent"),
		}
		if err := global.Db.Create(entry).Error; err != nil {
			log.Printf("[visitor] 记录访客日志失败: %v", err)
		}
	}
}

// visitorCategory 从公开博客路径解析出访问的内容分类。
// 挂在 /api/blog/* 上的路径形如 /api/blog/posts…、/api/blog/diaries…，据此区分内容类型。
func visitorCategory(path string) model.VisitorCategory {
	switch {
	case strings.HasPrefix(path, "/api/blog/posts") || strings.HasPrefix(path, "/api/blog/archive"):
		return model.VisitorCategoryPost
	case strings.HasPrefix(path, "/api/blog/diaries") || strings.HasPrefix(path, "/api/blog/diary"):
		return model.VisitorCategoryDiary
	case strings.HasPrefix(path, "/api/blog/portfolio"):
		return model.VisitorCategoryPortfolio
	case strings.HasPrefix(path, "/api/blog/tools"):
		return model.VisitorCategoryTool
	case strings.HasPrefix(path, "/api/blog/bookmarks") || strings.HasPrefix(path, "/api/blog/bookmark"):
		return model.VisitorCategoryBookmark
	case strings.HasPrefix(path, "/api/blog/"):
		return model.VisitorCategorySite
	default:
		return model.VisitorCategoryUnresolved
	}
}

// clientIP 提取真实访客 IP：优先反代设置的头，回退到连接地址。
func clientIP(c *gin.Context) string {
	if ip := strings.TrimSpace(c.GetHeader("X-Real-IP")); ip != "" {
		return ip
	}
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		if first := strings.TrimSpace(strings.Split(forwarded, ",")[0]); first != "" {
			return first
		}
	}
	return c.ClientIP()
}
