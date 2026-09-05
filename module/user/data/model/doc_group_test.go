package model

import (
	"sync"
	"testing"

	"gorm.io/gorm/schema"
)

// TestDocGroupSchemaHasRoomId 分组树接口按 room_id 过滤，模型缺少该字段会导致 SQL 引用不存在的列。
func TestDocGroupSchemaHasRoomId(t *testing.T) {
	s, err := schema.Parse(&DocGroup{}, &sync.Map{}, schema.NamingStrategy{TablePrefix: "td_", SingularTable: true})
	if err != nil {
		t.Fatalf("parse DocGroup schema failed: %v", err)
	}

	if s.Table != "td_doc_group" {
		t.Fatalf("DocGroup table = %s, want td_doc_group", s.Table)
	}

	field, ok := s.FieldsByName["RoomId"]
	if !ok {
		t.Fatal("DocGroup 缺少 RoomId 字段")
	}
	if field.DBName != "room_id" {
		t.Fatalf("RoomId column = %s, want room_id", field.DBName)
	}
}
