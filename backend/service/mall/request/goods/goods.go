package goods

import (
	"fastduck/treasure-doc/service/mall/request/common"
)

type FilterGoodsList struct {
	GoodsName string `json:"goodsName" form:"goodsName"`
	common.Pagination
	common.DataSort
}

type FilterGoodsDetail struct {
	GoodsId string `json:"goodsId"`
}
