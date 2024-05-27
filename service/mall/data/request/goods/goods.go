package goods

import (
	"fastduck/treasure-doc/service/mall/data/request/common"
)

type FilterGoodsList struct {
	GoodsName string `json:"goodsName" form:"goodsName"`
	common.Pagination
	common.DataSort
}

type FilterGoodsDetail struct {
	GoodsId int32 `json:"goodsId" form:"goodsId"`
}
