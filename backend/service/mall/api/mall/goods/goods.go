package goods

import (
	"fastduck/treasure-doc/service/mall/global"
	srvGoods "fastduck/treasure-doc/service/mall/internal/service/goods"
	reqGoods "fastduck/treasure-doc/service/mall/request/goods"
	"fastduck/treasure-doc/service/mall/response"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	resp := response.ListResponse{}
	var req reqGoods.FilterGoodsList
	if err := c.ShouldBindQuery(&req); err != nil {
		global.ZAPSUGAR.Infof("goods|List err:%+v", err)
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	if d, t, ok := srvGoods.GoodsList(c, req); ok != nil {
		global.ZAPSUGAR.Infof("goods|srvGoods.GoodsList err:%+v", ok)
		response.FailWithMessage(ok.Error(), c)
	} else {
		resp.List = d
		resp.Total = t
		response.OkWithData(resp, c)
	}
}

func Detail(c *gin.Context) {
	var req reqGoods.FilterGoodsDetail
	if err := c.ShouldBindQuery(&req); err != nil {
		global.ZAPSUGAR.Infof("goods|List err:%+v", err)
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	if d, ok := srvGoods.GoodsDetail(c, req); ok != nil {
		global.ZAPSUGAR.Infof("goods|srvGoods.GoodsList err:%+v", ok)
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}
