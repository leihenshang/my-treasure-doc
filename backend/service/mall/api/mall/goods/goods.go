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
		response.FailWithMessage(global.ErrResp(err), c)
		global.ZAPSUGAR.Infof("goods|List err:%+v", err)
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

}
