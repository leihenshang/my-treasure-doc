package common

type Pagination struct {
	Limit  uint32 `json:"limit" form:"limit"`
	Offset uint32 `json:"offset" form:"offset"`
}

type DataSort struct {
	SortField string `json:"sortField" form:"sortField" xml:"sortField"`
	IsDesc    bool   `json:"isDesc" form:"isDesc" xml:"isDesc"`
}

func NewPagination() *Pagination {
	return &Pagination{
		Limit:  20,
		Offset: 0,
	}
}
