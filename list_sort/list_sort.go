package listsort

import (
	"encoding/json"
	"errors"
	"strings"
)

type OrderStr string

var ErrInvalidSortParams error = errors.New("invalid sort params")

const (
	OrderAsc    OrderStr = "asc"
	OrderDesc   OrderStr = "desc"
	OrderAscUp           = "ASC"
	OrderDescUp          = "DESC"
)

type SortObj struct {
	SortParams SortParams
	RawSort    string
	Err        error
}

type SortParam struct {
	Field string
	Order OrderStr
}

type SortParams []*SortParam

func ParseSortParams(rawSort string, strType string, filter ...string) *SortObj {
	obj := &SortObj{RawSort: rawSort}
	if strType == "json" {
		err := json.Unmarshal([]byte(rawSort), &obj.SortParams)
		if err != nil {
			obj.Err = err
		}
	}

	filterMap := make(map[string]struct{})
	for _, v := range filter {
		if v == "" {
			continue
		}
		filterMap[v] = struct{}{}
	}

	for _, v := range strings.Split(rawSort, ",") {
		res := strings.Split(v, "_")
		if len(res) == 2 {
			if len(filterMap) > 0 {
				if _, ok := filterMap[res[0]]; ok {
					obj.SortParams = append(obj.SortParams, &SortParam{Field: res[0], Order: OrderStr(res[1])})
				}
			} else {
				obj.SortParams = append(obj.SortParams, &SortParam{Field: res[0], Order: OrderStr(res[1])})
			}
		}
	}
	if err := obj.SortParams.Validate(); err != nil {
		obj.Err = err
	}

	return obj
}

func (o SortObj) IsError() bool {
	return o.Err != nil
}

func (o SortObj) ShouldSort() bool {
	return len(o.SortParams) > 0
}

func (o SortObj) Exists(field string) bool {
	for _, v := range o.SortParams {
		if v.Field == field {
			return true
		}
	}
	return false
}

func (o SortParams) Validate() error {
	for _, param := range o {
		if err := param.Order.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (o OrderStr) Validate() error {
	if o == OrderAsc || o == OrderDesc || o == OrderAscUp || o == OrderDescUp {
		return nil
	}
	return errors.New("invalid order: " + string(o))
}
