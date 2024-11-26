package model

type DocHistory struct {
	BasicModel
	DocId   int64  `gorm:"column:doc_id;type:bigint(20) unsigned;" json:"docId"`
	UserId  int64  `gorm:"column:user_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"userId"` // 用户id
	Title   string `gorm:"column:title;type:varchar(100);NOT NULL" json:"title"`                     // 标题
	Content string `gorm:"column:content;type:text;NOT NULL" json:"content"`                         // 文档内容
}

type DocHistories []*DocHistory

func (m *DocHistory) TableName() string {
	return "td_doc_history"
}
