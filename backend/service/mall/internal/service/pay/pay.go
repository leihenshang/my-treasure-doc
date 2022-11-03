package pay

import (
	"context"
	"errors"
	payReq "fastduck/treasure-doc/service/mall/data/request/pay"
	payResp "fastduck/treasure-doc/service/mall/data/response/pay"
	"fastduck/treasure-doc/service/mall/global"
	orderDao "fastduck/treasure-doc/service/mall/internal/dao/order"
)

func Create(ctx context.Context, params payReq.ParamsPayCreate) (res *payResp.PayCreate, err error) {
	if params.OrderId <= 0 {
		err = errors.New("订单Id不能为空")
		return
	}
	if params.UserId <= 0 {
		err = errors.New("用户Id不能为空")
		return
	}

	orderF := &orderDao.GetOrderFilter{
		OrderId: params.OrderId,
		UserId:  params.UserId,
	}
	order, orderErr := orderDao.GetOrder(ctx, orderF)
	if orderErr != nil {
		err = errors.New("订单Id不能为空")
		global.ZapSugar.Errorf("[pay|Create]failed get order info,err:%+v,filter:%+v", orderErr, orderF)
		return
	}
	if order.Status != 1 {
		err = errors.New("订单状态错误")
		global.ZapSugar.Errorf("[pay|Create]failed get order info,status err.orderId:%+v", order.ID)
		return
	}

	//TODO 更新

	return
}
