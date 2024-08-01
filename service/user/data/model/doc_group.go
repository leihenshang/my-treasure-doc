package model

// DocGroup 文档分组
type DocGroup struct {
	BasicModel
	UserId        uint64 `gorm:"column:user_id;default:0;NOT NULL" json:"userId"`    // 用户Id
	Title         string `gorm:"column:title;NOT NULL" json:"title"`                 // 组名
	Icon          string `gorm:"column:icon;NOT NULL" json:"icon"`                   // 图标
	PId           uint64 `gorm:"column:p_id;default:0;NOT NULL" json:"pid"`          // 父级id
	Priority      int    `gorm:"column:priority;default:0;NOT NULL" json:"priority"` // 优先级
	ChildrenCount int    `gorm:"-:migration" json:"childrenCount"`
}

type DocGroups []*DocGroup

func (m *DocGroup) TableName() string {
	return "td_doc_group"
}

func (d DocGroups) GetIds() []uint64 {
	ids := make([]uint64, len(d))
	for k, v := range d {
		ids[k] = v.Id
	}
	return ids
}

func (d DocGroups) ToPidMap() map[uint64]*DocGroup {
	groupMap := make(map[uint64]*DocGroup, len(d))
	for _, v := range d {
		groupMap[v.PId] = v
	}
	return groupMap
}
