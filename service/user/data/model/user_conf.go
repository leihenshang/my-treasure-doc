package model

// UserConf 用户配置表
type UserConf struct {
	BaseModel
	UserId int64  `gorm:"column:user_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"userId"` // 用户id
	Key    string `gorm:"column:key;type:varchar(64);default:'';NOT NULL" json:"key"`
	Value  string `gorm:"column:value;type:text;NOT NULL" json:"value"`
}

func (m *UserConf) TableName() string {
	return "td_user_conf"
}

type EditorTheme string

const (
	EditorThemeDark  EditorTheme = "dark"
	EditorThemeLight EditorTheme = "light"
)

type EditorPreference string

const (
	EditorPreferenceVditor         EditorPreference = "vditor"
	EditorPreferenceCherryMarkdown EditorPreference = "cherryMarkdown"
)

type UserConfItem struct {
	EditorTheme      EditorTheme      `json:"editorTheme"`
	EditorPreference EditorPreference `json:"editorPreference"`
}

func NewUserConfItem() *UserConfItem {
	return &UserConfItem{
		EditorTheme:      EditorThemeLight,
		EditorPreference: EditorPreferenceVditor,
	}
}

func GetEditorList() []EditorPreference {
	return []EditorPreference{
		EditorPreferenceVditor,
		EditorPreferenceCherryMarkdown,
	}
}
