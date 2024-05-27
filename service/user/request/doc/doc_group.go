package doc

//CreateDocGroupRequest 创建文档分组
type CreateDocGroupRequest struct {
	Title string `json:"title" binding:"required,min=1,max=250"` // 标题
	PId   int    `json:"pId" binding:""`                         // 父级
	Icon  string `json:"icon" binding:""`                        // 图标
}

//UpdateDocGroupRequest 更新文档分组
type UpdateDocGroupRequest struct {
	Id    int    `json:"id" binding:"required"`
	Title string `json:"title" binding:"max=250"` // 标题
	PId   int    `json:"pId" binding:""`          // 父级
	Icon  string `json:"icon" binding:""`         // 图标
}

//GroupTreeRequest 文档分组树
type GroupTreeRequest struct {
	Pid     int `json:"pid" form:"pid"`         // 父级id
	WithDoc int `json:"withDoc" form:"withDoc"` // 是否带文档 0-否，1-是
	IsAll   int `json:"isAll" form:"isAll"`     // 是否返回全部 0-否，1-是
}
