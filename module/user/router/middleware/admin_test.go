package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/global"

	"github.com/gin-gonic/gin"
)

func TestRequireAdmin(t *testing.T) {
	tests := []struct {
		name       string
		user       *model.User
		wantStatus int
		wantCalled bool
	}{
		{name: "missing user", wantStatus: http.StatusUnauthorized},
		{name: "regular user", user: &model.User{UserType: model.UserTypeUser}, wantStatus: http.StatusForbidden},
		{name: "admin", user: &model.User{UserType: model.UserTypeAdmin}, wantStatus: http.StatusOK, wantCalled: true},
		{name: "root", user: &model.User{UserType: model.UserTypeRoot}, wantStatus: http.StatusOK, wantCalled: true},
	}

	gin.SetMode(gin.TestMode)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			called := false
			engine := gin.New()
			engine.Use(func(c *gin.Context) {
				if test.user != nil {
					c.Set(global.UserInfoKey, test.user)
				}
				c.Next()
			}, RequireAdmin())
			engine.GET("/protected", func(c *gin.Context) {
				called = true
				c.Status(http.StatusOK)
			})

			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/protected", nil))
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
			if called != test.wantCalled {
				t.Fatalf("handler called = %v, want %v", called, test.wantCalled)
			}
		})
	}
}
