package order

import (
	reqOrder "fastduck/treasure-doc/service/mall/data/request/order"
	"fastduck/treasure-doc/service/mall/data/response"
	"fastduck/treasure-doc/service/mall/global"
	srvOrder "fastduck/treasure-doc/service/mall/internal/service/order"
	"fastduck/treasure-doc/service/mall/middleware/auth"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	var req reqOrder.FilterOrderList
	if err := c.ShouldBindQuery(&req); err != nil {
		global.ZapSugar.Infof("order|List err:%+v", err)
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}
	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		global.ZapSugar.Infof("[order|srvOrder.OrderCreate] get user info err:%+v", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserId = int32(u.ID)

	if d, ok := srvOrder.OrderList(c, req); ok != nil {
		global.ZapSugar.Infof("[order|srvOrder.OrderList] err:%+v", ok)
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

func Detail(c *gin.Context) {
	var req reqOrder.FilterOrderDetail
	if err := c.ShouldBindQuery(&req); err != nil {
		global.ZapSugar.Infof("order|List err:%+v", err)
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	if d, ok := srvOrder.OrderDetail(c, req); ok != nil {
		global.ZapSugar.Infof("order|srvOrder.OrderList err:%+v", ok)
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}

func Create(c *gin.Context) {
	var req reqOrder.FilterOrderCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ZapSugar.Infof("[order|Create] parse user request data err:%+v", err)
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		global.ZapSugar.Infof("[order|srvOrder.OrderCreate] get user info err:%+v", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserId = int32(u.ID)

	if d, ok := srvOrder.OrderCreate(c, req); ok != nil {
		global.ZapSugar.Infof("[order|srvOrder.OrderCreate] err:%+v", ok)
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}
