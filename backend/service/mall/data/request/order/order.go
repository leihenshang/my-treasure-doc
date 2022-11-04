package order

import "fastduck/treasure-doc/service/mall/data/request/common"

type FilterOrderList struct {
	common.Pagination
	common.DataSort
	Status int32 `json:"status" from:"status"`
	UserId int32 `json:"userId" form:"userId"`
}

type FilterOrderDetail struct {
	OrderId int32 `json:"orderId" form:"orderId"`
	UserId  int32 `json:"userId" form:"userId"`
}

type ParamsOrderCreate struct {
	SkuId    int32 `json:"skuId"  form:"skuId"`
	Quantity int32 `json:"quantity" form:"quantity"`
	UserId   int32 `json:"userId" form:"userId"`
	//TODO 模拟超卖，记得后面删掉
	Exceed int32 `json:"exceed" form:"exceed"`
}
