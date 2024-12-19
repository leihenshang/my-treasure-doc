package model

type Note struct {
	BaseModel
	UserId   string   `gorm:"column:user_id;type:varchar(100);" json:"userId"`        // 用户id
	Title    string   `gorm:"column:title;type:varchar(100);" json:"title"`           // 标题
	Content  string   `gorm:"column:content;type:text;NOT NULL" json:"content"`       // 文档内容
	Color    string   `gorm:"column:color;type:varchar(100);" json:"color"`           // 颜色
	Icon     string   `gorm:"column:icon;type:varchar(100);" json:"icon"`             // 颜色
	IsTop    int      `gorm:"column:is_top;type:tinyint(4);NOT NULL" json:"isTop"`    // 是否置顶
	Priority int      `gorm:"column:priority;type:int(255);NOT NULL" json:"priority"` // 优先级
	DocId    string   `gorm:"column:doc_id;type:varchar(100);" json:"docId"`
	NoteType NoteType `gorm:"column:note_type;type:varchar(100);" json:"noteType"`
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

func (n Notes) GetDocIds() []string {
	ids := make([]string, 0)
	for _, note := range n {
		if note.DocId != "" {
			ids = append(ids, note.DocId)
		}
	}
	return ids
}
