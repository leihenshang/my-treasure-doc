package pay

import (
	"github.com/gin-gonic/gin"

	reqPay "fastduck/treasure-doc/service/mall/data/request/pay"
	"fastduck/treasure-doc/service/mall/data/response"
	"fastduck/treasure-doc/service/mall/global"
	srvPay "fastduck/treasure-doc/service/mall/internal/service/pay"
	"fastduck/treasure-doc/service/mall/middleware/auth"
)

func Create(c *gin.Context) {
	var req reqPay.ParamsPayCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ZapSugar.Infof("[pay|Create] parse user request data err:%+v", err)
		response.FailWithMessage(global.ErrResp(err), c)
		return
	}

	u, err := auth.GetUserInfoByCtx(c)
	if err != nil {
		global.ZapSugar.Infof("[pay|srvPay.PayCreate] get user info err:%+v", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserId = int32(u.ID)

	if d, ok := srvPay.Create(c, req); ok != nil {
		global.ZapSugar.Infof("[pay|srvPay.PayCreate] err:%+v", ok)
		response.FailWithMessage(ok.Error(), c)
	} else {
		response.OkWithData(d, c)
	}
}
