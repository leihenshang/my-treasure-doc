package middleware

import (
	"net"
	"net/http"
	"net/netip"
	"strings"

	commonresponse "fastduck/treasure-doc/module/common/response"

	"github.com/gin-gonic/gin"
)

// IPWhitelist 白名单校验：客户端 IP 命中 allowed 才放行，否则 403。
// allowed 支持精确 IP 与 CIDR（如 192.168.1.0/24）；为空时放行全部（仅靠令牌鉴权）。
// 客户端 IP 取自 gin ClientIP()（会考虑 X-Forwarded-For 等代理头），需在反代下正确配置信任代理。
func IPWhitelist(allowed []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(allowed) == 0 {
			c.Next()
			return
		}
		clientIP := c.ClientIP()
		if clientIP == "" {
			if host, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
				clientIP = host
			} else {
				clientIP = c.Request.RemoteAddr
			}
		}
		if ipAllowed(clientIP, allowed) {
			c.Next()
			return
		}
		commonresponse.Error(c, http.StatusForbidden, 40300, "来源 IP 不在白名单内")
		c.Abort()
	}
}

// ipAllowed 判断 IP 是否匹配任意规则（精确 IP 或 CIDR 前缀）。
func ipAllowed(clientIP string, allowed []string) bool {
	addr, err := netip.ParseAddr(strings.TrimSpace(clientIP))
	if err != nil {
		return false
	}
	for _, rule := range allowed {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		if prefix, err := netip.ParsePrefix(rule); err == nil {
			if prefix.Contains(addr) {
				return true
			}
			continue
		}
		if target, err := netip.ParseAddr(rule); err == nil && target == addr {
			return true
		}
	}
	return false
}
