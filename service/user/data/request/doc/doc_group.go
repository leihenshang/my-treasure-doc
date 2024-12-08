package doc

// CreateDocGroupRequest 创建文档分组
type CreateDocGroupRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=250"` // 标题
	PId      int64  `json:"pId" binding:""`                         // 父级
	Icon     string `json:"icon" binding:""`                        // 图标
	Priority int    `json:"priority" binding:""`
}

// UpdateDocGroupRequest 更新文档分组
type UpdateDocGroupRequest struct {
	Id    int64  `json:"id" binding:"required"`
	Title string `json:"title" binding:"max=250"` // 标题
	PId   int64  `json:"pId" binding:""`          // 父级
	Icon  string `json:"icon" binding:""`         // 图标
}

// GroupTreeRequest 文档分组树
type GroupTreeRequest struct {
	Pid     int64 `json:"pid" form:"pid"`         // 父级id
	WithDoc bool  `json:"withDoc" form:"withDoc"` // 是否返回子集
}
