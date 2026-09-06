package global

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	blogmodel "fastduck/treasure-doc/module/blog/data/model"
	usermodel "fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/config"

	"gorm.io/gorm"
)

// TestSQLiteSchemaAndJSONRoundTrip 用真实 SQLite 文件验证：
//  1. 跨方言模型（含 datatypes.JSON 与 timestamp）可正常 AutoMigrate；
//  2. JSON 字段能写回并读回；
//  3. 乐观锁 version + 1 在 SQLite 下生效。
func TestSQLiteSchemaAndJSONRoundTrip(t *testing.T) {
	dir := t.TempDir()
	dsn := "file:" + filepath.Join(dir, "smoke.db") + "?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on"

	db, err := openDatabaseWithConfig(&config.Config{
		Database: config.Database{Driver: config.DriverSQLite, Dsn: dsn},
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if sqlDb, err := db.DB(); err == nil {
		defer sqlDb.Close()
	}

	tables := append([]interface{}{&usermodel.User{}, &usermodel.UserToken{}}, blogmodel.Tables()...)
	if err := db.AutoMigrate(tables...); err != nil {
		t.Fatalf("automigrate: %v", err)
	}

	item := blogmodel.PortfolioItem{
		Slug:          "demo",
		Title:         "Demo",
		Summary:       "s",
		CategoryID:    "c",
		Cover:         "x",
		TechStack:     blogmodel.NewJSON([]string{"Go", "Vue"}),
		Links:         blogmodel.NewJSON([]string{"https://example.com"}),
		Content:       "c",
		PublishStatus: blogmodel.StatusDraft,
		PublishedOn:   time.Now(),
		PublishedAt:   time.Now(),
		Version:       1,
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("create portfolio: %v", err)
	}

	var got blogmodel.PortfolioItem
	if err := db.First(&got, "slug = ?", "demo").Error; err != nil {
		t.Fatalf("read portfolio: %v", err)
	}
	var tech []string
	if err := json.Unmarshal(got.TechStack, &tech); err != nil {
		t.Fatalf("unmarshal tech_stack: %v", err)
	}
	if len(tech) != 2 || tech[0] != "Go" {
		t.Fatalf("tech_stack round trip = %#v", tech)
	}

	// 乐观锁：version + 1
	if err := db.Model(&got).Update("version", gorm.Expr("version + 1")).Error; err != nil {
		t.Fatalf("update version: %v", err)
	}
	var after blogmodel.PortfolioItem
	if err := db.First(&after, "slug = ?", "demo").Error; err != nil {
		t.Fatalf("read after update: %v", err)
	}
	if after.Version != 2 {
		t.Fatalf("version after update = %d, want 2", after.Version)
	}
}
