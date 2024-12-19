package model

type Doc struct {
	BaseModel
	UserId    string    `gorm:"column:user_id;type:varchar(100);NOT NULL" json:"userId"`               // 用户id
	Title     string    `gorm:"column:title;type:varchar(100);NOT NULL" json:"title"`                  // 标题
	Content   string    `gorm:"column:content;type:text;NOT NULL" json:"content"`                      // 文档内容
	DocStatus int       `gorm:"column:doc_status;type:tinyint(4);default:1;NOT NULL" json:"docStatus"` // 1正常2审核中3禁用
	IsTop     int8      `gorm:"column:is_top;type:tinyint(4);NOT NULL" json:"isTop"`                   // 是否置顶
	Priority  int       `gorm:"column:priority;type:int(255);NOT NULL" json:"priority"`                // 优先级
	GroupId   string    `gorm:"column:group_id;type:varchar(100);NOT NULL" json:"groupId"`             // 分组id
	ReadOnly  int8      `gorm:"column:read_only;type:tinyint(4);NOT NULL" json:"readOnly"`             // 1-读写，2-只读
	Version   int       `gorm:"column:version;type:int(11);NOT NULL" json:"version"`
	GroupPath DocGroups `gorm:"-:all" json:"groupPath"`
	IsPin     int       `gorm:"-" json:"isPin"`
}

type Docs []*Doc

func (d *Doc) TableName() string {
	return "td_doc"
}

func (d Docs) GetGroupIds(unique bool) []string {
	var groupIds []string
	uniqueMap := make(map[string]struct{})
	for _, doc := range d {
		if unique {
			if _, ok := uniqueMap[doc.GroupId]; ok {
				continue
			}
		}
		groupIds = append(groupIds, doc.GroupId)
		if unique {
			uniqueMap[doc.GroupId] = struct{}{}
		}
	}
	return groupIds
}

func (d Docs) ToMap() map[string]*Doc {
	m := make(map[string]*Doc, len(d))
	for _, doc := range d {
		m[doc.Id] = doc
	}
	return m
}

func (d Docs) ToGroupIdMap() map[string]*Doc {
	m := make(map[string]*Doc, len(d))
	for _, doc := range d {
		m[doc.GroupId] = doc
	}
	return m
}

func (d *Doc) HiddenUnnecessary() *Doc {
	d.Content = ""
	return d
}
