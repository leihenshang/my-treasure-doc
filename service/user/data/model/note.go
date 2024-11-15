package model

type Note struct {
	BasicModel
	UserId   int64    `gorm:"column:user_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"userId"` // 用户id
	Content  string   `gorm:"column:content;type:text;NOT NULL" json:"content"`                         // 文档内容
	IsTop    int      `gorm:"column:is_top;type:tinyint(4);default:0;NOT NULL" json:"isTop"`            // 是否置顶
	Priority int      `gorm:"column:priority;type:int(255);default:0;NOT NULL" json:"priority"`         // 优先级
	NoteType NoteType `gorm:"column:note_type;type:varchar(100);default:'';NOT NULL" json:"noteType"`
}

type NoteType string

const (
	NoteTypeBookmark = `bookmark`
	NoteTypeTreeHole = `treeHole`
	NoteTypeTreeNote = `note`
)

type Notes []*Note

func (m *Note) TableName() string {
	return "td_note"
}
