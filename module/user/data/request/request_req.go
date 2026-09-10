package request

import (
	"errors"
	"strings"
)

type Pagination struct {
	Page     int   `json:"page" form:"page" xml:"page"`
	PageSize int   `json:"pageSize" form:"pageSize" xml:"pageSize"`
	Total    int64 `json:"total" form:"total" xml:"total"`
}

func (p Pagination) Offset() int {
	offset := (p.Page - 1) * p.PageSize
	if offset < 0 {
		offset = 1
	}
	return offset
}

type Sort struct {
	// OrderBy 形如 id_asc,name_desc
	OrderBy string `json:"orderBy" form:"orderBy" xml:"orderBy"`
}

func (l Sort) Sort(sortFields map[string]string) (string, error) {
	if sortFields == nil {
		return "", errors.New("sort field is empty")
	}

	sortSet := strings.Split(l.OrderBy, ",")
	var res []string
	for _, s := range sortSet {
		sortItem := strings.Split(s, "_")
		if len(sortItem) != 2 {
			return "", errors.New("sort error")
		}
		if !OrderByType(sortItem[1]).Check() {
			return "", errors.New("order by type error,allow asc or desc only")
		}
		sortItem[1] = OrderByType(sortItem[1]).UpperString()
		if v, ok := sortFields[sortItem[0]]; ok {
			sortItem[0] = v
			res = append(res, strings.Join(sortItem, " "))
		}
	}

	return strings.Join(res, ","), nil
}

type OrderByType string

const (
	OrderByDesc OrderByType = "DESC"
	OrderByAsc  OrderByType = "ASC"
)

// Check 判断排序方向是否为 asc/desc（大小写不敏感）。
func (o OrderByType) Check() bool {
	value := strings.ToUpper(string(o))
	return value == string(OrderByDesc) || value == string(OrderByAsc)
}

func (o OrderByType) UpperString() string {
	return strings.ToUpper(string(o))
}
