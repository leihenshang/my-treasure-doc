package model

// DocGroup 文档分组
type DocGroup struct {
	BasicModel
	UserId    int64     `gorm:"column:user_id;default:0;NOT NULL" json:"userId"`    // 用户Id
	Title     string    `gorm:"column:title;NOT NULL" json:"title"`                 // 组名
	Icon      string    `gorm:"column:icon;NOT NULL" json:"icon"`                   // 图标
	PId       int64     `gorm:"column:p_id;default:0;NOT NULL" json:"pid"`          // 父级id
	Priority  int       `gorm:"column:priority;default:0;NOT NULL" json:"priority"` // 优先级
	GroupPath string    `gorm:"column:group_path;default:'';NOT NULL" json:"groupPath"`
	GroupType GroupType `gorm:"-:all" json:"groupType"`
	IsLeaf    bool      `gorm:"-:all" json:"isLeaf"`
}

type GroupType string

const (
	GroupTypeGroup GroupType = "group"
	GroupTypeDoc   GroupType = "doc"
)

type DocGroups []*DocGroup

func (m *DocGroup) TableName() string {
	return "td_doc_group"
}

func (d DocGroups) GetIds() []int64 {
	ids := make([]int64, len(d))
	for k, v := range d {
		ids[k] = v.Id
	}
	return ids
}

func (d DocGroups) ToPidMap() map[int64]*DocGroup {
	groupMap := make(map[int64]*DocGroup, len(d))
	for _, v := range d {
		groupMap[v.PId] = v
	}
	return groupMap
}
