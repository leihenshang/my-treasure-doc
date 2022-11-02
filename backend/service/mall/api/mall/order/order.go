package order

import (
	reqOrder "fastduck/treasure-doc/service/mall/data/request/order"
	"fastduck/treasure-doc/service/mall/data/response"
	"fastduck/treasure-doc/service/mall/global"
	srvOrder "fastduck/treasure-doc/service/mall/internal/service/order"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	resp := response.ListResponse{}
	var req reqOrder.FilterOrderList
	if err := c.ShouldBindQuery(&req); err != nil {
		global.ZAPSUGAR.Infof("order|List err:%+v", err)
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	if d, t, ok := srvOrder.OrderList(c, req); ok != nil {
		global.ZAPSUGAR.Infof("order|srvOrder.OrderList err:%+v", ok)
		response.FailWithMessage(ok.Error(), c)
	} else {
		resp.List = d
		resp.Total = t
		response.OkWithData(resp, c)
	}
}

func Detail(c *gin.Context) {
	var req reqOrder.FilterOrderDetail
	if err := c.ShouldBindQuery(&req); err != nil {
		global.ZAPSUGAR.Infof("order|List err:%+v", err)
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	if d, ok := srvOrder.OrderDetail(c, req); ok != nil {
		global.ZAPSUGAR.Infof("order|srvOrder.OrderList err:%+v", ok)
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

func Create(c *gin.Context) {

}
