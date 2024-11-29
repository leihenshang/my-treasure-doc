package model

type Note struct {
	BasicModel
	UserId   int64    `gorm:"column:user_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"userId"` // 用户id
	Title    string   `gorm:"column:title;type:varchar(100);NOT NULL" json:"title"`                     // 标题
	Content  string   `gorm:"column:content;type:text;NOT NULL" json:"content"`                         // 文档内容
	Color    string   `gorm:"column:color;type:varchar(100);default:'';NOT NULL;" json:"color"`         // 颜色
	Icon     string   `gorm:"column:icon;type:varchar(100);NOT NULL;default:''" json:"icon"`            // 颜色
	IsTop    int      `gorm:"column:is_top;type:tinyint(4);default:0;NOT NULL" json:"isTop"`            // 是否置顶
	Priority int      `gorm:"column:priority;type:int(255);default:0;NOT NULL" json:"priority"`         // 优先级
	DocId    int64    `gorm:"column:doc_id;type:bigint(20);default:0;NOT NULL" json:"docId"`
	NoteType NoteType `gorm:"column:note_type;type:varchar(100);default:'';NOT NULL" json:"noteType"`
}

type NoteType string
type NoteTypes []string

const (
	NoteTypeBookmark = `bookmark`
	NoteTypeTreeHole = `treeHole`
	NoteTypeTreeNote = `note`
	NoteTypeDoc      = `doc`
)

type Notes []*Note

func (m *Note) TableName() string {
	return "td_note"
}

func (n Notes) GetDocIds() []int64 {
	ids := make([]int64, 0)
	for _, note := range n {
		if note.DocId > 0 {
			ids = append(ids, note.DocId)
		}
	}
	return ids
}
