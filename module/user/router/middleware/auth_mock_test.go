package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fastduck/treasure-doc/module/user/config"
	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

func TestConfiguredAdminMock(t *testing.T) {
	previous := global.GetConf()
	defer func() { global.Conf = previous }()

	tests := []struct {
		name       string
		config     *config.Config
		wantStatus int
		wantID     string
	}{
		{name: "dev default root", config: &config.Config{App: config.App{RunMode: config.GinModeDev}, Debug: config.Debug{EnableMockLogin: true}}, wantStatus: http.StatusOK, wantID: "9999999999"},
		{name: "dev custom root", config: &config.Config{App: config.App{RunMode: config.GinModeDev}, Debug: config.Debug{EnableMockLogin: true, MockUserId: "mock-admin"}}, wantStatus: http.StatusOK, wantID: "mock-admin"},
		{name: "release does not mock", config: &config.Config{App: config.App{RunMode: config.GinModeRelease}, Debug: config.Debug{EnableMockLogin: true}}, wantStatus: http.StatusUnauthorized},
	}

	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			global.Conf = test.config
			engine := gin.New()
			engine.Use(Auth(), RequireAdmin())
			engine.GET("/admin", func(c *gin.Context) {
				value, _ := c.Get(global.UserInfoKey)
				user := value.(*model.User)
				c.JSON(http.StatusOK, gin.H{"id": user.Id, "userType": user.UserType})
			})
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/admin", nil))
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if test.wantID != "" && !strings.Contains(recorder.Body.String(), `"id":"`+test.wantID+`"`) {
				t.Fatalf("mock ID missing: %s", recorder.Body.String())
			}
			if test.wantID != "" && !strings.Contains(recorder.Body.String(), `"userType":100`) {
				t.Fatalf("mock should be root: %s", recorder.Body.String())
			}
		})
	}
}
