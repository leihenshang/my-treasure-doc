package order

import "fastduck/treasure-doc/service/mall/data/request/common"

type FilterOrderList struct {
	common.Pagination
	common.DataSort
	UserId int32 `json:"userId" form:"userId"`
}

type FilterOrderDetail struct {
	OrderId int32 `json:"orderId" form:"orderId"`
	UserId  int32 `json:"userId" form:"userId"`
}

type FilterOrderCreate struct {
	SkuId    int32 `json:"skuId"  form:"skuId"`
	Quantity int32 `json:"quantity" form:"quantity"`
	UserId   int32 `json:"userId" form:"userId"`
}
