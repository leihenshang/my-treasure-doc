package model

type Doc struct {
	BasicModel
	UserId    int64  `gorm:"column:user_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"userId"` // 用户id
	Title     string `gorm:"column:title;type:varchar(100);NOT NULL" json:"title"`                     // 标题
	Content   string `gorm:"column:content;type:text;NOT NULL" json:"content"`                         // 文档内容
	DocStatus int    `gorm:"column:doc_status;type:tinyint(4);default:1;NOT NULL" json:"docStatus"`    // 1正常2审核中3禁用
	ViewCount int    `gorm:"column:view_count;type:int(11);default:0;NOT NULL" json:"viewCount"`       // 查看次数
	LikeCount int    `gorm:"column:like_count;type:int(11);default:0;NOT NULL" json:"likeCount"`       // 点赞次数
	IsTop     int    `gorm:"column:is_top;type:tinyint(4);default:0;NOT NULL" json:"isTop"`            // 是否置顶
	Priority  int    `gorm:"column:priority;type:int(255);default:0;NOT NULL" json:"priority"`         // 优先级
	GroupId   int64  `gorm:"column:group_id;type:bigint(20);default:0;NOT NULL" json:"groupId"`        // 分组id
	GroupPath string `gorm:"-:all" json:"groupPath"`
}

type Docs []*Doc

func (m *Doc) TableName() string {
	return "td_doc"
}

func (d Docs) GetGroupIds(unique bool) []int64 {
	var groupIds []int64
	uniqueMap := make(map[int64]struct{})
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
