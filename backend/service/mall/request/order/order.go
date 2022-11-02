package order

import "fastduck/treasure-doc/service/mall/request/common"

type FilterOrderList struct {
	common.Pagination
	common.DataSort
}

type FilterOrderDetail struct {
	OrderId int32 `json:"orderId" form:"orderId"`
}
