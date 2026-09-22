package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIPWhiteAllowed(t *testing.T) {
	cases := []struct {
		ip      string
		allowed []string
		want    bool
	}{
		{"192.168.1.5", []string{"192.168.1.0/24"}, true},
		{"192.168.2.5", []string{"192.168.1.0/24"}, false},
		{"10.0.0.1", []string{"10.0.0.1"}, true},
		{"10.0.0.2", []string{"10.0.0.1"}, false},
		{"2001:db8::1", []string{"2001:db8::/32"}, true},
		{"", []string{"192.168.1.0/24"}, false},
	}
	for _, tc := range cases {
		if got := ipAllowed(tc.ip, tc.allowed); got != tc.want {
			t.Fatalf("ipAllowed(%q, %v) = %v，想要 %v", tc.ip, tc.allowed, got, tc.want)
		}
	}
}

func TestIPWhitelistHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/x", IPWhitelist([]string{"127.0.0.0/8"}), func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	engine.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("白名单内 IP = %d，想要 200", resp.Code)
	}

	resp = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/x", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	engine.ServeHTTP(resp, req)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("白名单外 IP = %d，想要 403", resp.Code)
	}

	// 空白名单 = 放行全部
	engine2 := gin.New()
	engine2.GET("/y", IPWhitelist(nil), func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	resp = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/y", nil)
	req.RemoteAddr = "10.0.0.2:1"
	engine2.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("空白名单 = %d，想要 200", resp.Code)
	}
}
