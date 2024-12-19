package model

type DocHistory struct {
	BaseModel
	DocId   string `gorm:"column:doc_id;type:varchar(100);default:'';" json:"docId"`
	UserId  string `gorm:"column:user_id;type:varchar(100);default:'';" json:"userId"` // 用户id
	Title   string `gorm:"column:title;type:varchar(100);default:'';" json:"title"`    // 标题
	Content string `gorm:"column:content;type:text;NOT NULL" json:"content"`           // 文档内容
}

type DocHistories []*DocHistory

func (m *DocHistory) TableName() string {
	return "td_doc_history"
}
