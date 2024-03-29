package request

type IdRequest struct {
	Id int `json:"id" form:"id"  xml:"id"`
}

type ListRequest struct {
	Page      int    `json:"page" form:"page" xml:"page"`
	PageSize  int    `json:"pageSize" form:"pageSize" xml:"pageSize"`
	SortField string `json:"sortField" form:"sortField" xml:"sortField"`
	IsDesc    bool   `json:"isDesc" form:"isDesc" xml:"isDesc"`
	GroupId   int    `json:"groupId" form:"groupId" xml:"groupId"`
}
