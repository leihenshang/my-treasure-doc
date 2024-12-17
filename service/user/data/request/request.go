package request

import (
	"errors"
	"strings"

	"fastduck/treasure-doc/service/user/gid"
)

type IDReq struct {
	ID gid.Gid `json:"id" form:"id"  xml:"id"`
}

type ListPagination struct {
	Page     int   `json:"page" form:"page" xml:"page"`
	PageSize int   `json:"pageSize" form:"pageSize" xml:"pageSize"`
	Total    int64 `json:"total" form:"total" xml:"total"`
}

type ListSort struct {
	// OrderBy orderBy: id_asc,name_desc
	OrderBy string `json:"orderBy" form:"orderBy" xml:"orderBy"`
}

func (l ListSort) Sort(sortFields map[string]string) (string, error) {
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

const OrderByDesc OrderByType = "DESC"
const OrderByAsc OrderByType = "ASC"
const OrderByDescLower OrderByType = "desc"
const OrderByAscLower OrderByType = "asc"

func (o OrderByType) Check() bool {
	return strings.ToUpper(string(o)) != string(OrderByDesc) || strings.ToLower(string(o)) != string(OrderByAsc)
}

func (o OrderByType) UpperString() string {
	return strings.ToUpper(string(o))
}

func (p ListPagination) Offset() int {
	offset := (p.Page - 1) * p.PageSize
	if offset < 0 {
		offset = 1
	}
	return offset
}
